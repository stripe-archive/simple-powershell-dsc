package dsc

import (
	"context"

	"github.com/stripe-archive/simple-powershell-dsc/dsc/types"
)

// ConfigurationRepository is the interface that should be implemented in order
// to serve configuration and modules to DSC clients.
type ConfigurationRepository interface {
	// RegisterDscAgent is called by the client in order to register a
	// client with the ConfigurationRepository.
	RegisterDscAgent(ctx context.Context, req types.RegisterDscAgentRequest) (*types.RegisterDscAgentResponse, error)

	// GetConfiguration is called by the client in order to fetch a given
	// DSC configuration.
	GetConfiguration(ctx context.Context, req types.GetConfigurationRequest) (*types.GetConfigurationResponse, error)

	// GetModule is called by the client in order to fetch a given DSC
	// module.
	GetModule(ctx context.Context, req types.GetModuleRequest) (*types.GetModuleResponse, error)
}

// ReportServer is the interface that should be implemented in order to handle
// report submission and retrieval for DSC clients.
type ReportServer interface {
	// RegisterDscAgent is called by the client in order to register a
	// client with the ReportServer.
	RegisterDscAgent(ctx context.Context, req types.RegisterDscAgentRequest) (*types.RegisterDscAgentResponse, error)

	// SendReport is called by the client in order to send a report to the
	// ReportServer.
	SendReport(ctx context.Context, req types.SendReportRequest) (*types.SendReportResponse, error)

	// GetReports is called by the client in order to fetch a report from
	// the ReportServer.
	GetReports(ctx context.Context, req types.GetReportsRequest) (*types.GetReportsResponse, error)
}

// NodeStatus is the interface that should be implemented in order to store the
// status of a node.
type NodeStatus interface {
	// RegisterDscAgent is called by the client in order to register a
	// client with the NodeStatus.
	RegisterDscAgent(ctx context.Context, req types.RegisterDscAgentRequest) (*types.RegisterDscAgentResponse, error)

	// GetDscAction should return the action to perform for the provided
	// node. This method should compare the provided state to the known
	// state of the node and return an appropriate action to perform.
	GetDscAction(ctx context.Context, req types.GetDscActionRequest) (*types.GetDscActionResponse, error)
}
