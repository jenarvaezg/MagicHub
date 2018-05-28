package registry

import (
	"log"

	"github.com/jenarvaezg/MagicHub/interfaces"
)

type registry struct {
	services    map[string]interfaces.Service
	controllers map[string]interfaces.Controller
}

// NewRegistry returns a Registry for services to register themselves and get other services
func NewRegistry() interfaces.Registry {
	r := &registry{}
	r.services = make(map[string]interfaces.Service)
	r.controllers = make(map[string]interfaces.Controller)
	return r
}

func (r *registry) RegisterService(s interfaces.Service, name string) {
	r.services[name] = s
}

func (r *registry) RegisterController(s interfaces.Controller, name string) {
	r.controllers[name] = s
}

func (r *registry) GetService(name string) interfaces.Service {
	s, ok := r.services[name]
	if !ok {
		log.Panicf("Service %s not found in registry", name)
	}

	return s
}

func (r *registry) GetController(name string) interfaces.Controller {
	s, ok := r.controllers[name]
	if !ok {
		log.Panicf("Controller %s not found in registry", name)
	}

	return s
}

func (r *registry) AllServicesRegistered() {
	for _, s := range r.services {
		s.OnAllServicesRegistered(r)
	}
	log.Println("All services ready")
}

func (r *registry) AllControllersRegistered() {
	for _, c := range r.controllers {
		c.OnAllControllersRegistered(r)
	}
	log.Println("All controllers ready")
}
