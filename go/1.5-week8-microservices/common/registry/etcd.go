package registry

import (
	"context"
	"fmt"
	"sync"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
)

// Registry 封装了 etcd 的服务注册与发现功能
type Registry struct {
	client  *clientv3.Client
	leaseID clientv3.LeaseID
	key     string // 注册时写入的 key
	cancel  context.CancelFunc

	mu        sync.RWMutex
	instances map[string][]string // 缓存服务实例: serviceName -> []addr
}

// New 创建一个 Registry 实例并连接 etcd
func New(endpoints []string) (*Registry, error) {
	client, err := clientv3.New(clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 2 * time.Second,
	})
	if err != nil {
		return nil, fmt.Errorf("connect etcd: %w", err)
	}
	return &Registry{
		client:    client,
		instances: make(map[string][]string),
	}, nil
}

// Register 将服务注册到 etcd，并使用租约实现心跳保活
func (r *Registry) Register(serviceName, addr string, ttl int64) error {
	// 1. 创建租约
	lease, err := r.client.Grant(context.Background(), ttl)
	if err != nil {
		return fmt.Errorf("grant lease: %w", err)
	}
	r.leaseID = lease.ID

	// 2. 写入服务 key（前缀模式利于发现）
	r.key = fmt.Sprintf("/services/%s/%s", serviceName, addr)
	_, err = r.client.Put(context.Background(), r.key, addr, clientv3.WithLease(lease.ID))
	if err != nil {
		return fmt.Errorf("register: %w", err)
	}

	// 3. 后台心跳续约
	ctx, cancel := context.WithCancel(context.Background())
	r.cancel = cancel
	keepAliveCh, err := r.client.KeepAlive(ctx, lease.ID)
	if err != nil {
		cancel()
		return fmt.Errorf("keepalive: %w", err)
	}
	go func() {
		for range keepAliveCh {
			// 消费续约响应，保持租约活跃
		}
		fmt.Printf("[registry] %s 心跳通道关闭\n", serviceName)
	}()

	fmt.Printf("[registry] ✅ %s 注册成功 @ %s (ttl=%ds)\n", serviceName, addr, ttl)
	return nil
}

// Deregister 注销服务
func (r *Registry) Deregister() error {
	if r.cancel != nil {
		r.cancel()
	}
	if r.key != "" {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_, err := r.client.Delete(ctx, r.key)
		if err != nil {
			return err
		}
		fmt.Printf("[registry] 🗑️  已注销 key: %s\n", r.key)
	}
	return nil
}

// Discover 发现指定服务的所有实例地址（从 etcd 实时查询，失败返回 error）
func (r *Registry) Discover(serviceName string) ([]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	prefix := fmt.Sprintf("/services/%s/", serviceName)
	resp, err := r.client.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		// 失败时返回 error，调用方可用缓存兜底
		return nil, fmt.Errorf("discover %s: %w", serviceName, err)
	}

	addrs := make([]string, 0, len(resp.Kvs))
	for _, kv := range resp.Kvs {
		addrs = append(addrs, string(kv.Value))
	}

	r.mu.Lock()
	r.instances[serviceName] = addrs
	r.mu.Unlock()

	return addrs, nil
}

// DiscoverCached 发现服务实例：优先查 etcd，失败则用本地缓存兜底
func (r *Registry) DiscoverCached(serviceName string) []string {
	addrs, err := r.Discover(serviceName)
	if err == nil {
		return addrs
	}

	// etcd 不可用，返回缓存
	r.mu.RLock()
	cached := r.instances[serviceName]
	r.mu.RUnlock()

	if len(cached) > 0 {
		fmt.Printf("[registry] ⚠️  etcd 不可用，%s 使用缓存: %v\n", serviceName, cached)
	}
	return cached
}

// Watch 监听服务实例的上/下线变化
func (r *Registry) Watch(serviceName string, onChange func([]string)) {
	// 先拉取一次当前实例，避免错过 Watch 启动前的注册
	addrs := r.DiscoverCached(serviceName)
	onChange(addrs)

	prefix := fmt.Sprintf("/services/%s/", serviceName)
	watchCh := r.client.Watch(context.Background(), prefix, clientv3.WithPrefix())

	go func() {
		for range watchCh {
			addrs, err := r.Discover(serviceName)
			if err == nil {
				onChange(addrs)
			}
			// watch 失败（etcd 断连）时不做任何操作，
			// 保持当前缓存不变，等 etcd 恢复后 watch 会自动重连
		}
	}()
	fmt.Printf("[registry] 👀 开始监听 %s 实例变化 (当前: %v)\n", serviceName, addrs)
}

// Close 关闭 etcd 连接
func (r *Registry) Close() error {
	r.Deregister()
	return r.client.Close()
}
