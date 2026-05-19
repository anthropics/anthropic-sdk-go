package agenttoolset

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	anthropic "github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

// SetupSkills downloads the resolved agent's skills for sessionID into
// {e.Workdir}/skills/<name>/. For each skill it fetches the files via
// client.Beta.Skills.Versions.Download and extracts the archive (a zip or
// gzip/bzip2/plain tar archive) under a directory named after the skill. Archive
// members and skill names that would escape the workspace are refused; a failure
// on one skill is logged and does not block the others. Call this before
// starting the dispatcher (e.g. right after the workdir is ready).
//
// opts are applied to every request this makes (the session lookup and each
// skill version list/get/download). Self-hosted-environment callers must pass
// the environment key here — the session and skill endpoints are
// environment-scoped, and without it the requests fall back to the client's
// default credentials and fail. option.WithAuthToken alone only ADDS an
// Authorization header; the parent client's WithAPIKey middleware still
// emits X-Api-Key on every request, so both creds would land on the wire
// and the server rejects the dual auth. Pair the bearer with an explicit
// X-Api-Key delete:
//
//	opts := []option.RequestOption{
//		option.WithHeaderDel("X-Api-Key"),
//		option.WithAuthToken(environmentKey),
//	}
//	env.SetupSkills(ctx, client, sessionID, opts...)
func (e *AgentToolContext) SetupSkills(ctx context.Context, client anthropic.Client, sessionID string, opts ...option.RequestOption) error {
	log := slog.Default().With(slog.String("component", "tool-env"), slog.String("session_id", sessionID))
	session, err := client.Beta.Sessions.Get(ctx, sessionID, anthropic.BetaSessionGetParams{}, opts...)
	if err != nil {
		return fmt.Errorf("retrieve session %s: %w", sessionID, err)
	}
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

// Cleanup removes the per-session skill downloads [AgentToolContext.SetupSkills]
// created under the workdir ({Workdir}/skills). The EnvironmentWorker calls this
// when a work item is done so one session's skills do not leak into the next
// item served by the same worker. It is a no-op when no workdir is set.
func (e *AgentToolContext) Cleanup() error {
	if e.Workdir == "" {
		return nil
	}
	return os.RemoveAll(filepath.Join(e.Workdir, "skills"))
}

func (e *AgentToolContext) downloadSkill(ctx context.Context, client anthropic.Client, skillsRoot, skillID, skillVersion string, log *slog.Logger, opts ...option.RequestOption) error {
	versionID, err := resolveSkillVersion(ctx, client, skillID, skillVersion, opts...)
	if err != nil {
		return err
	}
	version, err := client.Beta.Skills.Versions.Get(ctx, versionID, anthropic.BetaSkillVersionGetParams{SkillID: skillID}, opts...)
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
	resp, err := client.Beta.Skills.Versions.Download(ctx, versionID, anthropic.BetaSkillVersionDownloadParams{SkillID: skillID}, opts...)
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
	if err := extractSkillArchive(tmpPath, dest); err != nil {
		return fmt.Errorf("extract skill: %w", err)
	}
	log.Info("downloaded skill",
		slog.String("skill_id", skillID),
		slog.String("version", versionID),
		slog.String("dest", dest))
	return nil
}

// resolveSkillVersion resolves version to the concrete numeric timestamp the
// /v1/skills/{id}/versions/{version} endpoints require. session.agent.skills[].version
// may be an alias such as "latest", which those endpoints reject — so list the
// skill's versions and pick the newest. Numeric versions are returned unchanged.
func resolveSkillVersion(ctx context.Context, client anthropic.Client, skillID, version string, opts ...option.RequestOption) (string, error) {
	if isNumericString(version) {
		return version, nil
	}
	var newest string
	pager := client.Beta.Skills.Versions.ListAutoPaging(ctx, skillID, anthropic.BetaSkillVersionListParams{}, opts...)
	for pager.Next() {
		v := pager.Current().Version
		if isNumericString(v) && (newest == "" || numericGreater(v, newest)) {
			newest = v
		}
	}
	if err := pager.Err(); err != nil {
		return "", fmt.Errorf("list skill versions: %w", err)
	}
	if newest == "" {
		return "", fmt.Errorf("skill %q has no concrete version to resolve %q against", skillID, version)
	}
	return newest, nil
}

func isNumericString(s string) bool {
	if s == "" {
		return false
	}
	for _, r := range s {
		if r < '0' || r > '9' {
			return false
		}
	}
	return true
}

// numericGreater reports whether decimal string a is numerically greater than b.
// Both must be non-empty digit strings without leading zeros (skill versions are
// Unix-epoch timestamps), so length-then-lexical ordering matches numeric order
// without risking integer overflow on very large values.
func numericGreater(a, b string) bool {
	if len(a) != len(b) {
		return len(a) > len(b)
	}
	return a > b
}
