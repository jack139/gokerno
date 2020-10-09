package terminal
 
import "unsafe"

type vidmem [900][1440]uint32

/* 
	通过地址获取内存指针 
*/
func get_vidMem(addr uint64) *vidmem {
	buff := (*vidmem)(unsafe.Pointer(uintptr(addr)))
	return buff
}

/*
	获取屏幕映射缓存地址 
*/
func Init() *vidmem{
	vidMem := get_vidMem(0xFFFF800000A00000)
	return vidMem
}

/*
	屏幕上指定位置输出一个字符
*/
func putchar(fb *vidmem, x int, y int, colorFR uint32, colorBK uint32, font byte) {
	var i, j int
	var testval byte
	var fontp *[16]byte

	testval = 0
	fontp = &font_ascii[font]

	for i=0; i<16; i++ {
		testval = 0x80
		for j=0; j<8; j++ {
			if ((*fontp)[i] & testval) != 0 {
				(*fb)[y+i][x+j] = colorFR
			} else {
				(*fb)[y+i][x+j] = colorBK
			}
			testval = testval >> 1
		}
	}
}
