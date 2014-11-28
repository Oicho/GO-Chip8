package graphics

import "github.com/nsf/termbox-go"

// PrintString print the given msg starting at the position x and y
// fg is the Text color
// bg is the back ground color
func PrintString(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, termbox.ColorWhite, bg)
		x++
	}
}

// PrintScreen Print a boolean array on the screen
func PrintScreen(screen [][]bool) {
	point := 0
	for x, superArr := range screen {
		for y, b := range superArr {
			if b {
				termbox.SetCell(x+point, y, ' ', termbox.ColorDefault, termbox.ColorWhite)
				termbox.SetCell(x+point+1, y, ' ', termbox.ColorDefault, termbox.ColorWhite)
			}
		}
		point++
	}
	termbox.Flush()
}
