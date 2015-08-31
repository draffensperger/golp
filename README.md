[![GoDoc](https://godoc.org/github.com/draffensperger/golp?status.svg)](https://godoc.org/github.com/draffensperger/golp) [![Build Status](https://travis-ci.org/draffensperger/golp.svg?branch=master)](https://travis-ci.org/draffensperger/golp)

Golp is a Golang wrapper for the [LPSolve](http://lpsolve.sourceforge.net/5.5/) linear (and integer) programming library.

## Installation

**Step 1: Get the golp Go code**

```
go get github.com/draffensperger/golp
```

**Step 2: Get the LPSolve library**

Golp is configured to dynamically link to LP solve and expects its files to be 
in the `lib/lp_solve` folder for your project.  You will need an LP Solve build 
suitable for your operating  system, which you can 
[get from SourceForge](http://sourceforge.net/projects/lpsolve/files/lpsolve/5.5.2.0/).

Here's how you could download and extract the LP Solve library for 64-bit Linux:

```
LP_URL=http://sourceforge.net/projects/lpsolve/files/lpsolve/5.5.2.0/lp_solve_5.5.2.0_dev_ux64.tar.gz
mkdir -p lib/lp_solve
wget -qO- $LP_URL | tar xvz -C lib/lp_solve
```

With some configuration changes, it would be possible to statically link to LP
Solve but that may have licensing implications for your project since LP Solve
is LGPL licensed.

## Usage 

Not all LPSolve functions are supported, but it's currently possible to run a 
simple linear and integer program using golp. Note that unlike the LP Solve C
library, the column indices are always zero based.

### Example with real-valued variables

The example below in an adaption of an example in the 
[LP Solve documentation.](http://lpsolve.sourceforge.net/5.5/formulate.htm)
for maximizing a farmer's profit.

```
package main

import "fmt"
import "github.com/draffensperger/golp"

func main() {
  lp := golp.NewLP(0, 2)
  lp.AddConstraint([]float64{110.0, 30.0}, golp.LE, 4000.0)
  lp.AddConstraint([]float64{1.0, 1.0}, golp.LE, 75.0)
  lp.SetObjFn([]float64{143.0, 60.0})
  lp.SetMaximize()

  lp.Solve()
  vars := lp.Variables()
  fmt.Printf("Plant %.3f acres of barley\n", vars[0])
  fmt.Printf("And  %.3f acres of wheat\n", vars[1])
  fmt.Printf("For optimal profit of $%.2f\n", lp.Objective())

  // No need to explicitly free underlying C structure as golp.LP finalizer will
}
```

Outputs:
```
Plant 21.875 acres of barley
And  53.125 acres of wheat
For optimal profit of $6315.62
```

### MIP (Mixed Integer Programming) example

LPSolve also supports setting variables to be integral or binary and uses the
branch-and-bound algorithm for such problems.


```
import "fmt"
import "github.com/draffensperger/golp"

func main() {
  lp := golp.NewLP(0, 4)
  lp.AddConstraintSparse([]golp.Entry{{0, 1.0}, {1, 1.0}}, golp.LE, 5.0)
  lp.AddConstraintSparse([]golp.Entry{{0, 2.0}, {1, -1.0}}, golp.GE, 0.0)
  lp.AddConstraintSparse([]golp.Entry{{0, 1.0}, {1, 3.0}}, golp.GE, 0.0)
  lp.AddConstraintSparse([]golp.Entry{{2, 1.0}, {3, 1.0}}, golp.GE, 0.5)
  lp.AddConstraintSparse([]golp.Entry{{2, 1.0}}, golp.GE, 1.1)
  lp.SetObjFn([]float64{-1.0, -2.0, 0.1, 3.0})
  lp.SetInt(2, true)
  lp.Solve()

  fmt.Printf("Objective value: %v\n", lp.Objective())
  vars := lp.Variables()
  fmt.Printf("Variable values:\n")
  for i := 0; i < lp.NumCols(); i++ {
    fmt.Printf("x%v = %v\n", i + 1, vars[i])
  }
}
```

Outputs:
```
Objective value: -8.133333333333333
Variable values:
x1 = 1.6666666666666665
x2 = 3.333333333333333
x3 = 2
x4 = 0
```

## Alternative Linear Programming Packages



## License

The golp golang code itself is MIT licensed. However, LPSolve itself is [licensed under the LGPL](http://lpsolve.sourceforge.net/5.5/LGPL.htm). The `stringbuilder.c` code is from [breckinloggins/libuseful](https://github.com/breckinloggins/libuseful).

