package main

import (
	"errors"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

type UploadResponse struct {
	Id            string `json:"id"`
	Name          string `json:"name"`
	Cid           string `json:"cid"`
	Size          int    `json:"size"`
	NumberOfFiles int    `json:"number_of_files"`
	MimeType      string `json:"mime_type"`
	UserId        string `json:"user_id"`
	IndexedAt     string `json:"indexed_at"`
	GroupId       string `json:"group_id,omitempty"`
}

type Options struct {
	GroupId string `json:"group_id"`
}

type Metadata struct {
	Name string `json:"name"`
}

type File struct {
	Id            string  `json:"id"`
	Name          string  `json:"name"`
	Cid           string  `json:"cid"`
	Size          int     `json:"size"`
	NumberOfFiles int     `json:"number_of_files"`
	MimeType      string  `json:"mime_type"`
	GroupId       *string `json:"group_id,omitempty"`
	CreatedAt     string  `json:"created_at"`
}

type ListFilesData struct {
	Files         []File `json:"files"`
	NextPageToken string `json:"next_page_token"`
}

type ListResponse struct {
	Data ListFilesData `json:"data"`
}

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
					&cli.IntFlag{
						Name:    "version",
						Aliases: []string{"v"},
						Value:   1,
						Usage:   "Set desired CID version to either 0 or 1. Default is 1.",
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
					groupId := ctx.String("groupId")
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
				Name:      "delete",
				Aliases:   []string{"d"},
				Usage:     "Delete a file by CID",
				ArgsUsage: "[CID of file]",
				Action: func(ctx *cli.Context) error {
					cid := ctx.Args().First()
					if cid == "" {
						return errors.New("no CID provided")
					}
					err := Delete(cid)
					return err
				},
			},
			{
				Name:    "list",
				Aliases: []string{"l"},
				Usage:   "List most recent files",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "cid",
						Aliases: []string{"c"},
						Value:   "null",
						Usage:   "Search files by CID",
					},
					&cli.StringFlag{
						Name:    "amount",
						Aliases: []string{"a"},
						Value:   "10",
						Usage:   "The number of files you would like to return, default 10 max 1000",
					},
					&cli.StringFlag{
						Name:    "name",
						Aliases: []string{"n"},
						Value:   "null",
						Usage:   "The name of the file",
					},
					&cli.StringFlag{
						Name:    "status",
						Aliases: []string{"s"},
						Value:   "pinned",
						Usage:   "Status of the file. Options are 'pinned', 'unpinned', or 'all'. Default: 'pinned'",
					},
					&cli.StringFlag{
						Name:    "pageOffset",
						Aliases: []string{"p"},
						Value:   "null",
						Usage:   "Allows you to paginate through files. If your file amount is 10, then you could set the pageOffset to '10' to see the next 10 files.",
					},
				},
				Action: func(ctx *cli.Context) error {
					cid := ctx.String("cid")
					amount := ctx.String("amount")
					name := ctx.String("name")
					status := ctx.String("status")
					offset := ctx.String("pageOffset")
					_, err := ListFiles(amount, cid, name, status, offset)
					return err
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
