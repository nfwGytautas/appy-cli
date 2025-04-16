package utils

var Verbose = false
var Console = console{
	dbg: &debugWriter{},
	err: &errorWriter{},
}
