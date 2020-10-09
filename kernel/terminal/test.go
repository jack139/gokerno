package terminal
 
/* 打印三个彩色条 */
func Print_test() {
	vidMem := Init()

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