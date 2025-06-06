package commands

import (
	"errors"
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/reconmap/shared-lib/pkg/api"
	"github.com/reconmap/shared-lib/pkg/configuration"
	"github.com/rodaine/table"
	"github.com/urfave/cli/v2"
)

func preActionChecks(c *cli.Context) error {
	if !configuration.HasConfig() {
		return errors.New("Rmap has not been configured. Please call the 'rmap config' command first.")
	}
	return nil
}

var CommandList []*cli.Command = []*cli.Command{
	{
		Name:   "login",
		Usage:  "Initiates session with the server",
		Flags:  []cli.Flag{},
		Before: preActionChecks,
		Action: func(c *cli.Context) error {
			err := Login()
			return err
		},
	},
	{
		Name:   "logout",
		Usage:  "Terminates session with the server",
		Flags:  []cli.Flag{},
		Before: preActionChecks,
		Action: func(c *cli.Context) error {
			err := Logout()
			return err
		},
	},
	{
		Name:  "config",
		Usage: "Creates a configuration file for Rmap",
		Flags: []cli.Flag{},
		Action: func(c *cli.Context) error {

			config := configuration.NewConfig()
			configurationFilePath, err := configuration.SaveConfig(config)
			if err != nil {
				return fmt.Errorf("error saving configuration: %w", err)
			}
			fmt.Printf("Configuration successfully saved to: %s\n", configurationFilePath)
			fmt.Println("You can now use the 'rmap login' command to authenticate with the server.")
			return nil
		},
	},
	{
		Name:    "command",
		Aliases: []string{"c"},
		Usage:   "Search and run commands",
		Before:  preActionChecks,
		Subcommands: []*cli.Command{
			{
				Name:  "search",
				Usage: "Search commands by keywords",
				Action: func(c *cli.Context) error {
					if c.Args().Len() == 0 {
						return errors.New("no keywords were entered after the search command")
					}
					var keywords string = strings.Join(c.Args().Slice(), " ")
					commands, err := api.GetCommandsByKeywords(keywords)
					if err != nil {
						return err
					}

					var numCommands int = len(*commands)
					fmt.Printf("%d commands matching '%s'\n", numCommands, keywords)

					if numCommands > 0 {
						fmt.Println()

						headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
						columnFmt := color.New(color.FgYellow).SprintfFunc()

						tbl := table.New("ID", "Name", "Description", "Output parser", "Executable type", "Executable path", "Arguments")
						tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt)

						for _, command := range *commands {
							tbl.AddRow(command.ID, command.Name, command.Description)

						}
						tbl.Print()
					}

					return err
				},
			},
			{
				Name:  "run",
				Usage: "Run a command and upload its output to the server",
				Flags: []cli.Flag{
					&cli.IntFlag{Name: "projectId", Aliases: []string{"pid"}, Required: false},
					&cli.IntFlag{Name: "commandUsageId", Aliases: []string{"cuid"}, Required: true},
					&cli.StringSliceFlag{Name: "var", Required: false},
				},
				Action: func(c *cli.Context) error {
					projectId := c.Int("projectId")
					commandUsageId := c.Int("cuid")

					usage, err := api.GetCommandUsageById(commandUsageId)
					if err != nil {
						return fmt.Errorf("unable to retrieve command usage with id=%d (%w)", commandUsageId, err)
					}
					err = RunCommand(projectId, usage, c.StringSlice("var"))
					if err != nil {
						return err
					}

					err = UploadResults(projectId, usage)
					return err
				},
			},
		},
	},
}
