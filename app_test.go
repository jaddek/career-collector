package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func getConfigsDataProvider() []string {
	collection := make([]string, 0)
	collection = append(collection, MTUP, SJS, HNZ, ASB)

	return collection
}

func TestCollectorResolver(t *testing.T) {

	collection := getConfigsDataProvider()

	for _, domain := range collection {
		config := &CollectorConfig{AllowedDomains: "", CacheDir: "", Env: "", Route: "", Domain: domain}
		collector, err := CollectorResolver(config)

		assert.Nil(t, err)
		assert.NotNil(t, collector)
	}

	config := &CollectorConfig{AllowedDomains: "", CacheDir: "", Env: "", Route: "", Domain: "test"}
	source, err := CollectorResolver(config)

	assert.NotNil(t, err)
	assert.Nil(t, source)
}
