package main

import (
	"errors"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:  "pinata",
		Usage: "A CLI for uploading files to Pinata! To get started make an API key at https://app.pinata.cloud/keys, then authorize the CLI with the auth command with your JWT",
		Commands: []*cli.Command{
			{
				Name:      "auth",
				Aliases:   []string{"a"},
				Usage:     "Authorize the CLI with your Pinata JWT",
				ArgsUsage: "[your Pinata JWT]",
				Action: func(ctx *cli.Context) error {
					jwt := ctx.Args().First()
					if jwt == "" {
						return errors.New("no jwt supplied")
					}
					err := SaveJWT(jwt)
					return err
				},
			},
			{
				Name:      "upload",
				Aliases:   []string{"u"},
				Usage:     "Upload a file or folder to Pinata",
				ArgsUsage: "[path to file]",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "group",
						Aliases: []string{"g"},
						Value:   "",
						Usage:   "Upload a file to a specific group by passing in the groupId",
					},
					&cli.StringFlag{
						Name:    "name",
						Aliases: []string{"n"},
						Value:   "nil",
						Usage:   "Add a name for the file you are uploading. By default it will use the filename on your system.",
					},
					&cli.BoolFlag{
						Name:  "cid-only",
						Usage: "Use if you only want the CID returned after an upload",
					},
				},
				Action: func(ctx *cli.Context) error {
					filePath := ctx.Args().First()
					groupId := ctx.String("group")
					name := ctx.String("name")
					cidOnly := ctx.Bool("cid-only")
					if filePath == "" {
						return errors.New("no file path provided")
					}
					_, err := Upload(filePath, groupId, name, cidOnly)
					return err
				},
			},
			{
				Name:    "groups",
				Aliases: []string{"g"},
				Usage:   "Interact with file groups",
				Subcommands: []*cli.Command{
					{
						Name:    "create",
						Aliases: []string{"c"},
						Usage:   "Create a new group",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:     "name",
								Aliases:  []string{"n"},
								Required: true,
								Usage:    "The name you want to give for a group",
							},
							&cli.BoolFlag{
								Name:    "public",
								Aliases: []string{"p"},
								Value:   false,
								Usage:   "Determine if the group should be public",
							},
						},
						Action: func(ctx *cli.Context) error {
							name := ctx.String("name")
							public := ctx.Bool("public")
							_, err := CreateGroup(name, public)
							return err
						},
					},
					{
						Name:    "list",
						Aliases: []string{"l"},
						Usage:   "List groups on your account",
						Flags: []cli.Flag{
							&cli.BoolFlag{
								Name:    "public",
								Aliases: []string{"p"},
								Usage:   "List only public groups",
							},
							&cli.StringFlag{
								Name:    "amount",
								Aliases: []string{"a"},
								Value:   "10",
								Usage:   "The number of groups you would like to return",
							},
						},
						Action: func(ctx *cli.Context) error {
							public := ctx.Bool("public")
							amount := ctx.String("amount")
							_, err := ListGroups(amount, public)
							return err
						},
					},
					{
						Name:      "delete",
						Aliases:   []string{"d"},
						Usage:     "Delete a group by ID",
						ArgsUsage: "[ID of group]",
						Action: func(ctx *cli.Context) error {
							groupId := ctx.Args().First()
							if groupId == "" {
								return errors.New("no ID provided")
							}
							err := DeleteGroup(groupId)
							return err
						},
					},
				},
			},
			{
				Name:    "files",
				Aliases: []string{"f"},
				Usage:   "Interact with your files on Pinata",
				Subcommands: []*cli.Command{
					{
						Name:      "delete",
						Aliases:   []string{"d"},
						Usage:     "Delete a file by ID",
						ArgsUsage: "[ID of file]",
						Action: func(ctx *cli.Context) error {
							fileId := ctx.Args().First()
							if fileId == "" {
								return errors.New("no CID provided")
							}
							err := DeleteFile(fileId)
							return err
						},
					},
					{
						Name:    "list",
						Aliases: []string{"l"},
						Usage:   "List most recent files",
						Flags: []cli.Flag{
							&cli.StringFlag{
								Name:    "name",
								Aliases: []string{"n"},
								Usage:   "Filter by name of the target file",
							},
							&cli.StringFlag{
								Name:    "cid",
								Aliases: []string{"c"},
								Usage:   "Filter results by CID",
							},
							&cli.StringFlag{
								Name:    "group",
								Aliases: []string{"g"},
								Usage:   "Filter results by group ID",
							},
							&cli.StringFlag{
								Name:    "mime",
								Aliases: []string{"m"},
								Usage:   "Filter results by file mime type",
							},
							&cli.StringFlag{
								Name:    "amount",
								Aliases: []string{"a"},
								Usage:   "The number of files you would like to return, default 10 max 1000",
							},
							&cli.StringFlag{
								Name:    "pageToken",
								Aliases: []string{"p"},
								Usage:   "Allows you to paginate through files.",
							},
							&cli.BoolFlag{
								Name:  "cidPending",
								Value: false,
								Usage: "Filter results based on whether or not the CID is pending",
							},
						},
						Action: func(ctx *cli.Context) error {
							amount := ctx.String("amount")
							pageToken := ctx.String("pageToken")
							name := ctx.String("name")
							cid := ctx.String("cid")
							group := ctx.String("group")
							mime := ctx.String("mime")
							cidPending := ctx.Bool("cidPending")
							_, err := ListFiles(amount, pageToken, cidPending, name, cid, group, mime)
							return err
						},
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
