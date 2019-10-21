package main

import (
	"encoding/json"
	"fmt"
	"github.com/gdamore/tcell"
	"os"
	"strings"
	"sync"
	"strconv"
	"time"
	"net"
)

const (
	ServerApplication Integer = iota
	ClientApplication
)

type Server struct {
}

//Integer type
type Integer int

//StatisticsPerLevel structure
type StatisticsPerLevel struct {
	CostPerLevel      Integer
	HealthPerLevel    Integer
	RedArmorPerLever  Integer
	BlueArmorPerLevel Integer
}

//Card structure
type Card struct {
	Name               string
	Cost               Integer
	Health             Integer
	RedArmor           Integer
	BlueArmor          Integer
	RedDamage          Integer
	BlueDamage         Integer
	AntiAttackSpeed    Integer
	Level              Integer
	Experience         Integer
	StatisticsPerLevel *StatisticsPerLevel
}

//NewWarriorCard creates a Warrior Card
func NewWarriorCard() *Card {
	return &Card{
		Name:            "Warrior",
		Cost:            1,
		RedDamage:       1,
		BlueDamage:      1,
		RedArmor:        2,
		BlueArmor:       1,
		AntiAttackSpeed: 4,
		Level:           1,
		Experience:      0,
		StatisticsPerLevel: &StatisticsPerLevel{
			CostPerLevel:      1,
			HealthPerLevel:    1,
			RedArmorPerLever:  2,
			BlueArmorPerLevel: 1,
		},
	}
}

