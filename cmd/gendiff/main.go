package main

import (
	"context"
	"fmt"
	"os"

	"code/code"
	urfaveCli "github.com/urfave/cli/v3"
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
		UsageText: "gendiff [--format stylish] <file1> <file2>",
		Flags: []urfaveCli.Flag{
			&urfaveCli.StringFlag{
				Name:    "format",
				Aliases: []string{"f"},
				Usage:   "output format (stylish, plain, json)",
				Value:   "stylish",
			},
		},
		Action: func(ctx context.Context, cmd *urfaveCli.Command) error {
			if cmd.Args().Len() != 2 {
				return urfaveCli.Exit("usage: gendiff [--format stylish] <file1> <file2>", 2)
			}
			f1 := cmd.Args().First()
			f2 := cmd.Args().Tail()[0]
			format := cmd.String("format")

			out, err := code.GenDiff(f1, f2, format)
			if err != nil {
				return urfaveCli.Exit(err.Error(), 1)
			}
			fmt.Println(out)
			return nil
		},
	}
}
