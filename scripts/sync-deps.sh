#!/usr/bin/env bash
set -euo pipefail

# Sync minerx deps to match versions used in minerx

cd "$(dirname "$0")/.."

mods=(
  "github.com/spf13/pflag@v1.0.6"
  "github.com/go-resty/resty/v2@v2.16.5"
  "github.com/golang-jwt/jwt/v4@v4.5.2"
  "github.com/onsi/ginkgo/v2@v2.22.2"
  "github.com/prometheus/common@v0.62.0"
  "github.com/google/go-cmp@v0.7.0"
  "github.com/google/gofuzz@v1.2.0"
  "github.com/go-kratos/kratos/v2@v2.8.3"
  "go.uber.org/automaxprocs@v1.6.0"
  "google.golang.org/grpc@v1.70.0"
  "google.golang.org/protobuf@v1.36.5"
  "google.golang.org/genproto/googleapis/api@v0.0.0-20250204164813-702378808489"
  "gorm.io/gen@v0.3.26"
  "gorm.io/gorm@v1.31.0"
  "k8s.io/apiserver@v0.33.2"
  "k8s.io/apimachinery@v0.33.2"
  "k8s.io/component-base@v0.33.2"
  "k8s.io/code-generator@v0.33.2"
  "k8s.io/gengo@v0.0.0-20250130153323-76c5745d3511"
  "k8s.io/kube-openapi@v0.0.0-20250318190949-c8a335a9a2ff"
  "k8s.io/klog/v2@v2.130.1"
  "golang.org/x/text@v0.28.0"
  "github.com/armon/go-socks5@v0.0.0-20160902184237-e75332964ef5"
)

echo "Syncing dependencies to minerx-aligned versions..."
for m in "${mods[@]}"; do
  echo "go get $m"
  go get "$m"
done

echo "Running go mod tidy..."
go mod tidy

echo "Done."

