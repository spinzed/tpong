package main

import "github.com/gdamore/tcell"

type Theme struct {
	Name       string
	Background tcell.Color
	// overlay is for the middle line and the score
	Overlay   tcell.Color
	Platforms tcell.Color
	Ball      tcell.Color
}

// Get background style.
func (t *Theme) GetBackgroundStyle() tcell.Style {
	return tcell.StyleDefault.Background(t.Background).Foreground(tcell.ColorWhite)
}

// Get overlay (middle line and number) style.
func (t *Theme) GetOverlayStyle() tcell.Style {
	return tcell.StyleDefault.Background(t.Overlay).Foreground(tcell.ColorWhite)
}

// Get platform style.
func (t *Theme) GetPlatformStyle() tcell.Style {
	return tcell.StyleDefault.Background(t.Platforms).Foreground(tcell.ColorWhite)
}

// Get ball style.
func (t *Theme) GetBallStyle() tcell.Style {
	return tcell.StyleDefault.Background(t.Ball).Foreground(tcell.ColorWhite)
}

// Get text style.
func (t *Theme) GetTextStyle() tcell.Style {
	// text foreground and background are constant
	text := tcell.ColorWhite
	textBg := tcell.NewRGBColor(100, 100, 100)
	return tcell.StyleDefault.Background(textBg).Foreground(text)
}

// Theme handler for themes.
type ThemeHandler struct {
	themes       []*Theme
	currentTheme int
	background   bool
}

func newThemeHandler(initialBg bool) *ThemeHandler {
	return &ThemeHandler{allThemes, 0, initialBg}
}

func (t *ThemeHandler) GetCurrent() *Theme {
	return t.themes[t.currentTheme]
}

func (t *ThemeHandler) Switch() {
	t.currentTheme++

	if t.currentTheme >= len(t.themes) {
		t.currentTheme = 0
	}
}

func (t *ThemeHandler) GetThemes() []*Theme {
	return t.themes
}

func (t *ThemeHandler) IsBgShown() bool {
	return t.background
}

func (t *ThemeHandler) ToggleBg() {
	t.background = !t.background
}

// All themes are stored in this array variable. They are for use with ThemeHandler
var allThemes = []*Theme{
	themeDefaultGrey,
	sapphire,
	purpur,
	colorful,
}

var themeDefaultGrey = &Theme{
	Name:       "Grey",
	Background: tcell.NewRGBColor(38, 38, 38),
	Overlay:    tcell.NewRGBColor(209, 209, 209),
	Platforms:  tcell.NewRGBColor(150, 150, 150),
	Ball:       tcell.NewRGBColor(255, 255, 255),
}

var sapphire = &Theme{
	Name:       "Sapphire",
	Background: tcell.NewRGBColor(0, 23, 31),
	Overlay:    tcell.NewRGBColor(100, 220, 255),
	Platforms:  tcell.NewRGBColor(0, 100, 175),
	Ball:       tcell.NewRGBColor(240, 240, 255),
}

var purpur = &Theme{
	Name:       "Purpur",
	Background: tcell.NewRGBColor(34, 25, 50),
	Overlay:    tcell.NewRGBColor(165, 100, 180),
	Platforms:  tcell.NewRGBColor(165, 50, 130),
	Ball:       tcell.NewRGBColor(230, 190, 240),
}

var colorful = &Theme{
	Name:       "Beach",
	Background: tcell.NewRGBColor(5, 79, 122),
	Overlay:    tcell.NewRGBColor(220, 200, 170),
	Platforms:  tcell.NewRGBColor(245, 180, 40),
	Ball:       tcell.NewRGBColor(245, 120, 100),
}
