// Package golp gives Go bindings for LPSolve.
package golp

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

// TestLP tests a real-valued linear programming example
func TestLP(t *testing.T) {
	lp := NewLP(0, 2)
	lp.SetVerboseLevel(NEUTRAL)
	lp.SetColName(0, "x")
	lp.SetColName(1, "y")
	assert.Equal(t, "x", lp.ColName(0))
	assert.Equal(t, "y", lp.ColName(1))

	require.NoError(t, lp.AddConstraint([]float64{120.0, 210.0}, LE, 15000))
	require.NoError(t, lp.AddConstraintSparse([]Entry{Entry{Col: 0, Val: 110.0}, Entry{Col: 1, Val: 30.0}}, LE, 4000))
	require.NoError(t, lp.AddConstraintSparse([]Entry{Entry{Col: 1, Val: 1.0}, Entry{Col: 0, Val: 1.0}}, LE, 75))

	lp.SetObjFn([]float64{143, 60})
	lp.SetMaximize()

	lpString := "/* Objective function */\nmax: +143 x +60 y;\n\n/* Constraints */\n+120 x +210 y <= 15000;\n+110 x +30 y <= 4000;\n+x +y <= 75;\n"
	assert.Equal(t, lpString, lp.WriteToString())

	solType := lp.Solve()
	assert.Equal(t, solType, OPTIMAL)

	delta := 0.000001
	assert.InDelta(t, 6315.625, lp.Objective(), delta)

	vars := lp.Variables()
	assert.Equal(t, len(vars), 2)
	assert.InDelta(t, 21.875, vars[0], delta)
	assert.InDelta(t, 53.125, vars[1], delta)

	assert.InDelta(t, 0.0, lp.DualResult(0), delta)
	assert.InDelta(t, 1.0375, lp.DualResult(1), delta)
}

// TestMIP tests a mixed-integer programming example
func TestMIP(t *testing.T) {
	lp := NewLP(0, 4)
	require.NoError(t, lp.AddConstraintSparse([]Entry{{0, 1.0}, {1, 1.0}}, LE, 5.0))
	require.NoError(t, lp.AddConstraintSparse([]Entry{{0, 2.0}, {1, -1.0}}, GE, 0.0))
	require.NoError(t, lp.AddConstraintSparse([]Entry{{0, 1.0}, {1, 3.0}}, GE, 0.0))
	require.NoError(t, lp.AddConstraintSparse([]Entry{{2, 1.0}, {3, 1.0}}, GE, 0.5))
	require.NoError(t, lp.AddConstraintSparse([]Entry{{2, 1.0}}, GE, 1.1))
	lp.SetObjFn([]float64{-1.0, -2.0, 0.1, 3.0})

	lp.SetInt(2, true)
	assert.Equal(t, lp.IsInt(2), true)

	solType := lp.Solve()
	assert.Equal(t, solType, OPTIMAL)

	delta := 0.000001
	assert.InDelta(t, -8.133333333, lp.Objective(), delta)

	vars := lp.Variables()
	assert.Equal(t, lp.NumCols(), 4)
	assert.Equal(t, len(vars), 4)
	assert.InDelta(t, 1.6666666666, vars[0], delta)
	assert.InDelta(t, 3.3333333333, vars[1], delta)
	assert.InDelta(t, 2.0, vars[2], delta)
	assert.InDelta(t, 0.0, vars[3], delta)

	duals := lp.Duals()
	assert.Len(t, duals, 5)
	assert.InDelta(t, -1.6666666666, duals[0], delta)
	assert.InDelta(t, 0.3333333333, duals[1], delta)
	assert.InDelta(t, 0.0, duals[2], delta)
	assert.InDelta(t, 0.0, duals[3], delta)
	assert.InDelta(t, 0.0, duals[4], delta)
}
