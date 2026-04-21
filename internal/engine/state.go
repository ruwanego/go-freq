package engine

type EngineState string

type RuntimeMode string

const (
	StateBooting    EngineState = "BOOTING"
	StateRecovering EngineState = "RECOVERING"
	StateReady      EngineState = "READY"
	StateFailed     EngineState = "FAILED"
)

const (
	ModeBacktest RuntimeMode = "BACKTEST"
	ModeLive     RuntimeMode = "LIVE"
)
