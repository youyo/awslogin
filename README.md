# awslogin

[![Go Report Card](https://goreportcard.com/badge/github.com/youyo/awslogin)](https://goreportcard.com/report/github.com/youyo/awslogin)

## Description

Login to the AWS management console.

## Install

- Brew

```
$ brew install youyo/tap/awslogin
```

Other platforms are download from [github release page](https://github.com/youyo/awslogin/releases).

## Usage

```bash
$ awslogin
Login to the AWS management console.

Usage:
  awslogin [flags]

Flags:
  -b, --browser string   Opens with the specified browse application
  -h, --help             help for awslogin
  -O, --output-url       output signin url
  -p, --profile string   use a specific profile from your credential file. (default "default")
  -S, --select-profile   interactive select profile
      --version          version for awslogin
```

### Login AWS management console.

```bash
$ awslogin
(open browser using default profile or $AWS_PROFILE)
```

### Login AWS management console using a specific profile.

```bash
$ awslogin --profile profile-1
(open browser using selected profile)
```

### Login AWS management console using interactive select.

```bash
$ awslogin --select-profile
(open browser using selected profile)
```

### Output SigninURL.

```bash
$ awslogin --output-url
https://signin.aws.amazon.com/federation?Action=...
```

## Author

[youyo](https://github.com/youyo)
