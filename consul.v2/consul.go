package consul

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/url"
	"strings"
	"sync"

	"github.com/chennqqi/goutils/utils"
	qgoutils "github.com/chennqqi/goutils/utils"
	consulapi "github.com/hashicorp/consul/api"
)

var (
	ErrNotExist = errors.New("NOT EXIST")
	HealthPath  = "health"
)

type ConsulOperator struct {
	*ConsulAppInfo

	//for check
	consul  *consulapi.Client `json:"-" yaml:"-"`
	once    sync.Once
	lockmap map[string]*consulapi.Lock
}

type ConsulAppInfo struct {
	ConsulHost string `json:"consul_host" yaml:"consul_host"`
	ConsulPort int    `json:"consul_port" yaml:"consul_port"`

	Config string     `json:"config" yaml:"config"`
	Values url.Values `json:"values" yaml:"values"`

	ServicePort int    `json:"service_port" yaml:"service_port"`
	ServiceIP   string `json:"service_ip" yaml:"service_ip"`
	ServiceName string `json:"service_name" yaml:"service_name"`

	CheckInterval string `json:"check_interval" yaml:"check_interval"`
	CheckHTTP     string `json:"check_http" yaml:"check_http"`
	CheckTCP      string `json:"check_tcp" yaml:"check_tcp"`
}

func ParseConsulUrl(consulUrl string) (*ConsulAppInfo, error) {
	var appinfo ConsulAppInfo
	appinfo.Values = make(url.Values)

	u, err := url.Parse(consulUrl)
	if err == nil {
		if u.Scheme != "consul" {
			return nil, fmt.Errorf(`expect scheme consul, not %v`, u.Scheme)
		}
		appinfo.Config = u.Path
		appinfo.ConsulHost = strings.Split(u.Host, ":")[0]
		fmt.Sscanf(u.Port(), "%d", &appinfo.ConsulPort)
		appinfo.Values = u.Query()
		querys := appinfo.Values

		appinfo.CheckInterval = querys.Get("check_interval")
		appinfo.CheckHTTP = querys.Get("check_http")
		appinfo.CheckTCP = querys.Get("check_tcp")

		appinfo.ServiceIP = querys.Get("service_ip")
		sPort := querys.Get("service_port")
		fmt.Sscanf(sPort, "%d", &appinfo.ServicePort)

		appinfo.ServiceName = querys.Get("service_name")
		return &appinfo, nil
	}

	return nil, err
}

func NewConsulOp(consulUrl string) (*ConsulOperator, error) {
	var c ConsulOperator
	c.lockmap = make(map[string]*consulapi.Lock)

	appinfo, err := ParseConsulUrl(consulUrl)
	if err != nil {
		return nil, err
	}
	c.ConsulAppInfo = appinfo
	return &c, nil
}

func (c *ConsulOperator) Fix() {
	if c.ConsulHost == "" {
		c.ConsulHost = "127.0.0.1"
	}
	if c.ConsulPort == 0 {
		c.ConsulPort = 8500
	}
	if c.ServicePort == 0 {
		c.ServicePort = 80
	}
	if c.ServiceIP == "" {
		c.ServiceIP, _ = qgoutils.GetHostIP()
		if c.ServiceIP == "" {
			c.ServiceIP, _ = qgoutils.GetInternalIP()
		}
	}
	if c.ServiceName == "" {
		c.ServiceName = utils.ApplicationName()
	}

	if c.CheckHTTP == "" && c.CheckTCP == "" {
		c.CheckHTTP = fmt.Sprintf("http://%v:%d/%v",
			c.ServiceIP, c.ServicePort, HealthPath)
		if c.CheckInterval == "" { // mix 10s
			c.CheckInterval = "10s"
		}
	}
}

