package main

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/nikooo777/ColorBlind/generator"
	"github.com/wailsapp/wails/v2/pkg/runtime"
	"golang.org/x/image/font/gofont/gobold"
)

const (
	appTitle = "HueCipher - Reverse Colorblind Plate Generator"

	defaultWindowWidth  = 1280
	defaultWindowHeight = 860
	minWindowWidth      = 960
	minWindowHeight     = 680

	plateWidth  = 700
	plateHeight = 650

	defaultDensity = 2000
	minDensity     = 10
	maxDensity     = 7000

	defaultTextDensity = 5500

	defaultShade = 170
	minShade     = 156
	maxShade     = 190

	defaultConfuserColor = "#569B9B"
	defaultSaveFilename  = "huecipher.png"

	shapeModeCircle = "Circle"
	shapeModeText   = "Text"
)

type App struct {
	ctx   context.Context
	mutex sync.Mutex
	plate *generator.Plate
}

type Defaults struct {
	ConfuserColor string `json:"confuserColor"`
	ShapeMode     string `json:"shapeMode"`
	Density       int    `json:"density"`
	MinDensity    int    `json:"minDensity"`
	MaxDensity    int    `json:"maxDensity"`
	TextDensity   int    `json:"textDensity"`
	Shade         int    `json:"shade"`
	MinShade      int    `json:"minShade"`
	MaxShade      int    `json:"maxShade"`
}

type PlateRequest struct {
	ConfuserColor string `json:"confuserColor"`
	ShapeMode     string `json:"shapeMode"`
	ShapeText     string `json:"shapeText"`
	Density       int    `json:"density"`
	Shade         int    `json:"shade"`
	ShowOutline   bool   `json:"showOutline"`
}

type PlateResponse struct {
	Image     string `json:"image"`
	Message   string `json:"message"`
	Placed    int    `json:"placed"`
	Requested int    `json:"requested"`
}

func NewApp() *App {
	return &App{}
}

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

func (a *App) GetDefaults() Defaults {
	return Defaults{
		Density:       defaultDensity,
		MinDensity:    minDensity,
		MaxDensity:    maxDensity,
		TextDensity:   defaultTextDensity,
		Shade:         defaultShade,
		MinShade:      minShade,
		MaxShade:      maxShade,
		ConfuserColor: defaultConfuserColor,
		ShapeMode:     shapeModeCircle,
	}
}

func (a *App) GeneratePlate(req PlateRequest) (PlateResponse, error) {
	req = normalizeRequest(req)
	plate, err := newPlate(req)
	if err != nil {
		return PlateResponse{}, err
	}

	a.mutex.Lock()
	defer a.mutex.Unlock()

	a.plate = plate
	return a.renderLocked(req)
}

func (a *App) RenderPlate(req PlateRequest) (PlateResponse, error) {
	req = normalizeRequest(req)

	a.mutex.Lock()
	defer a.mutex.Unlock()

	err := a.ensurePlateLocked(req)
	if err != nil {
		return PlateResponse{}, err
	}

	return a.renderLocked(req)
}

func (a *App) SaveCurrentPNG(req PlateRequest) (string, error) {
	req = normalizeRequest(req)

	pngBytes, err := a.renderCurrentPNG(req)
	if err != nil {
		return "", err
	}

	if a.ctx == nil {
		return "", fmt.Errorf("application is not ready")
	}

	path, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:           "Save PNG",
		DefaultFilename: defaultSaveFilename,
		Filters: []runtime.FileFilter{
			{DisplayName: "PNG Images (*.png)", Pattern: "*.png"},
		},
	})
	if err != nil {
		return "", err
	}
	if path == "" {
		return "", nil
	}
	if filepath.Ext(path) == "" {
		path += ".png"
	}

	err = os.WriteFile(path, pngBytes, 0o644)
	if err != nil {
		return "", err
	}

	return path, nil
}

func newPlate(req PlateRequest) (*generator.Plate, error) {
	factory, err := shapeFactory(req)
	if err != nil {
		return nil, err
	}

	return generator.NewPlate(plateWidth, plateHeight, req.Density, factory, circleProfile(req)), nil
}

func (a *App) ensurePlateLocked(req PlateRequest) error {
	if a.plate != nil {
		return nil
	}

	plate, err := newPlate(req)
	if err != nil {
		return err
	}

	a.plate = plate
	return nil
}

func (a *App) renderCurrentPNG(req PlateRequest) ([]byte, error) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	err := a.ensurePlateLocked(req)
	if err != nil {
		return nil, err
	}

	return a.renderPNGBytesLocked(req)
}

func normalizeRequest(req PlateRequest) PlateRequest {
	if !strings.EqualFold(req.ShapeMode, shapeModeText) {
		req.ShapeMode = shapeModeCircle
	} else {
		req.ShapeMode = shapeModeText
	}
	if req.Density < minDensity {
		req.Density = minDensity
	}
	if req.Density > maxDensity {
		req.Density = maxDensity
	}
	if req.Shade < minShade {
		req.Shade = minShade
	}
	if req.Shade > maxShade {
		req.Shade = maxShade
	}
	if req.ConfuserColor == "" {
		req.ConfuserColor = defaultConfuserColor
	}
	req.ShapeText = strings.TrimSpace(req.ShapeText)
	return req
}

func circleProfile(req PlateRequest) generator.CircleProfile {
	if req.ShapeMode == shapeModeText {
		return generator.CircleProfileFine
	}
	return generator.CircleProfileStandard
}

func shapeFactory(req PlateRequest) (func(w, h int) generator.Shape, error) {
	if req.ShapeMode == shapeModeText {
		if req.ShapeText == "" {
			return nil, fmt.Errorf("text required")
		}
		text := req.ShapeText
		return func(w, h int) generator.Shape {
			return generator.NewTextShape(text, gobold.TTF, w, h)
		}, nil
	}

	return func(w, h int) generator.Shape {
		return generator.NewCircleShape(float64(w), float64(h))
	}, nil
}

func (a *App) renderLocked(req PlateRequest) (PlateResponse, error) {
	pngBytes, err := a.renderPNGBytesLocked(req)
	if err != nil {
		return PlateResponse{}, err
	}

	placed := a.plate.CircleCount()
	encoded := base64.StdEncoding.EncodeToString(pngBytes)
	return PlateResponse{
		Image:     "data:image/png;base64," + encoded,
		Placed:    placed,
		Requested: req.Density,
		Message:   fmt.Sprintf("Placed %d / %d", placed, req.Density),
	}, nil
}

func (a *App) renderPNGBytesLocked(req PlateRequest) ([]byte, error) {
	confuserColor, err := parseHexColor(req.ConfuserColor)
	if err != nil {
		return nil, err
	}
	img := a.plate.Render(revealColor(req.Shade), &confuserColor, req.ShowOutline)

	var buf bytes.Buffer
	err = png.Encode(&buf, img)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func revealColor(shade int) color.NRGBA {
	return color.NRGBA{R: 255, G: 0, B: uint8(shade), A: 255}
}

func parseHexColor(value string) (color.NRGBA, error) {
	hex := strings.TrimPrefix(strings.TrimSpace(value), "#")
	if len(hex) != 6 {
		return color.NRGBA{}, fmt.Errorf("invalid color %q", value)
	}

	var r, g, b uint8
	_, err := fmt.Sscanf(hex, "%02x%02x%02x", &r, &g, &b)
	if err != nil {
		return color.NRGBA{}, fmt.Errorf("invalid color %q", value)
	}

	return color.NRGBA{R: r, G: g, B: b, A: 255}, nil
}
