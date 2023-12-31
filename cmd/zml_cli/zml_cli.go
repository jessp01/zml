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
var backgroundColor string
var width, height float64
var debug bool = false

func populateAppMetadata(app *cli.App) {
	cli.AppHelpTemplate = `NAME:
   {{.Name}} - {{.Usage}}

USAGE:
   {{.HelpName}} {{if .VisibleFlags}}[global options]{{end}} {{if .ArgsUsage}}{{.ArgsUsage}}{{else}}[input-file]{{end}}
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
			Name:        "font-dir, f",
			Usage:       "Path to font dir.\n",
			Destination: &fontDir,
		},
		cli.StringFlag{
			Name:        "title-font, t",
			Usage:       `font to use for titles; e.g: Roboto-Bold.ttf,30`,
			Destination: &titleFont,
		},
		cli.StringFlag{
			Name:        "element-font, e",
			Usage:       `font to use for node names; e.g: Roboto-Regular.ttf,15`,
			Destination: &elementFont,
		},
		cli.StringFlag{
			Name:        "label-font, l",
			Usage:       `font to use for connection labels; e.g: Roboto-Italic.ttf,15`,
			Destination: &labelFont,
		},
		cli.Float64Flag{
			Name:        "width, w",
			Usage:       "Image width.",
			Destination: &width,
			Value:       1024,
		},
		cli.Float64Flag{
			Name:        "height",
			Usage:       "Image height.",
			Destination: &height,
			Value:       1024,
		},
		cli.StringFlag{
			Name:        "background-color, b",
			Usage:       `Background colour`,
			Destination: &backgroundColor,
			Value:       "white",
		},
		cli.BoolFlag{
			Name:        "debug, d",
			Usage:       "Run in debug mode.",
			Destination: &debug,
		},
	}
}

func main() {
	app := cli.NewApp()
	populateAppMetadata(app)

	app.Action = func(c *cli.Context) error {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
		fmt.Printf("%f, %f, %s, %s\n", width, height, backgroundColor, fontDir)
		if c.NArg() < 1 {
			cli.ShowAppHelp(c)
			os.Exit(1)
		}

		fileName := c.Args().Get(0)

		fileBytes, err := ioutil.ReadFile(fileName)

		if err != nil {
			log.Fatal(err)
		}

		dia := zml.NewDiagram(fileName)
		if debug {
			dia.SetDebug(true)
		}
		dia.SetFontDir(fontDir)
		fontSize := 30.00
		if titleFont != "" {
			titleFontAttribs := strings.Split(titleFont, ",")
			if len(titleFontAttribs) > 1 {
				fontSize, _ = strconv.ParseFloat(titleFontAttribs[1], 64)
			}
			if err == nil {
				dia.SetTitleFont(zml.Font{Name: titleFontAttribs[0], Size: fontSize})
			}
		}
		fontSize = 15.00
		if labelFont != "" {
			labelFontAttribs := strings.Split(labelFont, ",")
			if len(labelFontAttribs) > 1 {
				fontSize, _ = strconv.ParseFloat(labelFontAttribs[1], 64)
			}
			if err == nil {
				dia.SetLabelFont(zml.Font{Name: labelFontAttribs[0], Size: fontSize})
			}
		}
		if elementFont != "" {
			elementFontAttribs := strings.Split(elementFont, ",")
			if len(elementFontAttribs) > 1 {
				fontSize, _ = strconv.ParseFloat(elementFontAttribs[1], 64)
			}
			if err == nil {
				dia.SetElementLabelFont(zml.Font{Name: elementFontAttribs[0], Size: fontSize})
			}
		}
		dia.ProcessData(fileBytes)
		dia.Render(width, height, backgroundColor)
		return nil
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
