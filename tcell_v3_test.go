package tview

import (
	"testing"

	"github.com/gdamore/tcell/v3"
)

func TestNormalizeKeyEventCompatibility(t *testing.T) {
	t.Run("shifted tab is backtab", func(t *testing.T) {
		got := normalizeKeyEvent(tcell.NewEventKey(tcell.KeyTab, "", tcell.ModShift))
		if got.Key() != tcell.KeyBacktab {
			t.Fatalf("key = %v, want KeyBacktab", got.Key())
		}
		if got.Modifiers()&tcell.ModShift != 0 {
			t.Fatalf("modifiers = %v, shifted-tab compatibility should consume ModShift", got.Modifiers())
		}
	})

	t.Run("backspace2 is backspace", func(t *testing.T) {
		got := normalizeKeyEvent(tcell.NewEventKey(tcell.KeyBackspace2, "", tcell.ModNone))
		if got.Key() != tcell.KeyBackspace {
			t.Fatalf("key = %v, want KeyBackspace", got.Key())
		}
	})

	t.Run("control modifier and key remain available", func(t *testing.T) {
		got := normalizeKeyEvent(tcell.NewEventKey(tcell.KeyCtrlA, "a", tcell.ModCtrl))
		if got.Key() != tcell.KeyCtrlA || got.Str() != "a" || got.Modifiers()&tcell.ModCtrl == 0 {
			t.Fatalf("Ctrl-A lost v3 key data: key=%v str=%q modifiers=%v", got.Key(), got.Str(), got.Modifiers())
		}
	})

	t.Run("Ctrl-C accepts legacy key and v3 rune forms", func(t *testing.T) {
		if !isCtrlC(tcell.NewEventKey(tcell.KeyCtrlC, "c", tcell.ModCtrl)) {
			t.Fatal("legacy Ctrl-C key was not recognized")
		}
		if !isCtrlC(tcell.NewEventKeyEx(tcell.KeyRune, "c", tcell.ModCtrl, true, 0, 1)) {
			t.Fatal("v3 Ctrl-C rune+modifier was not recognized")
		}
		if isCtrlC(tcell.NewEventKey(tcell.KeyRune, "c", tcell.ModNone)) {
			t.Fatal("plain c was treated as Ctrl-C")
		}
	})
}

func TestInputFieldPreservesAndDeletesGraphemeKey(t *testing.T) {
	field := NewInputField()
	handle := field.InputHandler()
	handle(tcell.NewEventKey(tcell.KeyRune, "👩🏽‍💻", tcell.ModNone), nil)
	if got := field.GetText(); got != "👩🏽‍💻" {
		t.Fatalf("typed grapheme = %q, want %q", got, "👩🏽‍💻")
	}
	handle(normalizeKeyEvent(tcell.NewEventKey(tcell.KeyBackspace2, "", tcell.ModNone)), nil)
	if got := field.GetText(); got != "" {
		t.Fatalf("text after backspace = %q, want empty string", got)
	}
}

func TestInputFieldAdvancedControlKeys(t *testing.T) {
	field := NewInputField().SetText("hello")
	handle := field.InputHandler()
	handle(tcell.NewEventKeyEx(tcell.KeyRune, "u", tcell.ModCtrl, true, 0, 1), nil)
	if got := field.GetText(); got != "" {
		t.Fatalf("advanced Ctrl-U must clear the field, got %q", got)
	}
	handle(tcell.NewEventKeyEx(tcell.KeyRune, " ", tcell.ModCtrl, true, 0, 1), nil)
	if got := field.GetText(); got != "" {
		t.Fatalf("Ctrl-Space must not insert a space, got %q", got)
	}
}
