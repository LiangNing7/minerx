#!/usr/bin/env bash

# Short-circuit if verify-generated.sh has already been sourced.
[[ $(type -t minerx::verify::generated::loaded) == function ]] && return 0

source "${PROJ_ROOT_DIR}/scripts/lib/init.sh"

# This function verifies whether generated files are up-to-date. The first two
# parameters are messages that get printed to stderr when changes are found,
# the rest are the function or command and its parameters for generating files
# in the work tree.
#
# Example: minerx::verify::generated "Mock files are out of date" "Please run 'scripts/update-mocks.sh'" scripts/update-mocks.sh
minerx::verify::generated() {
  ( # a subshell prevents environment changes from leaking out of this function
    local failure_header=$1
    shift
    local failure_tail=$1
    shift

    minerx::util::ensure_clean_working_dir

    # This sets up the environment, like GOCACHE, which keeps the worktree cleaner.
    minerx::golang::setup_env

    _tmpdir="$(minerx::realpath "$(mktemp -d -t "verify-generated-$(basename "$1").XXXXXX")")"
    git worktree add -f -q "${_tmpdir}" HEAD
    minerx::util::trap_add "git worktree remove -f ${_tmpdir}" EXIT
    cd "${_tmpdir}"

    # Update generated files.
    "$@"

    # Test for diffs
    diffs=$(git status --porcelain | wc -l)
    if [[ ${diffs} -gt 0 ]]; then
      if [[ -n "${failure_header}" ]]; then
        echo "${failure_header}" >&2
      fi
      git status >&2
      git diff >&2
      if [[ -n "${failure_tail}" ]]; then
        echo "" >&2
        echo "${failure_tail}" >&2
      fi
      return 1
    fi
  )
}

# Marker function to indicate verify-generated.sh has been fully sourced.
minerx::verify::generated::loaded() {
  return 0
}
