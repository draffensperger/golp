package lp

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLP(t *testing.T) {
	lp := NewLP(0, 2)
	lp.SetVerboseLevel(NEUTRAL)
	lp.SetColName(0, "x")
	lp.SetColName(1, "y")
	assert.Equal(t, "x", lp.GetColName(0))
	assert.Equal(t, "y", lp.GetColName(1))

	lp.AddConstraint([]float64{120.0, 210.0}, LE, 15000)
	lp.AddConstraintSparse([]Entry{Entry{col: 0, val: 110.0}, Entry{col: 1, val: 30.0}}, LE, 4000)
	lp.AddConstraintSparse([]Entry{Entry{col: 1, val: 1.0}, Entry{col: 0, val: 1.0}}, LE, 75)

	lp.SetObjFn([]float64{143, 60}, true)
	//lp.WriteToStdout()
	lp.Solve()

	delta := 0.000001
	assert.InDelta(t, 6315.625, lp.GetObjective(), delta)

	vars := lp.GetVariables()
	assert.Equal(t, len(vars), 2)
	assert.InDelta(t, 21.875, vars[0], delta)
	assert.InDelta(t, 53.125, vars[1], delta)
}

//	row := []C.double{120.0, 210.0}
//	C.add_constraintex(lp, j, &row[0], &colno[0], C.LE, 15000)

//	row = []C.double{110.0, 30.0}
//	C.add_constraintex(lp, j, &row[0], &colno[0], C.LE, 4000)

//	row = []C.double{1.0, 1.0}
//	C.add_constraintex(lp, j, &row[0], &colno[0], C.LE, 75)

//	C.set_add_rowmode(lp, 0)

//	row = []C.double{143.0, 60.0}
//	C.set_obj_fnex(lp, j, &row[0], &colno[0])
//	C.set_maxim(lp)
//	C.write_LP(lp, C.stdout)
//	C.set_verbose(lp, C.IMPORTANT)

//	ret := C.solve(lp)
//	fmt.Println("Solve returned: ")
//	fmt.Println(ret)

//	obj := C.get_objective(lp)
//	fmt.Println("Objective value: ")
//	fmt.Println(obj)

//	C.get_variables(lp, &row[0])
//	fmt.Println("x: ")
//	fmt.Println(row[0])

//	fmt.Println("y: ")
//	fmt.Println(row[1])

//	C.delete_lp(lp)

//	f := C.intFunc(C.fortytwo)
//	fmt.Println(int(C.bridge_int_func(f)))
//}
