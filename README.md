[![GoDoc](https://godoc.org/github.com/draffensperger/golp?status.svg)](https://godoc.org/github.com/draffensperger/golp) [![Build Status](https://travis-ci.org/draffensperger/golp.svg?branch=master)](https://travis-ci.org/draffensperger/golp) [![Code Climate](https://codeclimate.com/github/draffensperger/golp/badges/gpa.svg)](https://codeclimate.com/github/draffensperger/golp)

Golp is a Golang wrapper for the [LPSolve](http://lpsolve.sourceforge.net/5.5/) linear (and integer) programming library.

## Installation

**Step 1: Get the golp Go code**

```
go get github.com/draffensperger/golp
```

**Step 2: Get the LPSolve library**

Golp is configured to dynamically link to LPSolve and expects its files to be 
in the `lib/lp_solve` folder for your project.  You will need an LPSolve library 
suitable for your operating system, which you can 
[get from SourceForge](http://sourceforge.net/projects/lpsolve/files/lpsolve/5.5.2.0/).

Here's how you could download and extract the LPSolve library for 64-bit Linux:

```
LP_URL=http://sourceforge.net/projects/lpsolve/files/lpsolve/5.5.2.0/lp_solve_5.5.2.0_dev_ux64.tar.gz
mkdir -p lib/lp_solve
wget -qO- $LP_URL | tar xvz -C lib/lp_solve
```

With some configuration changes, it would be possible to statically link to 
LPSolve but that may have licensing/distribution implications for your project 
since LP Solve is [LGPL licensed](http://lpsolve.sourceforge.net/5.5/LGPL.htm). 

## Usage 

Not all LPSolve functions are supported, but it's currently possible to run a 
simple linear and integer program using golp. Note that unlike the LP Solve C
library, the column indices are always zero based. Feel free to open a GitHub
issue or pull request if you'd like more functions added.

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
branch-and-bound algorithm for such problems. This example is from the
[LPSolve integer variables documentation](http://lpsolve.sourceforge.net/5.5/integer.htm).


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

## Alternative linear / mixed integer solver libraries

There are also Go bindings for the GPL-licensed
[GNU Linear Programming Kit (GLPK)](http://www.gnu.org/software/glpk/) at
[github.com/lukpank/go-glpk](https://github.com/lukpank/go-glpk).

The Google [or-tools](https://github.com/google/or-tools) project provides a C++
SWIG compatible inteface for a number of other linear and mixed integer solvers
like CBC, CLP, GLOP, Gurobi, CPLEX, SCIP, and Sulum. 
There is Go support [for SWIG bindings](http://www.swig.org/Doc2.0/Go.html), so
it should be possible to write a wrapper that would connect to those other
solvers via the or-tools library as well.

## License

The LPSolve library this project depends on is
[LGPL licensed](http://lpsolve.sourceforge.net/5.5/LGPL.htm). 

The `stringbuilder.c` code is from [breckinloggins/libuseful](https://github.com/breckinloggins/libuseful).

The golp Go code is MIT licensed as follows:

The MIT License (MIT)

Copyright (c) 2015 David Raffensperger

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
THE SOFTWARE.
