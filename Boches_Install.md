## Bochs 2.6.11 安装 

### configure 配置
```shell
./configure --with-x11 --with-wx --enable-debugger --enable-disasm --enable-all-optimizations --enable-readline --enable-long-phy-address --enable-ltdl-install --enable-idle-hack --enable-plugins --enable-a20-pin --enable-x86-64 --enable-smp --enable-cpu-level=6 --enable-large-ramfile --enable-repeat-speedups --enable-fast-function-calls --enable-handlers-chaining --enable-trace-linking --enable-configurable-msrs --enable-show-ips --enable-cpp --enable-debugger-gui --enable-iodebug --enable-logging --enable-assert-checks --enable-fpu --enable-vmx=2 --enable-svm --enable-3dnow --enable-alignment-check --enable-monitor-mwait --enable-avx --enable-evex --enable-x86-debugger --enable-pci --enable-usb --enable-voodoo
```

### path 否则编译会报错
```c++
Description: Fix the build with SMP enabled
Origin: https://sourceforge.net/p/bochs/code/13778/

Index: bochs/bx_debug/dbg_main.cc
===================================================================
--- bochs/bx_debug/dbg_main.cc	(revision 13777)
+++ bochs/bx_debug/dbg_main.cc	(working copy)
@@ -1494,11 +1494,11 @@
 {
   char cpu_param_name[16];
 a
-  Bit32u index = BX_ITLB_INDEX_OF(laddr);		//这一行改成下面一行
+  Bit32u index = BX_CPU(dbg_cpu)->ITLB.get_index_of(laddr);
   sprintf(cpu_param_name, "ITLB.entry%d", index);
   bx_dbg_show_param_command(cpu_param_name, 0);
 
-  index = BX_DTLB_INDEX_OF(laddr, 0);		//同理
+  index = BX_CPU(dbg_cpu)->DTLB.get_index_of(laddr);
   sprintf(cpu_param_name, "DTLB.entry%d", index);
   bx_dbg_show_param_command(cpu_param_name, 0);
 }
```

### 编译安装
```shell
make
make install
```

### bochs 版本信息
```shell
bochs --version
========================================================================
                       Bochs x86 Emulator 2.6.11
              Built from SVN snapshot on January 5, 2020
                Timestamp: Sun Jan  5 08:36:00 CET 2020
========================================================================
Usage: bochs [flags] [bochsrc options]

  -n               no configuration file
  -f configfile    specify configuration file
  -q               quick start (skip configuration interface)
  -benchmark N     run bochs in benchmark mode for N millions of emulated ticks
  -dumpstats N     dump bochs stats every N millions of emulated ticks
  -r path          restore the Bochs state from path
  -log filename    specify Bochs log file name
  -unlock          unlock Bochs images leftover from previous session
  -rc filename     execute debugger commands stored in file
  -dbglog filename specify Bochs internal debugger log file name
  --help           display this help and exit
  --help features  display available features / devices and exit
  --help cpu       display supported CPU models and exit

For information on Bochs configuration file arguments, see the
bochsrc section in the user documentation or the man page of bochsrc.
00000000000p[      ] >>PANIC<< command line arg '--version' was not understood
00000000000e[SIM   ] notify called, but no bxevent_callback function is registered
00000000000e[SIM   ] notify called, but no bxevent_callback function is registered
========================================================================
Bochs is exiting with the following message:
[      ] command line arg '--version' was not understood
========================================================================
00000000000i[SIM   ] quit_sim called with exit code 1
```