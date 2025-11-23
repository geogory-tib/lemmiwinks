package main

import (
	"gerbil/ui"
	"log"
	"github.com/gdamore/tcell/v2"
)

const LONG_EXAMPLE_TEXT = `Person A: Hey, did you finish the homework yet?

Person B: Not yet, I’m still stuck on question three.

Person A: Which part is confusing you?

Person B: The subnetting section. It never clicks for me.

Person A: Want me to walk you through it?

Person B: Yes please, that would help a lot.

Person A: Okay, what IP range did they give you?

Person B: It starts with 192.168.20.0/26.

Person A: Alright, so first calculate how many hosts each subnet supports…

Person B: Ohhh, got it now—thanks!`

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
	ui.RenderNewLines = true
	input_box := ui.New_Input_Box(0,0,50,50,tcell.ColorLightGrey,tcell.ColorBlue)
	input_box.Input[0] = []rune("Hello This is a string\n")
	input_box.Input[1] = []rune("Hi")
	input_box.Render(app_screen)
	app_screen.Show()
	for {
		event := app_screen.PollEvent()
		switch a := event.(type) {
		case *tcell.EventKey:
			if a.Rune() == 'i'{
				input_box.Get_Input(app_screen)
			}
			return
		}
	}
}
