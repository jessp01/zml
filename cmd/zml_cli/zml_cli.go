package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jessp01/zml"
	"github.com/urfave/cli"
)

var fontDir string
var titleFont string
var labelFont string
var elementFont string
var debug bool = false

func populateAppMetadata(app *cli.App) {
	cli.AppHelpTemplate = `NAME:
   {{.Name}} - {{.Usage}}

USAGE:
   {{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}{{if .Commands}} command [command options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[input-file]{{end}}
   {{if len .Authors}}
COMMANDS:
{{range .Commands}}{{if not .HideHelp}}   {{join .Names ", "}}{{ "\t"}}{{.Usage}}{{ "\n" }}{{end}}{{end}}{{end}}{{if .VisibleFlags}}
GLOBAL OPTIONS:
   {{range .VisibleFlags}}{{.}}{{ "\n" }}
   {{end}}{{end}}{{if .Copyright }}
AUTHOR:
   {{range .Authors}}{{ . }}{{end}}
   {{end}}{{if .Commands}}
COPYRIGHT:
   {{.Copyright}}
   {{end}}
`
	app.Usage = "Diagram and flowchart tool"
	app.Version = "0.21.3"
	app.EnableBashCompletion = true
	cli.VersionFlag = cli.BoolFlag{
		Name:  "print-version, V",
		Usage: "print only the version",
	}
	app.Compiled = time.Now()
	app.Description = "Converts ZML text to PNG diagrams.\n"
	app.Authors = []cli.Author{
		{
			Name:  "Jesse Portnoy",
			Email: "jesse@packman.io",
		},
	}
	app.Copyright = "(c) packman.io"
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:        "font-dir, fd",
			Usage:       "Path to font dir.\n",
			Destination: &fontDir,
		},
		cli.StringFlag{
			Name:        "title-font, tf",
			Usage:       `font to use for titles; e.g: Roboto-Bold.ttf,30`,
			Destination: &titleFont,
		},
		cli.StringFlag{
			Name:        "element-font, ef",
			Usage:       `font to use for node names; e.g: Roboto-Regular.ttf,15`,
			Destination: &elementFont,
		},
		cli.StringFlag{
			Name:        "label-font, lf",
			Usage:       `font to use for connection labels; e.g: Roboto-Italic.ttf,15`,
			Destination: &labelFont,
		},
		cli.BoolFlag{
			Name:        "debug, d",
			Usage:       "Run in debug mode.\n",
			Destination: &debug,
		},
	}
}

func main() {
	app := cli.NewApp()
	populateAppMetadata(app)

	app.Action = func(c *cli.Context) error {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		if c.NArg() < 1 {
			fmt.Printf("USAGE: %s <filename> \n", os.Args[0])
			os.Exit(0)
		}

		fileName := c.Args().Get(0)

		fileBytes, err := ioutil.ReadFile(fileName)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		dia := zml.NewDiagram(fileName)
		if debug {
			dia.SetDebug(true)
		}
		dia.SetFontDir(fontDir)
		if titleFont != "" {
			titleFontAttribs := strings.Split(titleFont, ",")
			fontSize, _ := strconv.ParseFloat(titleFontAttribs[1], 64)
			if err == nil {
				dia.SetTitleFont(zml.Font{Name: titleFontAttribs[0], Size: fontSize})
			}
		}
		if labelFont != "" {
			labelFontAttribs := strings.Split(labelFont, ",")
			fontSize, _ := strconv.ParseFloat(labelFontAttribs[1], 64)
			if err == nil {
				dia.SetLabelFont(zml.Font{Name: labelFontAttribs[0], Size: fontSize})
			}
		}
		if elementFont != "" {
			elementFontAttribs := strings.Split(elementFont, ",")
			fontSize, _ := strconv.ParseFloat(elementFontAttribs[1], 64)
			if err == nil {
				dia.SetElementLabelFont(zml.Font{Name: elementFontAttribs[0], Size: fontSize})
			}
		}
		dia.ProcessData(fileBytes)
		return nil
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
