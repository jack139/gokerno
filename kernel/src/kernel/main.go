package main
 
import (
	"kernel/terminal"
	//"kernel/kfmt"
)


/*
	假入口，只为编译时不报错
*/
func main() {
	StartKerno()
}

/*
	kernel的go语言入口
*/
//go:noinline
//go:nosplit
func StartKerno() {

	terminal.InitTerminal()

	terminal.Print_test()
	//terminal.Printk_test()

	//kfmt.Printf_int("\n\n\nHello:\t%08x\n", 123)
	//kfmt.Printf_str("\nHello %s\n", "world")

	for { }
}
