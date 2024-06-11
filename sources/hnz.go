package sources

import (
	"encoding/json"
	"log"

	"github.com/gocolly/colly/v2"
)

type Hnz struct {
	SourceCollector
}

type HNZResponse struct {
	Rasp *HNZJobsCollection `json:"RASP"`
}

type HNZJobsCollection struct {
	Jobs []map[string]interface{} `json:"JOBS"`
}

type HNZJobsCollectionJobs struct {
	Jobs []*HNZJob `json:"jobs"`
}

type HNZJob struct {
	Link       string `json:"href"`
	Title      string `json:"positionTitle"`
	Island     string `json:"island"`
	Region     string `json:"region"`
	District   string `json:"district"`
	Closedate  string `json:"closedate"`
	Dateposted string `json:"dateposted"`
	Type       string `json:"jobtype"`
}

func (hnz *Hnz) normaliseJob(job *HNZJob) JobEntity {
	ajob := &HNZJobDecorator{Job: job}

	return JobEntity{
		Link:        ajob.GetLink(),
		Title:       ajob.GetTitle(),
		Description: "",
		Type:        ajob.GetType(),
	}
}

func (hnz *Hnz) getIndexCollector() *colly.Collector {
	collector := hnz.getCollector()

	collector.OnRequest(func(r *colly.Request) {
		log.Println("Visiting", r.URL.String())
	})

	collector.OnError(func(r *colly.Response, err error) {
		log.Println("Something went wrong:", err)
	})

	collector.OnScraped(func(r *colly.Response) {
		response := &HNZResponse{}
		err := json.Unmarshal(r.Body, response)

		if err != nil {
			log.Fatal(err)

			return
		}

		newJson, _ := json.Marshal(response.Rasp.Jobs[0])

		result := &HNZJobsCollectionJobs{}
		err2 := json.Unmarshal(newJson, result)

		if err2 != nil {
			log.Fatal(err)

			return
		}

		docs := result.Jobs

		for _, doc := range docs {
			doc.Link = "https://jobs.tewhatuora.govt.nz/" + (string)(doc.Link)

			job := hnz.normaliseJob(doc)
			hnz.Jobs = append(hnz.Jobs, job)
		}
	})

	return collector
}

func (hnz *Hnz) Collect() {
	indexCollector := hnz.getIndexCollector()

	indexCollector.Visit(hnz.config.Route)
}

func (hnz *Hnz) GetJobs() []JobEntity {
	return hnz.Jobs
}
