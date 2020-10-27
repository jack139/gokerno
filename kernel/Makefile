AS = x86_64-linux-gnu-as
GCC = x86_64-linux-gnu-gcc
GCCGO = x86_64-linux-gnu-gccgo
LD = x86_64-linux-gnu-ld
OBJCOPY = x86_64-linux-gnu-objcopy
GOFLAGS = -g -m64 -static -mcmodel=large -Werror -nostdinc -nostdlib
INCLUDE_DIRS = -I.
LDFLAGS = -m64 -mcmodel=large -nostdlib -nostartfiles -nodefaultlibs -ffreestanding #-static
LIBPATH = -L/usr/lib/x86_64-linux-gnu -L/usr/lib/gcc/x86_64-linux-gnu/6.3.0
LIBS = -lgcc -lgcc_eh -lgo -lc #-lpthread

GOARCH := amd64
GOOS := linux

GO_OBJ = terminal.o kfmt.o main.o

all: system
	$(OBJCOPY) -I elf64-x86-64 -S -R ".eh_frame" -R ".comment" -O binary system kernel.bin

system: head.o $(GO_OBJ)
	$(GCC) $(LDFLAGS) -o $@ $^ -T Kernel.lds $(LIBPATH) $(LIBS)
	#$(GCCGO) -mcmodel=large -o $@ $^ -static -Wl,-z,muldefs -Wl,-T,Kernel.lds
	#$(LD) $(LDFLAGS) -o $@ $^ -T Kernel.lds $(LIBPATH) $(LIBS)

head.o: head.S
	$(GCC) -E $^ > head.s
	$(AS) --64 -o $@ head.s

main.o:	main.go
	$(GCCGO) -c $(GOFLAGS) $(INCLUDE_DIRS) -o $@ -static $^

%.o: %/*.go
	$(GCCGO) -c $(GOFLAGS) $(INCLUDE_DIRS) -o $@ -static $^

clean:
	rm -rf *.a *.o *.s~ *.s *.S~ *.c~ *.h~ system  Makefile~ Kernel.lds~ kernel.bin 

