package interfaces

// Registry is a interface of all the methods required to be a Registry
type Registry interface {
	RegisterService(s Service, name string)
	RegisterController(s Controller, name string)
	GetService(name string) Service
	GetController(name string) Controller
	AllServicesRegistered()
	AllControllersRegistered()
}
