package main

import (
	"errors"
	"log"
	"os"
	"runtime"
	"strings"

	"github.com/gocolly/colly/v2"
)

const (
	SJS         = "sjs"
	ASB         = "asb"
	HNZ         = "hnz"
	MTUP        = "mtup"
	PROD        = "prod"
	DEV         = "dev"
	METHOD_POST = "POST"
	METHOD_GET  = "GET"
)

type ISource interface {
	GetRequestData() []byte
	GetRequestMethod() string
	GetQuery() string
	GetOnScrapedHandler(data []byte) []IEntity
}

type IJsonSource interface {
	ISource
	GetRequestData() []byte
}

type IHtmlSource interface {
	ISource
	GetContainer() string
	GetOnHtmlHandler(e *colly.HTMLElement)
}

type ISourceCollector interface {
	Collect()
	GetEntities() []IEntity
}

type IEntity interface {
	GetLink() string
	GetTitle() string
	GetDescription() string
	GetType() string
}

type EventEntity struct {
	Link        string `json:"link"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

func runDomain(collector ISourceCollector, channel chan<- []IEntity) {
	runtime.LockOSThread()
	defer wg.Done()

	collector.Collect()
	entities := collector.GetEntities()

	channel <- entities

	runtime.UnlockOSThread()
}

func getDomainConfig(domain string) CollectorConfig {
	prefix := strings.ToUpper(domain)

	config := CollectorConfig{
		AllowedDomains: os.Getenv(prefix + "_ALLOWED_DOMAINS"),
		CacheDir:       os.Getenv(prefix + "_CACHE_DIR"),
		Env:            os.Getenv("APP_ENV"),
		Route:          os.Getenv(prefix + "_ROUTE"),
		Domain:         domain,
	}

	if config.Route == "" {
		log.Printf("Route for domain %s is empty", domain)
	}

	return config
}

func run() []IEntity {
	domains, workers := getDomains()
	channel := make(chan []IEntity, workers)

	wg.Add(workers)

	for _, domain := range domains {
		collector, err := CollectorResolver(getDomainConfig(domain))

		if err != nil {
			log.Panic(err)

			continue
		}

		if collector == nil {
			log.Panic("Unsupported source")

			continue
		}

		go runDomain(collector, channel)
	}

	wg.Wait()

	close(channel)

	entities := make([]IEntity, 0)
	for j := range channel {
		entities = append(entities, j...)
	}

	return entities
}

func getDomains() ([]string, int) {
	activeDomains := os.Getenv("ACTIVE_DOMAINS")
	domains := strings.Split(activeDomains, ",")

	return domains, len(domains)
}

func CollectorResolver(collectorConfig CollectorConfig) (ISourceCollector, error) {
	switch collectorConfig.Domain {
	case MTUP:
		source := MtapMakeSource()

		return MakeSourceJsonCollector(source, collectorConfig), nil
	case SJS:
		source := SjsMakeSource()

		return MakeSourceJsonCollector(source, collectorConfig), nil
	case HNZ:
		source := HnzMakeSource()

		return MakeSourceJsonCollector(source, collectorConfig), nil
	case ASB:
		source := AsbMakeSource()

		return MakeSourcePaginationalHtmlCollector(source, collectorConfig), nil
	}

	return nil, errors.New("unsupported source")
}
