package main

import (
	"fmt"
	"github.com/itsmanjeet/app/internal/appimage"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	app := &cli.App{
		Name:  "appimage",
		Usage: "Manage Your AppImages",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "root-dir",
				Value: "/apps",
				Usage: "Specify Integration root directory",
			},
		},
		Action: func(context *cli.Context) error {
			if context.Args().Len() == 0 {
				return fmt.Errorf("no appimage provided")
			}
			file := context.Args().Get(0)
			log.Printf("APPIMAGE: %s\n", file)
			a, err := appimage.Load(file)
			if err != nil {
				return fmt.Errorf("failed to load appimage %v", err)
			}
			return a.Integrate(context.String("root-dir"))
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
