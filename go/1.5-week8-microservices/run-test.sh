#!/usr/bin/env bash
set -euo pipefail

ROOT="$(cd "$(dirname "$0")" && pwd)"

echo "=== 启动 etcd ==="
docker compose -f "$ROOT/docker-compose.yml" up -d

echo "=== 运行集成测试 ==="
cd "$ROOT"
go test -tags=integration -v -timeout 60s ./integration/

if [ "${KEEP_ETCD:-}" != "1" ]; then
    echo "=== 停止 etcd ==="
    docker compose -f "$ROOT/docker-compose.yml" down
fi

echo "=== 完成 ==="
