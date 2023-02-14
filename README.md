# exesizecmp

Tool for breaking down section size differences between two executables.

Sample usage:

```
$ go build -o regular.exe .
$ go build -gcflags=all=-l -o noinl.exe .
$ ./exesizecmp -i=regular.exe,noinl.exe
.debug_frame	36971	47336	10365	p=0.3%
.debug_info	324515	251907	-72608	p=-0.2%
.debug_line	183460	139005	-44455	p=-0.2%
.debug_loc	198301	176922	-21379	p=-0.1%
.debug_ranges	70153	26842	-43311	p=-0.6%
.go.buildinfo	528	544	16	p=0.0%
.gopclntab	579104	594024	14920	p=0.0%
.rodata		337062	250263	-86799	p=-0.3%
.strtab		68520	97775	29255	p=0.4%
.symtab		73224	97776	24552	p=0.3%
.text		859480	862859	3379	p=0.0%
.typelink	2132	2224	92	p=0.0%
$
```

