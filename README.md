# awslogin

[![CircleCI](https://circleci.com/gh/youyo/awslogin.svg?style=svg)](https://circleci.com/gh/youyo/awslogin)

## Description

Using AssumeRole, accept IAMRole and log in to the AWS management console.

## Usage

- Print Help.

```bash
$ awslogin -h
Using AssumeRole, accept IAMRole and login to the AWS management console.

Usage:
  awslogin [flags]
  awslogin [command]

Available Commands:
  help        Help about any command
  list        List profiles
  version     Show version

Flags:
  -a, --app string       Opens with the specified application.
  -h, --help             help for awslogin
  -p, --profile string   Use a specific profile.
  -e, --read-from-env    Use a specific profile read from the environment. [$AWS_PROFILE]

Use "awslogin [command] --help" for more information about a command.
```

- Login AWS management console.

```bash
$ awslogin
(open browser)
```

- Login AWS management console using a specific profile.

```bash
$ awslogin -p profile-1
(open browser)
```

- Login AWS management console using a specific profile read from the environment.

```bash
$ export AWS_PROFILE=profile-1
$ awslogin -e
(open browser)
```

- Print Arns.

```bash
$ awslogin list
test
profile-1
profile-2
```

---

## Install

- Configure AWS CLI [default profile]. http://docs.aws.amazon.com/streams/latest/dev/kinesis-tutorial-cli-installation.html#config-cli
- Configure using Assume role. http://docs.aws.amazon.com/cli/latest/userguide/cli-roles.html
- If use to MFA, please set `mfa_serial` parameter.
- Install awslogin command

```bash
$ brew tap youyo/awslogin
$ brew install awslogin
```

## Contribution

1. Fork ([https://github.com/youyo/awslogin/fork](https://github.com/youyo/awslogin/fork))
1. Create a feature branch
1. Setup Environment `make setup && make deps`
1. Write code
1. Run `gofmt -s`
1. Execute test `make test`
1. Commit your changes
1. Rebase your local changes against the master branch
1. Create a new Pull Request

## Author

[youyo](https://github.com/youyo)
