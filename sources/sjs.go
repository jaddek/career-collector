package sources

import (
	"encoding/json"
	"log"

	"github.com/gocolly/colly/v2"
)

type SjsRespone struct {
	TotalCount int       `json:"totalCount"`
	Documents  []*SjsJob `json:"documents"`
}

type SjsJob struct {
	Link             string
	Id               int       `json:"id"`
	JobNumber        int       `json:"jobNumber"`
	Name             string    `json:"Name"`
	ShortDescription string    `json:"ShortDescription"`
	Status           SjsStatus `json:"Status"`
	Type             SjsType   `json:"Type"`
}

type SjsType struct {
	Title string `json:"title"`
}

type SjsStatus struct {
	Title string `json:"title"`
}

type Sjs struct {
	SourceCollector
}

func (sjs *Sjs) getIndexCollector() *colly.Collector {

	collector := sjs.getCollector()

	collector.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL.String())
	})

	collector.OnError(func(r *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	collector.OnScraped(func(r *colly.Response) {
		result := &SjsRespone{}
		err := json.Unmarshal(r.Body, result)

		if err != nil {
			log.Fatal(err)
		}

		docs := result.Documents

		for _, doc := range docs {
			doc.Link = "https://www.sjs.co.nz/job/" + (string)(doc.Id)

			job := sjs.normaliseJob(doc)
			sjs.Jobs = append(sjs.Jobs, job)
		}
	})

	return collector
}

func (sjs *Sjs) normaliseJob(job *SjsJob) JobEntity {
	sjob := &SjsJobDecorator{Job: job}

	return JobEntity{
		Link:        sjob.GetLink(),
		Title:       sjob.GetTitle(),
		Description: sjob.GetDescription(),
		Type:        sjob.GetType(),
	}
}

func (sjs *Sjs) Collect() {
	collector := sjs.getIndexCollector()

	collector.Visit(sjs.config.Route)
}

func (sjs *Sjs) GetJobs() []JobEntity {
	return sjs.Jobs
}
