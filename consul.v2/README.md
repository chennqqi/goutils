# consul

## feature
1. register service
2. deregister service
3. set key/value
4. aquire lock


## samples

1.basics
```
	consulapi := consul.NewConsulOp()
	consulapi.Fix()
	consulapi.Ping()
	err := consulapi.Put("AAA",[]byte("1"))
	dat, err := consulapi.Get("AAA")
	err = consulapi.Delete("AAA")
	consulapi.Aquire("lock/test")
	consulapi.Release("lock/test")
	consulapi.RegisterService()
	consulapi.DeregisterService()	

```

2.consul app


	2.1 create consul app struct pointer
	2.2 create an service to consul
	2.3 register service :8081 as an profile 

```
	capp,err := consul.NewApp(":8081")

		service := newService()
		...//you code

	capp.Wait(func (stopSig os.Signal){
		service.Stop()
	})


```


3.consul app with cfg

	3.1 consul app get ${APPNAME}.yml from consul
	3.2 if not exist or error, load ${APPNAME}.yml from local
	3.3 if load ${APPNAME}.yml not exist on local, make default cfg ${APPNAME}.yml to local
	3.4 create an service to consul 
	consul app will auto register a service with appname
	

```

	var cfg YourAppCfg
	capp,err := consul.NewAppWithCfg(&cfg, "", ":8081")

		service := newService()
		...//you code

	capp.Wait(func (stopSig os.Signal){
		service.Stop()
	})


```