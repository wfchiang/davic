package main

import (
	"wfchiang/davic"
)
 
func main () {
	var value davic.Expr = davic.Expr{Operator:davic.OPT_WEBCALL}
	value.Eval()
}