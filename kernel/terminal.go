package terminal
 
import "unsafe"
 
/*
 * Map the text mode video memory into a multi-dimensional array that can be safely
 * used from Go.
 */


func get_vidMem(addr uint64) *[900][1440][4]byte {
	buff := (*[900][1440][4]byte)(unsafe.Pointer(uintptr(addr)))
	return buff
}

var vidMem *[900][1440][4]byte

func Init() {
	vidMem = get_vidMem(0xFFFF800000A00000)
}

func Print() {
	for x := 0; x < 1440; x++ {
		for y := 0; y< 20; y++ {
			(*vidMem)[y][x][0] = 0x00
			(*vidMem)[y][x][1] = 0x00
			(*vidMem)[y][x][2] = 0xff
			(*vidMem)[y][x][3] = 0x00
		}

		for y := 20; y< 40; y++ {
			(*vidMem)[y][x][0] = 0x00
			(*vidMem)[y][x][1] = 0xff
			(*vidMem)[y][x][2] = 0x00
			(*vidMem)[y][x][3] = 0x00
		}

		for y := 40; y< 60; y++ {
			(*vidMem)[y][x][0] = 0xff
			(*vidMem)[y][x][1] = 0x00
			(*vidMem)[y][x][2] = 0x00
			(*vidMem)[y][x][3] = 0x00
		}

		for y := 60; y< 80; y++ {
			(*vidMem)[y][x][0] = 0xff
			(*vidMem)[y][x][1] = 0xff
			(*vidMem)[y][x][2] = 0xff
			(*vidMem)[y][x][3] = 0x00
		}
	}
}