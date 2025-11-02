#!/usr/bin/env bash

# 本脚本功能：根据 scripts/environment.sh 配置，生成 MINERX 组件 YAML 配置文件。
# 示例：gen-config.sh scripts/environment.sh configs/miner-apiserver.yaml

env_file="$1"
template_file="$2"

PROJ_ROOT_DIR=$(dirname "${BASH_SOURCE[0]}")/..

source "${PROJ_ROOT_DIR}/scripts/lib/init.sh"

if [ $# -ne 2 ];then
    minerx::log::error "Usage: gen-config.sh manifests/env.local configs/minerx.service.tmpl"
    exit 1
fi

source "${env_file}"

declare -A envs

set +u
for env in $(sed -n 's/^[^#].*${\(.*\)}.*/\1/p' ${template_file})
do
    if [ -z "$(eval echo \$${env})" ];then
        minerx::log::error "environment variable '${env}' not set"
        missing=true
    fi
done

if [ "${missing}" ];then
    minerx::log::error 'You may run `source manifests/env.local` to set these environment'
    exit 1
fi

eval "cat << EOF
$(cat ${template_file})
EOF"
