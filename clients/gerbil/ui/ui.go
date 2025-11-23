
package ui
import (
	"github.com/gdamore/tcell/v2"
	"slices"
)

var RenderNewLines bool // will make a new line if a newline is present
type Vector2 struct {
	X, Y int
}

type Rectangle struct {
	Pos           Vector2
	Size          Vector2
	Outline_Color tcell.Color
	Inside_Color  tcell.Color
}

func (vec Vector2) add(vec2 Vector2) (ret Vector2) {
	ret.X = vec.X + vec2.X
	ret.Y = vec.Y + vec2.Y
	return
}

func (vec Vector2) addY(vec2 Vector2) (ret Vector2) {
	ret.X = vec.X
	ret.Y = vec.Y + vec2.Y
	return
}

func (vec Vector2) addX(vec2 Vector2) (ret Vector2) {
	ret.Y = vec.Y
	ret.X = vec.X + vec2.X
	return
}

func (rect Rectangle) Render(scr tcell.Screen) {
	inside_style := tcell.StyleDefault.Background(rect.Inside_Color)
	outside_style := tcell.StyleDefault.Background(rect.Outline_Color)
	if rect.Inside_Color == rect.Outline_Color {
		for y := rect.Pos.Y; y < (rect.Pos.Y + rect.Size.X); y++ {
			for x := rect.Pos.X; x < (rect.Pos.X + rect.Size.Y); x++ {
				scr.SetContent(x, y, ' ', nil, outside_style)
			}
		}
	} else {

		draw_bar := func(start_pos, end_pos Vector2) {
			for x := start_pos.X; x <= end_pos.X; x++ {
				scr.SetContent(x, start_pos.Y, ' ', nil, outside_style)
			}
		}
		draw_bar(rect.Pos, rect.Pos.addX(Vector2{rect.Size.Y, 0}))
		for y := rect.Pos.Y + 1; y < rect.Pos.addY(rect.Size).Y-1; y++ {
			scr.SetContent(rect.Pos.X, y, ' ', nil, outside_style)
			scr.SetContent(rect.Pos.addX(rect.Size).X, y, ' ', nil, outside_style)
		}
		draw_bar(Vector2{rect.Pos.X + rect.Pos.Y, rect.Pos.Y + rect.Size.X - 1}, rect.Pos.add(rect.Size))
		for y := rect.Pos.Y + 1; y < (rect.Pos.Y + 1 + rect.Size.X - 2); y++ {
			for x := rect.Pos.X + 1; x < (rect.Pos.X + 1 + rect.Size.Y - 1); x++ {
				scr.SetContent(x, y, ' ', nil, inside_style)
			}
		}

	}
}

type Displayed_String struct {
	Contents string
	Pos      Vector2
	Width    int
}

func (str Displayed_String) Render(scr tcell.Screen) {
	y := str.Pos.Y
	X := 0
	for _, ch := range str.Contents {
		if X > str.Width || ch == '\n' && RenderNewLines {
			y++
			X = 0
		}
		if ch != '\n'{
			scr.SetContent((X), y, ch, nil, tcell.StyleDefault)
			X++
		}
	}
}

type Message_Display_Box struct {
	Shape    Rectangle
	Contents []string
	Text_Color tcell.Color
	Scroll_Offset int
}
func (msg Message_Display_Box)Render(app_scrn tcell.Screen) {
	msg.Shape.Render(app_scrn)
	X := 1
	Y := 1
	for _,str := range msg.Contents[msg.Scroll_Offset:]{
		for _,ch := range str{
			if X >=  msg.Shape.Size.Y{
				X = 1
				Y++
			}
			if Y > msg.Shape.Size.X{
				return
			}
			app_scrn.SetContent((X + msg.Shape.Pos.X),(Y + msg.Shape.Pos.Y ),ch,nil,tcell.StyleDefault.Background(msg.Shape.Inside_Color))

			X++
			
		}
		X = 1
		Y++
	}

}
type Contacts_Bar struct {
	Shape    Rectangle
	Contacts []string
}
type Input_Box struct {
	Shape Rectangle
	Input [][]rune
	Cursor_pos Vector2
	ESC_PRESSED bool
	Scroll_Offset int
}

