package colors_package

import "github.com/fatih/color"

var (
	ColorBlue = color.New(color.FgBlue)
	ColorRed  = color.New(color.FgHiRed)
)

// custom FPrintf functions
var (
	WriteBlue = ColorBlue.FprintFunc()
	WriteRed  = ColorRed.FprintFunc()
)
