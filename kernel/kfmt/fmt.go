package kfmt

/*
	因为go的可变参数函数在调用时需要申请内存建立新切片，所以在没有内存管理时会导致内核崩掉，
	是用 interface{} 参数也有同样问题。
	因此，先只实现了 Printf_str() 和 Print_int() 两个固定参数的打印，
	后续实现内存管理后在实现通用的 Printf()

	这里的代码参考 gopher-os
*/

import (
	"terminal"
)

// maxBufSize defines the buffer size for formatting numbers.
const maxBufSize = 32

// 记录color_printk打印缓存区
type PrintkBuf struct {
	buff *(terminal.Buffer)
	length int
}

var (
	numFmtBuf [33]byte
	singleByte [1]byte
	printkBuff PrintkBuf
)

// 目前支持的替换符号:
//
// Strings:
//		%s the uninterpreted bytes of the string or byte slice
//
// Integers:
//              %o base 8
//              %d base 10
//              %x base 16, with lower-case letters for a-f
//
func Printf_int(format string, args uint64) {
	printkBuff.buff = terminal.GetPrintkBuf()
	printkBuff.length = 0
	fprintfInt(&printkBuff, format, args)
	terminal.ColorPrintk(terminal.WHITE, terminal.BLACK, printkBuff.length)
}

func Printf_str(format string, args string) {
	printkBuff.buff = terminal.GetPrintkBuf()
	printkBuff.length = 0
	fprintfStr(&printkBuff, format, args)
	terminal.ColorPrintk(terminal.WHITE, terminal.BLACK, printkBuff.length)
}


// 处理整数，写入color_printk缓存
func fprintfInt(w *PrintkBuf, format string, args uint64) {
	var (
		nextCh                       byte
		blockStart, blockEnd, padLen int
		fmtLen                       = len(format)
	)

	for blockEnd < fmtLen {
		nextCh = format[blockEnd]
		if nextCh != '%' {
			blockEnd++
			continue
		}

		if blockStart < blockEnd {
			// passing format[blockStart:blockEnd] to doWrite triggers a
			// memory allocation so we need to do this one byte at a time.
			for i := blockStart; i < blockEnd; i++ {
				singleByte[0] = format[i]
				doWrite(w, singleByte)
			}
		}

		// Scan til we hit the format character
		padLen = 0
		blockEnd++

	parseFmt:

		for ; blockEnd < fmtLen; blockEnd++ {
			nextCh = format[blockEnd]
			switch {
			case nextCh == '%':
				singleByte[0] = '%'
				doWrite(w, singleByte)
				break parseFmt
			case nextCh >= '0' && nextCh <= '9':
				padLen = (padLen * 10) + int(nextCh-'0')
				continue
			case nextCh == 'd' || nextCh == 'x' || nextCh == 'o':
				switch nextCh {
				case 'o':
					fmtInt(w, args, 8, padLen)
				case 'd':
					fmtInt(w, args, 10, padLen)
				case 'x':
					fmtInt(w, args, 16, padLen)
				}

				break parseFmt
			}

			// reached end of formatting string without finding a verb
			singleByte[0] = '?'
			doWrite(w, singleByte)
		}
		blockStart, blockEnd = blockEnd+1, blockEnd+1
	}

	if blockStart != blockEnd {
		// passing format[blockStart:blockEnd] to doWrite triggers a
		// memory allocation so we need to do this one byte at a time.
		for i := blockStart; i < blockEnd; i++ {
			singleByte[0] = format[i]
			doWrite(w, singleByte)
		}
	}

}

