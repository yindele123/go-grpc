package register

import (
	"github.com/hashicorp/consul/api"
)

type Register interface {
	Register(address string, port int, name string, tags []string, id string) error
	Deregister(serviceId string) error
	GetAllAervice() map[string]*api.AgentService
	FilterService(filter string) map[string]*api.AgentService
}
