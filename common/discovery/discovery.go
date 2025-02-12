package discovery

import (
	"context"
	"fmt"
	"math/rand"
)

type Registry interface {
	Register(ctx context.Context, instanceID, serverName, hostPort string) error
	Deregister(ctx context.Context, instanceID, serverName string) error
	Discover(ctx context.Context, serverName string) ([]string, error)
	HealthCheck(instanceID, serverName string) error
}

func GenerateInstanceID(serviceName string) string {
	return fmt.Sprintf("%s-%d", serviceName, rand.Intn(10000))
}
