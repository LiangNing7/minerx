#!/usr/bin/env bash

# This script sets up a go workspace locally and builds all go components.
# Usage: `scripts/build.sh <component name> <platform >`.
# Example: `scripts/build.sh onex-usercenter linux_amd64`

set -o errexit
set -o nounset
set -o pipefail

PROJ_ROOT_DIR=$(dirname "${BASH_SOURCE[0]}")/..
source "${PROJ_ROOT_DIR}/scripts/common.sh"

minerx::golang::build_binaries "$1" "$2"

