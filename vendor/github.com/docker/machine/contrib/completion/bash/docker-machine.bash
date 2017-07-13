#
# bash completion file for docker-machine commands
#
# This script provides completion of:
#  - commands and their options
#  - machine names
#  - filepaths
#
# To enable the completions either:
#  - place this file in /etc/bash_completion.d
#  or
#  - copy this file to e.g. ~/.docker-machine-completion.sh and add the line
#    below to your .bashrc after bash completion features are loaded
#    . ~/.docker-machine-completion.sh
#

# --- helper functions -------------------------------------------------------

_docker_machine_q() {
    docker-machine 2>/dev/null "$@"
}

_docker_machine_machines() {
    _docker_machine_q ls --format '{{.Name}}' "$@"
}

_docker_machine_value_of_option() {
    local pattern="$1"
    for (( i=2; i < ${cword}; ++i)); do
        if [[ ${words[$i]} =~ ^($pattern)$ ]] ; then
            echo ${words[$i + 1]}
            break
        fi
    done
}

# --- completion functions ---------------------------------------------------

_docker_machine_active() {
    case "${prev}" in
        --timeout|-t)
            return
            ;;
    esac

    if [[ "${cur}" == -* ]]; then
        COMPREPLY=($(compgen -W "--help --timeout -t" -- "${cur}"))
    fi
}

_docker_machine_config() {
    if [[ "${cur}" == -* ]]; then
        COMPREPLY=($(compgen -W "--help --swarm" -- "${cur}"))
    else
        COMPREPLY=($(compgen -W "$(_docker_machine_machines)" -- "${cur}"))
    fi
}

_docker_machine_create() {
    case "${prev}" in
        --driver|-d)
            COMPREPLY=($(compgen -W "
                amazonec2
                azure
                digitalocean
                exoscale
                generic
                google
                hyperv
                openstack
                rackspace
                softlayer
                virtualbox
                vmwarefusion
                vmwarevcloudair
                vmwarevsphere" -- "${cur}"))
            return
            ;;
    esac

    # driver specific options are only included in help output if --driver is given,
    # so we have to pass that option when calling docker-machine to harvest options.
    local driver="$(_docker_machine_value_of_option '--driver|-d')"
    local parsed_options="$(_docker_machine_q create ${driver:+--driver $driver} --help | grep '^   -' | sed 's/^   //; s/[^a-z0-9-].*$//')"
    if [[ ${cur} == -* ]]; then
        COMPREPLY=($(compgen -W "${parsed_options} -d --help" -- "${cur}"))
    fi
}

_docker_machine_env() {
    case "${prev}" in
        --shell)
            COMPREPLY=($(compgen -W "cmd fish powershell tcsh" -- "${cur}"))
            return
            ;;
    esac

    if [[ "${cur}" == -* ]]; then
	COMPREPLY=($(compgen -W "--help --no-proxy --shell --swarm --unset -u" -- "${cur}"))
    else
	COMPREPLY=($(compgen -W "$(_docker_machine_machines)" -- "${cur}"))
    fi
}

# See docker-machine-wrapper.bash for the use command
_docker_machine_use() {
    if [[ "${cur}" == -* ]]; then
        COMPREPLY=($(compgen -W "--help --swarm --unset" -- "${cur}"))
    else
        COMPREPLY=($(compgen -W "$(_docker_machine_machines)" -- "${cur}"))
    fi
}

_docker_machine_inspect() {
    case "${prev}" in
        --format|-f)
            return
            ;;
    esac

    if [[ "${cur}" == -* ]]; then
	COMPREPLY=($(compgen -W "--format -f --help" -- "${cur}"))
    else
	COMPREPLY=($(compgen -W "$(_docker_machine_machines)" -- "${cur}"))
    fi
}

_docker_machine_ip() {
    if [[ "${cur}" == -* ]]; then
        COMPREPLY=($(compgen -W "--help" -- "${cur}"))
    else
        COMPREPLY=($(compgen -W "$(_docker_machine_machines)" -- "${cur}"))
    fi
}

_docker_machine_kill() {
    if [[ "${cur}" == -* ]]; then
        COMPREPLY=($(compgen -W "--help" -- "${cur}"))
    else
        COMPREPLY=($(compgen -W "$(_docker_machine_machines)" -- "${cur}"))
    fi
}

_docker_machine_ls() {
    case "${prev}" in
        --filter|--format|-f|--timeout|-t)
            return
            ;;
    esac

    if [[ "${cur}" == -* ]]; then
        COMPREPLY=($(compgen -W "--filter --format -f --help --quiet -q --timeout -t" -- "${cur}"))
    fi
}

_docker_machine_regenerate_certs() {
    if [[ "${cur}" == -* ]]; then
        COMPREPLY=($(compgen -W "--force -f --help" -- "${cur}"))
    else
        COMPREPLY=($(compgen -W "$(_docker_machine_machines)" -- "${cur}"))
    fi
}

