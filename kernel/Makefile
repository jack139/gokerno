AS = x86_64-linux-gnu-as
GCC = x86_64-linux-gnu-gcc
GCCGO = x86_64-linux-gnu-gccgo
LD = x86_64-linux-gnu-ld
OBJCOPY = x86_64-linux-gnu-objcopy
GOFLAGS = -g -m64 -static -mcmodel=large -Werror -nostdinc -nostdlib
INCLUDE_DIRS = -I.
#LDFLAGS = -m64 -mcmodel=large -nostdlib -nostartfiles -nodefaultlibs -ffreestanding #-static
LDFLAGS = -n -static --no-ld-generated-unwind-info
LIBPATH = -L/usr/lib/x86_64-linux-gnu -L/usr/lib/gcc/x86_64-linux-gnu/6.3.0
LIBS = -lgcc -lgcc_eh -lgo -lc #-lpthread

GOARCH := amd64
GOOS := linux

GO_OBJ = terminal.o kfmt.o main.o

GO = go
BUILD_DIR := build
BUILD_ABS_DIR := $(CURDIR)/$(BUILD_DIR)
GOPATH := $(BUILD_ABS_DIR):$(shell pwd):$(GOPATH)
GOROOT := $(shell $(GO) env GOROOT)

all: system
	$(OBJCOPY) -I elf64-x86-64 -S -R ".eh_frame" -R ".comment" -O binary system kernel.bin

system: head.o go.o
	@#$(GCC) $(LDFLAGS) -o $@ $^ -T Kernel.lds #$(LIBPATH) $(LIBS)
	@#$(GCCGO) -mcmodel=large -o $@ $^ -static -Wl,-z,muldefs -Wl,-T,Kernel.lds
	@#$(LD) $(LDFLAGS) -o $@ $^ -T Kernel.lds $(LIBPATH) $(LIBS)
	$(LD) $(LDFLAGS) -T Kernel.lds -o $@ head.o build/go.o

go.o:
	@mkdir -p $(BUILD_DIR)

	@echo "[go] compiling go sources into a standalone .o file"
	GOARCH=$(GOARCH) GOOS=$(GOOS) GOPATH=$(GOPATH) $(GO) build -gcflags '$(GC_FLAGS)' -n kernel 2>&1 | sed \
	    -e "1s|^|set -e\n|" \
	    -e "1s|^|export GOOS=$(GOOS)\n|" \
	    -e "1s|^|export GOARCH=$(GOARCH)\n|" \
	    -e "1s|^|export GOROOT=$(GOROOT)\n|" \
	    -e "1s|^|export CGO_ENABLED=0\n|" \
	    -e "1s|^|alias pack='$(GO) tool pack'\n|" \
	    -e "/^mv/d" \
	    -e "/\/buildid/d" \
	    -e "s|-extld|-tmpdir='$(BUILD_ABS_DIR)' -linkmode=external -extldflags='-nostartfiles -nodefaultlibs -nostdlib -r' -extld|g" \
	    -e 's|$$WORK|$(BUILD_ABS_DIR)|g' \
            | sh 2>&1 |  sed -e "s/^/  | /g"

	@# build/go.o is a elf32 object file but all go symbols are unexported. Our
	@# asm entrypoint code needs to know the address to 'main.main' so we use
	@# objcopy to make that symbol exportable. Since nasm does not support externs
	@# with slashes we create a global symbol alias for kernel.Kmain
	@echo "[objcopy] create main.StartKerno alias"
	@echo "[objcopy] globalizing symbols {runtime.g0/m0/physPageSize}"
	objcopy \
		--add-symbol main.StartKerno=.text:0x`nm $(BUILD_DIR)/go.o | grep "main.StartKerno$$" | cut -d' ' -f1` \
		--globalize-symbol runtime.g0 \
		--globalize-symbol runtime.m0 \
		--globalize-symbol runtime.physPageSize \
		 $(BUILD_DIR)/go.o $(BUILD_DIR)/go.o

head.o: head.S
	$(GCC) -E $^ > head.s
	$(AS) --64 -o $@ head.s

clean:
	rm -rf *.a *.o *.s~ *.s *.S~ *.c~ *.h~ system  Makefile~ Kernel.lds~ kernel.bin 
	test -d $(BUILD_DIR) && rm -rf $(BUILD_DIR) || true