//Player structure
type Player struct {
	Health   Integer
	RedArmor Integer
	Energy   Integer
	Credit   Integer
	Deck     []Card
	Hand     []Card
	Board    [][]Card
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

//UpdateIndex function
func (bufferWidget *BufferWidget) UpdateIndex() {
	visualSize := bufferWidget.GetFullVisualArraySize()
	if visualSize >= bufferWidget.Height {
		bufferWidget.Index = visualSize - bufferWidget.Height
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
			separator = append(separator, tcell.RuneHLine)
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

//Reset function
func (inputWidget *InputWidget) Reset() {
	inputWidget.Line = nil
	inputWidget.Index = 0
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
		emptyRow = append(emptyRow, tcell.RuneHLine)
	}
	array = append(array, emptyRow)

	var inputRow []rune
	inputRow = append(inputRow, ' ', 'λ', ' ')
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
		emptyRow2 = append(emptyRow2, tcell.RuneHLine)
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

//AppManager structure
type AppManager struct {
	LocalIPs []net.IP
	ServerIP net.IP
	ServerPort Integer
	ServerAddress string
	Type           Integer
	UI             *UI
	Screen         tcell.Screen
	ScreenMutex    sync.Mutex
	Timer          *time.Timer
	CommandChannel chan []rune
}

//GetScreenWidth function
func (appManager *AppManager) GetScreenWidth() Integer {
	width, _ := appManager.GetScreenSize()
	return width
}

//GetScreenHeight function
func (appManager *AppManager) GetScreenHeight() Integer {
	_, height := appManager.GetScreenSize()
	return height
}

//FillScreen function
func (appManager *AppManager) FillScreen(r rune) {
	width, height := appManager.GetScreenSize()
	var y Integer
	for y = 0; y < height; y++ {
		var x Integer
		for x = 0; x < width; x++ {
			appManager.Screen.SetContent(int(x), int(y), r, nil, tcell.StyleDefault)
		}
	}
}

//FillScreenFromArray function
func (appManager *AppManager) FillScreenFromArray(array [][]rune) {
	height := appManager.GetScreenHeight()
	var y Integer
	for y = 0; y < height; y++ {
		if array[y] != nil {
			var x Integer
			for x = 0; x < Integer(len(array[y])); x++ {
				appManager.Screen.SetContent(int(x), int(y), array[y][x], nil, tcell.StyleDefault)
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
func (appManager *AppManager) UpdateScreen() {
	appManager.ScreenMutex.Lock()
	width, height := appManager.GetScreenSize()
	appManager.Screen.Clear()
	if (width < appManager.UI.MinimumWidth) || (height < appManager.UI.MinimumHeight) {
		appManager.FillScreen('X')
	} else {
		appManager.UI.Update(appManager.GetScreenSize())
		//PrettyLog(internal.UI)
		runeArray := appManager.UI.GetRuneArray()
		appManager.FillScreenFromArray(runeArray)
	}
	appManager.Screen.Sync()
	appManager.ScreenMutex.Unlock()
}

//GetScreenSize function
func (appManager *AppManager) GetScreenSize() (Integer, Integer) {
	width, height := appManager.Screen.Size()
	return Integer(width), Integer(height)
}

//Tick function
func (appManager *AppManager) Tick() {
	appManager.Timer.Reset(256 * time.Millisecond)
	appManager.UI.InputWidget.Tick()
	appManager.UpdateScreen()
}

//ResetTimer function
func (appManager *AppManager) ResetTimer() {
	appManager.Timer.Reset(256 * time.Millisecond)
	appManager.UI.InputWidget.shouldDrawIndex = true
}

//SetTimer function
func (appManager *AppManager) SetTimer() {
	appManager.Timer = time.NewTimer(256 * time.Millisecond)
	appManager.UI.InputWidget.shouldDrawIndex = true
}

//AskServerOrClient function
func (appManager *AppManager) AskServerOrClient() {
	appManager.WriteEntry("Which network role dou you want to become(server/client)?")
	appManager.UpdateScreen()
a:
	for {
		select {
		case command := <-appManager.CommandChannel:
			//PrettyLog(fmt.Sprint("Tick at", t))
			if strings.EqualFold(string(command), "server") {
				appManager.WriteEntryAndUpdate("Ok, now you are a server")
				appManager.Type = ServerApplication
				break a
			} else if strings.EqualFold(string(command), "client") {
				appManager.WriteEntryAndUpdate("Ok, now you are a client")
				appManager.Type = ClientApplication
				break a
			} else {
				appManager.WriteEntryAndUpdate("Incorrect answer(" + string(command) + "), choose between: server/client, and try again")
			}
		}
	}
}

//FindLocalIPs function
func (appManager *AppManager) FindLocalIPs() {
	ips := make([]net.IP, 0)
	ifaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}
	for _, i := range ifaces {
		//appManager.WriteEntryAndUpdate(StructToJSONPretty(i))
		addrs, err := i.Addrs()
		if err != nil {
			panic(err)
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
					ip = v.IP
					ips = append(ips, ip)
			case *net.IPAddr:
					ip = v.IP
					ips = append(ips, ip)
			}
		}
	}
	appManager.LocalIPs = ips
}

//AskForIP function
func (appManager *AppManager) AskForIP() {
	appManager.WriteEntryAndUpdate("Choose an address(type an index number):")
	var index Integer
	ips := appManager.LocalIPs
	for index = 0; index < Integer(len(ips)); index++ {
		appManager.WriteEntry(strconv.Itoa(int(index + 1)) + " " + string(tcell.RuneRArrow) + " " + ips[index].String())
	}
	appManager.UpdateScreen()
a:
	for {
		select {
		case command := <-appManager.CommandChannel:
			//PrettyLog(fmt.Sprint("Tick at", t))
			selectionInt, err := strconv.Atoi(string(command))
			selection := Integer(selectionInt)
			if err != nil {
				appManager.WriteEntryAndUpdate("Type an index number")
			} else {
				if selection <= 0 {
					appManager.WriteEntryAndUpdate("The index must be greater than 0")
				} else if selection > Integer(len(ips)) {
					appManager.WriteEntryAndUpdate("The index must be smaller than " + strconv.Itoa(len(ips) + 1))
				} else {
					appManager.ServerIP = ips[selection - 1]
					appManager.WriteEntryAndUpdate("Great, you have choosen address " + 
						appManager.ServerIP.String() + 
						" at index " + 
						strconv.Itoa(int(selection)))
					break a
				}
			}
			
		}
	}
}

//ListenForConnection function
func (appManager *AppManager) ListenForConnection() {
	var portInteger Integer
a:
	for portInteger = 2048; portInteger <= 32000; portInteger++ {
		//timeout := time.Second
		addressString := net.JoinHostPort(appManager.ServerIP.String(), strconv.Itoa(int(portInteger)))
		appManager.WriteEntry(addressString)
		listener, err := net.Listen("tcp", addressString)
		if err != nil {
			// handle error
			appManager.WriteEntry(fmt.Sprint("Server error: ", err))
		} else {
			appManager.WriteEntry("Server at " + addressString)
			connection, err := listener.Accept()
			if err != nil {
				appManager.WriteEntry(fmt.Sprint("tcp server accept error: ", err))
			} else {
				appManager.WriteEntry("Connection succesful")
				appManager.WriteEntry(StructToJSONPretty(connection))
				break a
			}
		}
	}
	appManager.UpdateScreen()
}

func (appManager *AppManager) DialServer() {
	appManager.WriteEntryAndUpdate("Connecting to server address: " + appManager.ServerAddress)
	connection, err := net.Dial("tcp", appManager.ServerAddress)
	if err != nil {
		appManager.WriteEntryAndUpdate(fmt.Sprint("Connection error: ", err))
	} else {
		appManager.WriteEntry("Connection succesful")
		appManager.WriteEntryAndUpdate(StructToJSONPretty(connection))

	}
}

func (appManager *AppManager) AskServerPort() {
	appManager.WriteEntryAndUpdate("Enter the server port (integer between 2048 and 32000 inclusive):")
a:
	for {
		select {
		case command := <-appManager.CommandChannel:
			//PrettyLog(fmt.Sprint("Tick at", t))
			portString := string(command)
			portInt, err := strconv.Atoi(portString)
			if err != nil {
				appManager.WriteEntryAndUpdate("Server port must be an integer between 2048 and 32000 ")
			} else {
				appManager.ServerPort = Integer(portInt)
				if appManager.ServerPort < 2048 || appManager.ServerPort > 32000 {
					appManager.WriteEntryAndUpdate("Server port must be an integer between 2048 and 32000 ")
				} else {
					appManager.ServerAddress = net.JoinHostPort(appManager.ServerIP.String(), portString)
					appManager.DialServer()
					break a
				}
			}
		}
	}
}

func (appManager *AppManager) ConnectToServer() {
	appManager.WriteEntryAndUpdate("Enter a server IP address:")
a:
	for {
		select {
		case command := <-appManager.CommandChannel:
			//PrettyLog(fmt.Sprint("Tick at", t))
			ipString := string(command)
			appManager.ServerIP = net.ParseIP(ipString)
			if appManager.ServerIP == nil {
				appManager.WriteEntryAndUpdate("Invalid server IP address, try again")
			} else {
				appManager.AskServerPort()
				break a
			}
		}
	}
}

//LogicLoop function
func (appManager *AppManager) LogicLoop() {
	appManager.AskServerOrClient()
	//appManager.FindLocalAddress()
	if appManager.Type == ServerApplication {
		appManager.FindLocalIPs()
		appManager.AskForIP()
		appManager.ListenForConnection()
	} else {
		appManager.ConnectToServer()
	}
	
	for {
		select {
		case command := <-appManager.CommandChannel:
			//PrettyLog(fmt.Sprint("Tick at", t))
			appManager.WriteEntryAndUpdate("Command received: " + string(command))
		}
	}
}

//SendCommand function
func (appManager *AppManager) SendCommand(command []rune) {
	if len(command) > 0 {
		commandCopy := make([]rune, len(command))
		copy(commandCopy, command)
		appManager.CommandChannel <- commandCopy
	}

}

//SendCommandFromInput function
func (appManager *AppManager) SendCommandFromInput() {
	appManager.SendCommand(appManager.UI.InputWidget.Line)
	appManager.UI.InputWidget.Reset()
}

//WriteEntryAndUpdate function
func (appManager *AppManager) WriteEntryAndUpdate(e string) {
	appManager.WriteEntry(e)
	appManager.UpdateScreen()
}

//WriteEntry function
func (appManager *AppManager) WriteEntry(e string) {
	appManager.UI.BufferWidget.AppendString(e)
	appManager.UI.BufferWidget.UpdateIndex()
}

//AppManagerLoop function
func AppManagerLoop(ui *UI) {
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
	var appManager AppManager
	appManager.Screen = screen
	appManager.UI = ui
	appManager.UpdateScreen()
	appManager.SetTimer()
	appManager.CommandChannel = make(chan []rune)
	go func() {
		for {
			select {
			case <-appManager.Timer.C:
				//PrettyLog(fmt.Sprint("Tick at", t))
				appManager.Tick()
			}
		}
	}()
	go appManager.LogicLoop()
loop:
	for {
		event := appManager.Screen.PollEvent()
		switch event := event.(type) {
		case *tcell.EventKey:
			switch event.Key() {
			case tcell.KeyEscape:
				//appManager.SendCommand([]rune("exit"))
				break loop
			case tcell.KeyUp:
				//PrettyLog("tcell.KeyUp")
				appManager.UI.BufferWidget.ScrollUp()
				appManager.UpdateScreen()
			case tcell.KeyDown:
				//PrettyLog("tcell.KeyDown")
				appManager.UI.BufferWidget.ScrollDown()
				appManager.UpdateScreen()
			case tcell.KeyLeft:
				//PrettyLog("tcell.KeyLeft")
				appManager.UI.InputWidget.ScrollLeft()
				appManager.ResetTimer()
				appManager.UpdateScreen()
			case tcell.KeyRight:
				//PrettyLog("tcell.KeyRight")
				appManager.UI.InputWidget.ScrollRight()
				appManager.ResetTimer()
				appManager.UpdateScreen()
			case tcell.KeyRune:
				appManager.UI.InputWidget.Typed(event.Rune())
				appManager.ResetTimer()
				appManager.UpdateScreen()
			case tcell.KeyBackspace, tcell.KeyBackspace2:
				//PrettyLog("tcell.KeyBackspace")
				appManager.UI.InputWidget.DeleteRune()
				appManager.ResetTimer()
				appManager.UpdateScreen()
			case tcell.KeyEnter:
				appManager.ResetTimer()
				//appManager.WriteEntry("tcell.KeyEnter")
				appManager.SendCommandFromInput()
				appManager.UpdateScreen()

			}
			//PrettyLog(event.Key())
			//PrettyLog(event.Name())
			//fmt.Println(width, height)
		case *tcell.EventResize:
			appManager.UpdateScreen()
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
	//PrettyLog(WarriorCard)
	bufferWidget := NewBufferWidget()
	bufferWidget.AppendString("Hola mundo me llamo José Manuel Martínez Quevedo")
	bufferWidget.AppendString("Hello world my name is PepeThePepe and this is GameTheGame")
	bufferWidget.AppendString(StructToJSONPretty(NewWarriorCard()))
	inputWidget := NewInputWidget()
	ui := NewUI(bufferWidget, inputWidget, 16, 8)
	//PrintStructPretty(bufferWidget)
	AppManagerLoop(ui)
}
