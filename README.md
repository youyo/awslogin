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
- Shell completion for bash and zsh
- Cross-platform: macOS, Linux, Windows (amd64/arm64)

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

# Open it in a browser
AWS_PROFILE=myprofile awslogin --open
```

## Usage

### Generate a sign-in URL (default)

Prints the URL to stdout. Pipe it, copy it, do whatever you want.

```bash
awslogin
awslogin | pbcopy  # copy to clipboard on macOS
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

### Show version

```bash
awslogin version
```

### Shell completion

```bash
# zsh
eval "$(awslogin completion zsh)"

# bash
eval "$(awslogin completion bash)"
```

Add the line to your `~/.zshrc` or `~/.bashrc` to persist it.

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
