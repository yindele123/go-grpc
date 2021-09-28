package register

type Register interface {
	Register(address string, port int, name string, tags []string, id string) error
	Deregister(serviceId string) error
}
