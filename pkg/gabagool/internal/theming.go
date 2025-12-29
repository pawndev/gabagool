package internal

import (
	"github.com/veandco/go-sdl2/sdl"
)

type Theme struct {
	HighlightColor       sdl.Color // Color1: Selected item background, footer button background
	AccentColor          sdl.Color // Color2: Pill backgrounds, status bar pill
	ButtonLabelColor     sdl.Color // Color3: Button label text (inside pills)
	TextColor            sdl.Color // Color4: Default text color
	HighlightedTextColor sdl.Color // Color5: Text on highlighted items
	HintColor            sdl.Color // Color6: Help text, status bar text
	BackgroundColor      sdl.Color // BGColor: Screen background
	FontPath             string
	BackgroundImagePath  string
}

var currentTheme Theme

func SetTheme(theme Theme) {
	currentTheme = theme
}

func GetTheme() Theme {
	return currentTheme
}
