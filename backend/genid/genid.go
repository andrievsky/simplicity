package genid

import (
	"github.com/bwmarrin/snowflake"
	"math/rand"
)

type Provider interface {
	Generate() string
	Validate(string) error
}

type SnowflakeProvider struct {
	instance *snowflake.Node
}

func NewSnowflakeProvider(nodeID int64) (*SnowflakeProvider, error) {
	node, err := snowflake.NewNode(nodeID)
	if err != nil {
		return nil, err
	}
	return &SnowflakeProvider{instance: node}, nil
}

func (p *SnowflakeProvider) Generate() string {
	return p.instance.Generate().String()
}

func (p *SnowflakeProvider) Validate(id string) error {
	_, err := snowflake.ParseString(id)
	if err != nil {
		return err
	}
	return nil
}

const dict = "abcdefghijklmnopqrstuvwxyz0123456789"

func GeneratePartialID(n int) string {
	s := make([]byte, n)
	for i := range s {
		s[i] = dict[rand.Intn(len(dict))]
	}
	return string(s)
}
