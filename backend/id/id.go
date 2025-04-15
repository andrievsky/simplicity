package id

import (
	"errors"
	"math/rand"
	"strings"
)

type Provider interface {
	Next() (string, error)
	Validate(id string) error
}

type RandomNameProvider struct {
	wordGroups  [][]string
	existIds    map[string]struct{}
	separator   string
	maxAttempts int
}

func NewRandomNameProvider(wordGroups [][]string, existIds map[string]struct{}) *RandomNameProvider {
	return &RandomNameProvider{wordGroups, existIds, "-", 100}
}

func (t *RandomNameProvider) generateRandomId() string {
	sb := strings.Builder{}
	first := true
	for _, group := range t.wordGroups {
		if first {
			first = false
		} else {
			sb.WriteString(t.separator)
		}
		sb.WriteString(group[rand.Intn(len(group))])
	}
	return sb.String()
}

func (t *RandomNameProvider) Next() (string, error) {
	id := t.generateRandomId()
	attempts := 1
	if t.Validate(id) != nil {
		if attempts >= t.maxAttempts {
			return "", errors.New("max attempts reached to generate a random id")
		}
		id = t.generateRandomId()
		attempts++
	}

	t.existIds[id] = struct{}{}
	return id, nil
}

func (t *RandomNameProvider) Validate(id string) error {
	if _, ok := t.existIds[id]; ok {
		return errors.New("id already exists")
	}
	return nil
}
