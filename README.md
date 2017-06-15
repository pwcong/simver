# SimVer
简单的验证服务。对外restful，服务端间rpc。

可简单对用户请求/验证次数进行限制，防止恶意请求。

## 安装方法
1. 首先克隆仓库 `go get github.com/pwcong/simver`
2. 到目录 `$GOROOT/github.com/pwcong/simver` 下执行 `go build ./main.go`
3. 运行Redis数据库服务，修改配置文件，执行 `./main`

## 验证流程
![process](https://raw.githubusercontent.com/pwcong/SnapShot/master/simver/check.png)

1. 用户访问客户端

2. 客户端向SimVer服务发起请求获取Key，响应内容如下：
    ```
    {
        "code": 200,
        "msg": "",
        "key": "xxx"
    }
    ```

3. 客户端携带Key向服务端发起请求

4. 服务端携带客户端的Key向SimVer服务发起验证请求，并获取验证结果。通过验证结果对客户端请求进行相应的处理

## 配置说明
**默认读取配置文件为 `./conf/simver.toml`，默认配置如下：**
```
[server]
ip = "0.0.0.0"                      # restful监听地址
port = 56789                        # restful监听端口

    [server.vertify]
    signingKey = "simver"           # key加密签名
    visitCounts = 10                # 周期内访问次数限制
    checkCounts = 1                 # 周期内验证次数限制
    expiredTime = 86400000000000    # 周期时长，单位为微秒

    [server.rpc]
    ip = "0.0.0.0"                  # rpc监听地址
    port = 56780                    # rpc监听端口

[databases]
    [databases.redis]
    ip = "127.0.0.1"                # redis数据库服务地址
    port = 6379                     # redis数据库服务端口
    password = ""                   # redis数据库服务连接密码
    db = 0                          # redis数据库服务指定数据库

[middlewares]

    [middlewares.cors]              # CORS中间件，允许跨域请求
    active = true                   # 是否启用该中间件
    allowOrigins = ["*"]            # 允许跨域来源
    allowMethods = ["GET", "HEAD", "PUT", "PATCH", "POST", "DELETE"] # 允许请求方法

    [middlewares.log]
    active = true
    format = "${time_rfc3339_nano} ${remote_ip} ${host} ${method} ${uri} ${status} ${latency_human} ${bytes_in} ${bytes_out}" # 日志打印格式
    output = "stdout"               # 可选值有："stdout"、"file"，"stdout" 为控制台输出，"file" 为日志记录到 ./log/server.log

```

## 服务端间RPC通讯说明
**服务端间RPC采用谷歌GRPC技术，通讯数据采用protobuf3格式，数据模型为 `./vertify/vertify.proto`，内容如下：**
```
syntax = "proto3";

package vertify;

service Vertify {
    rpc CheckKey(VertifyRequest) returns (VertifyResponse);
}

message VertifyRequest {
    string key = 1;
}

message VertifyResponse {
    bool checked = 1;
}

```