// 处理string，写入color_printk缓存
func fprintfStr(w *PrintkBuf, format string, args string) {
	var (
		nextCh                       byte
		blockStart, blockEnd, padLen int
		fmtLen                       = len(format)
	)

	for blockEnd < fmtLen {
		nextCh = format[blockEnd]
		if nextCh != '%' {
			blockEnd++
			continue
		}

		if blockStart < blockEnd {
			// passing format[blockStart:blockEnd] to doWrite triggers a
			// memory allocation so we need to do this one byte at a time.
			for i := blockStart; i < blockEnd; i++ {
				singleByte[0] = format[i]
				doWrite(w, singleByte)
			}
		}

		// Scan til we hit the format character
		padLen = 0
		blockEnd++

	parseFmt:

		for ; blockEnd < fmtLen; blockEnd++ {
			nextCh = format[blockEnd]
			switch {
			case nextCh == '%':
				singleByte[0] = '%'
				doWrite(w, singleByte)
				break parseFmt
			case nextCh >= '0' && nextCh <= '9':
				padLen = (padLen * 10) + int(nextCh-'0')
				continue
			case nextCh == 's':
				fmtString(w, args, padLen)

				break parseFmt
			}

			// reached end of formatting string without finding a verb
			singleByte[0] = '?'
			doWrite(w, singleByte)
		}
		blockStart, blockEnd = blockEnd+1, blockEnd+1
	}

	if blockStart != blockEnd {
		// passing format[blockStart:blockEnd] to doWrite triggers a
		// memory allocation so we need to do this one byte at a time.
		for i := blockStart; i < blockEnd; i++ {
			singleByte[0] = format[i]
			doWrite(w, singleByte)
		}
	}

}

// fmtString prints a formatted version of string value v, applying
// the padding specified by padLen.
func fmtString(w *PrintkBuf, v string, padLen int) {
	castedVal := v
	fmtRepeat(w, ' ', padLen-len(castedVal))
	// converting the string to a byte slice triggers a memory allocation
	// so we need to do this one byte at a time.
	for i := 0; i < len(castedVal); i++ {
		singleByte[0] = castedVal[i]
		doWrite(w, singleByte)
	}
}

// fmtRepeat writes count bytes with value ch.
func fmtRepeat(w *PrintkBuf, ch byte, count int) {
	singleByte[0] = ch
	for i := 0; i < count; i++ {
		doWrite(w, singleByte)
	}
}

// fmtInt prints out a formatted version of v in the requested base, applying
// the padding specified by padLen. This function supports all built-in signed
// and unsigned integer types and base 8, 10 and 16 output.
func fmtInt(w *PrintkBuf, v uint64, base, padLen int) {
	var (
		sval             int64
		uval             uint64
		divider          uint64
		remainder        uint64
		padCh            byte
		left, right, end int
	)

	if padLen >= maxBufSize {
		padLen = maxBufSize - 1
	}

	switch base {
	case 8:
		divider = 8
		padCh = '0'
	case 10:
		divider = 10
		padCh = ' '
	case 16:
		divider = 16
		padCh = '0'
	}

	// 64-bit
	uval = v

	// Handle signs
	if sval < 0 {
		uval = uint64(-sval)
	} else if sval > 0 {
		uval = uint64(sval)
	}

	for right < maxBufSize {
		remainder = uval % divider
		if remainder < 10 {
			numFmtBuf[right] = byte(remainder) + '0'
		} else {
			// map values from 10 to 15 -> a-f
			numFmtBuf[right] = byte(remainder-10) + 'a'
		}

		right++

		uval /= divider
		if uval == 0 {
			break
		}
	}

	// Apply padding if required
	for ; right-left < padLen; right++ {
		numFmtBuf[right] = padCh
	}

	// Apply negative sign to the rightmost blank character (if using enough padding);
	// otherwise append the sign as a new char
	if sval < 0 {
		for end = right - 1; numFmtBuf[end] == ' '; end-- {
		}

		if end == right-1 {
			right++
		}

		numFmtBuf[end+1] = '-'
	}

	// Reverse in place
	end = right
	for right = right - 1; left < right; left, right = left+1, right-1 {
		numFmtBuf[left], numFmtBuf[right] = numFmtBuf[right], numFmtBuf[left]
	}

	for i := 0; i < end; i++ {
		singleByte[0] = numFmtBuf[i]
		doWrite(w, singleByte)
	}
}

func doWrite(w *PrintkBuf, p [1]byte) {
	if w != nil {
		w.buff[w.length] = p[0]
		w.length++
	}
}
