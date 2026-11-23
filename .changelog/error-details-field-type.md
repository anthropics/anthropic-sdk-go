fix(api): change the type of ErrorDetailsUnion.Field to string

To compare it with a `WorkspacePolicyField` value, use `AsWorkspacePolicyViolation().Field`. The beta methods `Workspaces.Get`, `Workspaces.Archive` and `ComplianceSettings.Get` now take a params argument.
