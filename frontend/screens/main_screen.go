package screens

import (
	"fmt"
	"pdr/frontend/components"

	"gioui.org/layout"
	"gioui.org/widget/material"
	"github.com/ncruces/zenity"
)

type MainScreen struct {
	openBtn *components.Button
}

func NewMainScreen() *MainScreen {
	return &MainScreen{
		openBtn: components.NewButton("Открыть PDF", components.ButtonTypePrimary),
	}
}

func (s *MainScreen) Layout(gtx layout.Context, th *material.Theme) layout.Dimensions {
	if s.openBtn.Clicked(gtx) {

		filename, err := zenity.SelectFile(
			zenity.Title("Выбрать PDF"),
			zenity.FileFilters{{
				Name:     "PDF files",
				Patterns: []string{"*.pdf"},
			}},
		)
		if err == nil {
			fmt.Println("FILENAME: ", filename)
		} else {
			fmt.Println("err: ", err)
		}
		// открыть файл
	}

	return layout.Center.Layout(gtx, func(gtx layout.Context) layout.Dimensions {
		gtx.Constraints.Min.X = gtx.Dp(120)
		gtx.Constraints.Min.Y = gtx.Dp(40)
		return s.openBtn.Layout(gtx, th)
	})
}
