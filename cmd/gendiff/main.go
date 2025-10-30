package main

import (
	"code/code"
	"context"
	"fmt"
	urfaveCli "github.com/urfave/cli/v3"
	"os"
)

func main() {
	app := newApp()

	if err := app.Run(context.Background(), os.Args); err != nil {
		os.Exit(1)
	}
}

func newApp() *urfaveCli.Command {
	return &urfaveCli.Command{
		Name:      "gendiff",
		Usage:     "Compares two configuration files and shows a difference.",
		UsageText: "gendiff [global options]",
		Flags: []urfaveCli.Flag{
			&urfaveCli.BoolFlag{
				Name:    "format",
				Aliases: []string{"f"},
				Usage:   "string  output format (default: \"stylish\")",
			},
		},
		Action: func(_ context.Context, cmd *urfaveCli.Command) error {
			if cmd.Args().Len() != 2 {
				return urfaveCli.Exit("please provide only 2 arguments", 1)
			}

			configOne := cmd.Args().First()
			configTwo := cmd.Args().Tail()[0]

			res := code.GenDiff(configOne, configTwo)
			fmt.Printf("%+v\n", res)
			return nil
		},
	}
}
