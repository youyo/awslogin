# awslogin

[![Go Report Card](https://goreportcard.com/badge/github.com/youyo/awslogin)](https://goreportcard.com/report/github.com/youyo/awslogin)

## Description

Login to the AWS management console.

## Install

- Brew

```
$ brew tap youyo/tap
$ brew install awslogin
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
  -c, --cache            enable cache a credentials.
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

### Enable cache.

If you use mfa authentication, it may be difficult to authenticate each time. `--cache` option caches credentials and reuses it next time. Cache file is create to `~/.config/awslogin/cache/*` .

```bash
$ awslogin --cache --profile profile-1
Assume Role MFA token code: 000000
(open browser. require mfa.)

$ awslogin --cache --profile profile-1
(skip authentication and open browser)

$ awslogin --cache --profile profile-2
Assume Role MFA token code: 000000
(open browser. require mfa. because another profile.)
```

### Output SigninURL.

```bash
$ awslogin --output-url
https://signin.aws.amazon.com/federation?Action=...
```

## Author

[youyo](https://github.com/youyo)
