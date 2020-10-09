package main
 
import (
	"terminal"
	//"kfmt"
)


/*
	假入口，只为编译时不报错
*/
func main() {}


/*
	kernel的go语言入口
*/
func StartKerno() {
	
	terminal.InitTerminal()

	terminal.Print_test()
	terminal.Printk_test()

	for { }
}
