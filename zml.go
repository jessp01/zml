package zml

import (
	"fmt"
	"image/color"
	"log"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/fogleman/gg"
)

const (
	elemenetBoxWidth  = 100.0
	elemenetBoxHeight = 50.0
	elemenetsPadding  = 32

	rectangleStrokeWidth = 2.0
	lineStrokeWidth      = 1.0

	verticalSpaceBetweenEdges = 50

	width  = 1024
	height = 1000
	// TODO: expose on Diagram struct instead
	nodeBgColor    = "platered"
	nodeLabelColor = "white"
)

// Diagram represents a diagram
type Diagram struct {
	elemenets         []elemenet
	edges             []edge
	renderedElemenets []*elemenet
	elemenetsCoordMap map[string]elemenetCoord

	dc               *gg.Context
	title            string
	filename         string
	fontDir          string
	titleFont        Font
	labelFont        Font
	elementLabelFont Font
	debug            bool
}

// NewDiagram init function
func NewDiagram(filename string) *Diagram {
	coordMap := make(map[string]elemenetCoord)

	return &Diagram{
		elemenetsCoordMap: coordMap,
		filename:          filename,
	}
}

// Render generates an image from a `Diagram` object
func (dia *Diagram) Render(width, height float64, color string) {
	dia.dc = gg.NewContext(int(width), int(height))
	dia.dc.DrawRectangle(0, 0, width, height)
	dia.dc.SetRGB255(Colorlookup(color))
	dia.dc.Fill()

	dia.renderTitle()
	dia.renderElemenets()
	dia.renderConnections()

	imageOutputFile := fmt.Sprintf("%s.png", dia.filename)
	dia.dc.SavePNG(imageOutputFile)
	if dia.debug {
		log.Printf("Saved to %s\n", imageOutputFile)
	}
}

func (dia *Diagram) renderTitle() {
	s := dia.title
	if dia.fontDir != "" && dia.titleFont.Name != "" {
		if err := dia.dc.LoadFontFace(filepath.Join(dia.fontDir, dia.titleFont.Name), dia.titleFont.Size); err != nil {
			log.Printf(err.Error())
		}
	}
	textWidth, _ := dia.dc.MeasureString(s)
	centerX := float64(dia.dc.Width())/2.0 - float64(textWidth)/2.0
	dia.dc.SetColor(color.Black)
	dia.dc.DrawString(s, centerX, height*0.05)
	dia.dc.Stroke()
}

func (dia *Diagram) drawBorder(color string, rectangleStrokeWidth float64, startX, startY, endX, endY float64) {
	dia.dc.SetRGB255(Colorlookup(color))
	dia.dc.SetLineWidth(rectangleStrokeWidth)
	dia.dc.SetFillRule(gg.FillRuleEvenOdd)

	dia.dc.DrawLine(startX, startY, endX, startY)
	dia.dc.Stroke()

	dia.dc.DrawLine(startX, endY, endX, endY)
	dia.dc.Stroke()

	dia.dc.DrawLine(startX, startY, startX, endY)
	dia.dc.Stroke()

	dia.dc.DrawLine(endX, startY, endX, endY)
	dia.dc.Stroke()
}

func (dia *Diagram) drawDecisionNode(startX, startY float64, nodeBgColor, label string) {
	strWidth, strHeight := dia.dc.MeasureString(label)
	size := strWidth + 30
	centerStrWidth := startX - strWidth/2
	centerStrHeight := (startY + size/2) + (strHeight / 2)
	dia.dc.LineTo(startX+size/2, startY+size/2)
	dia.dc.LineTo(startX, startY+size)
	dia.dc.LineTo(startX-size/2, startY+size/2)
	dia.dc.LineTo(startX, startY)
	dia.dc.LineTo(startX+size/2, startY+size/2)
	dia.dc.SetRGB255(Colorlookup(nodeBgColor))
	dia.dc.FillPreserve()

	dia.dc.SetRGB255(Colorlookup(nodeLabelColor))
	dia.dc.DrawString(
		label,
		centerStrWidth,
		centerStrHeight,
	)
	dia.dc.Stroke()

	dia.dc.SetColor(color.Black)
}

