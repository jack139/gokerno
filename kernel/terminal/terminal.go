package terminal
 
import "unsafe"

type vidmem [900][1440][4]byte

/* 通过地址获取内存指针 */
func get_vidMem(addr uint64) *vidmem {
	buff := (*vidmem)(unsafe.Pointer(uintptr(addr)))
	return buff
}

/* 获取屏幕映射缓存地址 */
func Init() *vidmem{
	vidMem := get_vidMem(0xFFFF800000A00000)
	return vidMem
}
