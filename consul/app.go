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

func NewAppWithCustomCfg(cfg interface{}, confName, healthHost string) (*ConsulApp, error) {
	var capp ConsulApp
	appName := utils.ApplicationName()
	consulapi := NewConsulOp("")
	consulapi.Fix()
	capp.ConsulOperator = consulapi

	if err := consulapi.Ping(); err != nil {
		logrus.Error("[main] ping consul failed, try local")

		if confName == "" {
			confName = fmt.Sprintf("%s.yml", appName)
		}
		err := yamlconfig.Load(cfg, confName)
		if os.IsNotExist(err) {
			fmt.Println("configure not exist, make default")
			yamlconfig.Save(cfg, confName)
			return nil, err
		} else if err != nil {
			logrus.Error("[main:main] Load yml config error")
			return nil, err
		}
	} else {
		if confName == "" {
			confName = fmt.Sprintf("config/%s.yml", appName)
		}
		txt, err := consulapi.Get(confName)
		if err == nil {
			yaml.Unmarshal(txt, cfg)
		} else {
			logrus.Error("[main:main] Load yml from consul error ", err)
			err = yamlconfig.Load(cfg, "")
			if err != nil {
				fmt.Println("make empty local config")
				yamlconfig.Save(cfg, "")
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

func NewAppWithCfg(cfg interface{}, confName, healthHost string) (*ConsulApp, error) {
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
