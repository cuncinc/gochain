gochain
go开发的区块链系统

## 使用
1. 打开区块链中转服务器
```shell
./run.sh -s
```

2. 打开区块链节点
```shell
./run.sh -n [-server] [-ping]
```
打开多个节点
```shell
./run.sh -mn port nums
```

3. 打开区块链客户端
```shell
./run.sh -c
```