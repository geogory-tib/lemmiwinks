package main

import (
	"gerbil/ui"
	"github.com/gdamore/tcell/v2"
	"log"
)

func main() {
	app_screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatal(err)
	}
	defer app_screen.Fini()
	err = app_screen.Init()
	if err != nil {
		log.Fatal(err)
	}
	app_screen.SetStyle(tcell.StyleDefault)
	app_screen.Clear()
	rect := ui.New_Rectangle(0, 0, 50, 50, tcell.ColorDarkGreen, tcell.ColorDarkGreen)
	rect.Render(app_screen)
	app_screen.Show()
	for {
	}
}
