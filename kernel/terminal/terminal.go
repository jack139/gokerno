package terminal
 
import "unsafe"

/*
	屏幕信息
*/
const (
	X_CHAR_SIZE = 8
	Y_CHAR_SIZE = 16

	X_RESOLUTION = 1440
	Y_RESOLUTION = 900
)

/* 
	屏幕缓存地址类型 
*/
type vidmem [Y_RESOLUTION][X_RESOLUTION]uint32

/* 
	屏幕描述，记录当前输出位置 
*/
type Position struct {
	XResolution int
	YResolution int

	XPosition int
	YPosition int

	XCharSize int
	YCharSize int

	FB_addr *vidmem
	FB_length uint64
}

var (
	pos Position  // 屏幕位置信息
	buf [4096]byte  // pintk的缓存
)

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
func InitTerminal() *Position{
	pos.XResolution = X_RESOLUTION
	pos.YResolution = Y_RESOLUTION
	pos.XPosition = 0
	pos.YPosition = 0
	pos.XCharSize = X_CHAR_SIZE
	pos.YCharSize = Y_CHAR_SIZE
	pos.FB_addr = get_vidMem(0xFFFF800000A00000)
	// pos.FB_length = (XResolution * YResolution * 4 + PAGE_4K_SIZE - 1) & PAGE_4K_MASK

	return &pos
}

/*
	屏幕上指定位置输出一个字符
*/
func putchar(x int, y int, colorFR uint32, colorBK uint32, font byte) {
	var i, j int
	var testval byte
	var fontp *[16]byte

	fb := pos.FB_addr

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


/*
	打印缓存中的length个字符
*/
func color_printk(colorFR uint32, colorBK uint32, length int) int{
	var count, line, i int

	line = 0
	i = length

	for count=0; count<i || line>0; count++ {
		////	add \n \b \t
		if line>0 {
			count--

			line--
			putchar(pos.XPosition * pos.XCharSize, pos.YPosition * pos.YCharSize, colorFR, colorBK, ' ')
			pos.XPosition++

			goto Label_tab
		}

		if buf[count] == byte('\n') {
			pos.YPosition++
			pos.XPosition = 0
		} else if buf[count] == byte('\b'){
			pos.XPosition--
			if pos.XPosition < 0 {
				pos.XPosition = (pos.XResolution / pos.XCharSize - 1) * pos.XCharSize
				pos.YPosition--
				if pos.YPosition < 0 {
					pos.YPosition = (pos.YResolution / pos.YCharSize - 1) * pos.YCharSize
				}
			}	
			putchar(pos.XPosition * pos.XCharSize, pos.YPosition * pos.YCharSize, colorFR, colorBK, ' ')
		} else if buf[count] == byte('\t'){
			line = ((pos.XPosition + 8) & ^(8 - 1)) - pos.XPosition

			line--
			putchar(pos.XPosition * pos.XCharSize, pos.YPosition * pos.YCharSize, colorFR, colorBK, ' ')
			pos.XPosition++
		} else {
			putchar(pos.XPosition * pos.XCharSize, pos.YPosition * pos.YCharSize, colorFR, colorBK, buf[count])			
			pos.XPosition++
		}

Label_tab:

		if pos.XPosition >= (pos.XResolution / pos.XCharSize) {
			pos.YPosition++
			pos.XPosition = 0
		}
		if pos.YPosition >= (pos.YResolution / pos.YCharSize) {
			pos.YPosition = 0
		}

	}
	return i
}
