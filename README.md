# awslogin

[![Coverage Status](https://coveralls.io/repos/github/youyo/awslogin/badge.svg?branch=master)](https://coveralls.io/github/youyo/awslogin?branch=master)

## Description

Using AssumeRole, accept IAMRole and log in to the AWS management console.

## Usage

- Print Help.

```bash
$ awslogin -h
Usage of awslogin:
  -l	Print available ARN list and quit. (Short)
  -list
    	Print available ARN list and quit.
  -output
    	Print login url and quiet. It will not login automatically.
  -p string
    	Use Profile name. (Short)
  -profile string
    	Use Profile name
  -v	Print version information and quit. (Short)
  -version
    	Print version information and quit.
```

- Print Arns.

```bash
$ awslogin -list
test
```

- Login AWS management console.

```bash
$ awslogin -profile test
(open browser)
```

---


### For Zsh

Used with [peco](https://github.com/peco/peco).  
First install peco.  
Next install awslogin.  
Write to your `~/.zshrc` file.

```zsh
function al-src () {
    local selected_arn=$(awslogin -list | peco --query "$LBUFFER")
    BUFFER="awslogin -profile ${selected_arn}"
    zle accept-line
    zle clear-screen
}
zle -N al-src
bindkey '+_' al-src
```

Press '+_', you can select arn.

### For Bash

Used with bash-completion.  
First install bash-completion.  
Second install awslogin.  
Last create a config  file to /usr/local/etc/bash_completion.d/awslogin (For MacOS).

```bash
#!bash

_awslogin()
{
  local cur=${COMP_WORDS[COMP_CWORD]}
  CANDIDATE=`egrep "^\[profile" ~/.aws/config |perl -pe 's/]\n/ /g;s/\[profile//'`
  COMPREPLY=( $(compgen -W "$CANDIDATE" -- $cur) )
}

complete -F _awslogin awslogin
```

## Install

- Configure AWS CLI [default profile]. http://docs.aws.amazon.com/streams/latest/dev/kinesis-tutorial-cli-installation.html#config-cli
- Configure using Assume role. http://docs.aws.amazon.com/cli/latest/userguide/cli-roles.html
- If use to MFA, please set `mfa_serial` parameter.
- Install awslogin command

```bash
$ curl -sLO `curl -s 'http://grl.i-o.sh/youyo/awslogin?suffix=darwin_amd64.zip'`
$ unzip awslogin_darwin_amd64.zip
```

## Contribution

1. Fork ([https://github.com/youyo/awslogin/fork](https://github.com/youyo/awslogin/fork))
1. Create a feature branch
1. Setup Environment `make deps`
1. Write code
1. Run `gofmt -s`
1. Execute test `make test`
1. Commit your changes
1. Rebase your local changes against the master branch
1. Create a new Pull Request

## Author

[youyo](https://github.com/youyo)
