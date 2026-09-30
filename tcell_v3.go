package tview

import (
	"unicode/utf8"

	"github.com/gdamore/tcell/v3"
)

// normalizeKeyEvent keeps the key identities historically consumed by tview
// stable across tcell versions. Advanced key reporting in tcell v3 represents
// Shift-Tab as KeyTab+ModShift, while the legacy tview handlers expect
// KeyBacktab. Likewise, old backspace variants are normalized by the driver,
// but injected events should behave identically.
func normalizeKeyEvent(event *tcell.EventKey) *tcell.EventKey {
	if event == nil {
		return nil
	}
	if event.Key() == tcell.KeyRune && event.Modifiers()&tcell.ModCtrl != 0 {
		rr := []rune(event.Str())
		if len(rr) == 1 {
			r := rr[0]
			if r >= 'a' && r <= 'z' {
				return tcell.NewEventKey(tcell.KeyCtrlA+tcell.Key(r-'a'), event.Str(), event.Modifiers())
			}
			if r >= 'A' && r <= 'Z' {
				return tcell.NewEventKey(tcell.KeyCtrlA+tcell.Key(r-'A'), event.Str(), event.Modifiers())
			}
			if r == ' ' {
				return tcell.NewEventKey(tcell.KeyNUL, event.Str(), event.Modifiers())
			}
		}
	}
	switch event.Key() {
	case tcell.KeyTab:
		if event.Modifiers()&tcell.ModShift != 0 {
			return tcell.NewEventKey(tcell.KeyBacktab, "", event.Modifiers()&^tcell.ModShift)
		}
	case tcell.KeyBackspace2:
		return tcell.NewEventKey(tcell.KeyBackspace, "", event.Modifiers())
	}
	return event
}

func isCtrlC(event *tcell.EventKey) bool {
	return event != nil && (event.Key() == tcell.KeyCtrlC ||
		(event.Key() == tcell.KeyRune && event.Modifiers()&tcell.ModCtrl != 0 && event.Str() == "c"))
}

// eventRune provides the legacy single-rune view of a v3 key event. Most
// existing navigation handlers compare one rune; callers that accept text
// should use EventKey.Str() instead so composed graphemes are preserved.
func eventRune(event *tcell.EventKey) rune {
	if event == nil {
		return 0
	}
	r, _ := utf8.DecodeRuneInString(event.Str())
	return r
}

// decomposeStyle adapts the old tcell Style.Decompose API to the v3 getters.
func decomposeStyle(style tcell.Style) (fg, bg tcell.Color, attrs tcell.AttrMask) {
	return style.GetForeground(), style.GetBackground(), style.GetAttributes()
}

// getScreenContent adapts v3's grapheme-oriented Get API to the rune and
// combining-mark representation used by tview's existing drawing code.
func getScreenContent(screen tcell.Screen, x, y int) (main rune, combining []rune, style tcell.Style, width int) {
	str, style, width := screen.Get(x, y)
	if str == "" {
		return ' ', nil, style, width
	}
	runes := []rune(str)
	return runes[0], runes[1:], style, width
}
