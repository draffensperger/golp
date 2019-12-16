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
#cgo linux CFLAGS: -I${SRCDIR}/lpsolve
#cgo linux LDFLAGS: -L${SRCDIR}/lpsolve -llpsolve55 -Wl,-rpath=${SRCDIR}/lpsolve

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
	"fmt"
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

// VerboseLevel represents different verbose levels,
// see http://lpsolve.sourceforge.net/5.1/set_verbose.htm
type VerboseLevel int

// Verbose levels
const (
	NEUTRAL  VerboseLevel = iota // NEUTRAL == 0
	CRITICAL                     // CRITICAL == 1
	SEVERE
	IMPORTANT
	NORMAL
	DETAILED
	FULL
)

// Note that we can't use stringer because this does not work well with cgo
// yet: https://github.com/golang/go/issues/20358

func (level VerboseLevel) String() string {
	switch level {
	case NEUTRAL:
		return "NEUTRAL"
	case CRITICAL:
		return "CRITICAL"
	case SEVERE:
		return "SEVERE"
	case IMPORTANT:
		return "IMPORTANT"
	case NORMAL:
		return "NORMAL"
	case DETAILED:
		return "DETAILED"
	case FULL:
		return "FULL"
	default:
		return fmt.Sprintf("VerboseLevel(%d)", int(level))
	}
}

