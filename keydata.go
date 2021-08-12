package main

import "strings"

// ScreenKeyData struct contains info about key what events are
// dispatchable on the current screen.
// Keys is not of type KeyEvent because it is needed for the events
// to be quickly accessed via the key.
// AltKeys are used for special cases, per example, when the game is
// paused. Thus, it is optional (nil can be passed)
type ScreenKeyData struct {
	Keys    *KeyEventMap
	AltKeys *KeyEventMap
	Legend  *Legend
}

type Legend struct {
	// type can be either left or middle
	Type     string
	Keys     *KeyEventMap
	Selected int
}

func NewScreenKeyData(k *KeyEventMap, al *KeyEventMap, lk *KeyEventMap, lt string) *ScreenKeyData {
	legend := &Legend{lt, lk, 0}
	return &ScreenKeyData{k, al, legend}
}

var screens = map[string]*ScreenKeyData{
	screenStartMenu: NewScreenKeyData(&keysStart, nil, &legendKeysStart, "middle"),
	screenGame:      NewScreenKeyData(&keysGame, &altKeysGame, &legendKeysGame, "left"),
}

// Fetches the keys and event descriptions and formats them
// into a format suitable for showing on the screen as a legend.
// Mode dictates should keys be format for left or middle.
// May be moved to another file in the future.
func (g *Game) formatKeys() []string {
	mode := g.keyData.Legend.Type
	keys := g.keyData.Legend.Keys
	selected := g.keyData.Legend.Selected

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
	for _, keyev := range *keys {
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
	for i, keyevt := range *keys {
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
			if !g.aiActive {
				thstring += "IN"
			}
			thstring += "ACTIVE ]"
		case eventSwitchTheme:
			thstring += ": [ " + g.theme.GetCurrent().Name + " ] "
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
