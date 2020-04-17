package elasticsearch

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/elastic/go-elasticsearch/esapi"
	"github.com/elastic/go-elasticsearch/v8"
)

type Connection struct {
	conn *elasticsearch.Client
}

type Config struct {
	Addresses []string
}

type Document struct {
	ID    string
	Index string
	Body  interface{}
}

type ResultItem struct {
	ID     string
	Source map[string]interface{}
	Score  float64
	Type   string
}

type QueryResult struct {
	Hits []ResultItem
	Took float64
}

func New(cfg *Config) (*Connection, error) {
	es, err := elasticsearch.NewClient(elasticsearch.Config{
		Addresses: cfg.Addresses,
	})

	if err != nil {
		return nil, err
	}

	res, err := es.Info()
	log.Printf("Elasticsearch info: %s", res)

	return &Connection{conn: es}, err
}

func (c *Connection) Write(index string, data []Document) error {
	// Create a context object for the API calls
	ctx := context.Background()

	for _, item := range data {
		bod, err := json.Marshal(item.Body)

		if err != nil {
			return err
		}

		req := esapi.IndexRequest{
			Index:      item.Index,
			DocumentID: item.ID,
			Body:       strings.NewReader(string(bod)),
			Refresh:    "true",
		}

		res, err := req.Do(ctx, c.conn)
		if err != nil {
			return err
		}

		defer res.Body.Close()

		if res.IsError() {
			return fmt.Errorf("%s -error indexing document ID=%s", res.Status(), item.ID)
		}
	}

	return nil
}

func (c *Connection) Search(query string) (*QueryResult, error) {
	var buf bytes.Buffer

	buf.WriteString(query)

	res, err := c.conn.Search(
		c.conn.Search.WithContext(context.Background()),
		c.conn.Search.WithBody(&buf),
		c.conn.Search.WithTrackTotalHits(true),
		c.conn.Search.WithPretty(),
	)

	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.IsError() {
		var e map[string]interface{}
		if err := json.NewDecoder(res.Body).Decode(&e); err != nil {
			return nil, fmt.Errorf("Error parsing the response body: %s", err)
		} else {
			return nil, fmt.Errorf("[%s] %s: %s",
				res.Status(),
				e["error"].(map[string]interface{})["type"],
				e["error"].(map[string]interface{})["reason"],
			)
		}
	}

	var r map[string]interface{}

	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, fmt.Errorf("Error parsing the response body: %s", err)
	}

	log.Printf(
		"[%s] %d hits; took: %dms",
		res.Status(),
		int(r["hits"].(map[string]interface{})["total"].(map[string]interface{})["value"].(float64)),
		int(r["took"].(float64)),
	)

	ret := &QueryResult{
		Took: float64(r["took"].(float64)),
		Hits: []ResultItem{},
	}

	for _, hit := range r["hits"].(map[string]interface{})["hits"].([]interface{}) {

		ret.Hits = append(ret.Hits, ResultItem{
			ID:     string(hit.(map[string]interface{})["_id"].(string)),
			Source: map[string]interface{}(hit.(map[string]interface{})["_source"].(map[string]interface{})),
			Score:  float64(hit.(map[string]interface{})["_score"].(float64)),
			Type:   string(hit.(map[string]interface{})["_type"].(string)),
		})
	}

	return ret, err
}
