// NOTE: This code starts column index at 0, converts to starting at 1 for lpsolve library

package lp

/*
#cgo CFLAGS: -I./lib/lp_solve
#cgo LDFLAGS: -L./lib/lp_solve/ -llpsolve55 -Wl,-rpath=./lib/lp_solve
#include <stdlib.h>
#include "lp_lib.h"
*/
import "C"

import (
	"runtime"
	"unsafe"
)

type lp struct {
	ptr *C.lprec
}

func NewLP(rows, cols int) *lp {
	l := new(lp)
	l.ptr = C.make_lp(C.int(rows), C.int(cols))

	// Create the model by adding constraints by default
	C.set_add_rowmode(l.ptr, C.TRUE)

	runtime.SetFinalizer(l, deleteLP)
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

func (l lp) SetVerboseLevel(level int) {
	C.set_verbose(l.ptr, C.int(level))
}

func deleteLP(l *lp) {
	C.delete_lp(l.ptr)
}

func (l lp) SetColName(col int, name string) {
	cstrName := C.CString(name)
	C.set_col_name(l.ptr, C.int(col+1), cstrName)
	C.free(unsafe.Pointer(cstrName))
}

func (l lp) GetColName(col int) string {
	return C.GoString(C.get_col_name(l.ptr, C.int(col+1)))
}

type ConstraintType int

const ( // iota is reset to 0
	_  = iota // don't use 0
	LE        // LE == 1
	EQ        // EQ == 2
	GE        // GE == 3
)

func (l lp) AddConstraint(row []float64, ct ConstraintType, rightHand float64) error {
	cRow := make([]C.double, len(row)+1)
	cRow[0] = 0.0
	for i := 0; i < len(row); i++ {
		cRow[i+1] = C.double(row[i])
	}
	C.add_constraint(l.ptr, &cRow[0], C.int(ct), C.double(rightHand))
	return nil
}

type Entry struct {
	col int
	val float64
}

func (l lp) AddConstraintSparse(row []Entry, ct ConstraintType, rightHand float64) error {
	cRow := make([]C.double, len(row))
	cColNums := make([]C.int, len(row))
	for i, entry := range row {
		cRow[i] = C.double(entry.val)
		cColNums[i] = C.int(entry.col + 1)
	}
	C.add_constraintex(l.ptr, C.int(len(row)), &cRow[0], &cColNums[0], C.int(ct), C.double(rightHand))
	return nil
}

func (l lp) SetObjFn(row []float64, maximize bool) {
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

func (l lp) Solve() error {
	C.set_add_rowmode(l.ptr, C.FALSE)
	C.solve(l.ptr)
	return nil
}

func (l lp) WriteToStdout() {
	C.write_LP(l.ptr, C.stdout)
}

func (l lp) GetObjective() float64 {
	return float64(C.get_objective(l.ptr))
}

func (l lp) GetVariables() []float64 {
	numCols := int(C.get_Ncolumns(l.ptr))
	cRow := make([]C.double, numCols)
	C.get_variables(l.ptr, &cRow[0])
	row := make([]float64, numCols)
	for i := 0; i < numCols; i++ {
		row[i] = float64(cRow[i])
	}
	return row
}

//func formatLP(params TaskParams) {
//	/*
//		LP format:

//		step 1: break tasks up into hourly chunks, (later on re-weight the longer tasks as less rewarding)

//		variables are
//			date_hour_taskpart = 1.0 means do taskpart at date and hour
//			all of those date_hour_taskpart variables must be between 0 and 1
//			sum of all date_hour_taskpart = 1.0 (can only do one hour of total work per hour)

//		deadlines:
//			if a task part has a deadline, then the sum of all its work times before that deadline
//			must be 1.0

//		on or after:
//			if a task part has an on or after specified, the sum of work times before on or after must be 0

//		reward for each hour:
//			hour_reward = hour decay const * date_hour_taskpart * value/hr for task

//	*/

//	tasks := params.Tasks
//	//horizonHours := 22 * 8
//	horizonHours := 8

//	ncol := C.int(len(tasks) * horizonHours)
//	lp := C.make_lp(0, ncol)

//	for hour := 0; hour < horizonHours; hour++ {
//		for taskNum := 0; taskNum < len(tasks); taskNum++ {
//			nameStr := "h" + strconv.Itoa(hour) + "_t" + strconv.Itoa(taskNum)
//			name := C.CString(nameStr)
//			colNum := hour*len(tasks) + taskNum + 1
//			C.set_col_name(lp, C.int(colNum), name)
//			C.free(unsafe.Pointer(name))
//		}
//	}

//	C.set_add_rowmode(lp, 1)

//	// Each variable must be between 0 and 1

//	// Total tasks done in a hour must be <= 1
//	for hour := 0; hour < horizonHours; hour++ {
//		row := make([]C.double, len(tasks))
//		colNums := make([]C.int, len(tasks))
//		for taskNum := 0; taskNum < len(tasks); taskNum++ {
//			colNums[taskNum] = C.int(len(tasks)*hour + taskNum + 1)
//			row[taskNum] = 1.0
//		}
//		C.add_constraintex(lp, C.int(len(tasks)), &row[0], &colNums[0], C.LE, 1.0)
//	}

//	// Total amount done on each task must be <= task.EstimatedHours
//	for taskNum, task := range tasks {
//		row := make([]C.double, horizonHours)
//		colNums := make([]C.int, horizonHours)
//		for hour := 0; hour < horizonHours; hour++ {
//			colNums[hour] = C.int(len(tasks)*hour + taskNum + 1)
//			row[hour] = 1.0
//		}
//		C.add_constraintex(lp, C.int(horizonHours), &row[0], &colNums[0], C.LE, C.double(task.EstimatedHours))
//	}
//	C.set_add_rowmode(lp, 0)

//	// Objective function
//	decayRate := float32(0.95)
//	curHourValue := float32(1.0)
//	row := make([]C.double, len(tasks)*horizonHours+1)
//	for hour := 0; hour < horizonHours; hour++ {
//		for taskNum, task := range tasks {
//			row[len(tasks)*hour+taskNum+1] = C.double(curHourValue * task.Reward / task.EstimatedHours)
//		}
//		curHourValue *= decayRate
//	}

//	C.set_obj_fn(lp, &row[0])

//	C.set_maxim(lp)

//	fmt.Println("\n")
//	fmt.Println("LP formulation:")
//	C.write_LP(lp, C.stdout)

//	ret := C.solve(lp)
//	fmt.Printf("Solve returned: %v\n", ret)

//	obj := C.get_objective(lp)
//	fmt.Printf("Objective value: %v\n", obj)

//	C.get_variables(lp, &row[0])
//	for hour := 0; hour < horizonHours; hour++ {
//		for taskNum := 0; taskNum < len(tasks); taskNum++ {
//			nameStr := "h" + strconv.Itoa(hour) + "_t" + strconv.Itoa(taskNum)
//			val := row[len(tasks)*hour+taskNum]
//			fmt.Printf("%v: %v\n", nameStr, val)
//		}
//	}

//	C.delete_lp(lp)
//}
