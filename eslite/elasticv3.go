package eslite

import (
	"fmt"
	"log"
	"strings"

	"gopkg.in/olivere/elastic.v3"
)

//ElasticClientV3 object
type ElasticClientV3 struct {
	client *elastic.Client
	bkt    *elastic.BulkService
}

// open: connect with elasticsearch by user:pass@host:port
func (es *ElasticClientV3) Open(host string, port int, userName, pass string) error {
	url := fmt.Sprintf("http://%s:%d", host, port)
	if strings.HasPrefix(host, "http://") || strings.HasPrefix(host, "https://") {
		url = host
	}
	client, err := elastic.NewClient(elastic.SetURL(url), elastic.SetSniff(false))
	if err != nil {
		return err
	}
	info, code, err := client.Ping(url).Do()
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

//patch write elastic document
func (es *ElasticClientV3) Write(index string, id string,
	typ string, v interface{}) error {

	es.bkt.Add(elastic.NewBulkIndexRequest().Index(
		index).Type(typ).Id(id).Doc(v))

	return nil
}

//begin patch write
func (es *ElasticClientV3) Begin() error {
	return nil
}

//commit patch write
func (es *ElasticClientV3) Commit(pipeline string) error {
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

//close elasticsearch connection
func (es *ElasticClientV3) Close() {
	// Use the IndexExists service to check if a specified index exists.
}

//write a document directly
func (es *ElasticClientV3) WriteDirect(index string, id string,
	typ string, v interface{}) error {
	_, err := es.client.Index().Index(index).Type(typ).Id(id).BodyJson(v).Do()
	return err
}

//set elasticsearch pipeline, expect es version>5.0
func (es *ElasticClientV3) SetPipeline(pipeline string) error {
	return ErrNotSupportPipeline
}
