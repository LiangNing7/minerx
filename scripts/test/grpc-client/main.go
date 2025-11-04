//nolint:gosmopolitan
package main

import (
	"context"
	"fmt"
	"time"

	genericoptions "github.com/LiangNing7/goutils/pkg/options"

	"github.com/LiangNing7/minerx/internal/pkg/client/usercenter"
)

func main() {
	fmt.Println("=== gRPC 客户端构造流程演示 ===")

	// 演示1: 直连模式.
	// demo1DirectConn()

	// 演示2: 服务发现模式
	demo2ServiceDiscovery()
}

//nolint:unused
func demo1DirectConn() {
	fmt.Println("【演示1】直连模式")
	fmt.Println("=================")

	// step1. 创建配置选项.
	opts := &usercenter.UserCenterOptions{
		Server:  "127.0.0.1:50090",
		Timeout: 5 * time.Second,
	}
	fmt.Printf("✓ 配置选项: Server=%s, Timeout=%s\n", opts.Server, opts.Timeout)

	// step2. 创建 etcd 选项（直连模式不使用，但仍需提供）.
	etcdOpts := genericoptions.NewEtcdOptions()

	// step3. 创建客户端.
	client := usercenter.NewUserCenter(opts, etcdOpts)
	fmt.Printf("✓ 客户端创建成功（单例）\n")

	// step4. 调用服务.
	ctx := context.Background()
	// token 需登录，使用较新的token.
	token := "eyJhbGciOiJIUzUxMiIsImtpZCI6Ijk5MDEzZmM0LWNiOWQtNGZmMC1iMjJmLWY1Zjc0MTkwODY5MCIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJtaW5lcngtdXNlcmNlbnRlciIsInN1YiI6InVzZXItYnBkaG9sIiwiZXhwIjoxNzYyMjg2Mzc1LCJuYmYiOjE3NjIyNzkxNzUsImlhdCI6MTc2MjI3OTE3NX0.xxxxxxxxxxxxxxxx"
	userID, allowed, err := client.Authorize(ctx, token, "resource:demo", "read")

	// step5. 处理结果.
	if err != nil {
		fmt.Printf("✗ 调用失败: %v\n", err)
	} else {
		fmt.Printf("✓ 调用成功: userID=%s, allowed=%v\n", userID, allowed)
	}

	fmt.Println()
}

// 演示2: 服务发现模式（Etcd）.
func demo2ServiceDiscovery() {
	fmt.Println("【演示2】服务发现模式（Etcd）")
	fmt.Println("============================")

	// step1: 配置服务发现端点
	// 注意：使用 discovery:/// 前缀
	opts := &usercenter.UserCenterOptions{
		Server:  "discovery:///minerx-usercenter", // 服务发现格式
		Timeout: 10 * time.Second,
	}
	fmt.Printf("✓ 服务发现端点: %s\n", opts.Server)

	// step2: 配置 Etcd
	etcdOpts := &genericoptions.EtcdOptions{
		Endpoints: []string{"127.0.0.1:2379"}, // Etcd 地址
	}
	fmt.Printf("✓ Etcd 配置: Endpoints=%v\n", etcdOpts.Endpoints)

	// step3: 创建客户端
	// 内部流程:
	// 1. 检测到 discovery:/// 前缀
	// 2. 创建 Etcd 客户端
	// 3. 从 Etcd 查询服务实例
	// 4. 负载均衡选择实例
	// 5. 建立 gRPC 连接
	client := usercenter.NewUserCenter(opts, etcdOpts)
	fmt.Printf("✓ 客户端创建成功（通过服务发现）\n")

	// step4: 使用客户端
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	token := "eyJhbGciOiJIUzUxMiIsImtpZCI6Ijk5MDEzZmM0LWNiOWQtNGZmMC1iMjJmLWY1Zjc0MTkwODY5MCIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJtaW5lcngtdXNlcmNlbnRlciIsInN1YiI6InVzZXItYnBkaG9sIiwiZXhwIjoxNzYyMjg2Mzc1LCJuYmYiOjE3NjIyNzkxNzUsImlhdCI6MTc2MjI3OTE3NX0.-MleoEWoNdKW_Ya1gmG0c_rri0ftzkoc1bcjZoWWez7XxfIbOSA1vJhgLDEINMIZL0aEpN4sVrNCDeK1Xz1S2A"
	userID, allowed, err := client.Authorize(ctx, token, "resource:demo", "read")
	if err != nil {
		fmt.Printf("✗ 调用失败: %v\n", err)
	} else {
		fmt.Printf("✓ 调用成功: userID=%s, allowed=%v\n", userID, allowed)
	}

	fmt.Println()
}
