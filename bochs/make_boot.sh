#!/bin/bash

sudo mount boot.img /mnt/floppy/ -t vfat -o loop
sudo cp ../bootloader_bochs/loader.bin /mnt/floppy/
sudo cp ../kernel/kernel.bin /mnt/floppy/
sudo sync
sudo umount /mnt/floppy/
