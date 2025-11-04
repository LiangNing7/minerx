#!/bin/bash

# The root of the build/dist directory.
PROJ_ROOT_DIR=$(dirname "${BASH_SOURCE[0]}")/../..
# If common.sh has already been sourced, it will not be sourced again here.
[[ -z ${COMMON_SOURCED} ]] && source ${PROJ_ROOT_DIR}/scripts/installation/common.sh
# Set some environment variables.
MINERX_JAEGER_HOST=${MINERX_JAEGER_HOST:-127.0.0.1}
MINERX_JAEGER_PORT=${MINERX_JAEGER_PORT:-4317}

# Install Jaeger using containerization.
minerx::jaeger::docker::install()
{
  minerx::common::network
  docker run -d --name minerx-jaeger \
    --restart always \
    --network minerx \
    -e COLLECTOR_ZIPKIN_HOST_PORT=:9411 \
    -p ${MINERX_ACCESS_HOST}:6831:6831/udp \
    -p ${MINERX_ACCESS_HOST}:6832:6832/udp \
    -p ${MINERX_ACCESS_HOST}:5778:5778 \
    -p ${MINERX_ACCESS_HOST}:16686:16686 \
    -p ${MINERX_ACCESS_HOST}:${MINERX_JAEGER_PORT}:4317 \
    -p ${MINERX_ACCESS_HOST}:4318:4318 \
    -p ${MINERX_ACCESS_HOST}:14250:14250 \
    -p ${MINERX_ACCESS_HOST}:14268:14268 \
    -p ${MINERX_ACCESS_HOST}:14269:14269 \
    -p ${MINERX_ACCESS_HOST}:9411:9411 \
    jaegertracing/all-in-one:1.52

  sleep 2
  minerx::jaeger::status || return 1
  minerx::jaeger::info
  minerx::log::info "install jaeger successfully"
}

# Uninstall the docker container.
minerx::jaeger::docker::uninstall()
{
  docker rm -f minerx-jaeger &>/dev/null
  minerx::log::info "uninstall jaeger successfully"
}

# Install the jaeger step by step.
# sbs is the abbreviation for "step by step".
minerx::jaeger::sbs::install()
{
  minerx::jaeger::docker::install
  minerx::log::info "install jaeger successfully"
}

# Uninstall the jaeger step by step.
minerx::jaeger::sbs::uninstall()
{
  minerx::jaeger::docker::uninstall
  minerx::log::info "uninstall jaeger successfully"
}

# Start the jaeger container.
minerx::jaeger::docker::start()
{
  docker start minerx-jaeger &>/dev/null && minerx::log::info "jaeger started successfully" || minerx::log::error "failed to start jaeger"
}

# Stop the jaeger container.
minerx::jaeger::docker::stop()
{
  docker stop minerx-jaeger &>/dev/null && minerx::log::info "jaeger stopped successfully" || minerx::log::error "failed to stop jaeger"
}

# Restart the jaeger container.
minerx::jaeger::docker::restart()
{
  docker restart minerx-jaeger &>/dev/null && minerx::log::info "jaeger restarted successfully" || minerx::log::error "failed to restart jaeger"
}

# Print necessary information after docker or sbs installation.
minerx::jaeger::info()
{
  echo -e ${C_GREEN}Jaeger has been installed, here are some useful information:${C_NORMAL}
  cat << EOF | sed 's/^/  /'
OpenTelemetry Protocol (OTLP) over gRPC Endpoint: ${MINERX_JAEGER_HOST}:${MINERX_JAEGER_PORT}
EOF
}

# Status check after docker or sbs installation.
minerx::jaeger::status()
{
  minerx::util::telnet ${MINERX_JAEGER_HOST} ${MINERX_JAEGER_PORT} || return 1
}

if [[ "$*" =~ minerx::jaeger:: ]]; then
  eval $*
fi
