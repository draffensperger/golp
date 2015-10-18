/*
Package golp gives Go bindings for LPSolve, a Mixed Integer Linear
Programming (MILP) solver.

For usage examples, see https://github.com/draffensperger/golp#examples.

Not all LPSolve functions have bindings. Feel free to open an issue or
contact me if you would like more added.

One difference from the LPSolve C library, is that the golp columns are always
zero-based.

The Go code of golp is MIT licensed, but LPSolve itself is licensed under the
LGPL. This roughly means that you can include golp in a closed-source project
as long as you do not modify LPSolve itself and you use dynamic linking to
access LPSolve (and provide a way for someone to link your program to a
different version of LPSolve).
For the legal details: http://lpsolve.sourceforge.net/5.0/LGPL.htm
*/
package golp

/*
// For Mac, assume LPSolve installed via MacPorts
#cgo darwin CFLAGS: -I/opt/local/include/lpsolve
#cgo darwin LDFLAGS: -L/opt/local/lib -llpsolve55

// For Linux, assume LPSolve bundled in local lpsolve directory
#cgo linux CFLAGS: -I./lpsolve
#cgo linux LDFLAGS: -L./lpsolve -llpsolve55 -Wl,-rpath=./lpsolve

#include "lp_lib.h"
#include <stdlib.h>
#include "stringbuilder.h"

int write_lp_to_str_callback(void* userhandle, char* buf) {
	sb_append_str((stringbuilder*) userhandle, buf);
	return 0;
}

char* write_lp_to_str(lprec *lp) {
	stringbuilder* sb = sb_new();
	write_lpex(lp, sb, write_lp_to_str_callback);
	char* str = sb_cstring(sb);
	sb_destroy(sb, 0);
	return str;
}
*/
import "C"

import (
	"runtime"
	"unsafe"
)

// LP stores a linear (or mixed integer) programming problem
type LP struct {
	ptr *C.lprec
}

// NewLP create a new linear program structure with specified number of rows and
// columns. The underlying C data structure's memory will be freed in a Go
// finalizer, so there is no need to explicitly deallocate it.
func NewLP(rows, cols int) *LP {
	l := new(LP)
	l.ptr = C.make_lp(C.int(rows), C.int(cols))
	runtime.SetFinalizer(l, deleteLP)
	l.SetAddRowMode(true)
	l.SetVerboseLevel(IMPORTANT)
	return l
}

func deleteLP(l *LP) {
	C.delete_lp(l.ptr)
}

// NumRows returns the number of rows (constraints) in the linear program.
// See http://lpsolve.sourceforge.net/5.5/get_Nrows.htm
func (l *LP) NumRows() int {
	return int(C.get_Nrows(l.ptr))
}

// NumCols returns the number of columns (variables) in the linear program.
// See http://lpsolve.sourceforge.net/5.5/get_Ncolumns.htm
func (l *LP) NumCols() int {
	return int(C.get_Ncolumns(l.ptr))
}

// Verbose levels
const ( // iota is reset to 0
	NEUTRAL  = iota // NEUTRAL == 0
	CRITICAL        // CRITICAL == 1
	SEVERE
	IMPORTANT
	NORMAL
	DETAILED
	FULL
)

// SetVerboseLevel changes the output verbose level (golp defaults it to
// IMPORTANT).
// See http://lpsolve.sourceforge.net/5.1/set_verbose.htm
func (l *LP) SetVerboseLevel(level int) {
	C.set_verbose(l.ptr, C.int(level))
}

// SetColName changes a column name. Unlike the LPSolve C library, col is zero-based
func (l *LP) SetColName(col int, name string) {
	cstrName := C.CString(name)
	C.set_col_name(l.ptr, C.int(col+1), cstrName)
	C.free(unsafe.Pointer(cstrName))
}

// ColName gives a column name, index is zero-based.
func (l *LP) ColName(col int) string {
	return C.GoString(C.get_col_name(l.ptr, C.int(col+1)))
}

// SetInt specifies that the given column must take an integer value.
// This triggers LPSolve to use branch-and-bound instead of simplex to solve.
// See http://lpsolve.sourceforge.net/5.5/set_int.htm
func (l *LP) SetInt(col int, mustBeInt bool) {
	C.set_int(l.ptr, C.int(col+1), boolToUChar(mustBeInt))
}

// IsInt returns whether the given column must take an integer value
// See http://lpsolve.sourceforge.net/5.5/is_int.htm
func (l *LP) IsInt(col int) bool {
	return uCharToBool(C.is_int(l.ptr, C.int(col+1)))
}

// SetBinary specifies that the given column bust take a binary (0 or 1) value
// See http://lpsolve.sourceforge.net/5.5/set_binary.htm
func (l *LP) SetBinary(col int, mustBeBinary bool) {
	C.set_binary(l.ptr, C.int(col+1), boolToUChar(mustBeBinary))
}

// IsBinary returns whether the given column must take a binary (0 or 1) value
// See http://lpsolve.sourceforge.net/5.5/is_binary.htm
func (l *LP) IsBinary(col int) bool {
	return uCharToBool(C.is_binary(l.ptr, C.int(col+1)))
}

