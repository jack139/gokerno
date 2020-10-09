## 编译 Bootloader

### 编译 boot 和 loader
```shell
nasm boot.asm -o boot.bin
nasm loader.asm -o loader.bin
```

### 创建磁盘 image
```shell
bximage
```

### 制作启动 image
```shell
dd if=boot.bin of=../boot.img bs=512 count=1 conv=notrunc
mount ../boot.img /mnt/floppy/ -t vfat -o loop
cp loader.bin /mnt/floppy/
cp kernel.bin /mnt/floppy/
sync
umount /mnt/floppy/
```



## kernel 静态编译问题

因为gccgo相关库不是 -mcmodel=large 编译，所以链接会报错

```shell
ld -b elf64-x86-64 -z muldefs -static -nostdlib -o system head.o main.o -T Kernel.lds -L/usr/lib/x86_64-linux-gnu -L/usr/lib/gcc/x86_64-linux-gnu/6.3.0 --start-group -lpthread -lc -lgcc_eh -lgcc --end-group

ld: system: section .tdata lma 0xffff8000001b97f0 adjusted to 0xffff8000001b98e8
/usr/lib/gcc/x86_64-linux-gnu/6.3.0/libgcc.a(generic-morestack.o)：在函数‘__morestack_block_signals’中：
(.text+0x89a):  截断重寻址至相符: R_X86_64_PLT32 针对未定义的符号 pthread_sigmask
/usr/lib/gcc/x86_64-linux-gnu/6.3.0/libgcc.a(generic-morestack.o)：在函数‘__morestack_unblock_signals’中：
(.text+0x8f8):  截断重寻址至相符: R_X86_64_PLT32 针对未定义的符号 pthread_sigmask
/usr/lib/x86_64-linux-gnu/libc.a(assert.o)：在函数‘__assert_fail_base’中：
(.text+0x2d):  截断重寻址至相符: R_X86_64_PC32 针对未定义的符号 __pthread_setcancelstate
/usr/lib/x86_64-linux-gnu/libc.a(assert.o)：在函数‘__assert_fail_base’中：
(.text+0x35):  截断重寻址至相符: R_X86_64_32 针对 .rodata.str1.1
/usr/lib/x86_64-linux-gnu/libc.a(assert.o)：在函数‘__assert_fail_base’中：
(.text+0x4c):  截断重寻址至相符: R_X86_64_32 针对 .rodata.str1.1
/usr/lib/x86_64-linux-gnu/libc.a(assert.o)：在函数‘__assert_fail_base’中：
(.text+0x51):  截断重寻址至相符: R_X86_64_32 针对 .rodata.str1.1
/usr/lib/x86_64-linux-gnu/libc.a(assert.o)：在函数‘__assert_fail_base’中：
(.text+0x8b):  截断重寻址至相符: R_X86_64_32 针对 .rodata.str1.1
/usr/lib/x86_64-linux-gnu/libc.a(assert.o)：在函数‘__assert_fail_base’中：
(.text+0x11e):  截断重寻址至相符: R_X86_64_32 针对 .rodata.str1.16
/usr/lib/x86_64-linux-gnu/libc.a(assert.o)：在函数‘__assert_fail_base’中：
(.text+0x12f):  截断重寻址至相符: R_X86_64_32 针对 .rodata.str1.1
/usr/lib/x86_64-linux-gnu/libc.a(assert.o)：在函数‘__assert_fail’中：
(.text+0x153):  截断重寻址至相符: R_X86_64_32 针对 .rodata.str1.8
/usr/lib/x86_64-linux-gnu/libc.a(assert.o)：在函数‘__assert_fail’中：
(.text+0x15d):  从输出所省略的额外重寻址溢出
Makefile:14: recipe for target 'system' failed
make: *** [system] Error 1
```

不加 -static 时，可以链接成功

```shell
ld -b elf64-x86-64 -z muldefs -nostdlib -o system head.o main.o -T Kernel.lds -L/usr/lib/x86_64-linux-gnu -L/usr/lib/gcc/x86_64-linux-gnu/6.3.0 --start-group -lpthread -lc -lgcc_eh -lgcc --end-group
```

但是会需要运行时动态库

```shell
$ readelf -d system 

Dynamic section at offset 0x10afb0 contains 20 entries:
  标记        类型                         名称/值
 0x0000000000000001 (NEEDED)             共享库：[libpthread.so.0]
 0x0000000000000001 (NEEDED)             共享库：[libc.so.6]
 0x0000000000000001 (NEEDED)             共享库：[ld-linux-x86-64.so.2]
```

目前无影响，因为内核状态不需要相关动态库（C标准库 和 线程库）。

>参考： https://stackoverflow.com/questions/21796637/using-simple-linker-script-gives-relocation-truncated-to-fit-r-x86-64-32s
