package eslite

import (
	"fmt"
	"log"
	"strings"

	/*as doc v6 now use github master*/
	"github.com/olivere/elastic/v7"
	"golang.org/x/net/context"
)

//ElasticClientV7 object
type ElasticClientV7 struct {
	client   *elastic.Client
	bkt      *elastic.BulkService
	pipeline string
}

// open: connect with elasticsearch by user:pass@host:port
func (es *ElasticClientV7) Open(host string, port int, usrName, pass string) error {
	url := fmt.Sprintf("http://%s:%d", host, port)
	if strings.HasPrefix(host, "http://") || strings.HasPrefix(host, "https://") {
		url = host
	}
	fmt.Println(url)
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

//patch write elastic document
func (es *ElasticClientV7) Write(index string, id string,
	typ string, v interface{}) error {
	es.bkt.Add(elastic.NewBulkIndexRequest().Index(
		index).Id(id).Doc(v))
	return nil
}

//begin patch write
func (es *ElasticClientV7) Begin() error {
	return nil
}

//commit patch write
func (es *ElasticClientV7) Commit() error {
	//	log.Println("DOBEFORE bulkRequest:NumberOfActions", es.bkt.NumberOfActions())
	es.bkt.Pipeline(es.pipeline)
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

//close elasticsearch connection
func (es *ElasticClientV7) Close() {
	// Use the IndexExists service to check if a specified index exists.
}

//write a document directly
func (es *ElasticClientV7) WriteDirect(index, id, typ string,
	v interface{}) error {
	_, err := es.client.Index().Pipeline(es.pipeline).Index(index).Id(id).BodyJson(v).Do(context.Background())
	return err
}

//set elasticsearch pipeline
func (es *ElasticClientV7) SetPipeline(pipeline string) error {
	es.pipeline = pipeline
	return nil
}