// SetAddRowMode specifies whether adding by row (true) or by column (false)
// performs best. By default NewLP sets this for adding by row to perform best.
// See http://lpsolve.sourceforge.net/5.5/set_add_rowmode.htm
func (l *LP) SetAddRowMode(addRowMode bool) {
	C.set_add_rowmode(l.ptr, boolToUChar(addRowMode))
}

func boolToUChar(b bool) C.uchar {
	if b {
		return C.uchar(1)
	}
	return C.uchar(0)
}

func uCharToBool(c C.uchar) bool {
	return c != C.uchar(0)
}

// ConstraintType can be less than (golp.LE), greater than (golp.GE) or equal (golp.EQ)
type ConstraintType int

// Contraint type constants
const ( // iota is reset to 0
	_  = iota // don't use 0
	LE        // LE == 1
	GE        // GE == 2
	EQ        // EQ == 3
)

// AddConstraint adds a constraint to the linear program. This (unlike the
// LPSolve C function), exects the data in the row param to start at index 0
// for the first column.
// See http://lpsolve.sourceforge.net/5.5/add_constraint.htm
func (l *LP) AddConstraint(row []float64, ct ConstraintType, rightHand float64) error {
	cRow := make([]C.double, len(row)+1)
	cRow[0] = 0.0
	for i := 0; i < len(row); i++ {
		cRow[i+1] = C.double(row[i])
	}
	C.add_constraint(l.ptr, &cRow[0], C.int(ct), C.double(rightHand))
	return nil
}

// Entry is for sparse constraint or objective function rows
type Entry struct {
	Col int
	Val float64
}

// AddConstraintSparse adds a constraint row by specifying only the non-zero
// entries. Entries column indices are zero-based.
// See http://lpsolve.sourceforge.net/5.5/add_constraint.htm
func (l *LP) AddConstraintSparse(row []Entry, ct ConstraintType, rightHand float64) error {
	cRow := make([]C.double, len(row))
	cColNums := make([]C.int, len(row))
	for i, entry := range row {
		cRow[i] = C.double(entry.Val)
		cColNums[i] = C.int(entry.Col + 1)
	}
	C.add_constraintex(l.ptr, C.int(len(row)), &cRow[0], &cColNums[0], C.int(ct), C.double(rightHand))
	return nil
}

// SetObjFn changes the objective function. Row indices are zero-based.
// See http://lpsolve.sourceforge.net/5.5/set_obj_fn.htm
func (l *LP) SetObjFn(row []float64) {
	l.SetAddRowMode(false)

	cRow := make([]C.double, len(row)+1)
	cRow[0] = 0.0
	for i := 0; i < len(row); i++ {
		cRow[i+1] = C.double(row[i])
	}
	C.set_obj_fn(l.ptr, &cRow[0])
}

// SetMaximize will set the objective function  to maximize instead of
// minimizing by default.
// and http://lpsolve.sourceforge.net/5.5/set_maxim.htm
func (l *LP) SetMaximize() {
	C.set_maxim(l.ptr)
}

// SolutionType represents the result type
type SolutionType int

// Constacts for the solution result type
// See http://lpsolve.sourceforge.net/5.5/solve.htm
const ( // iota is reset to 0
	NOMEMORY   = -2
	OPTIMAL    = iota // don't use 0
	SUBOPTIMAL        // SUBOPTIMAL == 1
	INFEASIBLE        // INFEASIBLE == 2
	UNBOUNDED
	DEGENERATE
	NUMFAILURE
	USERABORT
	TIMEOUT
	PROCFAIL
	PROCBREAK
	FEASFOUND
	NOFEASFOUND
)

// Solve the linear (or mixed integer) program and return the solution type
// See http://lpsolve.sourceforge.net/5.5/solve.htm
func (l *LP) Solve() SolutionType {
	return SolutionType(C.solve(l.ptr))
}

// WriteToStdout writes a representation of the linear program to standard out
// See http://lpsolve.sourceforge.net/5.5/write_lp.htm
func (l *LP) WriteToStdout() {
	C.write_LP(l.ptr, C.stdout)
}

// WriteToString returns a representation of the linear program as a string
func (l *LP) WriteToString() string {
	cstr := C.write_lp_to_str(l.ptr)
	str := C.GoString(cstr)
	C.free(unsafe.Pointer(cstr))
	return str
}

// Objective gives the value of the objective function of the solved linear
// program.
// See http://lpsolve.sourceforge.net/5.5/get_objective.htm
func (l *LP) Objective() float64 {
	return float64(C.get_objective(l.ptr))
}

// Variables return the values for the variables of the solved linear program
// See http://lpsolve.sourceforge.net/5.5/get_variables.htm
func (l *LP) Variables() []float64 {
	numCols := int(C.get_Ncolumns(l.ptr))
	cRow := make([]C.double, numCols)
	C.get_variables(l.ptr, &cRow[0])
	row := make([]float64, numCols)
	for i := 0; i < numCols; i++ {
		row[i] = float64(cRow[i])
	}
	return row
}
