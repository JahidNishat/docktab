package main

import (
	"fmt"
	"os"

	"github.com/JahidNishat/docktab/internal/cli"
	"github.com/JahidNishat/docktab/internal/config"
	"github.com/JahidNishat/docktab/internal/docker"
	"github.com/JahidNishat/docktab/internal/logger"
	"github.com/JahidNishat/docktab/internal/table"
)

func main() {
	_ = config.LoadEnv()

	cfg := config.New()
	log := logger.New(cfg.Debug)
	renderer := table.NewRenderer()

	dockerClient, err := docker.NewClient(cfg.DockerHost)
	if err != nil {
		fmt.Fprintf(os.Stderr, "docktab: %v\n", err)
		os.Exit(1)
	}

	root := cli.NewRootCommand(dockerClient, renderer, log)
	if err := root.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "docktab: %v\n", err)
		os.Exit(1)
	}
}
