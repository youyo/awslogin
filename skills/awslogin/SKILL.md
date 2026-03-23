---
name: awslogin
description: LLM-first CLI for AWS console login. Generates federated login URLs via SSO/OIDC with full JSON output optimized for coding agents. Profile discovery, automated SSO authentication, session management, and browser integration all native to the CLI.
triggers:
  - awslogin
  - AWS login
  - AWS console
  - AWS federation
  - SSO login
---

# awslogin: LLM-First AWS Console Login CLI

`awslogin` is an LLM-first CLI for AWS console federation.

It is **not** a thin AWS SDK wrapper. Its main job is to turn AWS credential resolution and OIDC authentication into **stable, machine-readable JSON output** that works well for Claude Code, Codex, and other coding agents.

Full command:

```bash
awslogin
```

Default output format:

```bash
json
```

Supported output formats:

```bash
json
```

(JSON is the only output format. No alternate formats.)

Use `awslogin login` (implicit) to get federated console URLs. Use `awslogin list` for profile discovery and session context. Use `awslogin version` for diagnostic info.

---

## When to use this skill

Use `awslogin` when you need to:

- get a federated AWS console login URL
- discover available AWS profiles and their types (SSO, legacy credentials)
- check the current authenticated session
- trigger SSO device authorization flows for unattended/automated scenarios
- integrate AWS authentication into agent workflows (no manual browser steps)

Prefer `awslogin` over AWS console login web UI or aws-vault when:

- you need **LLM-friendly JSON output**
- you need **automatic OIDC device flow support**
- you want **non-interactive authentication** for agents
- you need **stderr event streaming** for real-time auth status

---

## Safety model

`awslogin` is intentionally minimal and read-heavy.

Supported operations:

- login (generate console URL)
- list (discover profiles)
- version (diagnostic info)
- completion (shell integration)

Explicitly unsupported:

- credential creation
- credential rotation
- secret injection
- profile modification
- environment manipulation (no `eval` output)

No credentials are ever printed to stdout. Credentials are resolved internally and used only for URL signing.

---

## Auth and configuration

Configuration:

```text
~/.aws/config       (profile definitions, SSO config)
~/.aws/credentials  (legacy credentials if present)
```

Token store:

```text
~/.aws/sso/cache/   (SSO session tokens, auto-managed by AWS SDK)
```

Resolution priority:

```text
AWS_PROFILE env var > default profile in config
```

Secrets are never stored in project-local files. Credentials are always read from `~/.aws/` or environment.

Primary auth modes:

- AWS SSO (via `~/.aws/config` with `[sso-session]` and `sso_start_url`)
- OIDC (auto-fallback when SSO session expires or is unavailable)
- Legacy IAM credentials (via `~/.aws/credentials`, not recommended)

Useful environment variables:

```text
AWS_PROFILE              (which profile to use; default: profile in config)
AWSLOGIN_OPEN           (set to true to auto-open browser; default: false)
AWSLOGIN_DURATION       (session duration in seconds; min 900, max 43200; default 3600)
AWS_REGION              (session region; overrides profile region)
AWS_DEFAULT_REGION      (fallback region)
AWS_SSO_CACHE_DIR       (advanced: override SSO token cache location)
```

**Important:** Do not pass `AWS_SECRET_ACCESS_KEY` or `AWS_ACCESS_KEY_ID` to `awslogin`. The SDK resolves these automatically from `~/.aws/credentials` or environment.

---

## Global behavior

Credentials are always resolved from:

1. Environment (`AWS_ACCESS_KEY_ID`, `AWS_SECRET_ACCESS_KEY`)
2. `~/.aws/credentials` file
3. SSO session (via `~/.aws/config`)

Profile selection:

- Use `AWS_PROFILE=dev awslogin` to select a non-default profile.
- No `--profile` flag exists. Profile selection is always via `AWS_PROFILE` environment variable or `[default]` section in `~/.aws/config`.

Region selection:

- Profile region from `~/.aws/config` is used by default.
- Override with `AWS_REGION=us-west-2 awslogin` or `AWS_DEFAULT_REGION`.

Output conventions:

- `stdout`: primary JSON result only, machine-readable.
- `stderr`: event stream (NDJSON), warnings, and diagnostic logs. **Do not ignore stderr.** SSO authentication events (like device flow URLs) are emitted here.

---

## Output contract

### Success (exit code 0)

All successful commands emit compact JSON to stdout:

```json
{
  "result": {
    "url": "https://signin.aws.amazon.com/federation?...",
    "region": "ap-northeast-1",
    "opened_in_browser": false
  }
}
```

### Error (exit code 1)

On complete failure, commands emit an error JSON shape to stdout:

