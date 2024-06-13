package main

import (
	"encoding/json"
	"fmt"
	"log"
)

type HNZResponseRoot struct {
	Rasp struct {
		Jobs []struct {
			Jobs []HNZJob `json:"jobs"`
		} `json:"JOBS"`
	} `json:"RASP"`
}

type Test struct {
	Jobs struct {
		Position string `json:"jobs"`
	} `json:"jobs"`
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

func (job HNZJob) GetLink() string {
	return job.Link
}

func (job HNZJob) GetTitle() string {
	return job.Title
}
func (job HNZJob) GetDescription() string {
	return ""
}
func (job HNZJob) GetType() string {
	return job.Type
}

type Hnz struct {
	Container     string
	RequestMethod string
}

func HnzMakeSource() *Hnz {
	source := &Hnz{}
	source.RequestMethod = METHOD_GET

	return source
}

func (hnz *Hnz) GetRequestMethod() string {
	return hnz.RequestMethod
}

func (hnz *Hnz) GetRequestData() []byte {
	return make([]byte, 0)
}

func (hnz *Hnz) GetQuery() string {
	return ""
}

func (hnz *Hnz) GetOnScrapedHandler(data []byte) []IEntity {
	parsedData := &HNZResponseRoot{}
	err := json.Unmarshal(data, parsedData)

	if err != nil {
		log.Fatal(err)
	}

	entities := make([]IEntity, 0)

	for _, job := range parsedData.Rasp.Jobs[0].Jobs {
		fmt.Println(job)
		job.Link = "https://jobs.tewhatuora.govt.nz" + job.Link
		entities = append(entities, job)
	}

	return entities
}
