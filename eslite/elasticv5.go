package eslite

import (
	"fmt"
	"log"
	"strings"

	"golang.org/x/net/context"
	"gopkg.in/olivere/elastic.v5"
)

type ElasticClientV5 struct {
	client *elastic.Client
	bkt    *elastic.BulkService
}

func (es *ElasticClientV5) Open(host string, port int, usrName, pass string) error {
	url := fmt.Sprintf("http://%s:%d", host, port)
	if strings.HasPrefix(host, "http://") || strings.HasPrefix(host, "https://") {
		url = host
	}
	client, err := elastic.NewClient(elastic.SetURL(url),
		elastic.SetBasicAuth(usrName, pass), elastic.SetSniff(false))

	if err != nil {
		return err
	}
	info, code, err := client.Ping(url).Do(context.TODO())
	if err != nil {
		// Handle error
		panic(err)
	}
	fmt.Printf("Elasticsearch returned with code %d and version %s\n", code, info.Version.Number)

	esversion, err := client.ElasticsearchVersion(url)
	if err != nil {
		// Handle error
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

func (es *ElasticClientV5) Write(index string, id string,
	typ string, v interface{}) error {

	es.bkt.Add(elastic.NewBulkIndexRequest().Index(
		index).Type(typ).Id(id).Doc(v))

	return nil
}

func (es *ElasticClientV5) Begin() error {
	return nil
}

func (es *ElasticClientV5) Commit() error {
	//	log.Println("DOBEFORE bulkRequest:NumberOfActions", es.bkt.NumberOfActions())

	bulkResponse, err := es.bkt.Do(context.Background())
	if err != nil {
		log.Println(err)
		return err
	}
	if bulkResponse == nil {
		log.Fatal("expected bulkResponse to be != nil; got nil")
	}
	//	log.Println("DOAFTER buolkRequest:NumberOfActions", es.bkt.NumberOfActions())
	return err
}

func (es *ElasticClientV5) Close() {
	// Use the IndexExists service to check if a specified index exists.
}

func (es *ElasticClientV5) WriteDirect(index string, id string,
	typ string, v interface{}) error {
	_, err := es.client.Index().Index(index).Type(typ).Id(id).BodyJson(v).Do(context.Background())
	return err
}
