/*
 Licensed under the Apache License, Version 2.0 (the "License");
 you may not use this file except in compliance with the License.
 You may obtain a copy of the License at

     https://www.apache.org/licenses/LICENSE-2.0

 Unless required by applicable law or agreed to in writing, software
 distributed under the License is distributed on an "AS IS" BASIS,
 WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 See the License for the specific language governing permissions and
 limitations under the License.
*/

package completion

import (
	"bytes"
	"fmt"

	"github.com/spf13/cobra"
)

const completionDesc = `
Generate autocompletion script for airshipctl for the specified shell (bash or zsh).

This command can generate shell autocompletion. e.g.

	$ airshipctl completion bash

Can be sourced as such

	$ source <(airshipctl completion bash)
`

var (
	completionShells = map[string]func(cmd *cobra.Command) error{
		"bash": runCompletionBash,
		"zsh":  runCompletionZsh,
	}
)

func NewCompletionCommand() *cobra.Command {
	shells := make([]string, 0, len(completionShells))
	for s := range completionShells {
		shells = append(shells, s)
	}

	cmd := &cobra.Command{
		Use:       "completion SHELL",
		Short:     "Generate autocompletions script for the specified shell (bash or zsh)",
		Long:      completionDesc,
		Args:      cobra.ExactArgs(1),
		RunE:      runCompletion,
		ValidArgs: shells,
	}

	return cmd
}

func runCompletion(cmd *cobra.Command, args []string) error {
	run, found := completionShells[args[0]]
	if !found {
		return fmt.Errorf("unsupported shell type %q", args[0])
	}

	return run(cmd)
}

func runCompletionBash(cmd *cobra.Command) error {
	return cmd.Root().GenBashCompletion(cmd.OutOrStdout())
}

