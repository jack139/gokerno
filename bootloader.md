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
sync
umount /mnt/floppy/
```
