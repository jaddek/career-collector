package main

import (
	"log"

	"github.com/gocolly/colly/v2"
)

type CollectorConfig struct {
	AllowedDomains string
	CacheDir       string
	Env            string
	Route          string
	Domain         string
}

func (cc *CollectorConfig) isDev() bool {
	return cc.Env == DEV
}

type SourceCollector struct {
	Source   ISource
	Config   *CollectorConfig
	Entities []IEntity
}

func (source *SourceCollector) getCollector() *colly.Collector {
	options := make([]colly.CollectorOption, 0)
	options = append(options, colly.AllowedDomains(source.Config.AllowedDomains))

	if source.Config.isDev() {
		options = append(options, colly.CacheDir(source.Config.CacheDir))
	}

	collector := colly.NewCollector(options...)

	collector.OnRequest(func(r *colly.Request) {
		log.Println("Visiting ", source.Config.Domain, r.URL.String())
	})

	collector.OnError(func(r *colly.Response, err error) {
		log.Println("Something went wrong:", source.Config.Domain, err)
	})

	return collector
}

func (source *SourceCollector) GetEntities() []IEntity {
	return source.Entities
}

type SourceJsonCollector struct {
	SourceCollector
	Source IJsonSource
}

func MakeSourceJsonCollector(source IJsonSource, collectorConfig *CollectorConfig) *SourceJsonCollector {
	collector := &SourceJsonCollector{}
	collector.Source = source
	collector.Config = collectorConfig

	return collector
}

func (source *SourceJsonCollector) Collect() {
	collector := source.getCollector()

	collector.OnScraped(func(r *colly.Response) {
		source.Entities = source.Source.GetOnScrapedHandler(r.Body)
	})

	switch source.Source.GetRequestMethod() {
	case METHOD_POST:
		source.applyPostStrategy(collector)
	case METHOD_GET:
		source.applyGetStrategy(collector)
	}
}

func (source *SourceJsonCollector) applyGetStrategy(collector *colly.Collector) {
	collector.Visit(source.Config.Route)
}

func (source *SourceJsonCollector) applyPostStrategy(collector *colly.Collector) {
	collector.PostRaw(source.Config.Route, source.Source.GetRequestData())
}

type SourceHtmlCollector struct {
	SourceCollector
	Source IHtmlSource
}

func MakeSourceHtmlCollector(source IHtmlSource, collectorConfig *CollectorConfig) *SourceHtmlCollector {
	collector := &SourceHtmlCollector{}
	collector.Source = source
	collector.Config = collectorConfig

	return collector
}

func (source *SourceHtmlCollector) Collect() {
	collector := source.getCollector()

	collector.OnHTML(source.Source.GetContainer(), func(e *colly.HTMLElement) {
		source.Source.GetOnHtmlHandler(e)
	})

	collector.OnScraped(func(r *colly.Response) {
		source.Entities = source.Source.GetOnScrapedHandler(r.Body)
	})

	switch source.Source.GetRequestMethod() {
	case "POST":
		source.applyPostStrategy(collector)
	case "GET":
		source.applyGetStrategy(collector)
	}
}

func (source *SourceHtmlCollector) applyGetStrategy(collector *colly.Collector) {
	collector.Visit(source.Config.Route)
}

func (source *SourceHtmlCollector) applyPostStrategy(collector *colly.Collector) {
	collector.PostRaw(source.Config.Route, source.Source.GetRequestData())
}

type SourcePaginationalHtmlCollector struct {
	SourceCollector
	Source IHtmlSource
}

func MakeSourcePaginationalHtmlCollector(
	source IHtmlSource,
	collectorConfig *CollectorConfig) *SourcePaginationalHtmlCollector {
	collector := &SourcePaginationalHtmlCollector{}
	collector.Source = source
	collector.Config = collectorConfig

	return collector
}

func (source *SourcePaginationalHtmlCollector) Collect() {
	collector := source.getCollector()
	collector.OnHTML(".page-links a", func(e *colly.HTMLElement) {
		collector.Visit(e.Request.AbsoluteURL(e.Attr("href")))
	})

	collector.OnHTML(source.Source.GetContainer(), func(e *colly.HTMLElement) {
		source.Source.GetOnHtmlHandler(e)
	})

	collector.OnScraped(func(r *colly.Response) {
		source.Entities = source.Source.GetOnScrapedHandler(r.Body)
	})

	switch source.Source.GetRequestMethod() {
	case "POST":
		source.applyPostStrategy(collector)
	case "GET":
		source.applyGetStrategy(collector)
	}
}

func (source *SourcePaginationalHtmlCollector) applyGetStrategy(collector *colly.Collector) {
	collector.Visit(source.Config.Route)
}

func (source *SourcePaginationalHtmlCollector) applyPostStrategy(collector *colly.Collector) {
	collector.PostRaw(source.Config.Route, source.Source.GetRequestData())
}
