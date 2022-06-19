package main

import (
	regex "goexpr/Regex"
)

func main() {
	restr := "abc?d+|cfuck*"
	re := regex.RE(restr)
	re.Match("abc")
	//true
}
