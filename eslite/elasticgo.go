package eslite

import (
	"bytes"
	"fmt"
	"log"
	"time"

	elastigo "github.com/wxiaomei/elastigo/lib"
)

type ElasticGoClient struct {
	client *elastigo.Conn
	bkt    *elastigo.BulkIndexer
}

func (es *ElasticGoClient) Open(host string, port int, userName, pass string) error {
	c := elastigo.NewConn()
	log.SetFlags(log.LstdFlags)
	c.Domain = host
	c.Port = fmt.Sprintf("%d", port)

	//	c.RequestTracer = func(method, url, body string) {
	//		log.Printf("Requesting %s %s", method, url)
	//		log.Printf("Request body: %s", body)
	//	}
	es.client = c
	return nil
}

func (es *ElasticGoClient) Write(index string, id string,
	typ string, v interface{}) error {
	err := es.bkt.Index(index, typ, id, "", "", nil, v)
	if err != nil {
		log.Println("ESGoClient ERR:", err)
	}
	return err
}

func (es *ElasticGoClient) Begin() error {
	indexer := es.client.NewBulkIndexer(10)
	indexer.BufferDelayMax = 60 * time.Second
	indexer.BulkMaxDocs = 1024
	indexer.BulkMaxBuffer = 1048576

	indexer.Sender = func(buf *bytes.Buffer) error {
		// @buf is the buffer of docs about to be written
		respJson, err := es.client.DoCommand("POST", "/_bulk", nil, buf)
		if err != nil {
			// handle it better than this
			fmt.Println(string(respJson))
		}
		return err
	}
	es.bkt = indexer
	es.bkt.Start()
	return nil
}

func (es *ElasticGoClient) Commit() error {
	es.bkt.Stop()
	return nil
}

func (es *ElasticGoClient) Close() {
	es.client.Close()
}

func (es *ElasticGoClient) WriteDirect(index string, id string,
	typ string, v interface{}) error {
	_, err := es.client.Index(index, typ, id, nil, v)
	return err
}
