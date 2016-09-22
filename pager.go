package main

import (
	// "fmt"
	"github.com/nsf/termbox-go"
	// "log"
	"regexp"
	// "strings"
)

const (
	QUIT      = 0
	NO_ACTION = 1
)

type Pager struct {
	str          string   // contents to display
	linenum        int      // num of lines in str
	lines []string
	File         string   // current file
	lineSelected string
	posY int
	offX, offY int
}

func (p *Pager) SetContent(s string) {
	p.str = s
	p.lines = regexp.MustCompile("\r\n|\n\r|\n|\r").Split(p.str, -1)
	p.offX = 0
	p.offY = 0
	p.posY = 0
	p.linenum = len(p.lines)
}

func (p *Pager) drawStatusLine(str string) {
	runes := []rune(str)
	for i := range(runes) {
		termbox.SetCell(i, 0, runes[i], termbox.ColorBlue, termbox.ColorWhite)
	}
}

func termLineNum() int {
	_, maxY := termbox.Size()
	return maxY
}

func termWidth() int {
	maxX, _:= termbox.Size()
	return maxX
}


func bodyLineNum() int {
	return termLineNum() -1
}



func (p *Pager) drawOneLine(y int) {
	line := p.lines[y]
	bodyOffsetY := 1

	screenY := y - p.offY + bodyOffsetY
	if screenY >= termLineNum()  || screenY < 1{
		return
	}


	color := termbox.ColorDefault
	backgroundColor := termbox.ColorDefault

	runes := []rune(line)

	for i := 0; i < min(termWidth(), len(runes)-p.offX); i++ {
		screenX := i
		termbox.SetCell(screenX,  screenY, runes[i+p.offX], color, backgroundColor)
	}
}

func (p *Pager) drawAllLines() {
	for j := p.offY; j < min(p.offY+bodyLineNum(), p.linenum); j++ {
		p.drawOneLine(j)
	}
}

func (p *Pager) Clear() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	termbox.Flush()
	termbox.Sync()
	p.Draw()
}
func (p *Pager) Draw() {
	p.DrawInternal(false)
}

func (p *Pager) DrawWithRefresh() {
	p.DrawInternal(true)
}

func (p *Pager) DrawInternal(needRefresh bool) {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	if needRefresh {
		termbox.Flush()
	}


	p.drawAllLines()
	// maxX, _ := termbox.Size()
	// empty := make([]byte, maxX)
	mode := ""
	file := ""
	nextFileUsage := ""
	if p.File != "" {
		file = " :: [file: " + p.File + " ]"
	}
	if file != "" {
		mode = file
	}
	// p.drawStatusLine( "USAGE [exit: ESC/q] [scroll: j,k/C-n,C-p] "+nextFileUsage+mode+string(empty))
	p.drawStatusLine( "USAGE [exit: ESC/q] [scroll: j,k/C-n,C-p] "+nextFileUsage+mode)
	termbox.SetCursor(0, p.posY-p.offY+1)
	termbox.Flush()
}




func (p *Pager) incrementLines(delta int) {
	prevIY := p.offY
	bodyLNum := bodyLineNum()
	p.offY = max(0, min(p.offY+delta, p.linenum -  bodyLNum ))
	p.posY = p.offY
	if p.offY - prevIY < delta {
		p.posY = p.linenum -1 
	}
}

func min(x, y int) int{
	if x < y {
		return x
	}
	return y
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}


func (p *Pager) viewModeKey(ev termbox.Event) int {
	switch ev.Key {
	case termbox.KeyEnter:
		defer termbox.Flush()
		return p.onEnter()
	case termbox.KeyEsc, termbox.KeyCtrlC:
		termbox.Flush()
		return QUIT
	case termbox.KeyCtrlL:
		termbox.Flush()
		termbox.Sync()
	case termbox.KeyArrowRight:
		p.scrollRight()
	case termbox.KeyArrowLeft:
		p.scrollLeft()
	case termbox.KeyCtrlN, termbox.KeyArrowDown:
		p.scrollDown()
	case termbox.KeyCtrlP, termbox.KeyArrowUp:
		p.scrollUp()
	case termbox.KeyCtrlD, termbox.KeySpace:
		p.incrementLines(29)
		p.DrawWithRefresh()
	case termbox.KeyCtrlU, termbox.KeyCtrlB:
		p.offY = max(0, p.offY - 29)
		p.posY = p.offY
		p.DrawWithRefresh()
	default:
		switch ev.Ch {
		case 'j':
			p.scrollDown()
		case 'k':
			p.scrollUp()
		case 'l':
			p.scrollRight()
		case 'h':
			p.scrollLeft()
		case 'q':
			termbox.Sync()
			return QUIT
		case 'b':
			p.offY = max(0, p.offY - 29)
			p.posY = p.offY
			p.DrawWithRefresh()
		case '<':
			p.offY = 0
			p.posY = 0
			p.offX = 0
			termbox.Sync()
			p.Draw()
		case '>':
			_, y := termbox.Size()
			p.offY = p.linenum - y + 1
			p.posY = p.linenum -1
			p.offX = 0
			if p.offY < 0 {
				p.offY = 0
			}
			termbox.Sync()
			p.Draw()
		default:
			p.Draw()
		}
	}
	return NO_ACTION
}

func (p *Pager) PollEvent() bool {
	p.Draw()
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			ret := p.viewModeKey(ev)
			if ret == QUIT {
				return true
			}
			p.Draw()
		default:
			p.Draw()
		}
	}
	return false
}


func (p *Pager) scrollDown() {
	bodyLNum := bodyLineNum()
	withRefresh := false
	
	p.posY = min(p.posY+1, p.linenum-1)
	if p.posY > p.offY + bodyLNum -1 && p.offY < p.linenum-bodyLNum  {
		p.offY++
		withRefresh = true
	}
	p.DrawInternal(withRefresh)
}

func (p *Pager) scrollUp() {
	withRefresh := false
	p.posY = max(p.posY-1, 0)
	if p.posY < p.offY {
		p.offY--
		withRefresh = true
	}
	p.DrawInternal(withRefresh)
}

func (p *Pager) scrollRight() {
	p.offX += 29
	p.DrawWithRefresh()
}

func (p *Pager) scrollLeft() {
	p.offX = max(0, p.offX - 29)
	p.DrawWithRefresh()
}

func (p *Pager) Init() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
}

func (p *Pager) Close() {
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	termbox.Flush()
	termbox.Sync()
	termbox.Close()
}


func (p *Pager) onEnter() int {
	if p.linenum ==0 {
		return NO_ACTION	
	}

	// log.Fatal("len %d, pos %d, linnum %d", len(p.lines), p.posY, p.linenum)	
	p.lineSelected = p.lines[p.posY]
	// log.Fatal("len %d, pos %d, linnum %d", len(p.lines), p.posY, p.linenum, p.lineSelected)	
	// p.lineSelected = "Test!:fuafuga"
	return QUIT
}





