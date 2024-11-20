package xrtmodels

// Schedule is used to register timed, polled reads of device resources
// The definition can refer to https://github.com/IOTechSystems/xrt-docs/blob/v2.0-branch/docs/mqtt-management/mqtt-management.md#schedule-format
type Schedule struct {
	Name     string             `json:"name"`
	Device   string             `json:"device"`
	Resource []string           `json:"resource"`
	Interval uint64             `json:"interval,omitempty"`
	OnChange bool               `json:"on_change"`
	Bounds   map[string]float64 `json:"bounds"`
	Publish  bool               `json:"publish"`
	Units    bool               `json:"units"`
	Options  interface{}        `json:"options,omitempty"`
}