// SetVerboseLevel changes the output verbose level (golp defaults it to
// IMPORTANT).
// See http://lpsolve.sourceforge.net/5.1/set_verbose.htm
func (l *LP) SetVerboseLevel(level VerboseLevel) {
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

// SetUnbounded specifies that the given column has a lower bound of -infinity
// and an upper bound of +infinity. (By default, columns have a lower bound of
// 0 and an upper bound of +infinity.)
// See http://lpsolve.sourceforge.net/5.5/set_unbounded.htm
func (l *LP) SetUnbounded(col int) {
	C.set_unbounded(l.ptr, C.int(col+1))
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

// SetBinary specifies that the given column must take a binary (0 or 1) value
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

// PresolveType specifies type of presolve,
// see http://lpsolve.sourceforge.net/5.5/set_presolve.htm
type PresolveType int

// Presolve types
const (
	NONE        PresolveType = 0
	ROWS                     = 1
	COLS                     = 2
	LINDEP                   = 4
	SOS                      = 32
	REDUCEMIP                = 64
	KNAPSACK                 = 128
	ELIMEQ2                  = 256
	IMPLIEDFREE              = 512
	REDUCEGCD                = 1024
	PROBEFIX                 = 2048
	PROBEREDUCE              = 4096
	ROWDOMANITE              = 8192
	COLDOMINATE              = 16384
	MERGEROWS                = 32768
	COLFIXDUAL               = 131072
	BOUNDS                   = 262144
	DUALS                    = 524288
	SENSDUALS                = 1048576
)

func (level PresolveType) String() string {
	switch level {
	case NONE:
		return "PRESOLVE_NONE"
	case ROWS:
		return "PRESOLVE_ROWS"
	case COLS:
		return "PRESOLVE_COLS"
	case LINDEP:
		return "PRESOLVE_LINDEP"
	case SOS:
		return "PRESOLVE_SOS"
	case REDUCEMIP:
		return "PRESOLVE_REDUCEMIP"
	case KNAPSACK:
		return "PRESOLVE_KNAPSACK"
	case ELIMEQ2:
		return "PRESOLVE_ELIMEQ2"
	case IMPLIEDFREE:
		return "PRESOLVE_IMPLIEDFREE"
	case REDUCEGCD:
		return "PRESOLVE_REDUCEGCD"
	case PROBEFIX:
		return "PRESOLVE_PROBEFIX"
	case PROBEREDUCE:
		return "PRESOLVE_PROBEREDUCE"
	case ROWDOMANITE:
		return "PRESOLVE_ROWDOMINATE"
	case COLDOMINATE:
		return "PRESOLVE_COLDOMINATE"
	case MERGEROWS:
		return "PRESOLVE_MERGEROWS"
	case COLFIXDUAL:
		return "PRESOLVE_COLFIXDUAL"
	case BOUNDS:
		return "PRESOLVE_BOUNDS"
	case DUALS:
		return "PRESOLVE_DUALS"
	case SENSDUALS:
		return "PRESOLVE_SENSDUALS"
	default:
		return fmt.Sprintf("PresolveType(%d)", int(level))
	}
}

// SetPresolve specifies whether pre solve should be used to try to simplify problem,
// by default it is set to not to perform pre solve, level specifies type of pre solve
// and maxLoops the maximum number of times pre solve may be done (use 0 to determine
// number of pre solve loops automatically by get_presolveloop()).
// For more info see: http://lpsolve.sourceforge.net/5.5/set_presolve.htm
func (l *LP) SetPresolve(level PresolveType, maxLoops int) {
	if maxLoops == 0 {
		maxLoops = l.GetPresolveLoops()
	}
	C.set_presolve(l.ptr, C.int(level), C.int(maxLoops))
}

// GetPresolveLoops determines optimal number of loops for pre solve.
// See: http://lpsolve.sourceforge.net/5.5/get_presolveloops.htm
func (l *LP) GetPresolveLoops() int {
	return int(C.get_presolveloops(l.ptr))
}

// ConstraintType can be less than (golp.LE), greater than (golp.GE) or equal (golp.EQ)
type ConstraintType int

// Contraint type constants
const ( // iota is reset to 0
	_  ConstraintType = iota // don't use 0
	LE                       // LE == 1
	GE                       // GE == 2
	EQ                       // EQ == 3
)

func (t ConstraintType) String() string {
	switch t {
	case LE:
		return "LE"
	case GE:
		return "GE"
	case EQ:
		return "EQ"
	default:
		return fmt.Sprintf("ConstraintType(%d)", int(t))
	}
}

// AddConstraint adds a constraint to the linear program. This (unlike the
// LPSolve C function), expects the data in the row param to start at index 0
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

// SolutionType represents the result type.
type SolutionType int

// Return values must not be enumerated from 0 in, many are not used
// any more and therefore there are gaps.
// Also lpsolve55 will not return PROCFAIL and other types any more,
// they're here for compatibility reasons.
// To make this clear we don't use iota but list the values.

// Constants for the solution result type.
// See http://lpsolve.sourceforge.net/5.5/solve.htm
const (
	NOMEMORY    SolutionType = -2
	OPTIMAL                  = 0
	SUBOPTIMAL               = 1
	INFEASIBLE               = 2
	UNBOUNDED                = 3
	DEGENERATE               = 4
	NUMFAILURE               = 5
	USERABORT                = 6
	TIMEOUT                  = 7
	PROCFAIL                 = 10
	PROCBREAK                = 11
	FEASFOUND                = 12
	NOFEASFOUND              = 13
)

func (t SolutionType) String() string {
	switch t {
	case NOMEMORY:
		return "NOMEMORY"
	case OPTIMAL:
		return "OPTIMAL"
	case SUBOPTIMAL:
		return "SUBOPTIMAL"
	case INFEASIBLE:
		return "INFEASIBLE"
	case UNBOUNDED:
		return "UNBOUNDED"
	case DEGENERATE:
		return "DEGENERATE"
	case NUMFAILURE:
		return "NUMFAILURE"
	case USERABORT:
		return "USERABORT"
	case TIMEOUT:
		return "TIMEOUT"
	case PROCFAIL:
		return "PROCFAIL"
	case PROCBREAK:
		return "PROCBREAK"
	case FEASFOUND:
		return "FEASFOUND"
	case NOFEASFOUND:
		return "NOFEASFOUND"
	default:
		return fmt.Sprintf("SolutionType(%d)", int(t))
	}
}

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