```json
{
  "error": {
    "code": "SSO_SESSION_EXPIRED",
    "message": "SSO session expired. Starting OIDC device authorization flow...",
    "details": "Access token missing or invalid"
  }
}
```

### stderr: Event stream (NDJSON)

Real-time events are emitted to stderr as newline-delimited JSON:

```
{"type":"sso_session_expired","message":"SSO session expired. Starting SSO login..."}
{"type":"sso_auth_required","verification_code":"ABCD-EFGH","verification_url":"https://device.amazonawsdev.com/?user_code=ABCD-EFGH","verification_url_base":"https://device.amazonawsdev.com/"}
{"type":"browser_opened","url":"https://device.amazonawsdev.com/?user_code=ABCD-EFGH"}
{"type":"sso_auth_complete","message":"SSO authentication complete. Fetching credentials..."}
```

**Important:** Agents must watch stderr and react to these events. For example:
- On `sso_auth_required`, agents should open the `verification_url` in a browser (using browser-use, playwright, or similar).
- awslogin will poll automatically until `sso_auth_complete` is emitted, then stdout will contain the login URL.

---

## Error codes

| Code | Meaning | Retryable | Action |
|------|---------|-----------|--------|
| `INVALID_ARGS` | Invalid CLI arguments | No | Check command syntax |
| `CONFIG_LOAD_FAILED` | AWS config or credentials file read failed | No | Check file permissions |
| `SSO_SESSION_EXPIRED` | SSO session token is missing or expired | Yes | Re-run; triggers OIDC device flow |
| `SSO_LOGIN_FAILED` | SSO authentication failed (user denied, etc.) | No | Check device authorization URL |
| `SSO_LEGACY_CONFIG` | Profile uses old SSO format (not supported) | No | Update `~/.aws/config` to use `sso_start_url` |
| `OIDC_REGISTER_FAILED` | OIDC client registration failed | No | Check network; retry after delay |
| `OIDC_DEVICE_AUTH_FAILED` | OIDC device auth flow failed | Yes | User denied authorization or timeout |
| `OIDC_TOKEN_EXPIRED` | OIDC token expired during flow | Yes | Re-run; restarts device flow |
| `FEDERATION_API_FAILED` | AWS federation API call failed | Yes | Check AWS account/role permissions |
| `BROWSER_OPEN_FAILED` | Could not auto-open browser (non-critical) | No | User must open URL manually (printed to stderr) |
| `INTERNAL_ERROR` | Unknown internal error | No | Check logs; file bug report |

---

## Command guide

### login (default)

Generate a federated AWS console URL.

```bash
awslogin
```

or explicitly:

```bash
awslogin login
```

This resolves the current AWS profile, authenticates via SSO/OIDC if needed, and returns a console login URL.

Example output:

```json
{
  "result": {
    "url": "https://signin.aws.amazon.com/federation?Action=login&Issuer=awslogin&Destination=https%3A%2F%2Fconsole.aws.amazon.com%2F&SigninToken=...",
    "region": "ap-northeast-1",
    "opened_in_browser": false
  }
}
```

**Behavior:**

- If SSO session is valid, authentication is instant.
- If SSO session has expired, `awslogin` automatically initiates OIDC device authorization flow and emits `sso_auth_required` event to stderr.
- Polling continues automatically until auth is complete or times out.
- Once authenticated, stdout contains the login URL.

**Useful variants:**

```bash
# Auto-open browser after URL generation (if supported)
AWSLOGIN_OPEN=true awslogin

# Use a specific profile
AWS_PROFILE=prod awslogin

# Override region
AWS_REGION=eu-west-1 awslogin

# Extend session duration (seconds, max 43200 = 12 hours)
AWSLOGIN_DURATION=7200 awslogin
```

---

### list

Discover available AWS profiles and show current session context.

```bash
awslogin list
```

Example output:

```json
{
  "result": {
    "profiles": [
      {
        "name": "default",
        "type": "sso",
        "sso_start_url": "https://my-sso-portal.awsapps.com/start/",
        "region": "ap-northeast-1"
      },
      {
        "name": "dev",
        "type": "sso",
        "sso_start_url": "https://my-sso-portal.awsapps.com/start/",
        "region": "us-east-1"
      },
      {
        "name": "legacy",
        "type": "iam",
        "region": "ap-northeast-1"
      }
    ],
    "current_session": {
      "aws_access_key_id": "ASIA...",
      "has_session_token": true,
      "region": "us-east-1",
      "source": "sso"
    }
  }
}
```

Use this command when:

- you need to discover which profiles are available
- you need to check the current authenticated profile
- an agent needs to present profile choices to a user
- you want to verify SSO session state before running `login`

---

### version

Show version and diagnostic info.

```bash
awslogin version
```

Example output:

```json
{
  "result": {
    "version": "v3.2.1"
  }
}
```

