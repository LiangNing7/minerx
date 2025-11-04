# 这是一个 UserCenter 服务的 gRPC 客户端，使用 Etcd 作为服务发现：
主要特点：
 - 使用 gRPC 协议连接 usercenter 服务
 - 支持 Etcd 服务发现 (newEtcdClient)
 - 默认服务地址：127.0.0.1:8081
 - 提供授权功能：Authorize(ctx, token, obj, act) 方法
 - 集成了链路追踪（tracing）中间件
 - 使用真实的 v1.UserCenterClient API