func (dia *Diagram) drawNode(lineEndY, startX, startY, endX float64, nodeBgColor, nodeLabelColor, label string) {
	endY := startY + elemenetBoxHeight
	strWidth, strHeight := dia.dc.MeasureString(label)
	centerStrWidth := startX + ((endX - startX) / 2) - strWidth/2
	centerStrHeight := (endY-startY)/2 + startY + (strHeight / 2)
	dia.dc.SetColor(color.Gray{Y: 230})
	dia.dc.DrawRoundedRectangle(
		startX,
		startY,
		elemenetBoxWidth,
		elemenetBoxHeight,
		5,
	)
	dia.dc.SetRGB255(Colorlookup(nodeBgColor))
	dia.dc.Fill()
	dia.dc.SetRGB255(Colorlookup(nodeLabelColor))

	dia.dc.DrawString(
		label,
		centerStrWidth,
		centerStrHeight,
	)
	dia.dc.Stroke()
	dia.dc.SetColor(color.Black)
}

func (dia *Diagram) renderElemenets() {
	for idx := range dia.elemenets {
		p := &dia.elemenets[idx]

		for rIdx := range dia.renderedElemenets {
			if dia.renderedElemenets[rIdx].Name == p.Name {
				return
			}
		}
		spacePerBlock := float64(dia.dc.Width() / len(dia.elemenets))
		startX := spacePerBlock*float64(len(dia.renderedElemenets)+1) - spacePerBlock/2 - elemenetsPadding
		endX := startX + elemenetBoxWidth
		startY := height * 0.1 // 10% from the top
		endY := startY + elemenetBoxHeight
		// dia.drawBorder("green", rectangleStrokeWidth, startX, startY, endX, endY)

		if dia.fontDir != "" && dia.elementLabelFont.Name != "" {
			if err := dia.dc.LoadFontFace(filepath.Join(dia.fontDir, dia.elementLabelFont.Name), dia.elementLabelFont.Size); err != nil {
				log.Printf(err.Error())
			}
		}

		dia.drawNode(endY, startX, startY, endX, nodeBgColor, nodeLabelColor, p.Name)
		dia.elemenetsCoordMap[p.Name] = elemenetCoord{
			X: startX,
			Y: startY,
		}

		// render vertical action line for each elemenet
		centerX := startX + (endX-startX)/2 - 2.5
		lineStartY := endY + 2.5
		lineEndY := float64(len(dia.edges)*(verticalSpaceBetweenEdges)) + lineStartY + verticalSpaceBetweenEdges // padding

		dia.dc.SetLineWidth(lineStrokeWidth)
		dia.dc.DrawLine(centerX, lineStartY, centerX, lineEndY)
		dia.dc.Stroke()

		startY = lineEndY + 1
		dia.drawNode(lineEndY, startX, startY, endX, nodeBgColor, nodeLabelColor, p.Name)
		dia.renderedElemenets = append(dia.renderedElemenets, p)

		// dia.drawDecisionNode(startX + 50, 300, "green", "A Decision Node")
	}
}

func (dia *Diagram) renderConnections() {
	renderedEdges := 0

	for idx := range dia.edges {
		e := &dia.edges[idx]
		fromCords := dia.elemenetsCoordMap[e.from.Name]
		toCords := dia.elemenetsCoordMap[e.to.Name]
		startX := fromCords.X + elemenetBoxWidth/2 - 2.5 // 2.5 = half of stroke width
		startY := fromCords.Y + elemenetBoxHeight + 2.5 + float64((1+renderedEdges)*verticalSpaceBetweenEdges)
		endX := toCords.X + elemenetBoxWidth/2 - 2.5
		isReverseEdge := endX < startX

		dia.dc.SetDash(6)
		dia.dc.DrawLine(
			startX,
			startY,
			endX,
			startY)
		dia.dc.Stroke()

		dia.dc.SetDash()

		if e.directional {
			arrowTipStartX := endX
			var arrowTipEndX float64

			if isReverseEdge {
				arrowTipEndX = arrowTipStartX + 10
			} else {
				arrowTipEndX = arrowTipStartX - 10
			}
			dia.dc.DrawLine(arrowTipStartX, startY, arrowTipEndX, startY-10)
			dia.dc.DrawLine(arrowTipStartX, startY, arrowTipEndX, startY+10)
			dia.dc.Stroke()
		}

		if e.Label != "" {
			if dia.fontDir != "" && dia.labelFont.Name != "" {
				if err := dia.dc.LoadFontFace(filepath.Join(dia.fontDir, dia.labelFont.Name), dia.labelFont.Size); err != nil {
					log.Printf(err.Error())
				}
			}
			textWidth, textHeight := dia.dc.MeasureString(e.Label)
			textY := startY + textHeight + 5
			textX := startX
			if isReverseEdge {
				textX -= elemenetsPadding / 2
				textX -= textWidth
			} else {
				textX += elemenetsPadding / 2
			}

			dia.dc.DrawString(e.Label, textX, textY)
		}

		renderedEdges++
	}
}