Use for:

- CI/CD diagnostic logs
- version checking before integration
- bug report context

---

### completion

Generate shell completion.

```bash
awslogin completion bash
awslogin completion zsh
```

Example (zsh):

```bash
autoload -Uz compinit && compinit
eval "$(awslogin completion zsh)"
```

**Important:** completion output is not JSON. It is shell code designed for `eval`. Do not JSON-parse it.

---

## Recommended patterns for coding agents

### 1. One-shot login for a known profile

Simple case: get console URL for the current or default profile.

```bash
awslogin
```

Parse JSON from stdout.

```json
{
  "result": {
    "url": "https://signin.aws.amazon.com/federation?...",
    "opened_in_browser": false
  }
}
```

If `opened_in_browser` is `false`, agent should open the URL using browser-use or similar.

### 2. Profile selection flow

When user has multiple profiles:

1. Query available profiles:

```bash
awslogin list
```

2. Present profile choices to user (from `result.profiles[]`).

3. Run login with selected profile:

```bash
AWS_PROFILE=prod awslogin
```

4. Return the `result.url` and auto-open (or show to user).

### 3. Handle SSO auth events from stderr

awslogin emits structured events to stderr during SSO authentication.

Recommended agent pattern:

1. Run `awslogin` and stream both stdout and stderr.
2. Watch stderr for NDJSON events.
3. On `sso_auth_required` event:
   - Extract `verification_url` from the event
   - Open it in browser (using browser-use skill or similar)
   - Do not block; let awslogin poll automatically
4. Monitor until `sso_auth_complete` event arrives
5. Read stdout JSON to get the console URL

Example stderr stream:

```
{"type":"sso_session_expired","message":"SSO session expired..."}
{"type":"sso_auth_required","verification_code":"ABCD-EFGH","verification_url":"https://..."}
{"type":"browser_opened","url":"https://..."}
{"type":"sso_auth_complete"}
```

stdout arrives after step 4 with the login URL.

### 4. Region and duration customization

```bash
AWS_PROFILE=dev AWS_REGION=eu-west-1 AWSLOGIN_DURATION=7200 awslogin
```

Use when agent knows the desired region or session duration.

### 5. Error recovery

If exit code is 1, parse the error JSON:

```bash
awslogin 2>/dev/null || true
```

Check `error.code`:

- `SSO_SESSION_EXPIRED` (retryable): auto-retry once
- `OIDC_DEVICE_AUTH_FAILED` (retryable): ask user to re-authenticate
- `BROWSER_OPEN_FAILED` (retryable): open URL manually or show to user
- Other codes (non-retryable): display error to user

---

## Anti-patterns

Avoid these:

- **Ignoring stderr.** SSO auth events are critical; agents must watch stderr.
- **Using `--profile` flag.** awslogin uses `AWS_PROFILE` env var only; there is no `--profile` flag.
- **Parsing `--help` output.** Help is human-readable; use this SKILL.md instead.
- **Assuming auto-browser opening.** Default is `AWSLOGIN_OPEN=false`; agents must handle URL opening via browser skill.
- **Passing credentials as arguments.** Credentials are resolved from config/environment only; never pass them via CLI args.
- **Using `aws cli` for federation.** `awslogin` is purpose-built for federation; aws CLI is not suitable for agent workflows.
- **Expecting non-JSON output for login/list/version.** JSON is the only format; use JSON parser.
- **Waiting indefinitely for SSO auth.** Set timeout (e.g., 5 minutes) on the awslogin process; user may deny or ignore the auth request.

---

## Design notes for agents

### Event-driven auth

awslogin uses stderr NDJSON for real-time event streaming. This allows agents to react immediately to auth state changes (e.g., opening a browser as soon as the verification URL is ready) without polling.

### No credential leakage

Credentials are never printed. The CLI internally signs AWS federation URLs using resolved credentials, but only the signed URL (safe to share) is output. Agents can safely log or display the URL.

### SSO-first, OIDC fallback

awslogin prefers SSO (via `~/.aws/config`). When SSO session expires, it automatically falls back to OIDC device authorization flow. This is transparent to agents.

### Session duration control

Agents can control console session duration via `AWSLOGIN_DURATION` (default 3600s = 1 hour, max 43200s = 12 hours). Useful for long-running tasks or short-lived debugging sessions.

### Region independence

awslogin respects AWS region config but always federates to the main console. Region is metadata in the response; the login URL itself is global.

---

## Minimal command set to remember

If you only remember a few commands, remember these:

```bash
awslogin                              # get console URL for current profile
AWS_PROFILE=prod awslogin              # get console URL for specific profile
awslogin list                          # discover profiles and current session
awslogin version                       # show version info
```

Monitor stderr for SSO auth events and open URLs in browser when needed.

