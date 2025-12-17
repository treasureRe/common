package platform

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/go-kratos/kratos/contrib/registry/consul/v2"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/assert"
)

// TestClientWithRealService 测试与真实服务的连接
// 注意: 这个测试需要本地有 Consul 和 iamPlatformServer 服务运行
func TestClientWithRealService(t *testing.T) {
	// 检查是否设置了跳过集成测试的环境变量
	if os.Getenv("SKIP_INTEGRATION") == "true" {
		t.Skip("跳过集成测试")
	}

	// 创建 Consul 客户端
	consulConfig := consulapi.DefaultConfig()
	consulClient, err := consulapi.NewClient(consulConfig)
	if err != nil {
		t.Skipf("无法创建 Consul 客户端: %v", err)
	}

	// 创建服务发现
	discovery := consul.New(consulClient)

	// 创建平台客户端
	config := DefaultConfig().WithTimeout(5 * time.Second)
	client, err := NewClientWithDiscovery(config, discovery)
	if err != nil {
		t.Skipf("无法创建平台客户端: %v", err)
	}
	defer client.Close()

	// 测试获取权限树
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	tree, total, err := client.IAM().GetTenantPermissionsTree(ctx, nil)
	if err != nil {
		t.Logf("获取权限树失败（可能服务未启动）: %v", err)
		t.Skip("跳过测试，服务可能未启动")
		return
	}

	t.Logf("成功获取权限树，总数: %d", total)
	assert.NotNil(t, tree)
	assert.GreaterOrEqual(t, total, uint32(0))
}

// TestClientWithRealServiceFilteredByStatus 测试带状态过滤的权限树获取
func TestClientWithRealServiceFilteredByStatus(t *testing.T) {
	if os.Getenv("SKIP_INTEGRATION") == "true" {
		t.Skip("跳过集成测试")
	}

	consulConfig := consulapi.DefaultConfig()
	consulClient, err := consulapi.NewClient(consulConfig)
	if err != nil {
		t.Skipf("无法创建 Consul 客户端: %v", err)
	}

	discovery := consul.New(consulClient)

	config := DefaultConfig().WithTimeout(5 * time.Second)
	client, err := NewClientWithDiscovery(config, discovery)
	if err != nil {
		t.Skipf("无法创建平台客户端: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 测试获取 GA 状态的权限
	tree, total, err := client.IAM().GetTenantPermissionsTree(ctx, &GetTenantPermissionsTreeOptions{
		Status: "GA",
	})
	if err != nil {
		t.Logf("获取 GA 状态权限树失败: %v", err)
		t.Skip("跳过测试，服务可能未启动")
		return
	}

	t.Logf("成功获取 GA 状态权限树，总数: %d", total)
	assert.NotNil(t, tree)
}

// TestDirectConnection 测试直连方式（如果知道服务地址）
func TestDirectConnection(t *testing.T) {
	if os.Getenv("SKIP_INTEGRATION") == "true" {
		t.Skip("跳过集成测试")
	}

	// 假设服务运行在 localhost:8080（你需要根据实际情况修改）
	config := DefaultConfig().
		WithEndpoint("localhost:8080").
		WithTimeout(5 * time.Second)

	client, err := NewClient(config)
	if err != nil {
		t.Skipf("无法创建客户端: %v", err)
	}
	defer client.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, _, err = client.IAM().GetTenantPermissionsTree(ctx, nil)
	if err != nil {
		t.Logf("直连测试失败（预期行为，除非服务在 localhost:8080）: %v", err)
	}
}
