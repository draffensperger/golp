[![GoDoc](https://godoc.org/github.com/draffensperger/golp?status.svg)](https://godoc.org/github.com/draffensperger/golp) [![Build Status](https://travis-ci.org/draffensperger/golp.svg?branch=master)](https://travis-ci.org/draffensperger/golp)

Golp is a Golang wrapper for the [LPSolve](http://lpsolve.sourceforge.net/5.5/) linear (and integer) programming library.

## Installation

To use golp, first you need to get the golp Go code:

```
go get https://github.com/draffensperger/golp
```

Then you will need to get LP solve itself. Golp is configured to dynamically
link to LP solve and expects its files to be in the `lib/lp_solve` folder for
your project. 

You will need an LP Solve build suitable for your operating system, which you
can [download from SourceForge](http://sourceforge.net/projects/lpsolve/files/lpsolve/5.5.2.0/).

Here's how you could download and extract the LP Solve library for 64-bit Linux:

```
wget -qO- http://sourceforge.net/projects/lpsolve/files/lpsolve/5.5.2.0/lp_solve_5.5.2.0_dev_ux64.tar.gz | tar xvz -C lib/lp_solve
```

With some configuration changes, it would be possible to statically link to LP
Solve but that may have licensing implications for your project since LP Solve
is LGPL licensed.

## Usage 

Not all LPSolve functions are supported, but it's currently possible to run a 
simple linear program using golp. See `lp_test.go` for an example. Note that 
the column indices are always zero based.

### Example with real-valued variables

### MIP (Mixed Integer Programming) example

## Linking to LPSolve

## Alternative Linear Programming Packages



## Licenses and acknowledgments

The golp golang code itself is MIT licensed. However, LPSolve itself is [licensed under the LGPL](http://lpsolve.sourceforge.net/5.5/LGPL.htm). The `stringbuilder.c` code is from [breckinloggins/libuseful](https://github.com/breckinloggins/libuseful).
