package config

// Queues config
var Queues = struct {
	Default string
	Delayed string
	Expired string
}{
	Default: "tasker:queue:default",
	Delayed: "tasker:queue:delayed",
	Expired: "tasker:queue:expired",
}
