package tests

import (
	"github.com/docker/docker/client"
	"go.uber.org/zap"
	"polar-sky/log"
	"polar-sky/task"
	"testing"
)

func TestContainer(t *testing.T) {
	testLogger, _ := zap.NewDevelopment()
	logger := testLogger.Sugar()
	log.Logger = *logger

	c := task.Config{
		Name:  "test-task-1",
		Image: "postgres:12",
		Env: []string{
			"POSTGRES_USER=cube",
			"POSTGRES_PASSWORD=secret",
		},
	}

	dc, _ := client.NewClientWithOpts(client.FromEnv, client.WithVersion("1.41"))
	d := task.Docker{
		Client: dc,
		Config: &c,
	}

	dr := d.Run()
	if dr.Error != nil {
		logger.Errorf("Failed run task %v with docker %v, error: %v\n", c, d, dr.Error)
	}

	dr = d.Stop(dr.ContainerID)
	if dr.Error != nil {
		logger.Errorf("Failed stop task %v with docker %v, error: %v\n", c, d, dr.Error)
	}
}
