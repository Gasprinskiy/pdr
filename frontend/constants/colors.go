package constants

import "image/color"

var (
	COLOR_DEFAULT = color.NRGBA{R: 255, G: 255, B: 255, A: 255}
	COLOR_PIMARY  = color.NRGBA{R: 99, G: 226, B: 183, A: 255}
	COLOR_INFO    = color.NRGBA{R: 112, G: 192, B: 232, A: 255}
	COLOR_WARNING = color.NRGBA{R: 242, G: 201, B: 125, A: 255}
	COLOR_ERROR   = color.NRGBA{R: 232, G: 128, B: 125, A: 255}
	//
	COLOR_PRIMARY_BACKGROUND   = color.NRGBA{R: 16, G: 16, B: 20, A: 255}
	COLOR_SECONDARY_BACKGROUND = color.NRGBA{R: 24, G: 24, B: 28, A: 255}
	//

	COLOR_PRIMARY_TEXT = color.NRGBA{R: 0, G: 0, B: 0, A: 255}
	COLOR_DEFAULT_TEXT = color.NRGBA{R: 212, G: 212, B: 213, A: 255}
)

var (
	ELEMENTS_RADIUS uint8 = 3
)
