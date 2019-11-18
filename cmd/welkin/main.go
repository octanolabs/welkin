package main

import (
	"fmt"
	"log"
	"os"
	"time"

	params "github.com/octanolabs/welkin/params"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:     "Welkin",
		Version:  params.VersionWithMeta,
		Compiled: time.Now(),
		Authors: []*cli.Author{
			&cli.Author{
				Name:  "Luke Williams",
				Email: "devnull@octano.dev",
			},
		},
		Copyright: "(c) 2019-present, Luke Williams",
		Commands: []*cli.Command{
			&cli.Command{
				Name:      "wallet",
				Category:  "WALLET COMMANDS",
				Usage:     "Manage master wallet",
				Aliases:   []string{"w"},
				ArgsUsage: "",
				Description: `
Manage master wallet, import a recovery phrase into a new
wallet, create a new wallet or update an existing wallet.
`,
				Subcommands: []*cli.Command{
					&cli.Command{
						Name:   "create",
						Usage:  "Create a new master wallet",
						Action: createNewWallet,
						// Action: utils.MigrateFlags(accountCreate),
						/* Flags: []cli.Flag{
							utils.DataDirFlag,
							utils.KeyStoreDirFlag,
							utils.PasswordFileFlag,
							utils.LightKDFFlag,
						}, */
						Description: `
    welkin wallet new
Creates a new master wallet and prints the recovery phrase.
`,
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func createNewWallet(c *cli.Context) error {
	newMnemonic := GenerateMnemonic()
	fmt.Println(newMnemonic)

	return nil
}
