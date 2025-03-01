package immem

import (
	"context"
	"errors"
	"sync"
	"time"
)

type Registry struct {
	sync.RWMutex
	addrs map[string]map[string]*serviceInstance
}

type serviceInstance struct {
	hostPort   string
	lastActive time.Time
}

func NewRegistry() *Registry {
	return &Registry{
		addrs: map[string]map[string]*serviceInstance{},
	}
}

func (r *Registry) Register(ctx context.Context, instanceID, serviceName, hostPort string) error {
	r.Lock()
	defer r.Unlock()

	if _, exists := r.addrs[serviceName]; !exists {
		r.addrs[serviceName] = map[string]*serviceInstance{}
	}

	r.addrs[serviceName][instanceID] = &serviceInstance{
		hostPort:   hostPort,
		lastActive: time.Now(),
	}
	return nil
}

func (r *Registry) Deregister(ctx context.Context, instanceID, serviceName string) error {
	r.Lock()
	defer r.Unlock()

	if _, exists := r.addrs[serviceName]; !exists {
		return nil
	}
	delete(r.addrs[serviceName], instanceID)

	return nil
}

func (r *Registry) Discover(ctx context.Context, serviceName string) ([]string, error) {
	r.Lock()
	defer r.Unlock()

	if len(r.addrs[serviceName]) == 0 {
		return nil, errors.New("no instances available")

	}

	var result []string

	for _, i := range r.addrs[serviceName] {
		result = append(result, i.hostPort)
	}
	return result, nil
}

func (r *Registry) ServiceAddress(ctx context.Context, serviceName string) ([]string, error) {
	r.RLock()
	defer r.RUnlock()

	if len(r.addrs[serviceName]) == 0 {
		return nil, errors.New("no service addressf found")
	}

	var result []string
	for _, i := range r.addrs[serviceName] {
		if i.lastActive.Before(time.Now().Add(-5 * time.Second)) {
			continue
		}
		result = append(result, i.hostPort)
	}

	return result, nil
}

func (r *Registry) HealthCheck(instanceID, serviceName string) error {
	r.Lock()
	defer r.Unlock()

	if _, ok := r.addrs[serviceName]; !ok {
		return errors.New("service is not registered yet")
	}

	if _, ok := r.addrs[serviceName][instanceID]; !ok {
		return errors.New("service instance is not registered yet")
	}

	r.addrs[serviceName][instanceID].lastActive = time.Now()

	return nil
}
