package main

import (
	g "github.com/AllenDang/giu"
	"image"
	"image/color"
	"math/rand"
	"time"
)

var (
	circleColour          = color.RGBA{200, 12, 12, 255}
	isCycling             = false
	cycleIntervals        = make(chan int32)
	cyclePeriod     int32 = 100
	seed                  = rand.NewSource(time.Now().UnixNano())
	randomGenerator       = rand.New(seed)
)

type CircleButtonWidget struct {
	g.Widget

	id      string
	clicked func()
}

func CircleButton(id string, clicked func()) *CircleButtonWidget {
	return &CircleButtonWidget{
		id:      id,
		clicked: clicked,
	}
}

func (c *CircleButtonWidget) Build() {
	var width float32 = 150
	var padding float32 = 8.0

	pos := g.GetCursorPos()

	radius := int(width/2 + padding*2)

	buttonWidth := float32(radius) * 2
	g.InvisibleButton().Size(buttonWidth, buttonWidth).OnClick(c.clicked).Build()

	center := pos.Add(image.Pt(radius, radius))

	canvas := g.GetCanvas()

	canvas.AddCircleFilled(center, float32(radius), &circleColour)
}

func setCyclePeriod() {
	cycleIntervals <- cyclePeriod
}

func onCircleClick() {
	isCycling = !isCycling
}

func cycleCircleColour() {
	interval := <-cycleIntervals
	ticker := time.NewTicker(time.Duration(interval) * time.Millisecond)
	for {
		select {
		case i := <-cycleIntervals:
			ticker.Stop()
			ticker = time.NewTicker(time.Duration(i) * time.Millisecond)

		case <-ticker.C:
			if !isCycling {
				continue
			}

			circleColour = color.RGBA{
				uint8(randomGenerator.Intn(255)),
				uint8(randomGenerator.Intn(255)),
				uint8(randomGenerator.Intn(255)),
				uint8(randomGenerator.Intn(120) + 135),
			}
			g.Update()
		}
	}
}

func uiLoop() {
	g.SingleWindow().Layout(
		g.Row(
			g.InputInt(&cyclePeriod).Label("Colour cycle period: "),
		),
		g.Row(
			g.Button("Update").OnClick(setCyclePeriod),
		),
		g.Align(g.AlignCenter).To(CircleButton("Circle Button", onCircleClick)),
	)
}

func main() {
	flags := g.MasterWindowFlagsFloating + g.MasterWindowFlagsNotResizable + g.MasterWindowFlagsFrameless + g.MasterWindowFlagsTransparent
	wnd := g.NewMasterWindow("Visual Distraction", 400, 300, flags)

	go cycleCircleColour()

	wnd.Run(uiLoop)
}
