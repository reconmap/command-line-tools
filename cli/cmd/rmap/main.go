package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/mail"
	"os"

	"github.com/reconmap/shared-lib/pkg/logging"

	"github.com/fatih/color"
	"github.com/reconmap/cli/internal/build"
	"github.com/reconmap/cli/internal/commands"
	"github.com/urfave/cli/v3"
)

func main() {
	logger := logging.GetLoggerInstance()
	defer logger.Sync()

	cli.VersionPrinter = func(c *cli.Command) {
		fmt.Printf("Version=%s\nBuildDate=%s\nGitCommit=%s\n", c.Version, build.BuildTime, build.BuildCommit)
	}

	app := cli.Command{}
	app.Flags = []cli.Flag{
		&cli.BoolFlag{
			Name:     "hide-banner",
			Usage:    "hide Reconmap's banner",
			Aliases:  []string{"b"},
			Required: false,
			Value:    false,
		},
	}
	app.Before = func(ctx context.Context, c *cli.Command) (context.Context, error) {
		if !c.Bool("hide-banner") {
			banner := "ICBfX19fICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICANCiB8ICBfIFwgX19fICBfX18gX19fICBfIF9fICBfIF9fIF9fXyAgIF9fIF8gXyBfXyAgDQogfCB8XykgLyBfIFwvIF9fLyBfIFx8ICdfIFx8ICdfIGAgXyBcIC8gX2AgfCAnXyBcIA0KIHwgIF8gPCAgX18vIChffCAoXykgfCB8IHwgfCB8IHwgfCB8IHwgKF98IHwgfF8pIHwNCiB8X3wgXF9cX19ffFxfX19cX19fL3xffCB8X3xffCB8X3wgfF98XF9fLF98IC5fXy8gDQogICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgfF98ICAgIA0KDQo="
			sDec, _ := base64.StdEncoding.DecodeString(banner)
			color.Set(color.FgHiRed)
			fmt.Print(string(sDec))
			color.Unset()
		}
		return nil, nil
	}
	app.Version = build.BuildVersion
	app.Copyright = "Apache License v2.0"
	app.Usage = "Reconmap's CLI"
	app.Description = "Reconmap's command line interface"
	app.Authors = []any{
		mail.Address{Name: "Reconmap", Address: "info@reconmap.com"},
	}
	app.Commands = commands.CommandList

	err := app.Run(context.Background(), os.Args)
	if err != nil {
		logger.Error(err)
	}
}
