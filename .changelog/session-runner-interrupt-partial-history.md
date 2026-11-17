fix(sessions): stop the session tool runner from waiting for the confirmation of an interrupted tool call when listing the session's events fails partway

After such a failure, an interrupt no longer starts the `MaxIdle` countdown from an idle that came before the failure. The runner waits for the next idle.
