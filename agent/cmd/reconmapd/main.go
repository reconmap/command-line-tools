package main

import (
	"context"
	"fmt"
	"os"
	"reconmap/agent/internal"
	"reconmap/agent/internal/configuration"

	shareconfig "github.com/reconmap/shared-lib/pkg/configuration"

	"github.com/urfave/cli/v3"
)

func ConfigAction(ctx context.Context, c *cli.Command) error {

	config := configuration.NewConfig()
	configurationFilePath, err := shareconfig.SaveConfig(config, configuration.ConfigFileName)
	if err != nil {
		return fmt.Errorf("error saving configuration: %w", err)
	}
	fmt.Printf("Configuration successfully saved to: %s\n", configurationFilePath)
	fmt.Println("You can now use the 'rmap login' command to authenticate with the server.")
	return nil
}

func RunAction(ctx context.Context, c *cli.Command) error {
	server := internal.NewApp()
	if err := server.Run(); err != nil {
		server.Logger.Error(*err)
	}
	return nil
}

func main() {
	app := cli.Command{
		Name:    "reconmapd",
		Usage:   "Reconmap's agent",
		Version: "1.0.0",
	}
	app.Copyright = "Apache License v2.0"
	app.Usage = "Reconmap's agent"
	app.Description = "Reconmap's agent for running scheduled commands"
	app.Commands = []*cli.Command{
		{
			Name:   "config",
			Usage:  "Creates a configuration file for Reconmapd",
			Flags:  []cli.Flag{},
			Action: ConfigAction,
		},
		{
			Name:   "run",
			Usage:  "Starts the Reconmapd server",
			Flags:  []cli.Flag{},
			Action: RunAction,
		},
	}

	err := app.Run(context.Background(), os.Args)
	if err != nil {
		panic(err)
	}
}