func (box Input_Box)Render(app_scrn tcell.Screen){
	box.Shape.Render(app_scrn)
	X := 1
	Y := 1
	render_text := box.Input[box.Scroll_Offset:]
	for _,str := range render_text{
		for _ ,ch := range str{
			app_scrn.SetContent((X + box.Shape.Pos.X),(Y + box.Shape.Pos.Y),ch,nil,tcell.StyleDefault.Background(box.Shape.Inside_Color))
			X++
		}
		Y++
		if Y == (box.Shape.Pos.Y + box.Shape.Size.X - 1){
			break
		}
	}
	app_scrn.Show()
}
func (box *Input_Box)Get_Input(scrn tcell.Screen) []rune{
	scrn.ShowCursor(box.Shape.Pos.X + 1, box.Shape.Pos.Y + 1)
	box.Render(scrn)
	for {
		ev := scrn.PollEvent()
		switch event :=  ev.(type){
			case *tcell.EventKey:
			key := event.Key()
			if key == tcell.KeyRune{
				ch := event.Rune()
				box.handle_rune(ch)
				
			}else{
				box.handle_control_keys(scrn,key)
				if box.ESC_PRESSED{
					return slices.Concat(box.Input...)
				}
				
			}
			scrn.ShowCursor(box.Shape.Pos.X + 1, box.Shape.Pos.Y + 1)
			box.Render(scrn)
		}
	}
}
func (box *Input_Box)handle_rune(ch rune){
	renderable_text := box.Input[box.Scroll_Offset:]
	current_line := renderable_text[box.Cursor_pos.Y]
	if box.Cursor_pos.X - 1 >= len(current_line){
		if len(current_line) + 1 >= box.Shape.Size.Y - 1{
			next_line := renderable_text[box.Cursor_pos.Y]
			next_line = slices.Insert(next_line,0,ch)
			box.Cursor_pos.X = 1
			box.Cursor_pos.Y++
			return
		}
		current_line = append(current_line,ch)
		box.Cursor_pos.X++ 
	}else if box.Cursor_pos.X - 1 < len(current_line){
		if len(current_line) + 1 >= box.Shape.Size.Y - 1{
			next_line := renderable_text[box.Cursor_pos.Y]
			next_line = slices.Insert(next_line,0,current_line[len(current_line) - 1])
			box.Cursor_pos.X = 1
			box.Cursor_pos.Y++
			return
		}
		current_line = slices.Insert(current_line,box.Cursor_pos.Y,ch)
	}
}
func (box *Input_Box)handle_control_keys(scrn tcell.Screen, key tcell.Key){
	renderable_text := box.Input[box.Scroll_Offset:]
	switch key{
		case tcell.KeyESC:
		{
			box.ESC_PRESSED = true
			scrn.HideCursor()
			return
		}
		case tcell.KeyEnter:
		{
			current_line := renderable_text[box.Cursor_pos.Y]
			if box.Cursor_pos.X >= len(current_line){
				current_line = append(current_line,'\n')
				box.Cursor_pos.X = 1
				box.Cursor_pos.Y++
			}else{
				next_line := renderable_text[box.Cursor_pos.Y]
				if len(next_line) == 0{
					next_line = append(next_line,current_line[box.Cursor_pos.X:]...)
					current_line = current_line[:box.Cursor_pos.X]
					current_line = append(current_line,'\n')
					box.Cursor_pos.X = 1
					box.Cursor_pos.Y++
				}else{
					next_line = slices.Insert(next_line,0,current_line[box.Cursor_pos.X:]...)
					current_line = current_line[:box.Cursor_pos.X]
					current_line = append(current_line,'\n')
					box.Cursor_pos.X = 1
					box.Cursor_pos.Y++
				}
			}
		}
		case tcell.KeyLeft:
		{
			box.Cursor_pos.X--
			if box.Cursor_pos.X == box.Shape.Pos.X{
				box.Cursor_pos.Y--
				box.Cursor_pos.X = len(renderable_text[box.Cursor_pos.Y])
			}
		}
		case tcell.KeyRight:
		{
			box.Cursor_pos.X++
			if box.Cursor_pos.X == box.Shape.Size.Y - 1{
				box.Cursor_pos.Y++
				box.Cursor_pos.X = 1
			}
		}
		case tcell.KeyUp:
		{
			box.Cursor_pos.Y--
			if box.Cursor_pos.Y <= box.Shape.Pos.Y && box.Scroll_Offset > 0{
				box.Scroll_Offset--
			}
		}
		case tcell.KeyDown:
		{
			box.Cursor_pos.Y++
			if box.Cursor_pos.Y >= (box.Shape.Pos.Y + box.Shape.Size.X) && box.Scroll_Offset < len(box.Input){
				box.Scroll_Offset--
			}
		}
	}
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

func New_DisplayString(x, y, W int, Contents string) (ret Displayed_String) {
	ret.Contents = Contents
	ret.Pos.Y = y
	ret.Pos.X = x
	ret.Width = W
	return
}

func New_Message_Display_Box(x,y,L,W int, outline,inside tcell.Color)(ret Message_Display_Box){
	ret.Contents = make([]string,0,400)
	ret.Shape.Pos.X = x
	ret.Shape.Pos.Y = y
	ret.Shape.Size.X = L
	ret.Shape.Size.Y = W
	ret.Shape.Inside_Color = inside
	ret.Shape.Outline_Color = outline
	return
}
func New_Input_Box(x,y,L,W int, outline,inside tcell.Color) (ret Input_Box){
	ret.Input = make([][]rune,100)
	for i := range len(ret.Input){
		ret.Input[i] = make([]rune,0,50)
	}
	ret.Shape.Pos.X = x
	ret.Shape.Pos.Y = y
	ret.Shape.Size.X = L
	ret.Shape.Size.Y = W
	ret.Shape.Outline_Color = outline
	ret.Shape.Inside_Color = inside
	return
}
