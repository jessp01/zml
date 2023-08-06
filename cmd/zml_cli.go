package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strings"

	"github.com/jessp01/zml"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Printf("USAGE: %s <filename> \n", os.Args[0])
		os.Exit(0)
	}

	fileName := os.Args[1]

	fileBytes, err := ioutil.ReadFile(fileName)

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	dia := zml.NewDiagram(fileName)
	dia.SetDebug(true)
	dia.SetFontDir("/usr/share/texlive/texmf-dist/fonts/truetype/google/roboto")
	dia.SetElementLabelFont(zml.Font{Name: "Roboto-Regular.ttf", Size: 15})
	dia.SetLabelFont(zml.Font{Name: "Roboto-Italic.ttf", Size: 15})
	dia.SetDebug(false)

	sliceData := strings.Split(string(fileBytes), "\n")
	firstLine := sliceData[0]
	titleRegexp := regexp.MustCompile(`^\[?title\]?\s*:\s*([\w\s.,-]+)`)
	matches := titleRegexp.FindStringSubmatch(string(firstLine))
	var title string
	if len(matches) > 1 {
		title = matches[1]
		sliceData = sliceData[1:]
		dia.SetTitle(title)
		dia.SetTitleFont(zml.Font{Name: "Roboto-Bold.ttf", Size: 21})
	}

	relationRegexp := regexp.MustCompile(`^\[?([A-Za-z\s]+)\]?([-]+>{0,2})\[?([A-Za-z\s]+)\]?:?(.*)?`)
	for _, line := range sliceData {
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//") {
			continue
		}
		parts := relationRegexp.FindStringSubmatch(string(line))
		if len(parts) > 4 {
			fromElemenet := strings.Trim(parts[1], " ")
			relationType := parts[2]
			toElemenet := strings.Trim(parts[3], " ")
			dia.AddElemenets(fromElemenet, toElemenet)
			label := ""
			if len(parts) == 5 && parts[4] != "" {
				label = parts[4]
			}
			if relationType[len(relationType)-1] == '>' {
				dia.AddDirectionalConnection(fromElemenet, toElemenet, label)
			} else {
				dia.AddConnection(fromElemenet, toElemenet, label)
			}
		}
	}
	dia.Render()

}
