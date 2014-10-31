package graphics

import "github.com/nsf/termbox-go"

// Tbprint print the given msg starting at the position x and y
// fg is the Text color
// bg is the back ground color
func Tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}
