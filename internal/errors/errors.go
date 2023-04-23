package errors

import "fmt"

var WrongDataFormat = fmt.Errorf("wrong data format received during data processing")
var EmptyRectangles = fmt.Errorf("no rectangles were given")
