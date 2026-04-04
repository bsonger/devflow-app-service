package main

import (
	"github.com/bsonger/devflow-app-service/pkg/router"
	"github.com/bsonger/devflow-app-service/platform/shared/bootstrap"
)

func main() {
	err := bootstrap.Run(bootstrap.Options{
		Name: "app-service",
		RouteOptions: router.Options{
			ServiceName:   "app-service",
			EnableSwagger: true,
			Modules: []router.Module{
				router.ModuleProject,
				router.ModuleApplication,
			},
		},
		PortEnv:        "APP_SERVICE_PORT",
		DefaultPort:    8081,
		MetricsPortEnv: "APP_SERVICE_METRICS_PORT",
		PprofPortEnv:   "APP_SERVICE_PPROF_PORT",
	})
	if err != nil {
		panic(err)
	}
}
