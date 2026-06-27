package main

import (
	"flag"
	"fmt"
	"go-web-demo/internal/router"

	"go-web-demo/internal/app"
)

func main() {
	env := parseEnvFlag()

	cfg, err := app.InitConfig(env)
	if err != nil {
		fmt.Printf("Init config failed: %v\n", err)
		return
	}

	if err := app.InitDatabase(cfg); err != nil {
		fmt.Printf("Init database failed: %v\n", err)
		return
	}

	h := app.CreateServer(cfg)
	router.SetupRoutes(h, cfg)
	h.Spin()
}

func parseEnvFlag() string {
	var env string
	flag.StringVar(&env, "env", "", "Environment: dev, test, prod (default: dev)")
	flag.Parse()
	return env
}
