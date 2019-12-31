# wallet - 钱包服务
wallet为钱包版本，目的是为了对项目需求探底并通过呀测获取监控指标，为集群版本的开发铺路。

# 目录结构
* assets：静态文件(存放配置等)
* cmd: 运行文件(存放服务启动文件)
* core: 核心代码
    * servant: 服务接口定义
    * storage: 数据接口定义
* docs: 开发文档说明
* logs: 日志文件
* pkg: 辅助插件
* router: 路由分配

# 技术栈
* 版本控制：gogs/git/github
* 开发语言: golang
* web 框架: gin
* grpc/secret_key
* 数据库: mysql
* 容器: docker/docker-compose/k8s  

## 安装/运行
```
mkdir -p $GOPATH/src/github.com/haleyrom/wallet
cd $GOPATH/src/github.com/haleyrom/wallet
git clone https://github.com/haleyrom/wallet.git
make server
```