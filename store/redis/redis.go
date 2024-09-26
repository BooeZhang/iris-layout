package redis

import (
	"context"
	"crypto/tls"
	"os"
	"sync"
	"time"

	"github.com/kataras/golog"
	"github.com/redis/go-redis/v9"

	"irir-layout/config"
)

var (
	r    redis.UniversalClient
	once sync.Once
)

func DialToRedis(cnf *config.Redis) {
	if cnf == nil {
		golog.Error("---> [REDIS] configuration files are empty")
		os.Exit(1)
	}
	golog.Debug("Creating new Redis connection pool")
	var (
		tlsConfig *tls.Config
		client    redis.UniversalClient
	)
	once.Do(func() {
		timeout := 5 * time.Second
		if cnf.Timeout > 0 {
			timeout = time.Duration(cnf.Timeout) * time.Second
		}
		// poolSize applies per cluster node and not for the whole cluster.
		poolSize := 500
		if cnf.MaxActive > 0 {
			poolSize = cnf.MaxActive
		}
		if cnf.UseSSL {
			tlsConfig = &tls.Config{
				InsecureSkipVerify: cnf.SSLInsecureSkipVerify,
			}
		}

		redisOption := &redis.UniversalOptions{
			Addrs:        cnf.Addrs,
			MasterName:   cnf.MasterName,
			Password:     cnf.Password,
			DB:           cnf.Database,
			DialTimeout:  timeout,
			ReadTimeout:  timeout,
			WriteTimeout: timeout,
			PoolSize:     poolSize,
			TLSConfig:    tlsConfig,
		}

		if cnf.MasterName != "" {
			golog.Info("---> [REDIS] Creating sentinel-backed failover client")
			client = redis.NewFailoverClient(redisOption.Failover())
		} else if cnf.EnableCluster {
			golog.Info("---> [REDIS] Creating cluster client")
			client = redis.NewClusterClient(redisOption.Cluster())
		} else {
			golog.Info("---> [REDIS] Creating single-node client")
			client = redis.NewClient(redisOption.Simple())
		}

		pong, err := client.Ping(context.Background()).Result()
		if err != nil {
			golog.Error("---> [REDIS] redis connect ping failed, err:", err)
			os.Exit(1)
		} else {
			golog.Info("---> [REDIS] redis connect ping response: ", pong)
		}
		r = client
	})

	if r == nil {
		golog.Errorf("---> [REDIS] failed to get redis store: %+v", r)
		os.Exit(1)
	}
}

// GetRedis 获取 redis session
func GetRedis() redis.UniversalClient {
	return r
}

func Close() {
	if r != nil {
		err := r.Close()
		if err != nil {
			golog.Error("---> [REDIS] redis close failed, err:", err)
		}
		golog.Info("---> [REDIS] redis closed")
	}
}
