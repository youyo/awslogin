# awslogin

## Description

Using AssumeRole, accept IAMRole and log in to the AWS management console.

## Usage

- Print Help.

```bash
$ awslogin -h
Usage:
  -l    Print available ARN list and quit. (Short)
  -list
    	Print available ARN list and quit.
  -r string
    	Use IAM-Role. (Short)
  -rolename string
    	Use IAM-Role.
  -v    Print version information and quit. (Short)
  -version
    	Print version information and quit.
```

- Print Arns.

```bash
$ awslogin -l
test
```

- Login AWS management console.

```bash
$ awslogin -r test
(open browser)
```

---

- Recommended usage.

Used with [peco](https://github.com/peco/peco).  
First install peco.  
Next install awslogin.  
Write to your `~/.zshrc` file.  

```zsh
function al-src () {
    local selected_arn=$(awslogin -l | peco --query "$LBUFFER")
    BUFFER="awslogin -r ${selected_arn}"
    zle accept-line
    zle clear-screen
}
zle -N al-src
bindkey '+_' al-src
```

Press '+_', you can select arn.


## Install

- Configure AWS CLI [default profile]. http://docs.aws.amazon.com/streams/latest/dev/kinesis-tutorial-cli-installation.html#config-cli
- Configure using Assume role. http://docs.aws.amazon.com/cli/latest/userguide/cli-roles.html
- Install awslogin command

```bash
$ wget `curl -s https://api.github.com/repos/youyo/awslogin/releases/latest | jq -r ".assets[0].browser_download_url"`
$ unzip awslogin_darwin_amd64.zip
```

## Contribution

1. Fork ([https://github.com/youyo/awslogin/fork](https://github.com/youyo/awslogin/fork))
1. Create a feature branch
1. Commit your changes
1. Rebase your local changes against the master branch
1. Run test suite with the `go test ./...` command and confirm that it passes
1. Run `gofmt -s`
1. Create a new Pull Request

## Author

[youyo](https://github.com/youyo)
