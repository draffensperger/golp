Golp is a Golang wrapper for the [LPSolve](http://lpsolve.sourceforge.net/5.5/) linear (and integer) programming library.

## Usage 

Not all LPSolve functions are supported, but it's currently possible to run a simple linear program using golp. See `lp_test.go` for an example. Note that the column indices are always zero based.

## Linking to LPSolve

You will need to copy `liblpsolve55.so` to the `./lib/lp_solve` folder in the directory where your final executable will be. The shared library included in this repo is for 64 bit Linux. You would need to get an appropriate library file for other systems from the LPSolve site. LPSolve source and binaries are at the [LPSolve SourceForge page](http://sourceforge.net/projects/lpsolve/). Static linking should be possible with a tweak to the cgo settings.

## Licenses and acknowledgments

The golp golang code itself is MIT licensed. However, LPSolve itself is [licensed under the LGPL](http://lpsolve.sourceforge.net/5.5/LGPL.htm). The `stringbuilder.c` code is from [breckinloggins/libuseful](https://github.com/breckinloggins/libuseful).
