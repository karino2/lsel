package main

import  (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)


func getpath(line string) string{
	arr := strings.Split(line, ":")
	return arr[0]
}


func main() {
	var p Pager
	p.Init()
	c, _ := ioutil.ReadAll(os.Stdin)
	p.SetContent(string(c))
	p.PollEvent()
	p.Close()


	if p.lineSelected != "" {
		path := getpath(p.lineSelected)
		fmt.Println(path)
		os.Exit(0)
	}
	os.Exit(1)
}
