[![GoDoc](https://godoc.org/github.com/draffensperger/golp?status.svg)](https://godoc.org/github.com/draffensperger/golp) [![Build Status](https://travis-ci.org/draffensperger/golp.svg?branch=master)](https://travis-ci.org/draffensperger/golp)

Golp is a Golang wrapper for the [LPSolve](http://lpsolve.sourceforge.net/5.5/) linear (and integer) programming library.

## Installation

Here's how to install it.

go get https://github.com/draffensperger/golp


Then get LP solve itself. Golp is configured to dynamically link to LP solve,
and it expects the 
and for instance if you are on 64-bit Linux, 

http://sourceforge.net/projects/lpsolve/files/lpsolve/5.5.2.0/

## Usage 

Not all LPSolve functions are supported, but it's currently possible to run a simple linear program using golp. See `lp_test.go` for an example. Note that the column indices are always zero based.

### Example with real-valued variables

### MIP (Mixed Integer Programming) example

## Linking to LPSolve

You will need to copy `liblpsolve55.so` to the `./lib/lp_solve` folder in the directory where your final executable will be. The shared library included in this repo is for 64 bit Linux. You would need to get an appropriate library file for other systems from the LPSolve site. LPSolve source and binaries are at the [LPSolve SourceForge page](http://sourceforge.net/projects/lpsolve/). Static linking should be possible with a tweak to the cgo settings.

## Alternative Linear Programming Packages



## Licenses and acknowledgments

The golp golang code itself is MIT licensed. However, LPSolve itself is [licensed under the LGPL](http://lpsolve.sourceforge.net/5.5/LGPL.htm). The `stringbuilder.c` code is from [breckinloggins/libuseful](https://github.com/breckinloggins/libuseful).
