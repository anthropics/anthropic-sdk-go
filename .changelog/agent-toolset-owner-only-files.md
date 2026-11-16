fix(tools): make agent toolset files and directories owner-only

The `write` tool now creates files as `0600` and directories as `0700` (they were `0644` and `0755`), and a downloaded skill's directory is `0700`. A process that runs as another user, such as a sidecar or a web server, can no longer read them. A file that `write` or `edit` replaces now keeps its permission bits, where before it was reset to `0644`. To share a file with another user, create it with the permissions you want before the agent writes to it, or run `chmod` afterwards.
