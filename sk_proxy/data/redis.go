package data

import (
	"github.com/astaxie/beego/logs"
	"github.com/garyburd/redigo/redis"
	"seckill02/sk_proxy/models"
	"strings"
	"time"
)

//redis config

type RedisPoolConf struct {
	RedisAddr        string
	RedisMaxIdle     int
	RedisMaxActive   int
	RedisIdleTimeout int
}

type RedisConf struct {
	RedisPool *redis.Pool
	RedisPoolConf
	ProxyToLayerQueueName string //队列名称
	LayerToProxyQueueName string //队列名称
	IdBlackListHash       string //用户黑名单hash表
	IpBlackListHash       string //IP黑名单hash表
	IdBlackListQueue      string //用户黑名单队列
	IpBlackListQueue      string //IP黑名单队列
}

//设置redisPool 参数
func (p *RedisPoolConf) SettingRedisConf(keys ...string) (err error) {
	for _, v := range keys {
		if strings.HasSuffix(v, "addr") {
			if err = models.AppConfigStringValue(&p.RedisAddr, v); err != nil {
				break
			}
		} else if strings.HasSuffix(v, "idle") {
			if err = models.AppConfigIntValue(&p.RedisMaxIdle, v); err != nil {
				break
			}
		} else if strings.HasSuffix(v, "active") {
			if err = models.AppConfigIntValue(&p.RedisMaxActive, v); err != nil {
				break
			}
		} else if strings.HasSuffix(v, "timeout") {
			if err = models.AppConfigIntValue(&p.RedisIdleTimeout, v); err != nil {
				break
			}
		}
	}
	return
}

//设置redis 其他的参数
func (p *RedisConf) SettingRedisConfOther(proxy2layerQueueName,
	layer2proxyQueueName,
	idBlackListHash,
	ipBlackListHash,
	idBlackListQueue,
	ipBlackListQueue string) (err error){
		if err = models.AppConfigStringValue(&p.ProxyToLayerQueueName,proxy2layerQueueName);err != nil {
			return
		}
		if err = models.AppConfigStringValue(&p.LayerToProxyQueueName,layer2proxyQueueName);err != nil {
			return
		}
		if err = models.AppConfigStringValue(&p.IdBlackListHash,idBlackListHash);err != nil {
			return
		}
		if err = models.AppConfigStringValue(&p.IpBlackListHash,ipBlackListHash);err != nil {
			return
		}
		if err = models.AppConfigStringValue(&p.IdBlackListQueue,idBlackListQueue);err != nil {
			return
		}
		if err = models.AppConfigStringValue(&p.IpBlackListQueue,ipBlackListQueue);err != nil {
			return
		}
		return
}

// 初始化redisPool
func (p RedisPoolConf) initRedisPool() (redisPool *redis.Pool, err error) {
	pool := &redis.Pool{
		MaxIdle:     p.RedisMaxIdle,
		MaxActive:   p.RedisMaxActive,
		IdleTimeout: time.Duration(p.RedisIdleTimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", p.RedisAddr)
		},
	}
	conn := pool.Get()
	defer conn.Close()

	_, err = conn.Do("ping")
	if err != nil {
		logs.Error("ping redis failed,err :%v", err)
		return
	}
	redisPool = pool
	return
}

func (p RedisPoolConf)InitRedisPoolValue(redisPool **redis.Pool) (err error) {
	pool,err := p.initRedisPool()
	if err != nil {
		logs.Error("init redis failed,err: %v,addr: %v",err,p.RedisAddr)
		return
	}
	*redisPool = pool
	return
}


//加载黑名单列表
func loadBlackList() {
	//用户ID

}

