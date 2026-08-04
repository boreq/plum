#!/bin/bash
set -e

cd "$(dirname "$0")/.."

# Build frontend
echo "Running yarn build"
cd plum-frontend
rm -rf dist
corepack yarn build

# Build backend
cd ..
echo "Running https://github.com/rakyll/statik"
go run github.com/rakyll/statik -f -src=./plum-frontend/dist
