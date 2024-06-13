package main

import "github.com/gocolly/colly/v2"

type AsbJob struct {
	Link        string
	Title       string `selector:".title"`
	Description string `selector:".description"`
	Details     string `selector:".details"`
}

func (job AsbJob) GetLink() string {
	return job.Link
}

func (job AsbJob) GetTitle() string {
	return job.Title
}

func (job AsbJob) GetDescription() string {
	return job.Description
}

func (job AsbJob) GetType() string {
	return job.Details
}

type Asb struct {
	Container     string
	RequestMethod string
	Jobs          []IEntity
}

func AsbMakeSource() *Asb {
	source := &Asb{}
	source.RequestMethod = METHOD_GET
	source.Container = ".job"

	return source
}

func (asb *Asb) GetContainer() string {
	return asb.Container
}

func (asb *Asb) GetRequestMethod() string {
	return asb.RequestMethod
}

func (asb *Asb) GetRequestData() []byte {
	return make([]byte, 0)
}

func (asb *Asb) GetQuery() string {
	return ""
}

func (asb *Asb) GetOnHtmlHandler(e *colly.HTMLElement) {
	job := AsbJob{}
	e.Unmarshal(&job)
	job.Link = e.Request.AbsoluteURL(e.Request.URL.Path)

	asb.Jobs = append(asb.Jobs, job)
}

func (asb *Asb) GetOnScrapedHandler(data []byte) []IEntity {
	return asb.Jobs
}
