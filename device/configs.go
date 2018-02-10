package device

type config struct {
	VendorID  int
	ProductID int
	Pins      []pin
	Buttons   []button
}

type pin struct {
	Number int
	Multi  bool
}

type button struct {
	Pin   int
	Value int
	Name  string
}

var configs = []config{{
	VendorID:  0x0079,
	ProductID: 0x0011,
	Pins: []pin{
		{Number: 3, Multi: false},
		{Number: 4, Multi: false},
		{Number: 5, Multi: true},
		{Number: 6, Multi: true},
	},
	Buttons: []button{
		{Pin: 5, Value: 32, Name: "a"},
		{Pin: 5, Value: 64, Name: "b"},
		{Pin: 5, Value: 128, Name: "y"},
		{Pin: 5, Value: 16, Name: "x"},
		{Pin: 6, Value: 2, Name: "r"},
		{Pin: 6, Value: 1, Name: "l"},
		{Pin: 6, Value: 32, Name: "start"},
		{Pin: 6, Value: 16, Name: "select"},
		{Pin: 4, Value: 0, Name: "up"},
		{Pin: 4, Value: 255, Name: "down"},
		{Pin: 3, Value: 0, Name: "left"},
		{Pin: 3, Value: 255, Name: "right"},
	},
}}
