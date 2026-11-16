fix(environments): stop EnvironmentWorker from running tools for a terminated session

For a terminated session, which includes an archived one, the worker now force-stops the work item right after it retrieves the session, before it downloads skills or memory stores. `HandleItem` no longer returns `ErrSessionMemoryNoToken` or a `SessionMemoryError` for such a session.
