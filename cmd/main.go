package main

import (
	"github.com/bsonger/devflow-app-service/pkg/infra/config"
	"github.com/bsonger/devflow-app-service/pkg/router"
	"github.com/bsonger/devflow-service-common/bootstrap"
	"github.com/bsonger/devflow-service-common/observability"
)

func main() {
	err := bootstrap.Run(bootstrap.Options[config.Config, router.Options, string]{
		Name: "app-service",
		RouteOptions: router.Options{
			ServiceName:   "app-service",
			EnableSwagger: true,
			Modules: []router.Module{
				router.ModuleProject,
				router.ModuleApplication,
			},
		},
		Load:        config.Load,
		InitRuntime: config.InitRuntime,
		NewRouter: func(opts router.Options) bootstrap.Runner {
			return router.NewRouterWithOptions(opts)
		},
		ResolveConfigPort: func(cfg *config.Config) int {
			if cfg != nil && cfg.Server != nil {
				return cfg.Server.Port
			}
			return 0
		},
		StartMetricsServer: observability.StartMetricsServer,
		StartPprofServer:   observability.StartPprofServer,
		PortEnv:            "APP_SERVICE_PORT",
		DefaultPort:        8081,
		MetricsPortEnv:     "APP_SERVICE_METRICS_PORT",
		PprofPortEnv:       "APP_SERVICE_PPROF_PORT",
	})
	if err != nil {
		panic(err)
	}
}
