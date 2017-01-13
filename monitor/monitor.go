package monitor

import "github.com/Knetic/govaluate"

type Monitor struct {
	Enabled    bool
	Name       string
	Expression string
}

type Monitors []Monitor

func (mon *Monitor) Evaluate(param map[string]interface{}) (bool, error) {
	exp, err := govaluate.NewEvaluableExpression(mon.Expression)
	if err != nil {
		return false, err
	}
	res, err := exp.Evaluate(param)
	evr := res.(bool)
	return evr, nil
}
