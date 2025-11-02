#!/bin/bash

# The root of the build/dist directory.
PROJ_ROOT_DIR=$(dirname "${BASH_SOURCE[0]}")/../..
# If common.sh has already been sourced, it will not be sourced again here.
[[ -z ${COMMON_SOURCED} ]] && source ${PROJ_ROOT_DIR}/scripts/installation/common.sh
# Set some environment variables.
MINERX_REDIS_HOST=${MINERX_REDIS_HOST:-127.0.0.1}
MINERX_REDIS_PORT=${MINERX_REDIS_PORT:-6379}
MINERX_PASSWORD=${MINERX_PASSWORD:-minerx(#)666}

# Install redis using containerization.
minerx::redis::docker::install()
{
  minerx::redis::pre_install

  minerx::common::network
  docker run -d --name minerx-redis \
    --restart always \
    --network minerx \
    -v ${MINERX_THIRDPARTY_INSTALL_DIR}/redis:/data \
    -p ${MINERX_ACCESS_HOST}:${MINERX_REDIS_PORT}:6379 \
    redis:7.2.3 \
    redis-server \
    --appendonly yes \
    --save 60 1 \
    --protected-mode no \
    --requirepass ${MINERX_PASSWORD} \
    --loglevel debug

  sleep 2
  minerx::redis::status || return 1
  minerx::redis::info
  minerx::log::info "install redis successfully"
}

# Uninstall the docker container.
minerx::redis::docker::uninstall()
{
  docker rm -f minerx-redis &>/dev/null
  minerx::util::sudo "rm -rf ${MINERX_THIRDPARTY_INSTALL_DIR}/redis"
  minerx::log::info "uninstall redis successfully"
}

# Install the redis step by step.
# sbs is the abbreviation for "step by step".
minerx::redis::sbs::install()
{
  minerx::redis::pre_install

  # 创建 `/var/lib/redis` 目录，否则 `redis-server` 命令启动时
  # 会报：`Can't chdir to '/var/lib/redis': No such file or directory` 错误
  minerx::util::sudo "mkdir -p /var/lib/redis"

  # 安装 Redis
  minerx::util::sudo "apt install -y -o Dpkg::Options::="--force-confmiss" --reinstall redis-server"

  # 配置 Redis
  # 修改 `/etc/redis/redis.conf` 文件，将 daemonize 由 no 改成 yes，表示允许 Redis 在后台启动
  redis_conf=/etc/redis/redis.conf
  # 注意：有的系统 redis 配置文件路径为 `/etc/redis.conf`
  [[ -f /etc/redis.conf ]] && redis_conf=/etc/redis.conf

  echo ${LINUX_PASSWORD} | sudo -S sed -i '/^daemonize/{s/no/yes/}' ${redis_conf}

  # 修改 Redis 端口为 ${MINERX_REDIS_PORT}
  echo ${LINUX_PASSWORD} | sudo -S sed -i "s/^port.*/port ${MINERX_REDIS_PORT}/g" ${redis_conf}

  # 在 `bind 127.0.0.1` 前面添加 `#` 将其注释掉，默认情况下只允许本地连接，注释掉后外网可以连接 Redis
  echo ${LINUX_PASSWORD} | sudo -S sed -i '/^bind .*127.0.0.1/s/^/# /' ${redis_conf}

  # 修改 requirepass 配置，设置 Redis 密码
  echo ${LINUX_PASSWORD} | sudo -S sed -i 's/^# requirepass.*$/requirepass '"${MINERX_REDIS_PASSWORD}"'/' ${redis_conf}

  # 因为我们上面配置了密码登录，需要将 protected-mode 设置为 no，关闭保护模式
  echo ${LINUX_PASSWORD} | sudo -S sed -i '/^protected-mode/{s/yes/no/}' ${redis_conf}

  # 为了能够远程连上 Redis，需要执行以下命令关闭防火墙，并禁止防火墙开机启动（如果不需要远程连接，可忽略此步骤）
  set +o errexit
  #minerx::util::sudo "systemctl stop firewalld.service"
  #minerx::util::sudo "systemctl disable firewalld.service"
  set -o errexit

  # 重启 Redis
  #minerx::util::sudo "redis-server ${redis_conf}"
  minerx::util::sudo "systemctl restart redis-server"

  minerx::redis::status || return 1
  minerx::redis::info
  minerx::log::info "install redis successfully"
}

minerx::redis::pre_install()
{
  minerx::util::sudo "apt install -y redis-tools"
}

# Uninstall the redis step by step.
minerx::redis::sbs::uninstall()
{
  # 先删除 redis-server 进程，否则 `systemctl stop redis-server` 可能会卡主
  redis_pid=$(pgrep -f redis-server)
  [[ ${redis_pid} != "" ]] && minerx::util::sudo "kill -9 ${redis_pid}"

  set +o errexit
  minerx::util::sudo "systemctl stop redis-server"
  minerx::util::sudo "systemctl disable redis-server"
  minerx::util::sudo "apt remove -y redis-server"
  minerx::util::sudo "rm -rf /var/lib/redis"
  set -o errexit
  minerx::log::info "uninstall redis successfully"
}

# Print necessary information after docker or sbs installation.
minerx::redis::info()
{
  echo -e ${C_GREEN}redis has been installed, here are some useful information:${C_NORMAL}
  cat << EOF | sed 's/^/  /'
Redis access endpoint is: ${MINERX_REDIS_HOST}:${MINERX_REDIS_PORT}
       Redis password is: ${MINERX_PASSWORD}
     Redis Login Command: redis-cli --no-auth-warning -h ${MINERX_REDIS_HOST} -p ${MINERX_REDIS_PORT} -a '${MINERX_REDIS_PASSWORD}'
EOF
}

# Status check after docker or sbs installation.
minerx::redis::status()
{
  minerx::util::telnet ${MINERX_REDIS_HOST} ${MINERX_REDIS_PORT} || return 1
  redis-cli --no-auth-warning -h ${MINERX_REDIS_HOST} -p ${MINERX_REDIS_PORT} -a "${MINERX_REDIS_PASSWORD}" --hotkeys || {
    minerx::log::error "can not login with ${MINERX_REDIS_USERNAME}, redis maybe not initialized properly."
    return 1
  }
}

if [[ "$*" =~ minerx::redis:: ]]; then
  eval $*
fi

