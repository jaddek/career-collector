package main

type NormalisedEntity struct {
	Link        string `json:"link"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Type        string `json:"type"`
}

func NormalizeSourceEntities(entities []IEntity) []*NormalisedEntity {
	list := make([]*NormalisedEntity, 0)

	for _, entity := range entities {
		list = append(list, &NormalisedEntity{
			Link:        entity.GetLink(),
			Title:       entity.GetTitle(),
			Description: entity.GetDescription(),
			Type:        entity.GetType(),
		})
	}

	return list
}
