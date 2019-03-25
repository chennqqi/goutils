package consul

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"strings"

	"github.com/Sirupsen/logrus"
	"github.com/chennqqi/goutils/closeevent"
	"github.com/chennqqi/goutils/utils"
	"github.com/chennqqi/goutils/yamlconfig"
	"gopkg.in/yaml.v2"
)

type ConsulApp struct {
	*ConsulOperator
}

func ReadTxt(c *ConsulOperator, file string) ([]byte, error) {
	if strings.HasPrefix(file, "consul://") {
		u, err := url.Parse(file)
		if gconsul == nil {
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

func NewAppWithCustomCfg(cfg interface{}, confName, healthHost string) (*ConsulApp, error) {
	var capp ConsulApp
	appName := utils.ApplicationName()
	consulapi := NewConsulOp("")
	consulapi.Fix()
	capp.ConsulOperator = consulapi

	if err := consulapi.Ping(); err != nil {
		if confName == "" {
			confName = fmt.Sprintf("%s.yml", appName)
		}

		logrus.Error("[consul/app.go]  ping consul failed, try local")
		err := yamlconfig.Load(cfg, confName)
		if os.IsNotExist(err) {
			fmt.Println("configure not exist, make default")
			yamlconfig.Save(cfg, confName)
			return nil, err
		} else if err != nil {
			logrus.Errorf("[consul/app.go] Load %v config from local error: %v", confName, err)
			return nil, err
		}
	} else {
		if confName == "" {
			confName = fmt.Sprintf("config/%s.yml", appName)
		} else if !strings.HasPrefix(confName, `config/`) {
			confName = fmt.Sprintf("config/%s", confName)
		}

		txt, err := consulapi.Get(confName)
		if err == nil {
			yaml.Unmarshal(txt, cfg)
		} else {
			logrus.Errorf("[consul/app.go] Load conf(%v) from consul error: %v", confName, err)
			err = yamlconfig.Load(cfg, confName)
			if err != nil {
				fmt.Println("make empty local config")
				yamlconfig.Save(cfg, confName)
				return nil, errors.New("make empty local config")
			}
		}
	}

	//post fix consul
	if healthHost == "" && cfg != nil {
		st := reflect.ValueOf(cfg).Elem()
		field := st.FieldByName("HealthHost")
		if !field.IsValid() {
			return nil, errors.New("cfg not contains`HealthHost")
		}
		healthHost = field.String()
	}
	if healthHost == "" {
		return nil, errors.New("cfg or HealthHost must be valid")
	}

	{
		consulapi.Name = appName
		v := strings.Split(healthHost, ":")
		if len(v) > 1 {
			fmt.Sscanf(v[1], "%d", &consulapi.Port)
		}
	}
	return &capp, nil
}

func NewAppEx(cfg interface{}, healthHost, consulUrl string) (*ConsulApp, error) {
	var capp ConsulApp

	host, port, confPath, err := ParseConsulUrl(consulUrl)
	if err != nil {
		return nil, err
	}
	appName := utils.ApplicationName()

	consulapi := NewConsulOp(host + ":" + port)
	consulapi.Fix()
	capp.ConsulOperator = consulapi

	if err := consulapi.Ping(); err != nil {
		if confName == "" {
			confName = fmt.Sprintf("%s.yml", appName)
		}

		logrus.Error("[consul/app.go]  ping consul failed, try local")
		err := yamlconfig.Load(cfg, confName)
		if os.IsNotExist(err) {
			fmt.Println("configure not exist, make default")
			yamlconfig.Save(cfg, confName)
			return nil, err
		} else if err != nil {
			logrus.Errorf("[consul/app.go] Load %v config from local error: %v", confName, err)
			return nil, err
		}
	} else { // consul is OK

		//try /${APPNAME}.yml
		//try /config/${APPNAME}.yml
		//try /config/${APPNAME}.yml
		var names []string
		confName = fmt.Sprintf("%s.yml", appName)
		names = append(names, confName)
		confName = appName
		names = append(names, confName)
		confName = fmt.Sprintf("config/%s.yml", appName)
		names = append(names, confName)
		confName = fmt.Sprintf("config/%s", confName)
		names = append(names, confName)

		//try /config/${APPNAME}.yml
		for i := 0; i < len(names); i++ {
			name := names[i]
			txt, err := consulapi.Get(confName)
			if err == nil {
			} else if err == ErrNotExist {

			} else {
				//RETURN ERROR
			}
		}

		txt, err := consulapi.Get(confName)
		if err == nil {
			yaml.Unmarshal(txt, cfg)
		} else {
			logrus.Errorf("[consul/app.go] Load conf(%v) from consul error: %v", confName, err)
			err = yamlconfig.Load(cfg, confName)
			if err != nil {
				fmt.Println("make empty local config")
				yamlconfig.Save(cfg, confName)
				return nil, errors.New("make empty local config")
			}
		}
	}

	//post fix consul
	if healthHost == "" && cfg != nil {
		st := reflect.ValueOf(cfg).Elem()
		field := st.FieldByName("HealthHost")
		if !field.IsValid() {
			return nil, errors.New("cfg not contains`HealthHost")
		}
		healthHost = field.String()
	}
	if healthHost == "" {
		return nil, errors.New("cfg or HealthHost must be valid")
	}

	{
		consulapi.Name = appName
		v := strings.Split(healthHost, ":")
		if len(v) > 1 {
			fmt.Sscanf(v[1], "%d", &consulapi.Port)
		}
	}
	return &capp, nil
}

func NewAppWithCfgEx(cfg interface{}, healthHost, agent string) (*ConsulApp, error) {
	return NewAppWithCustomCfgEx(cfg, "", healthHost, agent)
}

func NewAppWithCfg(cfg interface{}, healthHost string) (*ConsulApp, error) {
	return NewAppWithCustomCfg(cfg, "", healthHost)
}

func NewApp(healthHost string) (*ConsulApp, error) {
	//post fix consul
	if healthHost == "" {
		return nil, errors.New("healHost must be valid")
	}

	var capp ConsulApp
	appName := utils.ApplicationName()
	consulapi := NewConsulOp("")
	consulapi.Fix()
	capp.ConsulOperator = consulapi

	if err := consulapi.Ping(); err != nil {
		logrus.Error("[main] ping consul failed, try local")
		return nil, err
	}

	{
		consulapi.Name = appName
		v := strings.Split(healthHost, ":")
		if len(v) > 1 {
			fmt.Sscanf(v[1], "%d", &consulapi.Port)
		}
	}
	return &capp, nil
}

func NewAppEx(healthHost, agent string) (*ConsulApp, error) {
	//post fix consul
	if healthHost == "" {
		return nil, errors.New("healHost must be valid")
	}

	var capp ConsulApp
	appName := utils.ApplicationName()
	consulapi := NewConsulOp(agent)
	consulapi.Fix()
	capp.ConsulOperator = consulapi

	if err := consulapi.Ping(); err != nil {
		logrus.Error("[main] ping consul failed, try local")
		return nil, err
	}

	{
		consulapi.Name = appName
		v := strings.Split(healthHost, ":")
		if len(v) > 1 {
			fmt.Sscanf(v[1], "%d", &consulapi.Port)
		}
	}
	return &capp, nil
}

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
	logrus.Info("[main:main] quit, recv signal ", sig)
	if stopcall != nil {
		stopcall(sig)
	}
	c.DeregisterService()
}