_docker_machine_restart() {
    if [[ "${cur}" == -* ]]; then
        COMPREPLY=($(compgen -W "--help" -- "${cur}"))
    else
        COMPREPLY=($(compgen -W "$(_docker_machine_machines)" -- "${cur}"))
    fi
}

_docker_machine_rm() {
    if [[ "${cur}" == -* ]]; then
        COMPREPLY=($(compgen -W "--force -f --help -y" -- "${cur}"))
    else
	COMPREPLY=($(compgen -W "$(_docker_machine_machines)" -- "${cur}"))
    fi
}

_docker_machine_ssh() {
    if [[ "${cur}" == -* ]]; then
        COMPREPLY=($(compgen -W "--help" -- "${cur}"))
    else
        COMPREPLY=($(compgen -W "$(_docker_machine_machines)" -- "${cur}"))
    fi
}

_docker_machine_scp() {
    if [[ "${cur}" == -* ]]; then
        COMPREPLY=($(compgen -W "--help --recursive -r" -- "${cur}"))
    else
        _filedir
        # It would be really nice to ssh to the machine and ls to complete
        # remote files.
        COMPREPLY=($(compgen -W "$(_docker_machine_machines | sed 's/$/:/')" -- "${cur}") "${COMPREPLY[@]}")
    fi
}

_docker_machine_start() {
    if [[ "${cur}" == -* ]]; then
        COMPREPLY=($(compgen -W "--help" -- "${cur}"))
    else
        COMPREPLY=($(compgen -W "$(_docker_machine_machines --filter state=Stopped)" -- "${cur}"))
    fi
}

_docker_machine_status() {
    if [[ "${cur}" == -* ]]; then
        COMPREPLY=($(compgen -W "--help" -- "${cur}"))
    else
        COMPREPLY=($(compgen -W "$(_docker_machine_machines)" -- "${cur}"))
    fi
}

_docker_machine_stop() {
    if [[ "${cur}" == -* ]]; then
        COMPREPLY=($(compgen -W "--help" -- "${cur}"))
    else
        COMPREPLY=($(compgen -W "$(_docker_machine_machines --filter state=Running)" -- "${cur}"))
    fi
}

_docker_machine_upgrade() {
    if [[ "${cur}" == -* ]]; then
        COMPREPLY=($(compgen -W "--help" -- "${cur}"))
    else
        COMPREPLY=($(compgen -W "$(_docker_machine_machines)" -- "${cur}"))
    fi
}

_docker_machine_url() {
    if [[ "${cur}" == -* ]]; then
        COMPREPLY=($(compgen -W "--help" -- "${cur}"))
    else
        COMPREPLY=($(compgen -W "$(_docker_machine_machines)" -- "${cur}"))
    fi
}

_docker_machine_version() {
    if [[ "${cur}" == -* ]]; then
        COMPREPLY=($(compgen -W "--help" -- "${cur}"))
    else
        COMPREPLY=($(compgen -W "$(_docker_machine_machines)" -- "${cur}"))
    fi
}

_docker_machine_help() {
    if [[ "${cur}" == -* ]]; then
        COMPREPLY=($(compgen -W "--help" -- "${cur}"))
    else
        COMPREPLY=($(compgen -W "${commands[*]}" -- "${cur}"))
    fi
}

_docker_machine_docker_machine() {
    if [[ " ${wants_file[*]} " =~ " ${prev} " ]]; then
        _filedir
    elif [[ " ${wants_dir[*]} " =~ " ${prev} " ]]; then
        _filedir -d
    elif [[ "${cur}" == -* ]]; then
        COMPREPLY=($(compgen -W "${flags[*]} ${wants_dir[*]} ${wants_file[*]}" -- "${cur}"))
    else
        COMPREPLY=($(compgen -W "${commands[*]}" -- "${cur}"))
    fi
}

_docker_machine() {
    COMPREPLY=()
    local commands=(active config create env inspect ip kill ls regenerate-certs restart rm ssh scp start status stop upgrade url version help)

    local flags=(--debug --native-ssh --github-api-token --bugsnag-api-token --help --version)
    local wants_dir=(--storage-path)
    local wants_file=(--tls-ca-cert --tls-ca-key --tls-client-cert --tls-client-key)

    # Add the use subcommand, if we have an alias loaded
    if [[ ${DOCKER_MACHINE_WRAPPED} = true ]]; then
        commands=("${commands[@]}" use)
    fi

    local cur prev words cword
    _get_comp_words_by_ref -n : cur prev words cword
    local i
    local command=docker-machine

    for (( i=1; i < ${cword}; ++i)); do
        local word=${words[i]}
        if [[ " ${wants_file[*]} ${wants_dir[*]} " =~ " ${word} " ]]; then
            # skip the next option
            (( ++i ))
        elif [[ " ${commands[*]} " =~ " ${word} " ]]; then
            command=${word}
        fi
    done

    local completion_func=_docker_machine_"${command//-/_}"
    if declare -F "${completion_func}" > /dev/null; then
        ${completion_func}
    fi

    return 0
}

complete -F _docker_machine docker-machine docker-machine.exe