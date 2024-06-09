package sources

import (
	"log"

	"github.com/gocolly/colly/v2"
)

type Asb struct {
	SourceCollector
}

type AsbJob struct {
	Link        string
	Title       string `selector:"h1"`
	Description string `selector:".description"`
	Details     AsbJobDetail
}

type AsbJobDetail struct {
	Type string `selector:".jd tr.eq(3) td.eq(2)"`
}

func (asb *Asb) getDetailCollector(indexCollector *colly.Collector) *colly.Collector {
	collector := indexCollector.Clone()

	collector.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL.String())
	})

	collector.OnError(func(r *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	collector.OnHTML(".container", func(e *colly.HTMLElement) {
		asbJob := &AsbJob{}
		e.Unmarshal(asbJob)
		asbJob.Link = e.Request.AbsoluteURL(e.Request.URL.Path)

		job := asb.normaliseJob(asbJob)
		asb.Jobs = append(asb.Jobs, job)
	})

	return collector
}

func (asb *Asb) normaliseJob(job *AsbJob) JobEntity {
	ajob := &AsbJobDecorator{Job: job}

	return JobEntity{Link: ajob.GetLink(), Title: ajob.GetTitle()}
}

func (asb *Asb) getIndexCollector() *colly.Collector {
	collector := asb.getCollector()

	collector.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL.String())
	})

	collector.OnError(func(r *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	collector.OnHTML(".page-links a", func(e *colly.HTMLElement) {
		collector.Visit(e.Request.AbsoluteURL(e.Attr("href")))
	})

	return collector
}

func (asb *Asb) collectJobLinks(
	indexCollector *colly.Collector,
	detailCollector *colly.Collector) {

	indexCollector.OnHTML(".job a", func(e *colly.HTMLElement) {
		detailCollector.Visit(e.Request.AbsoluteURL(e.Attr("href")))
	})

	indexCollector.Visit(asb.config.Route)
}

func (asb *Asb) Collect() {
	indexCollector := asb.getIndexCollector()
	detailCollector := asb.getDetailCollector(indexCollector)

	asb.collectJobLinks(indexCollector, detailCollector)
}

func (asb *Asb) GetJobs() []JobEntity {
	return asb.Jobs
}
