package terminal

/* 打印三个彩色条 */
func Print_test() {
	vidMem := pos.FB_addr

	/* 输出红色彩条 */
	for x := 0; x < 1440; x++ {
		for y := 0; y< 20; y++ {
			(*vidMem)[y][x] = RED
		}
	}

	putchar(200, 200, WHITE, BLACK, byte('A')) // 'A'
}

func Printk_test() {
	buf[0] = 'H'
	buf[1] = 'e'
	buf[2] = '\t'
	buf[3] = 'l'
	buf[4] = 'o'
	buf[5] = '\n'
	buf[6] = 'A'
	buf[7] = 'B'
	buf[8] = '\b'
	buf[9] = 'C'

	color_printk(YELLOW, BLACK, 10)
}