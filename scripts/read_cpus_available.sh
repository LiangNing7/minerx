#!/bin/bash

set -e
set +u

if [ -z "$MAX_CPUS" ]; then
    MAX_CPUS=1

    case "$(uname -s)" in
    Darwin)
        MAX_CPUS=$(sysctl -n machdep.cpu.core_count)
        ;;
    Linux)
        # 使用 nproc 命令获取可用 CPU 数量
        MAX_CPUS=$(nproc --all)
        ;;
    *)
        # Unsupported host OS. Must be Linux or Mac OS X.
        echo "Unsupported OS"
        exit 1
        ;;
    esac
fi

echo "$MAX_CPUS"

