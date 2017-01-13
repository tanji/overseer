package monitor

import "testing"

func TestEvaluateMonitor(t *testing.T) {
	m := Monitor{true, "myname", "threads_connected > 32"}
	param := make(map[string]interface{}, 1)
	param["threads_connected"] = 33
	res, err := m.Evaluate(param)
	if err != nil {
		t.Fatal("Error evaluating expression:", err)
	}
	if res != true {
		t.Error("Expression did not evaluate")
	}
}
