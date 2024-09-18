#!/bin/bash

# 默认端口
DEFAULT_PORT=59101
SERVER_PORT=58080

# 启动 server 函数
start_server() {
    port=$1
    if [ -z "$port" ]; then
        port=$SERVER_PORT
    fi

    echo "Starting server on port $port..."
    cd src/test
    go run -tags=server server.go -port="$port"
}

start_tx_client() {
    echo "Starting tx client"
    cd src/test
    go run -tags=client txClient.go
}

# 启动 main 函数
start_node() {
    # port=$1
    # if [ -z "$port" ]; then
    #     port=$DEFAULT_PORT
    # fi

    shift  # 移除第一个参数 -m
    CMD="go run main.go $@"

    # 显示执行的命令
    echo "Executing: $CMD"

    cd src
    # 执行构建好的命令
    $CMD
}

# 启动多个 main 实例，端口递增
start_multiple_node() {
    port=$1
    count=$2
    if [ -z "$port" ]; then
        port=$DEFAULT_PORT
    fi

    if [ -z "$count" ]; then
        count=1
    fi

    echo "Starting $count main instances starting from port $port..."
    cd src

    # 使用一个子进程组，确保所有实例可以一起控制
    trap "kill 0" EXIT  # 当主进程被终止时，终止所有子进程

    for ((i=0; i<count; i++)); do
        new_port=$((port + i))
        echo "Starting main instance on port $new_port..."
        # go run main.go -port="$new_port" -ping -duration 1 &
        go run main.go -ping -duration 5 &
    done

    # 等待所有子进程完成
    wait
}


# # 启动多个 main 实例，端口递增
# start_multiple_node() {
#     port=$1
#     count=$2
#     if [ -z "$port" ]; then
#         port=$DEFAULT_PORT
#     fi

#     if [ -z "$count" ]; then
#         count=1
#     fi

#     echo "Starting $count main instances starting from port $port..."
#     cd src
#     for ((i=0; i<count; i++)); do
#         new_port=$((port + i))
#         echo "Starting main instance on port $new_port..."
#         go run main.go -port="$new_port" &
#     done
# }

# 显示帮助信息
show_help() {
    echo "Usage: $0 [options]"
    echo "Options:"
    echo "  -s, --server [PORT]    Start the server (default: 58080)"
    echo "  -n, --node [PORT]      Start the main on the specified port (default: 59101)"
    echo "  -mn, --mnode [PORT] [COUNT]  Start multiple main instances, starting from the specified port"
    echo "  -h, --help             Show this help message"
}

# 解析命令行参数
case $1 in
    -s|--server)
        start_server $2
        ;;
    -n|--node)
        start_node "$@"
        ;;
    -c|--client)
        start_tx_client "$@"
        ;;
    -mn|--mnode)
        start_multiple_node $2 $3
        ;;
    -h|--help|*)
        show_help
        ;;
esac
