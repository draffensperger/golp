// NOTE: This code starts column index at 0, converts to starting at 1 for lpsolve library

package golp

/*
#cgo CFLAGS: -I./lib/lp_solve
#cgo LDFLAGS: -L./lib/lp_solve/ -llpsolve55 -Wl,-rpath=./lib/lp_solve
#include <stdlib.h>
#include "lp_lib.h"
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

type LP struct {
	ptr *C.lprec
}

func NewLP(rows, cols int) *LP {
	l := new(LP)
	l.ptr = C.make_lp(C.int(rows), C.int(cols))
	runtime.SetFinalizer(l, deleteLP)
	l.SetAddRowMode(true)
	return l
}

const ( // iota is reset to 0
	NEUTRAL  = iota // NEUTRAL == 0
	CRITICAL        // CRITICAL == 1
	SEVERE
	IMPORTANT
	NORMAL
	DETAILED
	FULL
)

func (l *LP) SetVerboseLevel(level int) {
	C.set_verbose(l.ptr, C.int(level))
}

func deleteLP(l *LP) {
	C.delete_lp(l.ptr)
}

func (l *LP) SetColName(col int, name string) {
	cstrName := C.CString(name)
	C.set_col_name(l.ptr, C.int(col+1), cstrName)
	C.free(unsafe.Pointer(cstrName))
}

func (l *LP) GetColName(col int) string {
	return C.GoString(C.get_col_name(l.ptr, C.int(col+1)))
}

func (l *LP) SetAddRowMode(addRowMode bool) {
	C.set_add_rowmode(l.ptr, boolToUChar(addRowMode))
}

func boolToUChar(b bool) C.uchar {
	if b {
		return C.uchar(1)
	}
	return C.uchar(0)
}

type ConstraintType int

const ( // iota is reset to 0
	_  = iota // don't use 0
	LE        // LE == 1
	GE        // GE == 2
	EQ        // EQ == 3
)

func (l *LP) AddConstraint(row []float64, ct ConstraintType, rightHand float64) error {
	cRow := make([]C.double, len(row)+1)
	cRow[0] = 0.0
	for i := 0; i < len(row); i++ {
		cRow[i+1] = C.double(row[i])
	}
	C.add_constraint(l.ptr, &cRow[0], C.int(ct), C.double(rightHand))
	return nil
}

type Entry struct {
	Col int
	Val float64
}

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

func (l *LP) SetObjFn(row []float64, maximize bool) {
	l.SetAddRowMode(false)

	cRow := make([]C.double, len(row)+1)
	cRow[0] = 0.0
	for i := 0; i < len(row); i++ {
		cRow[i+1] = C.double(row[i])
	}
	C.set_obj_fn(l.ptr, &cRow[0])

	if maximize {
		C.set_maxim(l.ptr)
	}
}

type SolutionType int

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

func (l *LP) Solve() SolutionType {
	return SolutionType(C.solve(l.ptr))
}

func (l *LP) WriteToStdout() {
	C.write_LP(l.ptr, C.stdout)
}

func (l *LP) WriteToString() string {
	cstr := C.write_lp_to_str(l.ptr)
	str := C.GoString(cstr)
	C.free(unsafe.Pointer(cstr))
	return str
}

func (l *LP) GetObjective() float64 {
	return float64(C.get_objective(l.ptr))
}

func (l *LP) GetVariables() []float64 {
	numCols := int(C.get_Ncolumns(l.ptr))
	cRow := make([]C.double, numCols)
	C.get_variables(l.ptr, &cRow[0])
	row := make([]float64, numCols)
	for i := 0; i < numCols; i++ {
		row[i] = float64(cRow[i])
	}
	return row
}
