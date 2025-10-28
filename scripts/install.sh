#!/usr/bin/env bash

# The root of the build/dist directory
PROJ_ROOT_DIR=$(dirname "${BASH_SOURCE[0]}")/..
source "${PROJ_ROOT_DIR}/scripts/lib/init.sh"

function minerx::install::install_cfssl()
{
  mkdir -p $HOME/bin/
  wget https://github.com/cloudflare/cfssl/releases/download/v1.6.5/cfssl_1.6.5_linux_amd64 -O $HOME/bin/cfssl
  wget https://github.com/cloudflare/cfssl/releases/download/v1.6.5/cfssljson_1.6.5_linux_amd64 -O $HOME/bin/cfssljson
  wget https://github.com/cloudflare/cfssl/releases/download/v1.6.5/cfssl-certinfo_1.6.5_linux_amd64 -O $HOME/bin/cfssl-certinfo
  #wget https://pkg.cfssl.org/R1.2/cfssl_linux-amd64 -O $HOME/bin/cfssl
  #wget https://pkg.cfssl.org/R1.2/cfssljson_linux-amd64 -O $HOME/bin/cfssljson
  #wget https://pkg.cfssl.org/R1.2/cfssl-certinfo_linux-amd64 -O $HOME/bin/cfssl-certinfo
  chmod +x $HOME/bin/{cfssl,cfssljson,cfssl-certinfo}
  onex::log::info "install cfssl tools successfully"
}   

if [[ "$*" =~ minerx::install:: ]]; then
  eval $*
fi

