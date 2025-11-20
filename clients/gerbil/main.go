package main

import (
	"gerbil/ui"
	"log"
	"strings"
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
	message_box := ui.New_Message_Display_Box(0,0,60,60,tcell.ColorGreen,tcell.ColorGray)
	stuff_to_append := strings.Split(LONG_EXAMPLE_TEXT,"\n")
	message_box.Contents = append(message_box.Contents,stuff_to_append...)
	message_box.Scroll_Offset = 4
	message_box.Render(app_screen)
	app_screen.Show()
	for {
		event := app_screen.PollEvent()
		switch a := event.(type) {
		case *tcell.EventKey:
			a.Name()
			return
		}
	}
}
