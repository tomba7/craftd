package pods

import (
	"github.com/tomba7/craftd/pkg/client"
)

var Service podsServiceInterface = &podsService{
	lister: newPodLister(client.Client()),
}

type podsServiceInterface interface {
	Get(StatusFilters, *int) []*Pod
	Delete(StatusFilters, *int) []*Pod
}

type podsService struct{
	lister *podLister
}

func (s *podsService) Get(filters StatusFilters, minutes *int) []*Pod {
	return s.lister.get("default", filters, minutes)
}

func (s *podsService) Delete(filters StatusFilters, minutes *int) []*Pod {
	return s.lister.get("default", filters, minutes)
}
