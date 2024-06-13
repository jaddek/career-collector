package main

import (
	"encoding/json"
	"log"
	"strconv"
)

type SjsResponeRoot struct {
	TotalCount int      `json:"totalCount"`
	Documents  []SjsJob `json:"documents"`
}

type SjsJob struct {
	Link             string
	Id               int    `json:"id"`
	JobNumber        int    `json:"jobNumber"`
	Name             string `json:"Name"`
	ShortDescription string `json:"ShortDescription"`
	Status           struct {
		Title string `json:"title"`
	} `json:"Status"`
	Type struct {
		Title string `json:"title"`
	} `json:"Type"`
}

func (job SjsJob) GetLink() string {
	return job.Link
}

func (job SjsJob) GetTitle() string {
	return job.Name
}

func (job SjsJob) GetDescription() string {
	return job.ShortDescription
}

func (job SjsJob) GetType() string {
	return job.Type.Title
}

type Sjs struct {
	Container     string
	RequestMethod string
}

func SjsMakeSource() *Sjs {
	source := &Sjs{}
	source.RequestMethod = METHOD_GET

	return source
}

func (sjs *Sjs) GetRequestMethod() string {
	return sjs.RequestMethod
}

func (sjs *Sjs) GetRequestData() []byte {
	return make([]byte, 0)
}

func (sjs *Sjs) GetQuery() string {
	return ""
}

func (sjs *Sjs) GetOnScrapedHandler(data []byte) []IEntity {
	parsedData := &SjsResponeRoot{}
	err := json.Unmarshal(data, parsedData)

	if err != nil {
		log.Fatal(err)
	}

	jobs := make([]IEntity, 0)

	for _, job := range parsedData.Documents {
		job.Link = "https://www.sjs.co.nz/job/" + strconv.Itoa(job.Id)
		jobs = append(jobs, job)
	}

	return jobs
}
