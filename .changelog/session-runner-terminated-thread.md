fix(sessions): stop the session tool runner from running the tool calls of a terminated session thread

Once a `session.thread_status_terminated` event names a child session thread, the runner also no longer yields that thread's calls, including a denied call (`Confirmation="deny"`) and a call to a tool that it does not own. It no longer waits for a confirmation of any of that thread's calls either. If a tool is already running for one of the thread's calls when the session thread terminates, the runner cancels the tool's context, still yields the call with `Posted=false` and posts no result.
