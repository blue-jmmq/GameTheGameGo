package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gdamore/tcell"
)

//Integer type
type Integer int

//Card structure
type Card struct {
	Name       string
	Cost       Integer
	Health     Integer
	RedDamage  Integer
	BlueDamage Integer
	RedArmor   Integer
	BlueArmor  Integer
}

//WarriorCard variable
var WarriorCard Card = Card{Name: "Warrior", Cost: 1, RedDamage: 1, BlueDamage: 1, RedArmor: 1, BlueArmor: 1}

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

//PrettyPrint function
func PrettyPrint(structure interface{}) {
	fmt.Println(StructToJSONPretty(structure))
}

//PrettyLog function
func PrettyLog(structure interface{}) {
	fo, err := os.OpenFile("output.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	// close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()
	if _, err := fo.WriteString(StructToJSONPretty(structure)); err != nil {
		panic(err)
	}
	if _, err := fo.WriteString("\n"); err != nil {
		panic(err)
	}
}

//BufferWidget struct
type BufferWidget struct {
	Lines       [][]rune
	VisualArray [][]rune
	Width       Integer
	Height      Integer
	Index       Integer
}

//NewBufferWidget function creates a new BufferWidget
func NewBufferWidget() *BufferWidget {
	var bufferWidget BufferWidget
	bufferWidget.Lines = make([][]rune, 0)
	return &bufferWidget
}

//UpdateSize function
func (bufferWidget *BufferWidget) UpdateSize(width, height Integer) {
	bufferWidget.Width = width
	bufferWidget.Height = height
}

//ScrollDown function
func (bufferWidget *BufferWidget) ScrollDown() {
	lenght := bufferWidget.GetFullVisualArraySize()
	if bufferWidget.Index < lenght {
		bufferWidget.Index++
	}
}

//ScrollUp function
func (bufferWidget *BufferWidget) ScrollUp() {
	if bufferWidget.Index > 0 {
		bufferWidget.Index--
	}
}

//GetSize function
func (bufferWidget *BufferWidget) GetSize() (Integer, Integer) {
	return bufferWidget.Width, bufferWidget.Height
}

//Update function
func (bufferWidget *BufferWidget) Update(width, height Integer) {
	bufferWidget.UpdateSize(width, height)
	bufferWidget.VisualArray = bufferWidget.GetVisualArray()
}

//GetFullVisualArraySize function
func (bufferWidget *BufferWidget) GetFullVisualArraySize() Integer {
	fullArray := bufferWidget.GetFullVisualArray()
	return Integer(len(fullArray))
}

//GetSeparator function
func (bufferWidget *BufferWidget) GetSeparator() []rune {
	var separator []rune
	var index Integer
	for index = 0; index < bufferWidget.Width; index++ {
		if index == 0 {
			separator = append(separator, tcell.RuneDArrow)
		} else if index == bufferWidget.Width-1 {
			separator = append(separator, tcell.RuneDArrow)
		} else {
			separator = append(separator, '-')
		}
	}
	return separator
}

//GetFullVisualArray function
func (bufferWidget *BufferWidget) GetFullVisualArray() [][]rune {
	var array [][]rune
	var separator []rune = bufferWidget.GetSeparator()
	array = append(array, separator)
	var bufferEntries [][][]rune
	lines := bufferWidget.Lines
	for _, line := range lines {
		//PrettyLog(string(line))
		splittedLinesString := strings.Split(string(line), "\n")
		//PrettyLog(splittedLinesString)

		var splittedLines [][]rune
		for _, splittedLineString := range splittedLinesString {
			splittedLines = append(splittedLines, []rune(splittedLineString))
		}
		bufferEntries = append(bufferEntries, splittedLines)
	}
	//PrettyLog(bufferEntries)

	for _, bufferEntry := range bufferEntries {
		for _, line := range bufferEntry {
			drawableLines := bufferWidget.GetDrawableLines(line)
			for _, drawableLine := range drawableLines {
				array = append(array, drawableLine)
			}
			//array = append(array, []rune{' ', ' ', ' ', tcell.RuneRArrow})
		}
		array = append(array, separator)
	}
	//PrettyLog(array)
	return array
}

//GetVisualArray function
func (bufferWidget *BufferWidget) GetVisualArray() [][]rune {
	fullArray := bufferWidget.GetFullVisualArray()
	topIndex := bufferWidget.Index + bufferWidget.Height
	if Integer(len(fullArray)) < topIndex {
		return fullArray[bufferWidget.Index:]
	}
	return fullArray[bufferWidget.Index:topIndex]
}

//GetDrawableLines function
func (bufferWidget *BufferWidget) GetDrawableLines(line []rune) [][]rune {
	var drawableLines [][]rune
	preLine := []rune{' ', ' ', ' ', ' '}
	preLineSize := Integer(len(preLine))
	if Integer(len(line)) <= bufferWidget.Width-preLineSize {
		var drawableLine []rune
		drawableLine = append(drawableLine, preLine...)
		drawableLine = append(drawableLine, line...)
		//PrettyLog(dra)
		drawableLines = append(drawableLines, drawableLine)
	} else {
		var index Integer = 0
	a:
		for {
			var drawableLine []rune
			var dlIndex Integer
			var initialIndex Integer
			drawableLine = append(drawableLine, preLine...)
			initialIndex = preLineSize
		b:
			for dlIndex = initialIndex; dlIndex < bufferWidget.Width; dlIndex++ {
				drawableLine = append(drawableLine, line[index])
				index++
				if index >= Integer(len(line)) {
					break b
				}
			}
			drawableLines = append(drawableLines, drawableLine)
			if index >= Integer(len(line)) {
				break a
			}
		}
	}
	return drawableLines
}

//InputWidget struct
type InputWidget struct {
	Line            []rune
	VisualArray     [][]rune
	Width           Integer
	Height          Integer
	Index           Integer
	shouldDrawIndex bool
}

//NewInputWidget function creates a new InputWidget
func NewInputWidget() *InputWidget {
	var inputWidget InputWidget
	inputWidget.Line = make([]rune, 0)
	return &inputWidget
}

//ScrollLeft function
func (inputWidget *InputWidget) ScrollLeft() {
	if inputWidget.Index > 0 {
		inputWidget.Index--
	}
}

//ScrollRight function
func (inputWidget *InputWidget) ScrollRight() {
	if inputWidget.Index < Integer(len(inputWidget.Line)) {
		inputWidget.Index++
	}
}

//Tick function
func (inputWidget *InputWidget) Tick() {
	inputWidget.shouldDrawIndex = !inputWidget.shouldDrawIndex
	//PrettyLog(inputWidget.shouldDrawIndex)
}

//Typed function
func (inputWidget *InputWidget) Typed(r rune) {
	//PrettyLog(fmt.Sprint("Typed rune: ", r, ", ", string(r)))
	inputWidget.Line = append(inputWidget.Line, rune(0))
	copy(inputWidget.Line[inputWidget.Index+1:], inputWidget.Line[inputWidget.Index:])
	inputWidget.Line[inputWidget.Index] = r
	inputWidget.Index++
}

//UpdateSize function
func (inputWidget *InputWidget) UpdateSize(width, height Integer) {
	inputWidget.Width = width
	inputWidget.Height = height
}

//DeleteRune function
func (inputWidget *InputWidget) DeleteRune() {
	if inputWidget.Index > 0 && Integer(len(inputWidget.Line)) > 0 {
		inputWidget.Line = append(inputWidget.Line[:inputWidget.Index-1], inputWidget.Line[inputWidget.Index:]...)
		inputWidget.Index--
	}
}

//GetSize function
func (inputWidget *InputWidget) GetSize() (Integer, Integer) {
	return inputWidget.Width, inputWidget.Height
}

//Update function
func (inputWidget *InputWidget) Update(width, height Integer) {
	inputWidget.UpdateSize(width, height)
	inputWidget.VisualArray = inputWidget.GetVisualArray()
}

//IsLineSmall function
func (inputWidget *InputWidget) IsLineSmall(nRunesAvailable Integer) bool {
	return Integer(len(inputWidget.Line)) < nRunesAvailable
}

//IsLinePerfect function
func (inputWidget *InputWidget) IsLinePerfect(nRunesAvailable Integer) bool {
	return Integer(len(inputWidget.Line)) == nRunesAvailable
}

//LineIsPerfect function
func (inputWidget *InputWidget) LineIsPerfect(inputRow *[]rune) {
	var index Integer
	if inputWidget.Index == Integer(len(inputWidget.Line)) {
		for index = 1; index < Integer(len(inputWidget.Line)); index++ {
			*inputRow = append(*inputRow, inputWidget.Line[index])
		}
		if inputWidget.shouldDrawIndex {
			*inputRow = append(*inputRow, '█')
		}
	} else {
		inputWidget.LineIsSmall(inputRow)
	}
}

//LineIsSmall function
func (inputWidget *InputWidget) LineIsSmall(inputRow *[]rune) {
	var index Integer
	for index = 0; index < Integer(len(inputWidget.Line)); index++ {
		if inputWidget.shouldDrawIndex && index == inputWidget.Index {
			*inputRow = append(*inputRow, '█')
		} else {
			*inputRow = append(*inputRow, inputWidget.Line[index])
		}
	}
	if inputWidget.shouldDrawIndex && index == inputWidget.Index {
		*inputRow = append(*inputRow, '█')
	}
}

//LineIsBig function
func (inputWidget *InputWidget) LineIsBig(inputRow *[]rune, nRunesAvailable Integer) {
	var index Integer
	if inputWidget.Index < nRunesAvailable {
		for index = 0; index < nRunesAvailable; index++ {
			if inputWidget.shouldDrawIndex && index == inputWidget.Index {
				*inputRow = append(*inputRow, '█')
			} else {
				*inputRow = append(*inputRow, inputWidget.Line[index])
			}
		}
	} else {
		for index = inputWidget.Index - nRunesAvailable + 1; index < inputWidget.Index; index++ {
			*inputRow = append(*inputRow, inputWidget.Line[index])
		}
		if inputWidget.shouldDrawIndex {
			*inputRow = append(*inputRow, '█')
		} else if index < Integer(len(inputWidget.Line)) {
			*inputRow = append(*inputRow, inputWidget.Line[index])
		}
	}
}

//GetVisualArray function
func (inputWidget *InputWidget) GetVisualArray() [][]rune {
	var array [][]rune
	var emptyRow []rune
	var index Integer
	for index = 0; index < inputWidget.Width; index++ {
		emptyRow = append(emptyRow, '-')
	}
	array = append(array, emptyRow)

	var inputRow []rune
	inputRow = append(inputRow, 'λ', ' ')
	nConstantRunes := Integer(len(inputRow))
	nRunesAvailable := inputWidget.Width - nConstantRunes
	if inputWidget.IsLinePerfect(nRunesAvailable) {
		inputWidget.LineIsPerfect(&inputRow)
	} else {
		if inputWidget.IsLineSmall(nRunesAvailable) {
			inputWidget.LineIsSmall(&inputRow)
		} else {
			inputWidget.LineIsBig(&inputRow, nRunesAvailable)
		}
	}

	//inputRow = append(inputRow, '█')
	array = append(array, inputRow)

	var emptyRow2 []rune
	for index = 0; index < inputWidget.Width; index++ {
		emptyRow2 = append(emptyRow2, '-')
	}
	array = append(array, emptyRow2)

	return array
}

//AppendString function
func (bufferWidget *BufferWidget) AppendString(s string) {
	bufferWidget.Lines = append(bufferWidget.Lines, []rune(s))
}

//UI structure
type UI struct {
	BufferWidget  *BufferWidget
	InputWidget   *InputWidget
	MinimumWidth  Integer
	MinimumHeight Integer
	Width         Integer
	Height        Integer
}

//UpdateSize function
func (ui *UI) UpdateSize(width, height Integer) {
	ui.Width = width
	ui.Height = height
}

//GetSize function
func (ui *UI) GetSize() (Integer, Integer) {
	return ui.Width, ui.Height
}

//Update function
func (ui *UI) Update(width, height Integer) {
	ui.UpdateSize(width, height)
	ui.BufferWidget.Update(ui.Width, ui.Height-3)
	ui.InputWidget.Update(ui.Width, 3)
}

//Internal structure
type Internal struct {
	UI     *UI
	Screen tcell.Screen
	Timer  *time.Timer
}

//GetScreenWidth function
func (internal *Internal) GetScreenWidth() Integer {
	width, _ := internal.GetScreenSize()
	return width
}

//GetScreenHeight function
func (internal *Internal) GetScreenHeight() Integer {
	_, height := internal.GetScreenSize()
	return height
}

//FillScreen function
func (internal *Internal) FillScreen(r rune) {
	width, height := internal.GetScreenSize()
	var y Integer
	for y = 0; y < height; y++ {
		var x Integer
		for x = 0; x < width; x++ {
			internal.Screen.SetContent(int(x), int(y), r, nil, tcell.StyleDefault)
		}
	}
}

//FillScreenFromArray function
func (internal *Internal) FillScreenFromArray(array [][]rune) {
	height := internal.GetScreenHeight()
	var y Integer
	for y = 0; y < height; y++ {
		if array[y] != nil {
			var x Integer
			for x = 0; x < Integer(len(array[y])); x++ {
				internal.Screen.SetContent(int(x), int(y), array[y][x], nil, tcell.StyleDefault)
			}
		}
	}
}

//GetRuneArray function
func (ui *UI) GetRuneArray() [][]rune {
	var array [][]rune
	var y Integer
	for y = 0; y < ui.BufferWidget.Height; y++ {
		if y < Integer(len(ui.BufferWidget.VisualArray)) {
			array = append(array, ui.BufferWidget.VisualArray[y])
		} else {
			array = append(array, nil)
		}
	}
	for y = 0; y < ui.InputWidget.Height; y++ {
		array = append(array, ui.InputWidget.VisualArray[y])
	}
	return array
}

//UpdateScreen function updates the screen
func (internal *Internal) UpdateScreen() {
	width, height := internal.GetScreenSize()
	internal.Screen.Clear()
	if (width < internal.UI.MinimumWidth) || (height < internal.UI.MinimumHeight) {
		internal.FillScreen('X')
	} else {
		internal.UI.Update(internal.GetScreenSize())
		//PrettyLog(internal.UI)
		runeArray := internal.UI.GetRuneArray()
		internal.FillScreenFromArray(runeArray)
	}
	internal.Screen.Sync()
}

//GetScreenSize function
func (internal *Internal) GetScreenSize() (Integer, Integer) {
	width, height := internal.Screen.Size()
	return Integer(width), Integer(height)
}

//Tick function
func (internal *Internal) Tick() {
	internal.Timer.Reset(1000 * time.Millisecond)
	internal.UI.InputWidget.Tick()
	internal.UpdateScreen()
}

//ResetTimer function
func (internal *Internal) ResetTimer() {
	internal.Timer.Reset(1000 * time.Millisecond)
	internal.UI.InputWidget.shouldDrawIndex = true
}

//SetTimer function
func (internal *Internal) SetTimer() {
	internal.Timer = time.NewTimer(1000 * time.Millisecond)
	internal.UI.InputWidget.shouldDrawIndex = true
}

//InternalLoop handles TCell logic
func InternalLoop(ui *UI) {
	screen, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}
	if err = screen.Init(); err != nil {
		panic(err)
	}
	/*screen.SetStyle(tcell.StyleDefault.
	Foreground(tcell.ColorBlack).
	Background(tcell.ColorWhite))
	*/
	var internal Internal
	internal.Screen = screen
	internal.UI = ui
	internal.Screen.Clear()
	internal.UpdateScreen()
	internal.SetTimer()
	go func() {
		for {
			select {
			case <-internal.Timer.C:
				//PrettyLog(fmt.Sprint("Tick at", t))
				internal.Tick()
			}
		}
	}()
loop:
	for {
		event := internal.Screen.PollEvent()
		switch event := event.(type) {
		case *tcell.EventKey:
			switch event.Key() {
			case tcell.KeyEscape:
				break loop
			case tcell.KeyUp:
				//PrettyLog("tcell.KeyUp")
				internal.UI.BufferWidget.ScrollUp()
				internal.UpdateScreen()
			case tcell.KeyDown:
				//PrettyLog("tcell.KeyDown")
				internal.UI.BufferWidget.ScrollDown()
				internal.UpdateScreen()
			case tcell.KeyLeft:
				//PrettyLog("tcell.KeyLeft")
				internal.UI.InputWidget.ScrollLeft()
				internal.ResetTimer()
				internal.UpdateScreen()
			case tcell.KeyRight:
				//PrettyLog("tcell.KeyRight")
				internal.UI.InputWidget.ScrollRight()
				internal.ResetTimer()
				internal.UpdateScreen()
			case tcell.KeyRune:
				internal.UI.InputWidget.Typed(event.Rune())
				internal.ResetTimer()
				internal.UpdateScreen()
			case tcell.KeyBackspace, tcell.KeyBackspace2:
				//PrettyLog("tcell.KeyBackspace")
				internal.UI.InputWidget.DeleteRune()
				internal.ResetTimer()
				internal.UpdateScreen()
			}
			//PrettyLog(event.Key())
			//PrettyLog(event.Name())
			//fmt.Println(width, height)
		case *tcell.EventResize:
			internal.UpdateScreen()
		}
	}
	screen.Fini()
}

//RuneLinesToStringLines function
func RuneLinesToStringLines(runeLines [][]rune) []string {
	var stringLines []string
	for index := 0; index < len(runeLines); index++ {
		stringLines = append(stringLines, string(runeLines[index]))
	}
	return stringLines
}

//NewUI function
func NewUI(bufferWidget *BufferWidget, inputWidget *InputWidget, minimumWidth, minimunHeight Integer) *UI {
	var ui UI
	ui.BufferWidget = bufferWidget
	ui.InputWidget = inputWidget
	ui.MinimumWidth = minimumWidth
	ui.MinimumHeight = minimunHeight
	return &ui
}

func main() {
	PrettyLog(WarriorCard)
	bufferWidget := NewBufferWidget()
	bufferWidget.AppendString("Hola mundo me llamo José Manuel Martínez Quevedo")
	bufferWidget.AppendString("Hello world my name is PepeThePepe and this is GameTheGame")
	bufferWidget.AppendString(StructToJSONPretty(WarriorCard))
	inputWidget := NewInputWidget()
	ui := NewUI(bufferWidget, inputWidget, 16, 8)
	//PrintStructPretty(bufferWidget)
	InternalLoop(ui)
}
