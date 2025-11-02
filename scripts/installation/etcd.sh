#!/bin/bash

# The root of the build/dist directory.
PROJ_ROOT_DIR=$(dirname "${BASH_SOURCE[0]}")/../..
# If common.sh has already been sourced, it will not be sourced again here.
[[ -z ${COMMON_SOURCED} ]] && source ${PROJ_ROOT_DIR}/scripts/installation/common.sh
# Set some environment variables.
MINERX_ETCD_HOST=${MINERX_ETCD_HOST:-127.0.0.1}
MINERX_ETCD_CLIENT_PORT=${MINERX_ETCD_CLIENT_PORT:-2379}
MINERX_ETCD_PEER_PORT=${MINERX_ETCD_PEER_PORT:-2380}
MINERX_ETCD_VERSION=${ETCD_VERSION:-3.5.11}
MINERX_ETCD_DIR=${MINERX_ETCD_DIR:-${MINERX_THIRDPARTY_INSTALL_DIR}/etcd}

# Install etcd using containerization.
minerx::etcd::docker::install()
{
  minerx::common::network
  docker run -d --name minerx-etcd \
    --restart always \
    --network minerx \
    -v ${MINERX_ETCD_DIR}:/etcd-data \
    -p ${MINERX_ACCESS_HOST}:${MINERX_ETCD_CLIENT_PORT}:2379 \
    -p ${MINERX_ACCESS_HOST}:${MINERX_ETCD_PEER_PORT}:2380 \
    quay.io/coreos/etcd:v3.5.13 \
    /usr/local/bin/etcd \
    --advertise-client-urls http://0.0.0.0:2379 \
    --listen-client-urls http://0.0.0.0:2379 \
    --data-dir /etcd-data

  sleep 10
  minerx::etcd::status || return 1
  minerx::etcd::info
  minerx::log::info "install etcd successfully"
}

# Uninstall the docker container.
minerx::etcd::docker::uninstall()
{
  docker rm -f minerx-etcd &>/dev/null
  minerx::util::sudo "rm -rf ${MINERX_ETCD_DIR}"
  minerx::log::info "uninstall etcd successfully"
}

# Install the etcd step by step.
# sbs is the abbreviation for "step by step".
minerx::etcd::sbs::install()
{
  local os
  local arch

  os=$(minerx::util::host_os)
  arch=$(minerx::util::host_arch)

  download_file_name="etcd-v${MINERX_ETCD_VERSION}-${os}-${arch}"
  download_file="/tmp/${download_file_name}.tar.gz"
  url="https://github.com/coreos/etcd/releases/download/v${MINERX_ETCD_VERSION}/${download_file_name}.tar.gz"
  minerx::util::download_file "${url}" "${download_file}"
  tar -xvzf "${download_file}" -C /tmp/
  echo ${LINUX_PASSWORD} | sudo -S cp /tmp/${download_file_name}/{etcd,etcdctl,etcdutl} /usr/bin/
  rm "${download_file}"
  rm -rf /tmp/${download_file_name}

  # 创建 Etcd 配置文件
  minerx::util::sudo "mkdir -p /etc/etcd"
  # Etcd 会输出日志到 ${MINERX_LOG_DIR}，所以需要先创建该目录
  minerx::util::sudo "mkdir -p ${MINERX_LOG_DIR}"

  echo ${LINUX_PASSWORD} | sudo -S cat << EOF | sudo tee /etc/etcd/config.yaml
name: minerx
data-dir: ${MINERX_ETCD_DIR}
advertise-client-urls: http://0.0.0.0:${MINERX_ETCD_CLIENT_PORT}
listen-client-urls: http://0.0.0.0:${MINERX_ETCD_CLIENT_PORT}
initial-advertise-peer-urls: http://0.0.0.0:${MINERX_ETCD_PEER_PORT}
initial-cluster: minerx=http://0.0.0.0:${MINERX_ETCD_PEER_PORT}
log-outputs: [${MINERX_LOG_DIR}/etcd.log]
log-level: debug
EOF

  echo ${LINUX_PASSWORD} | sudo -S cat << EOF | sudo tee /lib/systemd/system/etcd.service
# Etcd systemd unit template from
# https://github.com/etcd-io/etcd/blob/main/contrib/systemd/etcd.service
[Unit]
Description=etcd key-value store # 指定了单元的描述，即 etcd 键值存储
Documentation=https://github.com/etcd-io/etcd # 提供了指向 etcd 项目文档的链接
After=network-online.target local-fs.target remote-fs.target time-sync.target # 指定了服务的启动顺序
Wants=network-online.target local-fs.target remote-fs.target time-sync.target # 指定了服务的启动依赖

[Service]
Type=notify # 指定了服务的类型。notify 类型表示服务会在准备就绪时发送通知
ExecStart=/usr/bin/etcd --config-file=/etc/etcd/config.yaml # 指定了服务启动时要执行的命令，这里是使用指定的配置文件启动 etcd
Restart=always # 指定了服务的重启行为。always 表示服务会在退出时总是被重启
RestartSec=10s # 指定了重启的间隔时间
LimitNOFILE=40000 # 指定了服务的文件描述符限制，这里设置为 40000

[Install]
WantedBy=multi-user.target # 指定了服务的安装目标，这里表示服务会被添加到 multi-user.target，以便在多用户模式下启动
EOF

  minerx::util::sudo "systemctl daemon-reload"
  minerx::util::sudo "systemctl start etcd"
  minerx::util::sudo "systemctl enable etcd"

  minerx::etcd::status || return 1
  minerx::etcd::info
  minerx::log::info "install etcd v${MINERX_ETCD_VERSION} successfully"
}

# Uninstall the etcd step by step.
minerx::etcd::sbs::uninstall()
{
  #set +o errexit
  #etcd_pids=$(pgrep -f $HOME/bin/etcd)
  #set -o errexit
  #if [[ ${etcd_pids} != "" ]];then
    #echo ${LINUX_PASSWORD} | sudo -S kill -9 ${etcd_pids} || true
  #fi

  # etcd, etcdctl, etcdutl 3 个文件这里不需要删除，因为下次安装时，文件会被覆盖，仍然是幂等安装
  minerx::util::sudo "systemctl stop etcd"
  minerx::util::sudo "rm -rf ${MINERX_ETCD_DIR}"
  minerx::log::info "uninstall etcd successfully"
}

# Print necessary information after docker or sbs installation.
minerx::etcd::info()
{
  minerx::color::green "etcd has been installed, here are some useful information:"
  cat << EOF | sed 's/^/  /'
Etcd endpoint is: ${MINERX_ETCD_HOST}:${MINERX_ETCD_CLIENT_PORT}
EOF
}

# Status check after docker or sbs installation.
minerx::etcd::status()
{
  minerx::util::telnet ${MINERX_ETCD_HOST} ${MINERX_ETCD_CLIENT_PORT} || return 1

  #echo "Waiting for etcd to come up."
  #minerx::util::wait_for_url "http://${MINERX_ETCD_HOST}:${MINERX_ETCD_CLIENT_PORT}/health" "etcd: " 0.25 80
}

if [[ "$*" =~ minerx::etcd:: ]]; then
  eval $*
fi

