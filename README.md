# awslogin

[![Test](https://github.com/youyo/awslogin/actions/workflows/test.yml/badge.svg)](https://github.com/youyo/awslogin/actions/workflows/test.yml)
[![Lint](https://github.com/youyo/awslogin/actions/workflows/lint.yml/badge.svg)](https://github.com/youyo/awslogin/actions/workflows/lint.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/youyo/awslogin)](https://goreportcard.com/report/github.com/youyo/awslogin)
[![Release](https://img.shields.io/github/v/release/youyo/awslogin)](https://github.com/youyo/awslogin/releases/latest)

[日本語](README.ja.md)

A CLI tool that generates AWS Management Console sign-in URLs from temporary credentials.

## Features

- Generate console sign-in URLs from AWS temporary credentials
- Open the console directly in your default browser with `--open`
- Customize session duration with `--duration`
- Configure defaults with environment variables (`AWSLOGIN_DURATION`, `AWSLOGIN_OPEN`)
- Shell completion for bash and zsh
- Cross-platform: macOS, Linux, Windows (amd64/arm64)
- JSON output on stdout/stderr — optimized for scripting and AI coding agents
- List available AWS profiles with `awslogin list`

## Install

### Homebrew

```bash
brew install youyo/tap/awslogin
```

### go install

```bash
go install github.com/youyo/awslogin@latest
```

### GitHub Releases

Download a binary for your OS and architecture from the [Releases page](https://github.com/youyo/awslogin/releases).

## Quick Start

```bash
# Get a sign-in URL using a named profile
AWS_PROFILE=myprofile awslogin
# {"result":{"url":"https://signin.aws.amazon.com/federation?...","region":"ap-northeast-1","opened_in_browser":false}}

# Open it in a browser
AWS_PROFILE=myprofile awslogin --open
# {"result":{"url":"https://signin.aws.amazon.com/federation?...","region":"ap-northeast-1","opened_in_browser":true}}
```

## Usage

### Generate a sign-in URL (default)

Prints JSON to stdout with the sign-in URL. Pipe it, extract it with jq, do whatever you want.

```bash
awslogin
# {"result":{"url":"https://signin.aws.amazon.com/federation?...","region":"ap-northeast-1","opened_in_browser":false}}

# Extract URL with jq
awslogin | jq -r '.result.url'

# Copy URL to clipboard on macOS
awslogin | jq -r '.result.url' | pbcopy
```

### Open in browser (`--open` / `-o`)

```bash
awslogin --open
awslogin -o
```

### Set session duration (`--duration` / `-d`)

Default is 3600 seconds (1 hour).

```bash
awslogin --duration 7200   # 2 hours
awslogin -d 7200
```

### Switch AWS profile

Use the `AWS_PROFILE` environment variable, same as the AWS CLI.

```bash
AWS_PROFILE=production awslogin
AWS_PROFILE=staging awslogin -o
```

### Environment variables

Set defaults so you don't have to pass the same flags every time.

| Variable | Description | Example |
|----------|-------------|---------|
| `AWSLOGIN_DURATION` | Session duration in seconds (900-43200) | `export AWSLOGIN_DURATION=7200` |
| `AWSLOGIN_OPEN` | Open URL in browser (`true`/`false`) | `export AWSLOGIN_OPEN=true` |

Command-line flags always take precedence over environment variables.

```bash
# Always use 2-hour sessions and open in browser
export AWSLOGIN_DURATION=7200
export AWSLOGIN_OPEN=true
awslogin

# Override for a one-off
awslogin -d 900
```

### List profiles (`list`)

Shows all configured AWS profiles and any active session from environment variables.

```bash
awslogin list
# {"result":{"profiles":[{"name":"dev","type":"sso","sso_start_url":"https://...","region":"ap-northeast-1"},{"name":"prod","type":"credentials","region":"us-east-1"}],"current_session":null}}
```

### Show version

```bash
awslogin version
# {"result":{"version":"v3.2.1"}}
```

### Shell completion

```bash
# zsh
eval "$(awslogin completion zsh)"

# bash
eval "$(awslogin completion bash)"
```

Add the line to your `~/.zshrc` or `~/.bashrc` to persist it.

## SSO Profile Support

awslogin supports AWS SSO profiles configured with the modern `sso-session` format.

When your SSO session has expired, awslogin automatically detects the `InvalidGrantException` and starts the OIDC device authorization flow:

1. A browser window opens automatically
2. An authorization code is displayed — confirm it in the browser
3. After successful authentication, awslogin retries and generates the console URL

**Only the modern `[sso-session]` format is supported.** Legacy profiles with a bare `sso_start_url` key will receive a migration error.

## JSON Output Format

All commands (except `completion`) output JSON to stdout. Events and progress are output as NDJSON to stderr.

### stdout (result)

```json
{"result": {"url": "...", "region": "...", "opened_in_browser": false}}
```

### stdout (error, exit code 1)

```json
{"error": {"code": "SSO_SESSION_EXPIRED", "message": "SSO session expired", "details": "..."}}
```

### stderr (events, NDJSON)

```json
{"type": "sso_auth_required", "verification_code": "ABCD-EFGH", "verification_url": "https://..."}
```

### Example `~/.aws/config`

```ini
[profile my-sso]
sso_session = my-sso
sso_account_id = 123456789012
sso_role_name = AdministratorAccess
region = ap-northeast-1

[sso-session my-sso]
sso_start_url = https://my-org.awsapps.com/start
sso_region = ap-northeast-1
sso_registration_scopes = sso:account:access
```

```bash
# First run or after session expiry: browser opens, then JSON output
AWS_PROFILE=my-sso awslogin
# stderr: {"type":"sso_session_expired","message":"SSO session expired. Starting SSO login..."}
# stderr: {"type":"sso_auth_required","verification_code":"ABCD-EFGH","verification_url":"https://..."}
# stderr: {"type":"sso_auth_complete"}
# stdout: {"result":{"url":"https://signin.aws.amazon.com/federation?...","region":"ap-northeast-1","opened_in_browser":false}}
```

## Migrating from v2

v3.0.0 includes breaking changes.

| v2 | v3 | Why |
|----|-----|-----|
| Opens browser by default | Prints URL to stdout by default | Easier to compose with pipes and scripts |
| `--output-url` (`-O`) to print URL | Default behavior (no flag needed) | URL output is the primary use case |
| `--profile` (`-p`) | `AWS_PROFILE` env var | Follows the AWS SDK convention |
| `--select-profile` (`-S`) | Removed | Interactive profile picker dropped |
| `--browser` (`-b`) | Removed | Only the default browser is supported |
| `--version` flag | `awslogin version` subcommand | Matches the Kong CLI framework convention |
| Plain text output | JSON output | Machine-readable, optimized for scripting and AI agents |

### What changed under the hood

- **CLI framework**: Cobra + Viper replaced with [Kong](https://github.com/alecthomas/kong)
- **AWS SDK**: v1 replaced with v2
- **MFA/SSO**: Delegated to the AWS SDK v2 credential chain (custom implementation removed)
- **Shell completion**: Static file (`_awslogin`) replaced with `awslogin completion` subcommand

## Development

```bash
go build -o awslogin .
go test ./...
golangci-lint run
```

## License

[MIT](LICENSE)

## Author

[youyo](https://github.com/youyo)
