package main

import (
	"fmt"
	"wfchiang/davic"
)
 
func main () {
	fmt.Println("==== Test Expr.Eval() ====")
	var value davic.Expr = davic.Expr{Operator:davic.OPT_WEBCALL}
	value.Eval()

	fmt.Println("==== Test JNode ====")
	var json_bytes = []byte("{\"abc\":123, \"xyz\":{\"111\":\"222\"}}")
	var jnode davic.JNode = davic.JNode{}
	jnode.InterpretBytes(json_bytes)
	fmt.Println(jnode.GetKeys())
	fmt.Println(jnode.GetValue([]string{"xyz", "111"}))
}