package main

import (
	"strings"

	"github.com/gdamore/tcell"
)

type UIHandler struct {
	screen     tcell.Screen
	ball       *Ball
	theme      *ThemeHandler
	menuSelect int
}

func NewUIHandler(bgHidden bool) (*UIHandler, error) {
	ui := &UIHandler{}
	var err error

	if ui.screen, err = initScreen(); err != nil {
		return nil, err
	}

	// game asset init
	w, _ := ui.screen.Size()

	// half of the width needs to be divisible by 2 because of ball
	if w/2%2 == 1 {
		w += 2
	}
	ui.ball = newBall((w-ballDiam)/2, 0, ballDiam, 1, 1)

	ui.theme = newThemeHandler(!bgHidden)

	return ui, nil
}

// Create a screen ready to use
func initScreen() (tcell.Screen, error) {
	s, err := tcell.NewScreen()

	if err != nil {
		return s, err
	}
	if err = s.Init(); err != nil {
		return s, err
	}

	// polling events from terminal has been substituted for listening for events from
	// /dev/input/event*, but this has to be left dangling like this for g.screen.Fini
	// to function properly
	go func() {
		for {
			s.PollEvent()
		}
	}()

	return s, nil
}

// Move selected legend item up.
func (ui *UIHandler) MoveMenuSelectedUp(state GameState) {
	ui.menuSelect--
	legend := screenData[state].LegendKeys

	if ui.menuSelect < 0 {
		ui.menuSelect = len(legend) - 1
	}
}

// Move selected legend item down.
func (ui *UIHandler) MoveMenuSelectedDown(state GameState) {
	ui.menuSelect++
	legend := screenData[state].LegendKeys

	if ui.menuSelect > len(legend)-1 {
		ui.menuSelect = 0
	}
}

// Dispatches an action according to the selected legend item.
func (ui *UIHandler) GetSelectedMenuAction(state GameState) Event {
	keysarr := screenData[state].LegendKeys

	// extract the event that has be dispatched and return it
	return keysarr[ui.menuSelect].Event
}

func (ui *UIHandler) Reset() {
	ui.screen.Clear()
}

func (ui *UIHandler) Paint() {
	ui.screen.Show()
}

func (ui *UIHandler) ScreenSize() (int, int) {
	return ui.screen.Size()
}

func (ui *UIHandler) Destroy() {
	if ui.screen != nil {
		ui.screen.Fini()
	}
}

// Fetches the keys and event descriptions and formats them
// into a format suitable for showing on the screen as a legend.
// Mode dictates should keys be format for left or middle.
func (ui *UIHandler) formatKeys(mode string, keys KeyEventMap, aiActive bool) []string {
	selected := ui.menuSelect

	if mode != "middle" && mode != "left" {
		panic("Invalid format mode:" + mode)
	}

	var maxKeylen int

	// switch default key names with these
	alternate := map[string]string{
		"Up":   "ArrowUp",
		"Down": "ArrowDown",
	}

	// cycle all event names and get their keys and descs
	for _, keyev := range keys {
		key := keyev.Key.Name
		realKey := key

		if alternate[key] != "" {
			realKey = alternate[key]
		}

		if len(key) > maxKeylen {
			maxKeylen = len(realKey)
		}
	}

	var final []string

	// format every line and append it to the final line list
	for i, keyevt := range keys {
		key := keyevt.Key.Name
		var thstring string

		if mode == "left" {
			thstring += strings.Repeat(" ", (maxKeylen-len(key))/2-1)
			if i == selected {
				thstring += ">> "
			} else {
				thstring += "   "
			}
			thstring += key
			thstring += strings.Repeat(" ", maxKeylen-len(key)-((maxKeylen-len(key))/2)+1)
			thstring += "- " + keyevt.Event.Description + " "
		}
		if mode == "middle" {
			thstring += "  "
			if i == selected {
				thstring += ">> "
			}
			thstring += key + " - " + keyevt.Event.Description
		}

		// add key event specific things
		switch keyevt.Event {
		case eventToggleAI:
			thstring += ": [ "
			if aiActive {
				thstring += "IN"
			}
			thstring += "ACTIVE ]"
		case eventSwitchTheme:
			thstring += ": [ " + ui.theme.GetCurrent().Name + " ] "
		}

		// add the second part of the selector when the text is centered
		if mode == "middle" {
			if i == selected {
				thstring += " <<"
			}
			thstring += " "
			if len(thstring)%2 == 1 {
				thstring += " "
			}
		}

		final = append(final, thstring)
	}

	return final
}