func runCompletionZsh(cmd *cobra.Command) error {
	out := cmd.OutOrStdout()
	zshInitialization := `#compdef airshipctl

__airshipctl_bash_source() {
	alias shopt=':'
	alias _expand=_bash_expand
	alias _complete=_bash_comp
	emulate -L sh
	setopt kshglob noshglob braceexpand
	source "$@"
}
__airshipctl_type() {
	# -t is not supported by zsh
	if [ "$1" == "-t" ]; then
		shift
		# fake Bash 4 to disable "complete -o nospace". Instead
		# "compopt +-o nospace" is used in the code to toggle trailing
		# spaces. We don't support that, but leave trailing spaces on
		# all the time
		if [ "$1" = "__airshipctl_compopt" ]; then
			echo builtin
			return 0
		fi
	fi
	type "$@"
}
__airshipctl_compgen() {
	local completions w
	completions=( $(compgen "$@") ) || return $?
	# filter by given word as prefix
	while [[ "$1" = -* && "$1" != -- ]]; do
		shift
		shift
	done
	if [[ "$1" == -- ]]; then
		shift
	fi
	for w in "${completions[@]}"; do
		if [[ "${w}" = "$1"* ]]; then
			echo "${w}"
		fi
	done
}
__airshipctl_compopt() {
	true # don't do anything. Not supported by bashcompinit in zsh
}
__airshipctl_declare() {
	if [ "$1" == "-F" ]; then
		whence -w "$@"
	else
		builtin declare "$@"
	fi
}
__airshipctl_ltrim_colon_completions()
{
	if [[ "$1" == *:* && "$COMP_WORDBREAKS" == *:* ]]; then
		# Remove colon-word prefix from COMPREPLY items
		local colon_word=${1%${1##*:}}
		local i=${#COMPREPLY[*]}
		while [[ $((--i)) -ge 0 ]]; do
			COMPREPLY[$i]=${COMPREPLY[$i]#"$colon_word"}
		done
	fi
}
__airshipctl_get_comp_words_by_ref() {
	cur="${COMP_WORDS[COMP_CWORD]}"
	prev="${COMP_WORDS[${COMP_CWORD}-1]}"
	words=("${COMP_WORDS[@]}")
	cword=("${COMP_CWORD[@]}")
}
__airshipctl_filedir() {
	local RET OLD_IFS w qw
	__debug "_filedir $@ cur=$cur"
	if [[ "$1" = \~* ]]; then
		# somehow does not work. Maybe, zsh does not call this at all
		eval echo "$1"
		return 0
	fi
	OLD_IFS="$IFS"
	IFS=$'\n'
	if [ "$1" = "-d" ]; then
		shift
		RET=( $(compgen -d) )
	else
		RET=( $(compgen -f) )
	fi
	IFS="$OLD_IFS"
	IFS="," __debug "RET=${RET[@]} len=${#RET[@]}"
	for w in ${RET[@]}; do
		if [[ ! "${w}" = "${cur}"* ]]; then
			continue
		fi
		if eval "[[ \"\${w}\" = *.$1 || -d \"\${w}\" ]]"; then
			qw="$(__airshipctl_quote "${w}")"
			if [ -d "${w}" ]; then
				COMPREPLY+=("${qw}/")
			else
				COMPREPLY+=("${qw}")
			fi
		fi
	done
}
__airshipctl_quote() {
	if [[ $1 == \'* || $1 == \"* ]]; then
		# Leave out first character
		printf %q "${1:1}"
	else
		printf %q "$1"
	fi
}
autoload -U +X bashcompinit && bashcompinit
# use word boundary patterns for BSD or GNU sed
LWORD='[[:<:]]'
RWORD='[[:>:]]'
if sed --help 2>&1 | grep -q GNU; then
	LWORD='\<'
	RWORD='\>'
fi
__airshipctl_convert_bash_to_zsh() {
	sed \
	-e 's/declare -F/whence -w/' \
	-e 's/_get_comp_words_by_ref "\$@"/_get_comp_words_by_ref "\$*"/' \
	-e 's/local \([a-zA-Z0-9_]*\)=/local \1; \1=/' \
	-e 's/flags+=("\(--.*\)=")/flags+=("\1"); two_word_flags+=("\1")/' \
	-e 's/must_have_one_flag+=("\(--.*\)=")/must_have_one_flag+=("\1")/' \
	-e "s/${LWORD}_filedir${RWORD}/__airshipctl_filedir/g" \
	-e "s/${LWORD}_get_comp_words_by_ref${RWORD}/__airshipctl_get_comp_words_by_ref/g" \
	-e "s/${LWORD}__ltrim_colon_completions${RWORD}/__airshipctl_ltrim_colon_completions/g" \
	-e "s/${LWORD}compgen${RWORD}/__airshipctl_compgen/g" \
	-e "s/${LWORD}compopt${RWORD}/__airshipctl_compopt/g" \
	-e "s/${LWORD}declare${RWORD}/__airshipctl_declare/g" \
	-e "s/\\\$(type${RWORD}/\$(__airshipctl_type/g" \
	-e 's/aliashash\["\(.\{1,\}\)"\]/aliashash[\1]/g' \
	-e 's/FUNCNAME/funcstack/g' \
	<<'BASH_COMPLETION_EOF'
`
	if _, err := out.Write([]byte(zshInitialization)); err != nil {
		return fmt.Errorf("could not write zsh completion file: %s", err.Error())
	}

	buf := new(bytes.Buffer)
	if err := cmd.Root().GenBashCompletion(buf); err != nil {
		return err
	}

	if _, err := out.Write(buf.Bytes()); err != nil {
		return fmt.Errorf("could not write zsh completion file: %s", err.Error())
	}

	zshTail := `
BASH_COMPLETION_EOF
}
__airshipctl_bash_source <(__airshipctl_convert_bash_to_zsh)
`
	if _, err := out.Write([]byte(zshTail)); err != nil {
		return fmt.Errorf("could not write zsh completion file: %s", err.Error())
	}
	return nil
}
