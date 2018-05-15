package registry

import (
	"log"
	"sync"

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
	log.Println("All services ready!", r.services)
	var wg sync.WaitGroup
	wg.Add(len(r.services))
	for _, service := range r.services {
		go func(s interfaces.Service) {
			s.OnAllServicesRegistered(r)
			wg.Done()
		}(service)
	}
	wg.Wait()
	log.Println("All services done doing their thing")
}

func (r *registry) AllControllersRegistered() {
	log.Println("All controllers ready!", r.controllers)
	var wg sync.WaitGroup
	wg.Add(len(r.controllers))
	for _, controller := range r.controllers {
		go func(c interfaces.Controller) {
			c.OnAllControllersRegistered(r)
			wg.Done()
		}(controller)
	}
	wg.Wait()
	log.Println("All controlers done doing their thing")
}
