# Common utilities, variables and checks for all build scripts.
set -eEuo pipefail

# The root of the build/dist directory
PROJ_ROOT_DIR=$(dirname "${BASH_SOURCE[0]}")/../..

# 所有的 MinerX 都会统一加载 scripts/common.sh 脚本
source "${PROJ_ROOT_DIR}/scripts/common.sh"

# 设置 MINERX_ENV_FILE（重要）
MINERX_ENV_FILE=${MINERX_ENV_FILE:-${PROJ_ROOT_DIR}/manifests/env/env.local}
# 加载本地安装环境变量（非常重要的一步，后面很多步骤都依赖于env.local中的变量设置）
source ${MINERX_ENV_FILE}

COMMON_SOURCED=true # Sourced flag

# 设置本地/容器化安装的环境变量，主要是为了避免端口冲突
export MINERX_CONTROLLER_MANAGER_METRICS_PORT=59081
export MINERX_CONTROLLER_MANAGER_HEALTHZ_PORT=59082
export MINERX_MINERSET_CONTROLLER_METRICS_PORT=60081
export MINERX_MINERSET_CONTROLLER_HEALTHZ_PORT=60082
export MINERX_MINER_CONTROLLER_METRICS_PORT=61081
export MINERX_MINER_CONTROLLER_HEALTHZ_PORT=61082

# 确保 minerx 容器网络存在。
# 在 uninstall 时，可不删除 minerx 容器网络，可以作为一个无害的无用数据
minerx::common::network()
{
  docker network ls |grep -q minerx || docker network create minerx
}
