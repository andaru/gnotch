#!/bin/bash
set -euo pipefail

protos="gnotch"

for p in $protos; do
    protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/$p.proto
done
