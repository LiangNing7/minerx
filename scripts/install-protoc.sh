#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

PROJ_ROOT_DIR=$(dirname "${BASH_SOURCE[0]}")/..
source "${PROJ_ROOT_DIR}/scripts/lib/init.sh"
source "${PROJ_ROOT_DIR}/scripts/lib/protoc.sh"

minerx::protoc::install

