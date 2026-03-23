package cmd

import (
	"fmt"
	"io"
	"os"
)

// CompletionCmd はシェル補完スクリプトを生成するサブコマンド
type CompletionCmd struct {
	Shell  string    `arg:"" enum:"bash,zsh" help:"Shell type (bash or zsh)."`
	Writer io.Writer `kong:"-"`
}

// Run は指定されたシェル用の補完スクリプトを stdout に出力する
func (c *CompletionCmd) Run() error {
	w := c.Writer
	if w == nil {
		w = os.Stdout
	}

	switch c.Shell {
	case "bash":
		_, err := fmt.Fprint(w, bashCompletionScript)
		return err
	case "zsh":
		_, err := fmt.Fprint(w, zshCompletionScript)
		return err
	default:
		return fmt.Errorf("unsupported shell: %s", c.Shell)
	}
}

const bashCompletionScript = `_awslogin_completions() {
    local cur prev commands flags
    COMPREPLY=()
    cur="${COMP_WORDS[COMP_CWORD]}"
    prev="${COMP_WORDS[COMP_CWORD-1]}"

    commands="list version completion"
    flags="--open -o --duration -d --help"

    case "${prev}" in
        completion)
            COMPREPLY=( $(compgen -W "bash zsh" -- "${cur}") )
            return 0
            ;;
        --duration|-d)
            return 0
            ;;
    esac

    if [[ "${cur}" == -* ]]; then
        COMPREPLY=( $(compgen -W "${flags}" -- "${cur}") )
        return 0
    fi

    COMPREPLY=( $(compgen -W "${commands}" -- "${cur}") )
    return 0
}

complete -F _awslogin_completions awslogin
`

const zshCompletionScript = `#compdef awslogin

_awslogin() {
    local -a commands
    commands=(
        'list:List available AWS profiles and current session'
        'version:Show version information'
        'completion:Generate shell completion script'
    )

    _arguments \
        '(-o --open)'{-o,--open}'[Open URL in default browser]' \
        '(-d --duration)'{-d,--duration}'[Session duration in seconds]:duration' \
        '--help[Show help]' \
        '1:command:->commands' \
        '*::arg:->args'

    case "$state" in
        commands)
            _describe -t commands 'awslogin commands' commands
            ;;
        args)
            case "${words[1]}" in
                completion)
                    _values 'shell' bash zsh
                    ;;
            esac
            ;;
    esac
}

compdef _awslogin awslogin
`
