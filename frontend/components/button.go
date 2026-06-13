package components

import (
	"image"
	"image/color"
	"pdr/frontend/constants"

	"gioui.org/io/semantic"
	"gioui.org/layout"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
	"gioui.org/widget/material"
)

type ButtonType uint8

const (
	ButtonTypeDefault ButtonType = iota
	ButtonTypePrimary
	ButtonTypeInfo
	ButtonTypeWarning
	ButtonTypeError
)

func (bt ButtonType) isDefault() bool {
	return bt == ButtonTypeDefault
}

var buttonColorMap = map[ButtonType]color.NRGBA{
	ButtonTypeDefault: constants.COLOR_DEFAULT,
	ButtonTypePrimary: constants.COLOR_PIMARY,
	ButtonTypeInfo:    constants.COLOR_INFO,
	ButtonTypeWarning: constants.COLOR_WARNING,
	ButtonTypeError:   constants.COLOR_ERROR,
}

type Button struct {
	click  widget.Clickable
	label  string
	bg     color.NRGBA
	fg     color.NRGBA
	radius unit.Dp
	bType  ButtonType
}

func NewButton(label string, bt ButtonType) *Button {
	var (
		textColor = constants.COLOR_PRIMARY_TEXT
	)

	if bt.isDefault() {
		textColor = constants.COLOR_DEFAULT_TEXT
	}

	return &Button{
		label:  label,
		bg:     buttonColorMap[bt],
		fg:     textColor,
		radius: unit.Dp(constants.ELEMENTS_RADIUS),
		bType:  bt,
	}
}

func (b *Button) Clicked(gtx layout.Context) bool {
	return b.click.Clicked(gtx)
}

func (b *Button) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	semantic.Button.Add(gtx.Ops)

	if b.click.Pressed() {
		b.bg.A = b.bg.A - 55
	} else {
		b.bg = buttonColorMap[b.bType]
	}

	return b.click.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		// фон кнопки
		rr := gtx.Dp(b.radius)
		bounds := clip.RRect{
			Rect: image.Rectangle{Max: gtx.Constraints.Min},
			SE:   rr,
			SW:   rr,
			NW:   rr,
			NE:   rr,
		}
		defer bounds.Push(gtx.Ops).Pop()

		paint.ColorOp{Color: b.bg}.Add(gtx.Ops)
		paint.PaintOp{}.Add(gtx.Ops)

		// текст по центру
		return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
			l := material.Label(th, unit.Sp(14), b.label)
			l.Color = b.fg
			return l.Layout(gtx)
		})
	})
}
