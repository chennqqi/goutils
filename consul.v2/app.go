package consul

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/url"
	"os"
	"os/signal"
	"reflect"
	"strings"

	"github.com/chennqqi/goutils/closeevent"
	"github.com/chennqqi/goutils/utils"
	"github.com/chennqqi/goutils/yamlconfig"
	"gopkg.in/yaml.v2"
)

type ConsulApp struct {
	*ConsulOperator
}

// load an conf file text form file(FQDN)
// try consul first, if not exist, try local file
func ReadTxt(c *ConsulOperator, file string) ([]byte, error) {
	if strings.HasPrefix(file, "consul://") {
		u, err := url.Parse(file)
		if c == nil {
			return nil, errors.New("consul not set")
		} else if err != nil {
			return nil, err
		} else {
			return c.Get(u.Path)
		}
	} else {
		return ioutil.ReadFile(file)
	}
}

// create a consul app with load cfg, using 127.0.0.1:8500
// cfg is required parameter, the cfg address
// arg heathHost is must parameter, a health http address of app
// arg consulUrl is an FQDN format string which contains consul host,port,configpath, if configpath is empty,
// ${appname}.yml, config/${appname} will be tried to load in order
func NewConsulAppWithCfg(cfg interface{}, consulUrl string) (*ConsulApp, error) {
	var capp ConsulApp

	appinfo, err := ParseConsulUrl(consulUrl)
	if err != nil {
		return nil, err
	}
	appName := utils.ApplicationName()
	consulapi, err := NewConsulOp(consulUrl)
	if err != nil {
		return nil, err
	}

	consulapi.Fix()
	capp.ConsulOperator = consulapi

	//try /${APPNAME}.yml
	//try /config/${APPNAME}.yml
	//try /config/${APPNAME}.yml
	var names []string
	var defaultName string
	var confName = appinfo.Config
	if confName == "" {
		confName = fmt.Sprintf("%s.yml", appName)
		defaultName = confName
		names = append(names, confName)
		confName = appName
		names = append(names, confName)
		confName = fmt.Sprintf("config/%s.yml", appName)
		names = append(names, confName)
		confName = fmt.Sprintf("config/%s", confName)
		names = append(names, confName)
	} else {
		names = append(names, confName)
		defaultName = appinfo.Config
	}

	if err := consulapi.Ping(); err != nil {
		log.Println("[consul/app.go] ping consul failed, try local")
		var exist bool
		for i := 0; i < len(names); i++ {
			log.Printf("[consul/app.go] ping consul failed, try load %v", names[i])
			err := yamlconfig.Load(cfg, names[i])
			if err == nil {
				exist = true
				break
			}
		}
		if !exist {
			log.Printf("[consul/app.go] all try load failed, make default %v", defaultName)
			yamlconfig.Save(cfg, defaultName)
			return nil, ErrNotExist
		}
	} else { // consul is OK
		//try /config/${APPNAME}.yml
		var exist bool
		for i := 0; i < len(names); i++ {
			txt, err := consulapi.Get(names[i])
			if err == nil {
				log.Printf("[consul/app.go] successfully get consul kv: %v", names[i])
				yaml.Unmarshal(txt, cfg)
				exist = true
				break
			} else {
				log.Printf("[consul/app.go] failed get consul kv(%v), error(%v)", names[i], err)
			}
		}
		if !exist {
			log.Printf("[consul/app.go] all try load failed, make default %v", defaultName)
			yamlconfig.Save(cfg, defaultName)
			return nil, ErrNotExist
		}
	}

	//post fix consul
	if cfg != nil {
		//health
		var healthHost string
		var health string

		{
			st := reflect.ValueOf(cfg).Elem()
			field := st.FieldByName("Health")
			if field.IsValid() {
				health = field.String()
				//health is an url
				if strings.HasPrefix(health, "tcp") {
					consulapi.CheckTCP = health
					consulapi.CheckHTTP = ""
				} else if strings.HasPrefix(health, "http") {
					consulapi.CheckHTTP = health
					consulapi.CheckTCP = ""
				}
			}
		}
		if health == "" {
			st := reflect.ValueOf(cfg).Elem()
			field := st.FieldByName("HealthHost")
			if field.IsValid() {
				healthHost = field.String()
				v := strings.Split(healthHost, ":")
				if len(v) > 1 {
					fmt.Sscanf(v[1], "%d", &consulapi.ServicePort)
				}
				fmt.Println("healthHost:", v, consulapi.ServicePort)
				if v[0] != "" && v[0] != "127.0.0.1" && v[0] != "##1" {
					ip := net.ParseIP(v[0])
					if ip != nil {
						consulapi.ServiceIP = ip.String()
					}
				}
			}
		}
	}
	return &capp, nil
}

// create a consul app
// arg heathHost is must parameter, a health http address of app
// arg agent is option, if empty using default 127.0.0.1:8500
func NewConsulApp(consulUrl string) (*ConsulApp, error) {
	consulapi, err := NewConsulOp(consulUrl)
	if err != nil {
		return nil, err
	}
	return &ConsulApp{consulapi}, nil
}

// wait for main function return and register app to service to consul
func (c *ConsulApp) Wait(stopcall func(os.Signal), signals ...os.Signal) {
	quitChan := make(chan os.Signal, 1)
	defer close(quitChan)
	if len(signals) > 0 {
		signal.Notify(quitChan, signals...)
	} else {
		closeevent.CloseNotify(quitChan)
	}

	c.RegisterService()
	sig := <-quitChan
	log.Println("[main:main] quit, recv signal ", sig)
	if stopcall != nil {
		stopcall(sig)
	}
	c.DeregisterService()
}
