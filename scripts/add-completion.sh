#!/usr/bin/env bash

cmd=$1
shell=$2

function add() 
{
  cat << EOF >> ${HOME}/.bashrc

# ${cmd} shell completion 
if [ -f \${HOME}/.${cmd}-completion.bash ]; then 
    source \${HOME}/.${cmd}-completion.bash
fi 
EOF
}

${cmd} completion ${shell} > ${HOME}/.${cmd}-completion.bash

if ! grep -q "# ${cmd} shell completion" ${HOME}/.bashrc;then
  add
fi
