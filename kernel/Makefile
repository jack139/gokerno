CFLAGS = -mcmodel=large -m64 -nostdlib -nostartfiles -nodefaultlibs
#LDFLAGS = -z muldefs -static -nostdlib
LDFLAGS = -z muldefs -nostdlib
LIBPATH = -L/usr/lib/x86_64-linux-gnu -L/usr/lib/gcc/x86_64-linux-gnu/6.3.0
LIBS = -lgcc -lgcc_eh -lgo -lc #-lpthread

GOARCH := amd64
GOOS := linux

all: system
	objcopy -I elf64-x86-64 -S -R ".eh_frame" -R ".comment" -O binary system kernel.bin

system:	head.o terminal.o main.o
	ld -b elf64-x86-64 $(LDFLAGS) -o $@ $^ -T Kernel.lds $(LIBPATH) $(LIBS)

terminal.o:	terminal.go
	gccgo $(CFLAGS) -c $^ -o $@

main.o:	main.go
	gccgo $(CFLAGS) -c $^ -o $@

head.o:	head.S
	gcc -E $^ > head.s
	as --64 -o $@ head.s

clean:
	rm -rf *.a *.o *.s~ *.s *.S~ *.c~ *.h~ system  Makefile~ Kernel.lds~ kernel.bin 

