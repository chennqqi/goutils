package consul

import (
	"encoding/json"
	"fmt"

	"sync"

	"github.com/Sirupsen/logrus"
	qgoutils "github.com/chennqqi/goutils/utils"
	"github.com/chennqqi/goutils/yamlconfig"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/pkg/errors"
)

var (
	ErrNotExist = errors.New("NOT EXIST")
)

const (
	CONSUL_HEALTH_PATH = "health"
)

type ConsulOperator struct {
	Agent    string `json:"agent,omitempty" yaml:"agent,omitempty"`
	IP       string `json:"ip" yaml:"ip"`
	Port     int    `json:"port" yaml:"port"`
	Name     string `json:"Name" yaml:"Name"`
	Path     string `json:"path,omitempty" yaml:"path,omitempty"`
	Interval string `json:"interval,omitempty" yaml:"interval,omitempty"`

	//for check
	consul  *consulapi.Client `json:"-" yaml:"-"`
	once    sync.Once
	lockmap map[string]*consulapi.Lock
}

func NewConsulOp(agent string) *ConsulOperator {
	var c ConsulOperator
	c.lockmap = make(map[string]*consulapi.Lock)
	c.Agent = agent
	return &c
}

func (c *ConsulOperator) Fix() {
	if c.Agent == "" {
		c.Agent = "localhost:8500"
	}
	if c.Path == "" {
		c.Path = CONSUL_HEALTH_PATH
	}
	if c.Port == 0 {
		c.Port = 80
	}
	if c.IP == "" {
		c.IP, _ = qgoutils.GetHostIP()
		if c.IP == "" {
			c.IP, _ = qgoutils.GetInternalIP()
		}
	}
	if c.Interval == "" {
		c.Interval = "10s"
	}
}

func (c *ConsulOperator) Ping() error {
	var retErr error
	c.once.Do(func() {
		consulCfg := consulapi.DefaultConfig()
		consulCfg.Address = c.Agent
		consul, err := consulapi.NewClient(consulCfg)
		retErr = err
		if err != nil {
			logrus.Error("New consul client error: ", err)
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
			logrus.Error("consul Require Lockkey error ", err)
			return err
		}
		c.lockmap[key] = lock
	}
	_, err = lock.Lock(stopChan)
	if err != nil {
		logrus.Error("consul Require lock.Lock error ", err)
		return err
	}
	return nil
}

func (c *ConsulOperator) Release(key string) error {
	lock, exist := c.lockmap[key]
	if !exist {
		return errors.Errorf("%v lock not exist", key)
	}
	err := lock.Unlock()
	if err != nil {
		logrus.Error("consul Require lock.Lock error ", err)
		return err
	}
	return nil
}

func (c *ConsulOperator) RegisterService() error {
	consul := c.consul
	agent := consul.Agent()
	check := consulapi.AgentServiceCheck{
		Interval:                       c.Interval,
		HTTP:                           fmt.Sprintf("http://%s:%d/%s", c.IP, c.Port, c.Path),
		DeregisterCriticalServiceAfter: "1m",
	}

	service := &consulapi.AgentServiceRegistration{
		ID:      c.Name,
		Name:    c.Name,
		Check:   &check,
		Address: c.IP,
		Port:    c.Port,
	}
	txt, _ := json.MarshalIndent(*service, " ", "\t")
	fmt.Println("register service:", string(txt))
	return agent.ServiceRegister(service)
}

func (c *ConsulOperator) DeregisterService() error {
	consul := c.consul
	agent := consul.Agent()
	return agent.ServiceDeregister(c.Name)
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

func (c *ConsulOperator) Save() {
	yamlconfig.Save(*c, "consul.yml")
}
