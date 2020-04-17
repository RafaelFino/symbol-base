package symbol

import (
	"github.com/RafaelFino/symbol-base/internal/storage/elasticsearch"
	"github.com/RafaelFino/symbol-base/internal/storage/redis"
)

type Symbol struct {
	es *elasticsearch.Connection
	rd *redis.Connection
}

type Config struct {
	EsConfig    *elasticsearch.Config
	RedisConfig *elasticsearch.Config
}

func New(cfg *Config) *Symbol {
	es, err := elasticsearch.New(cfg.EsConfig)

	if err != nil {
		panic(err)
	}

	rd, err :=	redis.New(cfg.RedisConfig)

	if err != nil {
		panic(err)
	}

	return &Symbol{
		es: es, 
		rd: rd,
	}
}

func (s *Symbol) Put(id string, data interface{}) error {
	s.rd.Set(id, data, 0)
}

func (s *Symbol) Search() ([]interface{}, error) {
	ids, err := s.es.Search()
}

func (s *Symbol) GetInfo(ids []string) {map[string]interface{}, error} {

}

func (s *Symbol) GetPrice(ids []string) {map[string]interface{}, error} {

}
