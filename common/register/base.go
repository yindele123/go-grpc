package register

type Register interface {
	Register(address string, port int, name string, tags interface{}, id string) error
	Deregister(serviceId string) error
}
