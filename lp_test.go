package golp

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
	lp.AddConstraintSparse([]Entry{Entry{Col: 0, Val: 110.0}, Entry{Col: 1, Val: 30.0}}, LE, 4000)
	lp.AddConstraintSparse([]Entry{Entry{Col: 1, Val: 1.0}, Entry{Col: 0, Val: 1.0}}, LE, 75)

	lp.SetObjFn([]float64{143, 60}, true)

	lpString := "/* Objective function */\nmax: +143 x +60 y;\n\n/* Constraints */\n+120 x +210 y <= 15000;\n+110 x +30 y <= 4000;\n+x +y <= 75;\n"
	assert.Equal(t, lpString, lp.WriteToString())

	lp.Solve()

	delta := 0.000001
	assert.InDelta(t, 6315.625, lp.GetObjective(), delta)

	vars := lp.GetVariables()
	assert.Equal(t, len(vars), 2)
	assert.InDelta(t, 21.875, vars[0], delta)
	assert.InDelta(t, 53.125, vars[1], delta)
}
