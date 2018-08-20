package eslite

import (
	"fmt"
	"log"
	"net/http"
	"strings"

	"gopkg.in/olivere/elastic.v1"
)

type ElasticClientV1 struct {
	client *elastic.Client
	bkt    *elastic.BulkService
}

func (es *ElasticClientV1) Open(host string, port int, userName, pass string) error {
	url := fmt.Sprintf("http://%s:%d", host, port)
	if strings.HasPrefix(host, "http://") || strings.HasPrefix(host, "https://") {
		url = host
	}
	client, err := elastic.NewClient(http.DefaultClient, url)
	if err != nil {
		return err
	}
	info, code, err := client.Ping().URL(url).Do()
	if err != nil {
		// Handle error
		log.Println(err)
		panic(err)
	}
	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

	esversion, err := client.ElasticsearchVersion(url)
	if err != nil {
		// Handle error
		log.Println(err)
		panic(err)
	}
	fmt.Printf("Elasticsearch version %s\n", esversion)
	es.client = client

	/*
		BulkService will be reset after each Do call.
		In other words, you can reuse BulkService to send many batches.
		You do not have to create a new BulkService for each batch.
	*/
	es.bkt = es.client.Bulk()
	return nil
}

func (es *ElasticClientV1) Write(index string, id string,
	typ string, v interface{}) error {

	es.bkt.Add(elastic.NewBulkIndexRequest().Index(
		index).Type(typ).Id(id).Doc(v))

	return nil
}

func (es *ElasticClientV1) Begin() error {
	return nil
}

func (es *ElasticClientV1) Commit() error {
	log.Println("DOBEFORE bulkRequest:NumberOfActions", es.bkt.NumberOfActions())

	bulkResponse, err := es.bkt.Do()
	if err != nil {
		log.Println(err)
		return err
	}
	if bulkResponse == nil {
		log.Fatal("expected bulkResponse to be != nil; got nil")
	}
	log.Println("DOAFTER buolkRequest:NumberOfActions", es.bkt.NumberOfActions())
	return err
}

func (es *ElasticClientV1) Close() {

}

func (es *ElasticClientV1) WriteDirect(index string, id string,
	typ string, v interface{}) error {
	_, err := es.client.Index().Index(index).Type(typ).Id(id).BodyJson(v).Do()
	return err
}

func (es *ElasticClientV1) SetPipeline(pipeline string) error {
	return ErrNotSupportPipeline
}
