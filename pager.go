package main

import (
	// "fmt"
	"github.com/gdamore/tcell"
	// "github.com/zyedidia/tcell"
	"github.com/mattn/go-runewidth"
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
	posY, posX int
	offX, offY int
	screen tcell.Screen
}

func (p *Pager) SetContent(s string) {
	p.str = s
	p.lines = regexp.MustCompile("\r\n|\n\r|\n|\r").Split(p.str, -1)
	p.offX = 0
	p.offY = 0
	p.posX = 0
	p.posY = 0
	p.linenum = len(p.lines)
	if p.linenum > 0 && p.lines[p.linenum -1] == "" {
		p.linenum -= 1
	}
}


func putln(s tcell.Screen, row int, str string, style tcell.Style) {
	puts(s, style, 0, row, str)
}

func puts(s tcell.Screen, style tcell.Style, x, y int, str string) {
	i := 0
	var deferred []rune
	dwidth := 0
	for _, r := range str {
		switch runewidth.RuneWidth(r) {
		case 0:
			if len(deferred) == 0 {
				deferred = append(deferred, ' ')
				dwidth = 1
			}
		case 1:
			if len(deferred) != 0 {
				s.SetContent(x+i, y, deferred[0], deferred[1:], style)
				i += dwidth
			}
			deferred = nil
			dwidth = 1
		case 2:
			if len(deferred) != 0 {
				s.SetContent(x+i, y, deferred[0], deferred[1:], style)
				i += dwidth
			}
			deferred = nil
			dwidth = 2
		}
		deferred = append(deferred, r)
	}
	if len(deferred) != 0 {
		s.SetContent(x+i, y, deferred[0], deferred[1:], style)
		i += dwidth
	}
}



func (p *Pager) drawStatusLine(str string) {
	putln(p.screen, 0, str, tcell.StyleDefault.Foreground(tcell.ColorBlue).Background(tcell.ColorWhite))
}


func (p *Pager) Size() (int, int) {
	x, y := p.screen.Size()
	// return x, y-1
	return x, y
}

func (p *Pager) termLineNum() int {
	_, maxY := p.Size()
	return maxY
}

func (p *Pager) termWidth() int {
	maxX, _:= p.Size()
	return maxX
}


func (p *Pager) bodyLineNum() int {
	return p.termLineNum() -1
}



func (p *Pager) drawOneLine(y int) {
	line := p.lines[y]
	bodyOffsetY := 1

	screenY := y - p.offY + bodyOffsetY
	if screenY >= p.termLineNum()  || screenY < 1{
		return
	}


	

	runes := []rune(line)

	for i := 0; i < min(p.termWidth(), len(runes)-p.offX); i++ {
		screenX := i
		p.screen.SetContent(screenX,  screenY, runes[i+p.offX], nil, tcell.StyleDefault)
	}
}

func (p *Pager) drawAllLines() {
	p.screen.SetStyle(tcell.StyleDefault)
	for j := p.offY; j < min(p.offY+p.bodyLineNum(), p.linenum); j++ {
		p.drawOneLine(j)
	}
}

func (p *Pager) Clear() {
	p.screen.Clear()
	p.Draw()
}
func (p *Pager) Draw() {
	p.DrawInternal(false)
}

func (p *Pager) DrawWithRefresh() {
	p.DrawInternal(true)
}

func (p *Pager) DrawInternal(needRefresh bool) {
	p.screen.Clear()
	if needRefresh {
		// termbox.Flush()
		// p.screen.Sync()
		// p.screen.Show()
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
	p.screen.ShowCursor(p.posX-p.offX, p.posY-p.offY+1)
	p.screen.Show()
}




func (p *Pager) incrementLines(delta int) {
	prevIY := p.offY
	bodyLNum := p.bodyLineNum()
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


func (p *Pager) viewModeKey(ev *tcell.EventKey) int {
	switch ev.Key() {
	case tcell.KeyEnter:
		// defer p.screen.Flush()
		return p.onEnter()
	case tcell.KeyEsc, tcell.KeyCtrlC:
		// termbox.Flush()
		return QUIT
	case tcell.KeyCtrlL:
		p.screen.Sync()
		/*
		termbox.Flush()
		termbox.Sync()
		*/
	case tcell.KeyRight:
		p.scrollRight()
	case tcell.KeyLeft:
		p.scrollLeft()
	case tcell.KeyCtrlN, tcell.KeyDown:
		p.scrollDown()
	case tcell.KeyCtrlP, tcell.KeyUp:
		p.scrollUp()
	case tcell.KeyCtrlD:
		p.incrementLines(29)
		p.DrawWithRefresh()
	case tcell.KeyCtrlU, tcell.KeyCtrlB:
		p.offY = max(0, p.offY - 29)
		p.posY = p.offY
		p.DrawWithRefresh()
	default:
		switch ev.Rune() {
		case 'j':
			p.scrollDown()
		case 'k':
			p.scrollUp()
		case 'l':
			p.scrollRight()
		case 'h':
			p.scrollLeft()
		case 'q':
			p.screen.Sync()
			return QUIT
		case 'b':
			p.offY = max(0, p.offY - 29)
			p.posY = p.offY
			p.DrawWithRefresh()
		case ' ':
			p.incrementLines(29)
			p.DrawWithRefresh()
		case '<':
			p.offY = 0
			p.posY = 0
			p.offX = 0
			p.screen.Sync()
			p.Draw()
		case '>':
			_, y := p.screen.Size()
			p.offY = p.linenum - y + 1
			p.posY = p.linenum -1
			p.offX = 0
			if p.offY < 0 {
				p.offY = 0
			}
			p.screen.Sync()
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
		ev := p.screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
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
	bodyLNum := p.bodyLineNum()
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
	withRefresh := false
	p.posX += 1
	width := p.termWidth()
	if p.posX - p.offX >= width {
		withRefresh = true
		p.offX += width/2
	}
	p.DrawInternal(withRefresh)
}

func (p *Pager) scrollLeft() {
	withRefresh := false
	p.posX = max(0, p.posX -1)

	if p.posX < p.offX {
		withRefresh = true
		p.offX = max(0, p.offX - p.termWidth()/2)
	}
	p.DrawInternal(withRefresh)
}

func (p *Pager) Init() {
	s, e := tcell.NewScreen()
	if e != nil {
		panic(e)
	}
	p.screen = s
	e = p.screen.Init()
	if e != nil {
		panic(e)
	}
}

func (p *Pager) Close() {
	p.screen.Fini()
}


func (p *Pager) onEnter() int {
	if p.linenum ==0 {
		return NO_ACTION	
	}

	// log.Fatal("len %d, pos %d, linnum %d", len(p.lines), p.posY, p.linenum, p.lines[p.posY])	
	p.lineSelected = p.lines[p.posY]
	// log.Fatal("len %d, pos %d, linnum %d", len(p.lines), p.posY, p.linenum, p.lineSelected)	
	// p.lineSelected = "Test!:fuafuga"
	return QUIT
}





