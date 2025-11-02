#!/bin/bash

# The root of the build/dist directory.
PROJ_ROOT_DIR=$(dirname "${BASH_SOURCE[0]}")/../..
# If common.sh has already been sourced, it will not be sourced again here.
[[ -z ${COMMON_SOURCED} ]] && source ${PROJ_ROOT_DIR}/scripts/installation/common.sh
# Set some environment variables.
MINERX_MYSQL_HOST=${MINERX_MYSQL_HOST:-127.0.0.1}
MINERX_MYSQL_PORT=${MINERX_MYSQL_PORT:-3306}
MINERX_PASSWORD=${MINERX_PASSWORD:-minerx(#)666}

# Install mariadb using containerization.
minerx::mariadb::docker::install()
{
  # 安装客户端工具，访问 MariaDB
  minerx::util::sudo "apt install -y mariadb-client"

  minerx::common::network
  docker run -d --name minerx-mariadb \
    --restart always \
    --network minerx \
    -v ${MINERX_THIRDPARTY_INSTALL_DIR}/mariadb:/var/lib/mysql \
    -p ${MINERX_ACCESS_HOST}:${MINERX_MYSQL_PORT}:3306 \
    -e MYSQL_ROOT_PASSWORD=${MINERX_PASSWORD} \
    mariadb:11.2.2

  echo "Sleeping to wait for all minerx-mariadb container to complete startup ..."
  sleep 10
  minerx::mariadb::status || return 1

  minerx::mariadb::info
  minerx::log::info "install mariadb successfully"
}

# Uninstall the docker container.
minerx::mariadb::docker::uninstall()
{
  docker rm -f minerx-mariadb &>/dev/null
  minerx::util::sudo "rm -rf ${MINERX_THIRDPARTY_INSTALL_DIR}/mariadb"
  minerx::log::info "uninstall mariadb successfully"
}

# Install the mariadb step by step.
# sbs is the abbreviation for "step by step".
minerx::mariadb::sbs::install()
{
  # 本机 apt 安装后 MySQL 端口固定位 3306
  export MINERX_MYSQL_PORT=3306

  # 从指定的 URL中获取 MariaDB 的发布密钥。这个密钥用于验证 MariaDB
  # 软件包的签名，确保软件包在下载和安装过程中的完整性和安全性
  echo ${LINUX_PASSWORD} | sudo -S apt-key adv --fetch-keys 'https://mariadb.org/mariadb_release_signing_key.asc'
  # 配置 MariaDB 11.2.2 apt 源（docker install 和 sbs install 版本都要保持一致）
  echo ${LINUX_PASSWORD} | sudo -S echo "deb [arch=amd64,arm64] https://mirrors.aliyun.com/mariadb/repo/11.2.2/debian/ $(lsb_release -cs) main" | sudo tee /etc/apt/sources.list.d/mariadb-11.2.2.list

  # 注意：一定要执行 `apt update`，否则可能安装的还是旧的软件包
  minerx::util::sudo "apt update"

  # 需要先创建 /var/lib/mysql/ 目录，否则 `systemctl start mariadb` 时可能会报错
  minerx::util::sudo "mkdir -p /var/lib/mysql"

  # 执行以下命令，防止uninstall后，出现：`update-alternatives: error: alternative path /etc/mysql/mariadb.cnf doesn't exist` 错误
  # 安装 MariaDB 客户端和 MariaDB 服务端
  minerx::util::sudo "apt install -y -o Dpkg::Options::="--force-confmiss" --reinstall mariadb-client mariadb-server"

  # 启动 MariaDB，并设置开机启动
  minerx::util::sudo "systemctl enable mariadb"

  # 为了方便你访问 MySQL，这里我们设置 MySQL 允许从所有机器网卡访问
  echo ${LINUX_PASSWORD} | sudo -S sed -i 's/^bind-address.*/bind-address = 0.0.0.0/g' /etc/mysql/mariadb.conf.d/50-server.cnf

  minerx::util::sudo "systemctl restart mariadb"

  #  设置 root 初始密码
  minerx::util::sudo "mysqladmin -u${MINERX_MYSQL_ADMIN_USERNAME} password ${MINERX_MYSQL_ADMIN_PASSWORD}"

  minerx::mariadb::status || return 1
  minerx::mariadb::info
  minerx::log::info "install mariadb successfully"
}

# Uninstall the mariadb step by step.
minerx::mariadb::sbs::uninstall()
{
  # `|| true` 实现幂等
  minerx::util::sudo "systemctl stop mariadb" || true
  minerx::util::sudo "systemctl disable mariadb" || true
  minerx::util::sudo "apt remove -y mariadb-client mariadb-server" || true

  # 删除配置文件和数据目录，以及其他关联安装文件
  minerx::util::sudo "rm -rvf /var/lib/mysql"
  minerx::util::sudo "rm -rvf /etc/mysql"
  minerx::util::sudo "rm -rvf /usr/share/keyrings/mariadb.gpg"
  minerx::util::sudo "rm -vf /etc/apt/sources.list.d/mariadb-11.2.2.list"
  minerx::log::info "uninstall mariadb successfully"
}

# Print necessary information after docker or sbs installation.
minerx::mariadb::info()
{
  minerx::color::green "mariadb has been installed, here are some useful information:"
  cat << EOF | sed 's/^/  /'
MySQL access endpoint is: ${MINERX_MYSQL_HOST}:${MINERX_MYSQL_PORT}
        root password is: ${MINERX_PASSWORD}
# `mysql` will be deprecated in the future, so here use `mariadb` instead.
Access command: mariadb -h ${MINERX_MYSQL_HOST} -P ${MINERX_MYSQL_PORT} -u root -p'${MINERX_PASSWORD}'
EOF
}

# Status check after docker or sbs installation.
minerx::mariadb::status()
{
  sleep 20
  # 基础检查：检查端口，基础检查
  minerx::util::telnet ${MINERX_MYSQL_HOST} ${MINERX_MYSQL_PORT} || return 1

  # 终态检查：检查 MySQL 是否成功运行
  echo mariadb -h${MINERX_MYSQL_HOST} -P${MINERX_MYSQL_PORT} -u${MINERX_MYSQL_ADMIN_USERNAME} -p${MINERX_MYSQL_ADMIN_PASSWORD} -e quit
  mariadb -h${MINERX_MYSQL_HOST} -P${MINERX_MYSQL_PORT} -u${MINERX_MYSQL_ADMIN_USERNAME} -p${MINERX_MYSQL_ADMIN_PASSWORD} -e quit &>/dev/null || {
    minerx::log::error "can not login with root, mariadb maybe not initialized properly."
    return 1
  }
}

if [[ "$*" =~ minerx::mariadb:: ]]; then
  eval $*
fi