// AddElemenets sets the `elemenet` array on the Diagram object
func (dia *Diagram) AddElemenets(name ...string) {
	skipAdd := false
	for _, n := range name {
		skipAdd = false
		for i := range dia.elemenets {
			if dia.elemenets[i].Name == n {
				if dia.debug {
					log.Printf("AddElemenets(): skipping %s\n", n)
				}
				skipAdd = true
				break
			}
		}
		if !skipAdd {
			if dia.debug {
				log.Printf("AddElemenets(): adding %s\n", n)
			}
			dia.elemenets = append(dia.elemenets, elemenet{Name: n})
		}
	}
}

// AddDirectionalConnection adds a connection (renders as an arrowed line) between two elemenets
func (dia *Diagram) AddDirectionalConnection(from, to string, label string) error {
	var fromPar *elemenet
	var toPar *elemenet
	for i := range dia.elemenets {
		if dia.elemenets[i].Name == from {
			fromPar = &dia.elemenets[i]
		}
		if dia.elemenets[i].Name == to {
			toPar = &dia.elemenets[i]
		}
	}
	if fromPar == nil {
		panic(fmt.Sprintf("elemenet \"%s\" not found", from))
	}
	if toPar == nil {
		panic(fmt.Sprintf("elemenet \"%s not found", to))
	}

	if dia.debug {
		log.Printf("{from: %s, to: %s, Label: %s, directional: true}\n", fromPar.Name, toPar.Name, label)
	}
	dia.edges = append(dia.edges, edge{from: *fromPar, to: *toPar, Label: label, directional: true})
	return nil
}

// AddConnection adds a connection (renders as a line) between two elemenets
func (dia *Diagram) AddConnection(from, to string, label string) error {
	var fromPar *elemenet
	var toPar *elemenet
	for i := range dia.elemenets {
		if dia.elemenets[i].Name == from {
			fromPar = &dia.elemenets[i]
		}
		if dia.elemenets[i].Name == to {
			toPar = &dia.elemenets[i]
		}
	}
	if fromPar == nil || toPar == nil {
		return fmt.Errorf("elemenet not found")
	}

	if dia.debug {
		log.Printf("{from: %s, to: %s, Label: %s, directional: false}\n", fromPar.Name, toPar.Name, label)
	}
	dia.edges = append(dia.edges, edge{from: *fromPar, to: *toPar, Label: label, directional: false})
	return nil
}

// SetTitle sets the diagram's title
func (dia *Diagram) SetTitle(title string) {
	dia.title = title
	if dia.debug {
		log.Printf("title: %s", dia.title)
	}
}

// SetFontDir path to font dir
func (dia *Diagram) SetFontDir(dir string) {
	dia.fontDir = dir
	if dia.debug {
		log.Printf("fontDir: %s", dia.fontDir)
	}
}

// SetTitleFont sets font to use for the title
func (dia *Diagram) SetTitleFont(font Font) {
	dia.titleFont = font
	if dia.debug {
		log.Printf("Title font: %s, %f", dia.titleFont.Name, dia.titleFont.Size)
	}
}

// SetLabelFont sets font to use for labels
func (dia *Diagram) SetLabelFont(font Font) {
	dia.labelFont = font
	if dia.debug {
		log.Printf("Label font: %s, %f", dia.labelFont.Name, dia.labelFont.Size)
	}
}

// SetElementLabelFont sets font to use for labels
func (dia *Diagram) SetElementLabelFont(font Font) {
	dia.elementLabelFont = font
	if dia.debug {
		log.Printf("Element font: %s, %f", dia.elementLabelFont.Name, dia.elementLabelFont.Size)
	}
}

// SetDebug set debug value
func (dia *Diagram) SetDebug(debug bool) {
	dia.debug = debug
}

// ProcessData generates image from ZML data
func (dia *Diagram) ProcessData(data []byte) {
	sliceData := strings.Split(string(data), "\n")
	firstLine := sliceData[0]
	titleRegexp := regexp.MustCompile(`^\[?title\]?\s*:\s*([\w\s.,-]+)`)
	matches := titleRegexp.FindStringSubmatch(string(firstLine))
	var title string
	if len(matches) > 1 {
		title = matches[1]
		sliceData = sliceData[1:]
		dia.SetTitle(title)
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
}
