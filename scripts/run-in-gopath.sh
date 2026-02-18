#!/usr/bin/env bash

# This script sets up a temporary Kubernetes GOPATH and runs an arbitrary
# command under it. Go tooling requires that the current directory be under
# GOPATH or else it fails to find some things, such as the vendor directory for
# the project.
# Usage: `hack/run-in-gopath.sh <command>`.

set -o errexit
set -o nounset
set -o pipefail

PROJ_ROOT_DIR=$(dirname "${BASH_SOURCE[0]}")/..
source "${PROJ_ROOT_DIR}/scripts/lib/init.sh"

# This sets up a clean GOPATH and makes sure we are currently in it.
minerx::golang::setup_env

# Run the user-provided command.
"${@}"

