#!/bin/bash

# 生成 proto 文件的 Go 代码
# 用法: ./scripts/gen-proto.sh

set -e

PROJECT_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
PROTO_DIR="$PROJECT_ROOT/api/proto/v1"

# 检测 protoc 路径
if [ -n "$PROTOC" ]; then
    PROTOC_BIN="$PROTOC"
elif command -v protoc &> /dev/null; then
    PROTOC_BIN="protoc"
elif [ -f "/usr/local/bin/protoc" ]; then
    PROTOC_BIN="/usr/local/bin/protoc"
elif [ -f "/usr/bin/protoc" ]; then
    PROTOC_BIN="/usr/bin/protoc"
else
    echo "Error: protoc not found. Please install protoc first."
    echo "Download: https://github.com/protocolbuffers/protobuf/releases"
    exit 1
fi

# 获取 protoc include 路径
PROTOC_INCLUDE=""
case "$(uname)" in
    Darwin)
        if [ -f "/usr/local/opt/protobuf/include/google/protobuf" ]; then
            PROTOC_INCLUDE="/usr/local/opt/protobuf/include"
        elif [ -f "/opt/homebrew/opt/protobuf/include" ]; then
            PROTOC_INCLUDE="/opt/homebrew/opt/protobuf/include"
        fi
        ;;
    Linux)
        if [ -f "/usr/include/google/protobuf" ]; then
            PROTOC_INCLUDE="/usr/include"
        elif [ -f "/usr/local/include/google/protobuf" ]; then
            PROTOC_INCLUDE="/usr/local/include"
        fi
        ;;
    MINGW*|MSYS*|CYGWIN*)
        # Windows
        PROTOC_BIN="C:\Code_Tools\protoc-32.1-win64\bin\protoc.exe"
        PROTOC_INCLUDE="C:\Code_Tools\protoc-32.1-win64\include"
        ;;
esac

echo "Using protoc: $PROTOC_BIN"
echo "Proto dir: $PROTO_DIR"

# 生成 Go 代码
generate_go() {
    local proto_file="$1"
    local proto_name=$(basename "$proto_file")
    
    echo "Generating Go code for: $proto_name"
    
    $PROTOC_BIN \
        --proto_path="$PROTO_DIR" \
        --proto_path="$PROTOC_INCLUDE" \
        --go_out="$PROTO_DIR" \
        --go_opt=paths=source_relative \
        --go-grpc_out="$PROTO_DIR" \
        --go-grpc_opt=paths=source_relative \
        "$proto_file"
}

# 生成所有 proto 文件
if [ $# -eq 0 ]; then
    for proto_file in "$PROTO_DIR"/*.proto; do
        if [ -f "$proto_file" ]; then
            generate_go "$proto_file"
        fi
    done
else
    for proto_file in "$@"; do
        if [ -f "$proto_file" ]; then
            generate_go "$proto_file"
        else
            echo "Warning: File not found: $proto_file"
        fi
    done
fi

echo "Done!"