func (c *ConsulOperator) Ping() error {
	var retErr error
	c.once.Do(func() {
		consulCfg := consulapi.DefaultConfig()
		consulCfg.Address = fmt.Sprintf("%v:%d", c.ConsulHost, c.ConsulPort)
		consul, err := consulapi.NewClient(consulCfg)
		retErr = err
		if err != nil {
			log.Println("New consul client error:", err)
			return
		}
		c.consul = consul
	})
	return retErr
}

func (c *ConsulOperator) Get(name string) ([]byte, error) {
	consul := c.consul
	kv := consul.KV()

	pair, _, err := kv.Get(name, nil)
	if err != nil {
		return nil, err
	}
	if pair == nil {
		return nil, ErrNotExist
	}
	return pair.Value, nil
}

func (c *ConsulOperator) GetEx(name string) ([]byte, uint64, error) {
	consul := c.consul
	kv := consul.KV()

	pair, _, err := kv.Get(name, nil)
	if err != nil {
		return nil, 0, err
	}
	if pair == nil {
		return nil, 0, ErrNotExist
	}
	return pair.Value, pair.ModifyIndex, nil
}

func (c *ConsulOperator) Put(name string, value []byte) error {
	consul := c.consul
	kv := consul.KV()
	pair := &consulapi.KVPair{
		Key:   name,
		Value: value,
	}
	_, err := kv.Put(pair, nil)
	return err
}

func (c *ConsulOperator) Delete(name string) error {
	consul := c.consul
	kv := consul.KV()
	_, err := kv.Delete(name, nil)
	return err
}

func (c *ConsulOperator) Acquire(key string, stopChan <-chan struct{}) error {
	lock, exist := c.lockmap[key]
	var err error
	if !exist {
		lock, err = c.consul.LockKey(key)
		if err != nil {
			log.Println("consul Acquire Lock key error ", err)
			return err
		}
		c.lockmap[key] = lock
	}
	_, err = lock.Lock(stopChan)
	if err != nil {
		log.Println("consul Acquire lock.Lock error ", err)
		return err
	}
	return nil
}

func (c *ConsulOperator) Release(key string) error {
	lock, exist := c.lockmap[key]
	if !exist {
		return fmt.Errorf("%v lock not exist", key)
	}
	err := lock.Unlock()
	if err != nil {
		log.Println("consul Release lock.Lock error ", err)
		return err
	}
	return nil
}

func (c *ConsulOperator) RegisterService() error {
	consul := c.consul
	agent := consul.Agent()
	check := consulapi.AgentServiceCheck{
		Interval:                       c.CheckInterval,
		HTTP:                           c.CheckHTTP,
		TCP:                            c.CheckTCP,
		DeregisterCriticalServiceAfter: "1m",
	}

	service := &consulapi.AgentServiceRegistration{
		ID:      c.ServiceName,
		Name:    c.ServiceName,
		Check:   &check,
		Address: c.ServiceIP,
		Port:    c.ServicePort,
	}
	txt, _ := json.MarshalIndent(*service, " ", "\t")
	fmt.Println("register service:", string(txt))
	return agent.ServiceRegister(service)
}

func (c *ConsulOperator) DeregisterService() error {
	consul := c.consul
	agent := consul.Agent()
	return agent.ServiceDeregister(c.ServiceName)
}

func (c *ConsulOperator) PrintServices(name string) error {
	consul := c.consul
	catalog := consul.Catalog()
	services, _, err := catalog.Service(name, "", nil)
	if err != nil {
		return err
	}
	fmt.Println("LIST services:")
	for _, v := range services {
		txt, _ := json.MarshalIndent(v, " ", "\t")
		fmt.Println(string(txt))
	}
	return err
}

func (c *ConsulOperator) ListService(name string) ([]*consulapi.CatalogService, error) {
	consul := c.consul
	catalog := consul.Catalog()
	services, _, err := catalog.Service(name, "", nil)
	return services, err
}

func (c *ConsulOperator) ListServices() (map[string][]string, error) {
	consul := c.consul
	catalog := consul.Catalog()
	services, _, err := catalog.Services(nil)
	return services, err
}
