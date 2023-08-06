package zml

import (
	"fmt"
	"image/color"
	"log"
	"path/filepath"

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
	height = 768
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
func (dia *Diagram) Render() {
	dia.dc = gg.NewContext(width, height)
	dia.dc.DrawRectangle(0, 0, width, height)
	dia.dc.SetColor(color.White)
	dia.dc.Fill()

	dia.renderTitle()
	dia.renderElemenets()
	dia.renderConnections()

	dia.dc.SavePNG(fmt.Sprintf("%s.png", dia.filename))
}

func (dia *Diagram) renderTitle() {
	s := dia.title
	if err := dia.dc.LoadFontFace(filepath.Join(dia.fontDir, dia.titleFont.Name), dia.titleFont.Size); err != nil {
		log.Printf(err.Error())
	}
	textWidth, _ := dia.dc.MeasureString(s)
	centerX := float64(dia.dc.Width())/2.0 - float64(textWidth)/2.0
	log.Printf("title: %s", dia.title)
	dia.dc.SetColor(color.Black)
	dia.dc.DrawString(s, centerX, height*0.05)
	dia.dc.Stroke()
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
		// startX := float64(elemenetsPadding + (len(dia.renderedelemenets) * (elemenetBoxWidth + 1000/(len(dia.elemenets)))))
		endX := startX + elemenetBoxWidth
		startY := 1000 * 0.1 // 10% from the top
		endY := startY + elemenetBoxHeight
		// draw the border
		/* dia.dc.SetColor(color.Black)
		dia.dc.SetLineWidth(rectangleStrokeWidth)
		dia.dc.SetFillRule(gg.FillRuleEvenOdd)

		dia.dc.DrawLine(startX, startY, endX, startY)
		dia.dc.Stroke()

		dia.dc.DrawLine(startX, endY, endX, endY)
		dia.dc.Stroke()

		dia.dc.DrawLine(startX, startY, startX, endY)
		dia.dc.Stroke()

		dia.dc.DrawLine(endX, startY, endX, endY)
		dia.dc.Stroke()*/

		if err := dia.dc.LoadFontFace(filepath.Join(dia.fontDir, dia.elementLabelFont.Name), dia.elementLabelFont.Size); err != nil {
			log.Printf(err.Error())
		}

		dia.dc.SetColor(color.Gray{Y: 230})
		dia.dc.DrawRoundedRectangle(
			startX,
			startY,
			elemenetBoxWidth,
			elemenetBoxHeight,
			5,
		)
		dia.dc.SetRGB255(Colorlookup("platered"))
		dia.dc.Fill()
		dia.dc.SetColor(color.White)
		strWidth, strHeight := dia.dc.MeasureString(p.Name)
		centerStrWidth := startX + ((endX - startX) / 2) - strWidth/2
		centerStrHeight := (endY-startY)/2 + startY + (strHeight / 2)

		dia.dc.DrawString(
			p.Name,
			centerStrWidth,
			centerStrHeight,
		)
		dia.dc.Stroke()
		dia.dc.SetColor(color.Black)
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
		// new
		dia.dc.SetColor(color.Gray{Y: 230})
		startY = lineEndY + 1
		endY = startY + elemenetBoxHeight
		dia.dc.DrawRoundedRectangle(
			startX,
			startY,
			elemenetBoxWidth,
			elemenetBoxHeight,
			5,
		)
		dia.dc.SetRGB255(Colorlookup("platered"))
		dia.dc.Fill()
		dia.dc.SetColor(color.White)

		centerStrHeight = (endY-startY)/2 + startY + (strHeight / 2)
		dia.dc.DrawString(
			p.Name,
			centerStrWidth,
			centerStrHeight,
		)
		dia.dc.Stroke()
		dia.dc.SetColor(color.Black)
		// end new
		dia.renderedElemenets = append(dia.renderedElemenets, p)

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
			if err := dia.dc.LoadFontFace(filepath.Join(dia.fontDir, dia.labelFont.Name), dia.labelFont.Size); err != nil {
				log.Printf(err.Error())
			}
			textWidth, textHeight := dia.dc.MeasureString(e.Label)
			textY := startY + textHeight
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
				// log.Printf("AddElemenets(): skipping %s\n", n)
				skipAdd = true
				break
			}
		}
		if !skipAdd {
			// log.Printf("AddElemenets(): adding %s\n", n)
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

	dia.edges = append(dia.edges, edge{from: *fromPar, to: *toPar, Label: label, directional: false})
	return nil
}

// SetTitle sets the diagram's title
func (dia *Diagram) SetTitle(s string) {
	dia.title = s
}

// SetFontDir path to font dir
func (dia *Diagram) SetFontDir(dir string) {
	dia.fontDir = dir
}

// SetTitleFont sets font to use for the title
func (dia *Diagram) SetTitleFont(font Font) {
	dia.titleFont = font
}

// SetLabelFont sets font to use for labels
func (dia *Diagram) SetLabelFont(font Font) {
	dia.labelFont = font
}

// SetElementLabelFont sets font to use for labels
func (dia *Diagram) SetElementLabelFont(font Font) {
	dia.elementLabelFont = font
}
