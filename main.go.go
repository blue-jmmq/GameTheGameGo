package main

import (
	"encoding/json"
	"fmt"

	"github.com/gdamore/tcell"
)

//Weapon structure
type Weapon struct {
	Hands int64
	Power int64
}

//Armor structure
type Armor struct {
	Head int64
	Body int64
	Legs int64
}

//Shield structure
type Shield struct {
	Power int64
}

//Figther structure
type Figther struct {
	Name   string
	ID     int64
	Health int64
	Armor  Armor
	Weapon Weapon
	Shield Shield
}

//StructToJSON function
func StructToJSON(structure interface{}) string {
	JSONBytes, _ := json.Marshal(structure)
	return string(JSONBytes)
}

//StructToJSONPretty function
func StructToJSONPretty(structure interface{}) string {
	JSONBytes, _ := json.MarshalIndent(structure, "", "    ")
	return string(JSONBytes)
}

//PrintStruct function
func PrintStruct(structure interface{}) {
	fmt.Println(StructToJSON(structure))
}

//PrintStructPretty function
func PrintStructPretty(structure interface{}) {
	fmt.Println(StructToJSONPretty(structure))
}

func main() {
	screen, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}
	if err = screen.Init(); err != nil {
		panic(err)
	}
	screen.SetStyle(tcell.StyleDefault.
		Foreground(tcell.ColorBlack).
		Background(tcell.ColorWhite))
	screen.Clear()
loop:
	for {
		event := screen.PollEvent()
		switch event := event.(type) {
		case *tcell.EventKey:
			switch event.Key() {
			case tcell.KeyEscape, tcell.KeyEnter:
				break loop
			}
			width, height := screen.Size()
			for y := 0; y < height; y++ {
				for x := 0; x < width; x++ {
					screen.SetContent(x, y, event.Rune(), nil, tcell.StyleDefault)
				}
			}
			screen.Sync()
			//fmt.Println(width, height)
		case *tcell.EventResize:
			screen.Sync()
		}
	}
	screen.Fini()

	fmt.Println("Welcome to 'Battle'")
	var fighter Figther = Figther{Health: 512}
	PrintStructPretty(fighter)
}
