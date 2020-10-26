GCC = x86_64-linux-gnu-gcc
GCCGO = x86_64-linux-gnu-gccgo
LD = x86_64-linux-gnu-ld
CFLAGS = -m64 -mcmodel=large -nostdinc
LDFLAGS = -m64 -z muldefs -nostdlib -ffreestanding #-static
LIBPATH = -L/usr/lib/x86_64-linux-gnu -L/usr/lib/gcc/x86_64-linux-gnu/6.3.0
LIBS = -lgcc -lgcc_eh -lgo -lc #-lpthread

GOARCH := amd64
GOOS := linux

GO_OBJ = terminal.o kfmt.o main.o

all: system
	objcopy -I elf64-x86-64 -S -R ".eh_frame" -R ".comment" -O binary system kernel.bin

head.o: head.S
	gcc -E $^ > head.s
	as --64 -o $@ head.s

system: head.o $(GO_OBJ)
	$(GCC) $(LDFLAGS) -o $@ $^ -T Kernel.lds $(LIBPATH) $(LIBS)

terminal.o:	terminal/*.go
	$(GCCGO) -c -g -m64  -o $@ -static $^

kfmt.o:	kfmt/*.go
	$(GCCGO) -c -g -m64  -o $@ -static $^

main.o:	main.go
	$(GCCGO) -c -g -m64  -o $@ -static $^

clean:
	rm -rf *.a *.o *.s~ *.s *.S~ *.c~ *.h~ system  Makefile~ Kernel.lds~ kernel.bin 

