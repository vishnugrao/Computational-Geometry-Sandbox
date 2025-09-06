// Interactive GUI for generating and visualizing random 2D points
// Uses Fyne for cross-platform GUI functionality

package main

import (
	"fmt"
	"image/color"
	"math/rand"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

// Point represents a 2D point with X and Y coordinates
type Point struct {
	X, Y float64
}

// GUIConfig holds configuration for the GUI application
type GUIConfig struct {
	NumPoints    int
	MinX, MaxX   float64
	MinY, MaxY   float64
	CanvasWidth  float32
	CanvasHeight float32
	PointSize    float32
	PointColor   color.Color
}

// PointVisualizer manages the GUI application
type PointVisualizer struct {
	config       *GUIConfig
	points       []Point
	canvas       *fyne.Container
	pointObjects []fyne.CanvasObject
	app          fyne.App
	window       fyne.Window
}

// NewPointVisualizer creates a new point visualizer instance
func NewPointVisualizer() *PointVisualizer {
	pv := &PointVisualizer{
		config: &GUIConfig{
			NumPoints:    30,
			MinX:         0.0,
			MaxX:         100.0,
			MinY:         0.0,
			MaxY:         100.0,
			CanvasWidth:  600,
			CanvasHeight: 400,
			PointSize:    6,
			PointColor:   color.RGBA{R: 70, G: 130, B: 255, A: 255}, // Steel blue
		},
		points: make([]Point, 0),
	}

	// Initialize random seed
	rand.Seed(time.Now().UnixNano())

	return pv
}

// GeneratePoints creates random points based on current configuration
func (pv *PointVisualizer) GeneratePoints() {
	pv.points = make([]Point, pv.config.NumPoints)

	for i := 0; i < pv.config.NumPoints; i++ {
		x := pv.config.MinX + rand.Float64()*(pv.config.MaxX-pv.config.MinX)
		y := pv.config.MinY + rand.Float64()*(pv.config.MaxY-pv.config.MinY)
		pv.points[i] = Point{X: x, Y: y}
	}
}

// UpdateCanvas redraws all points on the canvas
func (pv *PointVisualizer) UpdateCanvas() {
	// Clear existing objects
	pv.canvas.Objects = nil
	pv.pointObjects = nil

	// Add background
	bg := canvas.NewRectangle(color.RGBA{R: 248, G: 249, B: 250, A: 255})
	bg.Resize(fyne.NewSize(pv.config.CanvasWidth, pv.config.CanvasHeight))
	pv.canvas.Add(bg)

	// Add border
	border := canvas.NewRectangle(color.Transparent)
	border.StrokeColor = color.RGBA{R: 200, G: 200, B: 200, A: 255}
	border.StrokeWidth = 2
	border.Resize(fyne.NewSize(pv.config.CanvasWidth, pv.config.CanvasHeight))
	pv.canvas.Add(border)

	// Calculate scaling factors
	padding := float32(20)
	scaleX := (pv.config.CanvasWidth - 2*padding) / float32(pv.config.MaxX-pv.config.MinX)
	scaleY := (pv.config.CanvasHeight - 2*padding) / float32(pv.config.MaxY-pv.config.MinY)

	// Add points
	for _, point := range pv.points {
		// Transform coordinates
		canvasX := padding + float32(point.X-pv.config.MinX)*scaleX - pv.config.PointSize/2
		canvasY := padding + float32(pv.config.MaxY-point.Y)*scaleY - pv.config.PointSize/2

		// Create point circle
		circle := canvas.NewCircle(pv.config.PointColor)
		circle.Move(fyne.NewPos(canvasX, canvasY))
		circle.Resize(fyne.NewSize(pv.config.PointSize, pv.config.PointSize))

		pv.canvas.Add(circle)
		pv.pointObjects = append(pv.pointObjects, circle)
	}

	pv.canvas.Refresh()
}

// CreateControlPanel creates the control panel with sliders and buttons
func (pv *PointVisualizer) CreateControlPanel() *container.Scroll {
	// Number of points slider
	pointsLabel := widget.NewLabel("Number of Points:")
	pointsSlider := widget.NewSlider(5, 200)
	pointsSlider.SetValue(float64(pv.config.NumPoints))
	pointsValue := widget.NewLabel(fmt.Sprintf("%d", pv.config.NumPoints))

	pointsSlider.OnChanged = func(value float64) {
		pv.config.NumPoints = int(value)
		pointsValue.SetText(fmt.Sprintf("%d", pv.config.NumPoints))
	}

	// X range controls
	xMinLabel := widget.NewLabel("X Min:")
	xMinEntry := widget.NewEntry()
	xMinEntry.SetText(fmt.Sprintf("%.1f", pv.config.MinX))

	xMaxLabel := widget.NewLabel("X Max:")
	xMaxEntry := widget.NewEntry()
	xMaxEntry.SetText(fmt.Sprintf("%.1f", pv.config.MaxX))

	// Y range controls
	yMinLabel := widget.NewLabel("Y Min:")
	yMinEntry := widget.NewEntry()
	yMinEntry.SetText(fmt.Sprintf("%.1f", pv.config.MinY))

	yMaxLabel := widget.NewLabel("Y Max:")
	yMaxEntry := widget.NewEntry()
	yMaxEntry.SetText(fmt.Sprintf("%.1f", pv.config.MaxY))

	// Point size slider
	sizeLabel := widget.NewLabel("Point Size:")
	sizeSlider := widget.NewSlider(2, 15)
	sizeSlider.SetValue(float64(pv.config.PointSize))
	sizeValue := widget.NewLabel(fmt.Sprintf("%.0f", pv.config.PointSize))

	sizeSlider.OnChanged = func(value float64) {
		pv.config.PointSize = float32(value)
		sizeValue.SetText(fmt.Sprintf("%.0f", pv.config.PointSize))
	}

	// Color selection buttons
	colorLabel := widget.NewLabel("Point Color:")

	blueBtn := widget.NewButton("Blue", func() {
		pv.config.PointColor = color.RGBA{R: 70, G: 130, B: 255, A: 255}
		pv.UpdateCanvas()
	})

	redBtn := widget.NewButton("Red", func() {
		pv.config.PointColor = color.RGBA{R: 255, G: 99, B: 71, A: 255}
		pv.UpdateCanvas()
	})

	greenBtn := widget.NewButton("Green", func() {
		pv.config.PointColor = color.RGBA{R: 60, G: 179, B: 113, A: 255}
		pv.UpdateCanvas()
	})

	purpleBtn := widget.NewButton("Purple", func() {
		pv.config.PointColor = color.RGBA{R: 138, G: 43, B: 226, A: 255}
		pv.UpdateCanvas()
	})

	// Generate button
	generateBtn := widget.NewButton("🎲 Generate New Points", func() {
		// Update config from entries
		if val, err := strconv.ParseFloat(xMinEntry.Text, 64); err == nil {
			pv.config.MinX = val
		}
		if val, err := strconv.ParseFloat(xMaxEntry.Text, 64); err == nil {
			pv.config.MaxX = val
		}
		if val, err := strconv.ParseFloat(yMinEntry.Text, 64); err == nil {
			pv.config.MinY = val
		}
		if val, err := strconv.ParseFloat(yMaxEntry.Text, 64); err == nil {
			pv.config.MaxY = val
		}

		pv.GeneratePoints()
		pv.UpdateCanvas()
	})
	generateBtn.Importance = widget.HighImportance

	// Info display
	infoLabel := widget.NewRichTextFromMarkdown("**Random Point Generator**\n\nAdjust parameters and click generate to create new scattered points!")
	infoLabel.Wrapping = fyne.TextWrapWord

	// Layout controls
	rangeForm := container.NewGridWithColumns(2,
		xMinLabel, xMinEntry,
		xMaxLabel, xMaxEntry,
		yMinLabel, yMinEntry,
		yMaxLabel, yMaxEntry,
	)

	pointControls := container.NewVBox(
		pointsLabel,
		container.NewBorder(nil, nil, nil, pointsValue, pointsSlider),
	)

	sizeControls := container.NewVBox(
		sizeLabel,
		container.NewBorder(nil, nil, nil, sizeValue, sizeSlider),
	)

	colorControls := container.NewVBox(
		colorLabel,
		container.NewGridWithColumns(2, blueBtn, redBtn, greenBtn, purpleBtn),
	)

	controlPanel := container.NewVBox(
		infoLabel,
		widget.NewSeparator(),
		pointControls,
		widget.NewSeparator(),
		widget.NewLabel("Coordinate Ranges:"),
		rangeForm,
		widget.NewSeparator(),
		sizeControls,
		widget.NewSeparator(),
		colorControls,
		widget.NewSeparator(),
		generateBtn,
	)

	return container.NewScroll(controlPanel)
}

// Run starts the GUI application
func (pv *PointVisualizer) Run() {
	pv.app = app.NewWithID("point-visualizer")
	pv.app.SetIcon(theme.ComputerIcon())
	pv.window = pv.app.NewWindow("🎯 Random Point Visualizer")
	pv.window.Resize(fyne.NewSize(1000, 600))
	pv.window.CenterOnScreen()

	// Create canvas
	pv.canvas = container.NewWithoutLayout()
	canvasContainer := container.NewScroll(pv.canvas)
	canvasContainer.SetMinSize(fyne.NewSize(pv.config.CanvasWidth, pv.config.CanvasHeight))

	// Create control panel
	controlPanel := pv.CreateControlPanel()
	controlPanel.Resize(fyne.NewSize(300, 600))

	// Initial points generation
	pv.GeneratePoints()
	pv.UpdateCanvas()

	// Create main layout
	content := container.NewHSplit(
		canvasContainer,
		controlPanel,
	)
	content.SetOffset(0.65) // Give more space to canvas

	pv.window.SetContent(content)
	pv.window.ShowAndRun()
}

func main() {
	visualizer := NewPointVisualizer()
	visualizer.Run()
}
