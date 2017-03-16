#eslite

es write wrapper

##interface

	type ESLite interface {
		Open(host string, port int) error
		Close()
		Begin() error
		Write(index string, id string,
			typ string, v interface{}) error
	
		WriteDirect(index string, id string,
			typ string, v interface{}) error
	
		Commit() error
	}

##demo

	var es eslite.ESLite
	var ESPort int
	var ESEngine,EsHost string
	//...
	//configure ...
	//...
	switch ESEngine {
	case "ElasticClientV":
		es = &eslite.ElasticClientV1{}
	case "ElasticClientV2":
		es = &eslite.ElasticClientV2{}
	case "ElasticClientV3":
		es = &eslite.ElasticClientV3{}
	case "ElasticClientV5":
		es = &eslite.ElasticClientV5{}
	case "ElasticGoClient":
		es = &eslite.ElasticGoClient{}
	default:
		es = &eslite.ElasticGoClient{}
	}
	if err := es.Open(cfg.EsHost, cfg.ESPort, "", ""); err != nil {
		logrus.Error("Open es error ", err)
		return
	}
	//bacth write
	es.Begin()

	var idx int
	for {
		es.Write(xxxxxxxx)
		idx++
		if idx%512==0{
			es.Commit()
			es.Begin()
		}
	}
