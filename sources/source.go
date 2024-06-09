package sources

import (
	"github.com/gocolly/colly/v2"
)

const (
	SJS  = "sjs"
	ASB  = "asb"
	PROD = "prod"
	DEV  = "dev"
)

type Collector interface {
	Collect()
	GetJobs() []JobEntity
}

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
	config CollectorConfig
	Jobs   []JobEntity
}

func (source *SourceCollector) getCollector() *colly.Collector {
	options := make([]colly.CollectorOption, 0)
	options = append(options, colly.AllowedDomains(source.config.AllowedDomains))

	if source.config.isDev() {
		options = append(options, colly.CacheDir(source.config.CacheDir))
	}

	collector := colly.NewCollector(options...)

	return collector
}

type Job interface {
	GetLink() string
	GetTitle() string
	GetDescription() string
	GetType() string
}

type JobEntity struct {
	Link        string `json:"link"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

type AsbJobDecorator struct {
	Job *AsbJob
}

func (job *AsbJobDecorator) GetLink() string {
	return job.Job.Link
}

func (job *AsbJobDecorator) GetTitle() string {
	return job.Job.Title
}
func (job *AsbJobDecorator) GetDescription() string {
	return job.Job.Description
}
func (job *AsbJobDecorator) GetType() string {
	return job.Job.Details.Type
}

type SjsJobDecorator struct {
	Job *SjsJob
}

func (job *SjsJobDecorator) GetLink() string {
	return job.Job.Link
}

func (job *SjsJobDecorator) GetTitle() string {
	return job.Job.Name
}
func (job *SjsJobDecorator) GetDescription() string {
	return job.Job.ShortDescription
}
func (job *SjsJobDecorator) GetType() string {
	return job.Job.Type.Title
}

func SourceResolver(config CollectorConfig) Collector {
	switch config.Domain {
	case ASB:
		return &Asb{SourceCollector{config: config}}
	case SJS:
		return &Sjs{SourceCollector{config: config}}
	default:
		return nil
	}
}
