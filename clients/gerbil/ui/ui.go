package ui

import (
	"gerbil/util"

	"github.com/gdamore/tcell/v2"
)

type Vector2 struct {
	X, Y int
}

type Rectangle struct {
	Pos           Vector2
	Size          Vector2
	Outline_Color tcell.Color
	Inside_Color  tcell.Color
}

func (rect Rectangle) Render(scr tcell.Screen) {
	//	inside_style := tcell.StyleDefault.Background(rect.Inside_Color)
	outside_style := tcell.StyleDefault.Background(rect.Outline_Color)
	if rect.Inside_Color == rect.Outline_Color {
		for y := rect.Pos.Y; y < (rect.Pos.Y + rect.Size.X); y++ {
			for x := rect.Pos.X; x < (rect.Pos.X + rect.Pos.Y); x++ {
				scr.SetContent(x, y, ' ', nil, outside_style)
			}
		}
	} else {
		util.Todo()
	}

}

type Displayed_String struct {
	Contents string
	Pos      Vector2
	Width    int
}
type Message_Display_Box struct {
	Shape    Rectangle
	Contents []string
}
type Contacts_Bar struct {
	Shape    Rectangle
	Contacts []string
}
type Input_Box struct {
	Shape Rectangle
	Input string
}

func New_Rectangle(x, y, L, W int, outline, inside tcell.Color) (ret Rectangle) {
	ret.Pos.X = x
	ret.Pos.Y = y
	ret.Size.X = L
	ret.Size.Y = W
	ret.Inside_Color = inside
	ret.Outline_Color = outline
	return
}
