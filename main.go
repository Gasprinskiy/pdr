package main

import (
	"log"
	"os"
	"pdr/frontend/constants"
	"pdr/frontend/screens"

	"gioui.org/app"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget/material"
)

func init() {
	// runtime.GOMAXPROCS(1)
}

func main() {
	go func() {
		w := new(app.Window)
		w.Option(
			app.Title("PDR"),
			app.Size(unit.Dp(800), unit.Dp(600)),
		)
		if err := loop(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()
	app.Main()
}

func loop(w *app.Window) error {
	var ops op.Ops

	th := material.NewTheme()
	screen := screens.NewMainScreen()

	for {
		e := w.Event()
		switch e := e.(type) {
		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)

			paint.Fill(gtx.Ops, constants.COLOR_PRIMARY_BACKGROUND)
			screen.Layout(gtx, th)

			e.Frame(gtx.Ops)

		case app.DestroyEvent:
			return e.Err
		}
	}
}
