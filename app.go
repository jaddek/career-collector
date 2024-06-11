package main

import (
	"log"
	"os"
	"runtime"
	"strings"

	source "github.com/jaddek/tapnngo/sources"
)

func runDomain(source source.Collector, channel chan<- []source.JobEntity) {
	runtime.LockOSThread()
	defer wg.Done()

	source.Collect()
	jobs := source.GetJobs()

	channel <- jobs

	runtime.UnlockOSThread()
}

func getDomainConfig(domain string) source.CollectorConfig {
	prefix := strings.ToUpper(domain)

	config := source.CollectorConfig{
		AllowedDomains: os.Getenv(prefix + "_ALLOWED_DOMAINS"),
		CacheDir:       os.Getenv(prefix + "_CACHE_DIR"),
		Env:            os.Getenv("APP_ENV"),
		Route:          os.Getenv(prefix + "_ROUTE"),
		Domain:         domain,
	}

	return config
}

func run() []source.JobEntity {
	domains, workers := getDomains()
	channel := make(chan []source.JobEntity, workers)

	wg.Add(workers)

	for _, domain := range domains {
		source := source.SourceResolver(getDomainConfig(domain))

		if source == nil {
			log.Panic("Unsupported source")
			continue
		}

		go runDomain(source, channel)
	}

	wg.Wait()

	close(channel)

	jobs := make([]source.JobEntity, 0)
	for j := range channel {
		jobs = append(jobs, j...)
	}

	return jobs
}

func getDomains() ([]string, int) {
	activeDomains := os.Getenv("ACTIVE_DOMAINS")
	domains := strings.Split(activeDomains, ",")

	return domains, len(domains)
}
