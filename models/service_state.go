package models

// service state
const (
	_ ServiceState = iota
	ServiceStateAlive
	ServiceStateDead
)

// ServiceState indicates current service state
type ServiceState uint32
