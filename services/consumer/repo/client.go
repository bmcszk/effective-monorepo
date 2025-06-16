package repo

import (
	"context"
	"os"
	"strings"
	"time"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/namespace"
)

type Config struct {
	URIs   []string
	Prefix string
}

func NewConfig() *Config {
	prefix := os.Getenv("ETCD_PREFIX")
	if prefix != "" && !strings.HasSuffix(prefix, "/") {
		prefix += "/"
	}
	return &Config{
		URIs:   strings.Split(os.Getenv("ETCD_URIS"), "|"),
		Prefix: prefix,
	}
}

type Repo struct {
	client *clientv3.Client
}

func NewRepo() (*Repo, error) {
	config := NewConfig()
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   config.URIs,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		return nil, err
	}
	if config.Prefix != "" {
		cli.KV = namespace.NewKV(cli.KV, config.Prefix)
		cli.Watcher = namespace.NewWatcher(cli.Watcher, config.Prefix)
		cli.Lease = namespace.NewLease(cli.Lease, config.Prefix)
	}
	return &Repo{
		client: cli,
	}, nil
}

func (r *Repo) Close() error {
	return r.client.Close()
}

func (r *Repo) Get(ctx context.Context, key string) (string, error) {
	resp, err := r.client.Get(ctx, key)
	if err != nil {
		return "", err
	}
	if resp.Count == 0 {
		return "", nil
	}
	return string(resp.Kvs[0].Value), nil
}

func (r *Repo) Put(ctx context.Context, key string, value string) error {
	_, err := r.client.Put(ctx, key, value)
	return err
}
