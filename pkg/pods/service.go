package pods

import (
	"github.com/tomba7/craftd/pkg/client"
)

var Service podsServiceInterface = &podsService{
	lister: newPodLister(client.Client()),
}

type podsServiceInterface interface {
	Get(StatusFilters, *int) []*Pod
}

type podsService struct{
	lister *podLister
}

func (s *podsService) Get(filters StatusFilters, minutes *int) []*Pod {
	return s.lister.get("default", filters, minutes)
}
