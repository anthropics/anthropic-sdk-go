package agenttoolset

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	anthropic "github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

// SetupSkillsFromSession downloads the session agent's skills into
// {e.Workdir}/skills/<name>/. It reads the resolved agent off session, and for
// each skill fetches the files via client.Beta.Skills.Versions.Download and
// extracts the archive (a zip or gzip/bzip2/plain tar archive) under a
// directory named after the skill. Archive members and skill names that would
// escape the workspace are refused; a failure on one skill is logged and does
// not block the others. Call this before starting the dispatcher (e.g. right
// after the workdir is ready).
//
// Pass the session. A session's resources cannot change while it runs, so the
// caller fetches it once and shares that snapshot with the memory-store
// download — the two can then never disagree about the attached resources.
// Callers holding only an id have the deprecated [AgentToolContext.SetupSkills].
//
// opts are applied to every request this makes (each skill version
// list/get/download). Self-hosted-environment callers must pass the
// environment key here — the skill endpoints are environment-scoped, and
// without it the requests fall back to the client's default credentials and
// fail. option.WithAuthToken alone only ADDS an Authorization header; the
// parent client's WithAPIKey middleware still emits X-Api-Key on every
// request, so both creds would land on the wire and the server rejects the
// dual auth. Pair the bearer with an explicit X-Api-Key delete:
//
//	opts := []option.RequestOption{
//		option.WithHeaderDel("X-Api-Key"),
//		option.WithAuthToken(environmentKey),
//	}
//	session, err := client.Beta.Sessions.Get(ctx, sessionID, anthropic.BetaSessionGetParams{}, opts...)
//	// ...
//	env.SetupSkillsFromSession(ctx, client, session, opts...)
func (e *AgentToolContext) SetupSkillsFromSession(ctx context.Context, client anthropic.Client, session *anthropic.BetaManagedAgentsSession, opts ...option.RequestOption) error {
	log := slog.Default().With(slog.String("component", "tool-env"), slog.String("session_id", session.ID))
	skillsRoot, err := filepath.Abs(filepath.Join(e.Workdir, "skills"))
	if err != nil {
		return fmt.Errorf("resolve skills dir: %w", err)
	}
	for _, skill := range session.Agent.Skills {
		if err := e.downloadSkill(ctx, client, skillsRoot, skill.SkillID, skill.Version, log, opts...); err != nil {
			log.Warn("failed to download skill", slog.String("skill_id", skill.SkillID), slog.Any("error", err))
		}
	}
	return nil
}

// SetupSkills is [AgentToolContext.SetupSkillsFromSession] for a caller that holds
// only the session's id: it fetches the session — opts authorize that lookup
// as well — and delegates.
//
// Deprecated: prefer [AgentToolContext.SetupSkillsFromSession]. This entry point costs an
// extra session fetch on every call; a caller that uses it
// and then downloads the session's memory stores fetches the session twice,
// and two fetches can disagree about the attached resources. It remains
// supported for callers written before the session-taking form existed.
func (e *AgentToolContext) SetupSkills(ctx context.Context, client anthropic.Client, sessionID string, opts ...option.RequestOption) error {
	slog.Default().With(slog.String("component", "tool-env"), slog.String("session_id", sessionID)).
		Warn("SetupSkills(sessionID) is deprecated and costs an extra session fetch; " +
			"fetch the session once and pass it to SetupSkillsFromSession instead")
	session, err := client.Beta.Sessions.Get(ctx, sessionID, anthropic.BetaSessionGetParams{}, opts...)
	if err != nil {
		return fmt.Errorf("retrieve session %s: %w", sessionID, err)
	}
	return e.SetupSkillsFromSession(ctx, client, session, opts...)
}

// Cleanup removes the per-skill directories
// [AgentToolContext.SetupSkillsFromSession] downloaded under {Workdir}/skills
// — exactly those, and nothing else. The skills directory itself, and
// anything else the caller placed in it, survive: the directory lives inside
// the caller's workdir and its other contents are the caller's, not ours to
// delete. (A caller directory whose name collides with a downloaded skill is
// already replaced at download time, so there is nothing of the caller's left
// for Cleanup to spare.) The EnvironmentWorker calls this when a work item is done so one
// session's skills do not leak into the next item served by the same worker.
func (e *AgentToolContext) Cleanup() error {
	var errs []error
	for _, dir := range e.downloadedSkillDirs {
		if err := os.RemoveAll(dir); err != nil {
			errs = append(errs, err)
		}
	}
	e.downloadedSkillDirs = nil
	return errors.Join(errs...)
}

func (e *AgentToolContext) downloadSkill(ctx context.Context, client anthropic.Client, skillsRoot, skillID, skillVersion string, log *slog.Logger, opts ...option.RequestOption) error {
	version, err := client.Beta.Skills.Versions.Get(ctx, skillVersion, anthropic.BetaSkillVersionGetParams{SkillID: skillID}, opts...)
	if err != nil {
		return fmt.Errorf("retrieve skill version: %w", err)
	}
	// The directory is the skill's name, reduced to a single safe path
	// component so a hostile name can't escape skillsRoot.
	dirname := filepath.Base(strings.TrimSpace(version.Name))
	if dirname == "" || dirname == "." || dirname == ".." || strings.ContainsAny(dirname, `/\`) {
		dirname = skillID
	}
	dest := filepath.Join(skillsRoot, dirname)
	if dest != skillsRoot && !strings.HasPrefix(dest, skillsRoot+string(os.PathSeparator)) {
		return fmt.Errorf("skill name %q escapes the skills dir", version.Name)
	}
	// skillVersion may be the alias "latest", which only the retrieve
	// endpoint resolves; download by the concrete id it returned.
	resp, err := client.Beta.Skills.Versions.Download(ctx, version.ID, anthropic.BetaSkillVersionDownloadParams{SkillID: skillID}, opts...)
	if err != nil {
		return fmt.Errorf("download skill: %w", err)
	}
	defer resp.Body.Close()

	// Stream the archive to a temp file rather than buffering it whole in
	// memory: a skill bundle can be large, and the zip extractor needs random
	// access over the file anyway.
	tmp, err := os.CreateTemp("", "skill-archive-*")
	if err != nil {
		return fmt.Errorf("create temp file for skill archive: %w", err)
	}
	tmpPath := tmp.Name()
	defer os.Remove(tmpPath)
	if _, err := io.Copy(tmp, resp.Body); err != nil {
		tmp.Close()
		return fmt.Errorf("stream skill archive to disk: %w", err)
	}
	if err := tmp.Close(); err != nil {
		return fmt.Errorf("flush skill archive: %w", err)
	}

	if err := os.RemoveAll(dest); err != nil {
		return fmt.Errorf("clear %s: %w", dest, err)
	}
	// Recorded before the extract, like a successful one: even a
	// half-extracted directory is this call's to remove, not the caller's.
	e.downloadedSkillDirs = append(e.downloadedSkillDirs, dest)
	if err := extractSkillArchive(tmpPath, dest); err != nil {
		return fmt.Errorf("extract skill: %w", err)
	}
	log.Info("downloaded skill",
		slog.String("skill_id", skillID),
		slog.String("version", version.ID),
		slog.String("dest", dest))
	return nil
}
