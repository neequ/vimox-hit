# Troubleshooting

Run `vimox-hit doctor` to check API connectivity, then use `vimox-hit ai-log list` to locate the request. `show` reads details visible to the current account; `diagnose` creates a redacted diagnostic bundle; `copy` places the same bundle on the clipboard. Credentials and signed URL query parameters are removed.

If `doctor` reports `authenticated: false`, project and AI log commands cannot run until CLI OAuth authorization is available. Use the Hit web app for account data; do not copy access tokens from browser storage.

Retry read-only queries a limited number of times. Before retrying a task, check its server-side state. Every model operation and final video generation requires explicit user approval. Request new approval after a retry or any change to the model, inputs, parameters, or estimated cost.

See the [Simplified Chinese troubleshooting guide](zh-CN/troubleshooting.md).
