#!/bin/bash
set -e

cd "$(dirname "$0")/.."

# Build frontend
echo "Running yarn build"
cd plum-frontend
rm -rf dist
corepack yarn build

# Build backend
cd ../plum-backend
echo "Running https://github.com/rakyll/statik"
# -m omits modification times so that identical assets always produce an
# identical statik.go
go run github.com/rakyll/statik -f -m -dest=./entrypoints/http -src=../plum-frontend/dist
