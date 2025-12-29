package cannoli

import (
	"github.com/BrandonKowalski/gabagool/v2/pkg/gabagool/internal"
)

func InitCannoliTheme(fontPath string) internal.Theme {
	return internal.Theme{
		HighlightColor:       internal.HexToColor(0xFFFFFF),
		AccentColor:          internal.HexToColor(0x008080),
		ButtonLabelColor:     internal.HexToColor(0x000000),
		HintColor:            internal.HexToColor(0x000000),
		TextColor:            internal.HexToColor(0xFFFFFF),
		HighlightedTextColor: internal.HexToColor(0x000000),
		BackgroundColor:      internal.HexToColor(0xFFFFFF),
		FontPath:             fontPath,
	}
}
