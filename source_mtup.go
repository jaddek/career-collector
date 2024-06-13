package main

import (
	"encoding/json"
	"log"
)

type MtupResponseRoot struct {
	Data struct {
		Result struct {
			Edges []MtupEdge `json:"edges"`
		} `json:"result"`
	} `json:"data"`
}

type MtupEdge struct {
	Node struct {
		DateTime    string `json:"dateTime"`
		Description string `json:"description"`
		EventType   string `json:"eventType"`
		Id          string `json:"id"`
		Title       string `json:"title"`
		Link        string `json:"eventUrl"`
		Venue       struct {
			Name string `json:"name"`
		} `json:"venue"`
	} `json:"node"`
}

func (edge MtupEdge) GetLink() string {
	return edge.Node.Link
}

func (edge MtupEdge) GetTitle() string {
	return edge.Node.Title
}

func (edge MtupEdge) GetDescription() string {
	return edge.Node.Description
}

func (edge MtupEdge) GetType() string {
	return edge.Node.Venue.Name
}

type Mtup struct {
	Container     string
	RequestMethod string
}

func MtapMakeSource() *Mtup {
	source := &Mtup{}
	source.RequestMethod = METHOD_POST

	return source
}

func (mtup *Mtup) GetRequestMethod() string {
	return mtup.RequestMethod
}

func (mtup *Mtup) GetRequestData() []byte {
	requestData := map[string]interface{}{
		"operationName": "recommendedEventsWithSeries",
		"variables": map[string]interface{}{
			"first":                   20,
			"lat":                     "-41.279998779296875",
			"lon":                     "174.77999877929688",
			"startDateRange":          "2024-06-11T18:06:55-04:00[US/Eastern]",
			"numberOfEventsForSeries": 5,
			"seriesStartDate":         "2024-06-11",
			"sortField":               "DATETIME",
			"doConsolidateEvents":     true,
			"doPromotePaypalEvents":   true,
			"indexAlias":              "popular_events_nearby_current",
		},
		"extensions": map[string]map[string]interface{}{
			"persistedQuery": {
				"version":    1,
				"sha256Hash": "0f0332e9a4b01456580c1f669f26edc053d50382b3e338d5ca580f194a27feab",
			},
		},
	}

	raw, err := json.Marshal(requestData)

	if err != nil {
		log.Fatal(err)
	}

	return raw
}

func (mtup *Mtup) GetQuery() string {
	return ""
}

func (mtup *Mtup) GetOnScrapedHandler(data []byte) []IEntity {
	parsedData := &MtupResponseRoot{}
	err := json.Unmarshal(data, parsedData)

	if err != nil {
		log.Fatal(err)
	}

	edges := make([]IEntity, 0)

	for _, edge := range parsedData.Data.Result.Edges {
		edges = append(edges, edge)
	}

	return edges
}
