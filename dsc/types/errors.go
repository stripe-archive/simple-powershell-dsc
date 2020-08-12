package types

import (
	//"errors"
	"fmt"
)

type ConfigurationNotFoundError struct {
	AgentID string
	Name    string
}

func (e ConfigurationNotFoundError) Error() string {
	return fmt.Sprintf("dsc: configuration %q not found", e.Name)
}

type ModuleNotFoundError struct {
	AgentID string
	Name    string
	Version string
}

func (e ModuleNotFoundError) Error() string {
	return fmt.Sprintf("dsc: module %q (version: %q) not found", e.Name, e.Version)
}

type ReportNotFoundError struct {
	AgentID string
	JobID   string
}

func (e ReportNotFoundError) Error() string {
	return fmt.Sprintf("dsc: job %q for agent %q not found", e.JobID, e.AgentID)
}

type AgentNotRegisteredError struct {
	AgentID string
}

func (e AgentNotRegisteredError) Error() string {
	return fmt.Sprintf("dsc: agent %q not registered", e.AgentID)
}
