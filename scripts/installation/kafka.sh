#!/bin/bash

# The root of the build/dist directory.
PROJ_ROOT_DIR=$(dirname "${BASH_SOURCE[0]}")/../..
# If common.sh has already been sourced, it will not be sourced again here.
[[ -z ${COMMON_SOURCED} ]] && source ${PROJ_ROOT_DIR}/scripts/installation/common.sh
# Set some environment variables.
MINERX_KAFKA_HOST=${MINERX_KAFKA_HOST:-127.0.0.1}
MINERX_KAFKA_PORT=${MINERX_KAFKA_PORT:-4317}

# Install kafka using containerization.
# Refer to https://www.baeldung.com/ops/kafka-docker-setup
minerx::kafka::docker::install()
{
  minerx::common::network
  docker run -d --restart always --name minerx-zookeeper --network minerx -p 2181:2181 -t wurstmeister/zookeeper
  docker run -d --name minerx-kafka --link minerx-zookeeper:zookeeper \
    --restart always \
    --network minerx \
    --restart=always \
    -v /etc/localtime:/etc/localtime \
    -p ${MINERX_KAFKA_HOST}:${MINERX_KAFKA_PORT}:9092 \
    --env KAFKA_ZOOKEEPER_CONNECT=zookeeper:2181 \
    --env KAFKA_ADVERTISED_HOST_NAME=${MINERX_KAFKA_HOST} \
    --env KAFKA_ADVERTISED_PORT=${MINERX_KAFKA_PORT} \
    wurstmeister/kafka

  echo "Sleeping to wait for minerx-kafka container to complete startup ..."
  sleep 5
  minerx::kafka::status || return 1
  minerx::kafka::info
  minerx::log::info "install kafka successfully"
}

# Uninstall the docker container.
minerx::kafka::docker::uninstall()
{
  docker rm -f minerx-zookeeper &>/dev/null
  docker rm -f minerx-kafka &>/dev/null
  minerx::log::info "uninstall kafka successfully"
}

# Install the kafka step by step.
# sbs is the abbreviation for "step by step".
# Refer to https://kafka.apache.org/documentation/#quickstart
minerx::kafka::sbs::install()
{
  minerx::kafka::docker::install
  minerx::log::info "install kafka successfully"
}

# Uninstall the kafka step by step.
minerx::kafka::sbs::uninstall()
{
  minerx::kafka::docker::uninstall
  minerx::log::info "uninstall kafka successfully"
}

# Print necessary information after docker or sbs installation.
minerx::kafka::info()
{
  echo -e ${C_GREEN}kafka has been installed, here are some useful information:${C_NORMAL}
  cat << EOF | sed 's/^/  /'
Kafka brokers is: ${MINERX_KAFKA_HOST}:${MINERX_KAFKA_PORT}
EOF
}

# Status check after docker or sbs installation.
minerx::kafka::status()
{
  minerx::util::telnet ${MINERX_KAFKA_HOST} ${MINERX_KAFKA_PORT} || return 1
}

if [[ $* =~ minerx::kafka:: ]]; then
  eval $*
fi

