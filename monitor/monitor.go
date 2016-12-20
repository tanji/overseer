package monitor

type Monitor struct {
	Enabled    bool
	Name       string
	Expression string
}

type Monitors []Monitor
