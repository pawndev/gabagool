package gabagool

import (
	"sync/atomic"
	"time"

	"github.com/BrandonKowalski/gabagool/v2/pkg/gabagool/internal"
	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

// TimeFormat specifies 12-hour or 24-hour clock display
type TimeFormat int

const (
	TimeFormat24Hour TimeFormat = iota
	TimeFormat12Hour
)

// DynamicStatusBarIcon allows goroutines to safely update icon content.
// Use SetText to update from any goroutine.
type DynamicStatusBarIcon struct {
	text atomic.Value // stores string
}

// NewDynamicStatusBarIcon creates a new dynamic icon with initial text
func NewDynamicStatusBarIcon(initialText string) *DynamicStatusBarIcon {
	d := &DynamicStatusBarIcon{}
	d.text.Store(initialText)
	return d
}

// SetText updates the icon text (goroutine-safe)
func (d *DynamicStatusBarIcon) SetText(s string) {
	d.text.Store(s)
}

// GetText returns the current text
func (d *DynamicStatusBarIcon) GetText() string {
	if v := d.text.Load(); v != nil {
		return v.(string)
	}
	return ""
}

// StatusBarIcon represents a single icon slot in the status bar
type StatusBarIcon struct {
	Text    string                // Icon text (font glyph/symbol)
	Dynamic *DynamicStatusBarIcon // If set, reads from this instead of static Text
}

// StatusBarOptions configures the status bar appearance and behavior
type StatusBarOptions struct {
	Enabled    bool
	ShowTime   bool
	TimeFormat TimeFormat
	Icons      []StatusBarIcon // Max 3 icons
}

// DefaultStatusBarOptions returns sensible defaults with the status bar disabled
func DefaultStatusBarOptions() StatusBarOptions {
	return StatusBarOptions{
		Enabled:    false,
		ShowTime:   true,
		TimeFormat: TimeFormat24Hour,
		Icons:      nil,
	}
}

// renderStatusBar renders the status bar in the top-right corner of the component
func renderStatusBar(
	renderer *sdl.Renderer,
	font *ttf.Font,
	symbolFont *ttf.Font,
	options StatusBarOptions,
	margins internal.Padding,
) {
	if !options.Enabled {
		return
	}

	scaleFactor := internal.GetScaleFactor()
	window := internal.GetWindow()
	windowWidth, _ := window.Window.GetSize()

	padding := int32(float32(15) * scaleFactor)
	iconSpacing := int32(float32(10) * scaleFactor)

	// Start from right edge
	currentX := windowWidth - margins.Right - padding
	contentY := margins.Top + padding

	// Get time text height for vertical alignment
	var timeHeight int32
	if options.ShowTime {
		timeText := formatCurrentTime(options.TimeFormat)
		surface, err := font.RenderUTF8Blended(timeText, internal.GetTheme().MainColor)
		if err == nil && surface != nil {
			timeHeight = surface.H
			surface.Free()
		}
	}

	// 1. Render time (rightmost element)
	if options.ShowTime {
		timeText := formatCurrentTime(options.TimeFormat)
		currentX = renderStatusBarTime(renderer, font, timeText, currentX, contentY)
		currentX -= iconSpacing
	}

	// 2. Render icons (up to 3, right to left), vertically centered with time
	maxIcons := 3
	if len(options.Icons) < maxIcons {
		maxIcons = len(options.Icons)
	}

	// Icons render right-to-left (last icon closest to time)
	for i := maxIcons - 1; i >= 0; i-- {
		icon := options.Icons[i]
		currentX = renderStatusBarIcon(renderer, symbolFont, icon, currentX, contentY, timeHeight)
		if i > 0 {
			currentX -= iconSpacing
		}
	}
}

func formatCurrentTime(format TimeFormat) string {
	now := time.Now()
	switch format {
	case TimeFormat12Hour:
		return now.Format("3:04 PM")
	case TimeFormat24Hour:
		return now.Format("15:04")
	default:
		return now.Format("15:04")
	}
}

func renderStatusBarTime(
	renderer *sdl.Renderer,
	font *ttf.Font,
	timeText string,
	rightX, y int32,
) int32 {
	textColor := internal.GetTheme().MainColor

	surface, err := font.RenderUTF8Blended(timeText, textColor)
	if err != nil || surface == nil {
		return rightX
	}
	defer surface.Free()

	texture, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		return rightX
	}
	defer texture.Destroy()

	textX := rightX - surface.W
	rect := sdl.Rect{X: textX, Y: y, W: surface.W, H: surface.H}
	renderer.Copy(texture, nil, &rect)

	return textX
}

func renderStatusBarIcon(
	renderer *sdl.Renderer,
	font *ttf.Font,
	icon StatusBarIcon,
	rightX, y, lineHeight int32,
) int32 {
	// Resolve text (dynamic takes priority)
	var text string
	if icon.Dynamic != nil {
		text = icon.Dynamic.GetText()
	} else {
		text = icon.Text
	}

	if text == "" {
		return rightX
	}

	textColor := internal.GetTheme().MainColor
	surface, err := font.RenderUTF8Blended(text, textColor)
	if err != nil || surface == nil {
		return rightX
	}
	defer surface.Free()

	texture, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		return rightX
	}
	defer texture.Destroy()

	// Position text at rightX, vertically centered with line height
	textX := rightX - surface.W
	textY := y + (lineHeight-surface.H)/2
	rect := sdl.Rect{X: textX, Y: textY, W: surface.W, H: surface.H}
	renderer.Copy(texture, nil, &rect)
	return textX
}
