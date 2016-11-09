package simpleredis

import (
	"errors"
	"log"
	"strings"

	"github.com/garyburd/redigo/redis"
)

type Storage interface {
	Save(op string, key string, value string) error
	Open() error
	Close()
}

type StorageRedis struct {
	pool *redis.Pool
	URL  string
}

func (sr *StorageRedis) Open() error {
	if sr.pool != nil {
		return errors.New("pool not nil, maybe opened already")
	}
	sr.pool = newPool(sr.URL)
	if err := recover(); err != nil {
		log.Println(err)
		return errors.New("New Pool failed")
	}

	return nil
}

func (sr *StorageRedis) Close() {
	sr.pool.Close()
}

func (sr *StorageRedis) Save(op string, key string, value string) error {
	if sr.pool == nil {
		panic("Redis pool not opened")
	}

	//startTime := time.Now()

	// 从连接池里面获得一个连接
	c := sr.pool.Get()
	// 连接完关闭，其实没有关闭，是放回池里，也就是队列里面，等待下一个重用
	defer c.Close()

	if ok, err := redis.Bool(c.Do(op, key, value)); ok {
	} else {
		log.Print(err)
		return err
	}
	return nil
}

// 重写生成连接池方法
func newPool(host string) *redis.Pool {
	if strings.HasPrefix(host, "redis://") {
		return &redis.Pool{
			MaxIdle:   80,
			MaxActive: 12000, // max number of connections
			Dial: func() (redis.Conn, error) {
				c, err := redis.DialURL(host)
				if err != nil {
					panic(err.Error())
				}
				return c, err
			},
		}
	} else {
		return &redis.Pool{
			MaxIdle:   80,
			MaxActive: 12000, // max number of connections
			Dial: func() (redis.Conn, error) {
				c, err := redis.Dial("tcp", host)
				if err != nil {
					panic(err.Error())
				}
				return c, err
			},
		}
	}
}
