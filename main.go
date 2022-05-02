package main

import (
	"fmt"
	"image/color"
	"io/fs"
	"io/ioutil"
	"os"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
)

func init() {
	fontPath := "C:/Windows/Fonts/"

	fontPaths := paths(fontPath)
	for _, path := range fontPaths {
		os.Setenv("FYNE_FONT", fontPath+path.Name())
		//楷体:simkai.ttf
		//黑体:simhei.ttf
		if strings.Contains(path.Name(), "simkai.ttf") {
			os.Setenv("FYNE_FONT", fontPath+path.Name())
			break
		}
	}
	fmt.Println("=============")
}

func paths(path string) []fs.FileInfo {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}
	return files
}

func main() {
	myApp := app.New()
	w := myApp.NewWindow("Text")

	text := canvas.NewText("中国", color.Black)
	text.Alignment = fyne.TextAlignTrailing
	text.TextStyle = fyne.TextStyle{Italic: true}
	w.SetContent(text)

	w.ShowAndRun()
	os.Unsetenv("FYNE_FONT")
}
