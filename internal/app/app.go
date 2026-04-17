package app

import (
	"github.com/meshenka/npcgenerator/internal/app/query"
)

type Application struct {
	Queries Queries
}

type Queries struct {
	GetDescriptors query.GetDescriptorsHandler
}

func NewApplication() Application {
	return Application{
		Queries: Queries{
			GetDescriptors: query.NewGetDescriptorsHandler(),
		},
	}
}
