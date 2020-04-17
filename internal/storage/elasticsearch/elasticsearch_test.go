package elasticsearch

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/oklog/ulid"
)

func prepare() *Connection {
	cfg := &Config{
		Addresses: []string{"http://localhost:9200"},
	}

	es, err := New(cfg)

	if err != nil {
		panic(err)
	}

	return es
}

type testDocument struct {
	ID       string
	When     time.Time
	Data     string
	ToReturn string
}

func TestInsert(t *testing.T) {
	es := prepare()

	docQty := 200
	newDocs := make([]Document, docQty)
	index := "test_index"

	for i := 0; i < docQty; i++ {
		id := GetULID()

		toReturn := "false"
		if i%2 == 0 {
			toReturn = "true"
		}

		newDocs[i] = Document{
			ID: id,
			Body: testDocument{
				ID:       id,
				When:     time.Now(),
				Data:     fmt.Sprintf("data count: %d", i),
				ToReturn: toReturn,
			},
			Index: index,
		}
	}

	err := es.Write(index, newDocs)

	if err != nil {
		t.Error(err)
	}
}

const query string = `
	{
		"query": {
			"match": {
				"ToReturn": "true"
			}
		}
	}
	`

func TestSearch(t *testing.T) {
	es := prepare()

	result, err := es.Search(query)

	if err != nil {
		t.Error(err)
	}

	fmt.Printf("\tResult: %v", result)

	if len(result.Hits) == 0 {
		t.Errorf("no hits on query: %s", query)
	}

	for _, hit := range result.Hits {
		fmt.Printf("\tID: %s Score: %f Type: %s Source: %v\n", hit.ID, hit.Score, hit.Type, hit.Source)
	}
}

func GetULID() string {
	t := time.Now()
	entropy := rand.New(rand.NewSource(t.UnixNano()))

	return ulid.MustNew(ulid.Timestamp(t), entropy).String()
}
