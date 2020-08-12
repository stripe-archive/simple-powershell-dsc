package types

import (
	"io"
)

// 3.6: The GetConfiguration request SHOULD<10> get the configuration from the
// server.
type GetConfigurationRequest struct {
	AgentID           string
	ConfigurationName string
}

// GetConfigurationRequest is a response to a GetConfigurationRequest
type GetConfigurationResponse struct {
	// 3.6.5.2.2: In the response body, the configuration represents a BLOB.
	//
	// 3.6.5.2.3: The client gets the configuration from the server as
	// content-type application/octet-stream in the response body for the
	// GetConfiguration request
	Content io.Reader

	// The server MUST send the checksum in the response headers as
	// specified in section 2.2.2.2.
	Checksum string

	// The server MUST send the ChecksumAlgorithm in the response headers
	// as specified in section 2.2.2.3.
	ChecksumAlgorithm string

	//ConfigurationName string
}

// 3.7: GetModule request SHOULD<11> get the module from the server.
type GetModuleRequest struct {
	AgentID string
	Name    string
	Version string
}

// GetModuleResponse is a response to a GetModuleRequest
type GetModuleResponse struct {
	// 3.7.5.1.1.2: ModuleData represents a BLOB.
	//
	// 3.7.5.1.1.3: The client gets the module from the server as
	// content-type application/octet-stream in the response body for the
	// GetModule request.
	Content io.Reader

	// The server MUST send the checksum in the response headers as
	// specified in section 2.2.2.2.
	Checksum string

	// The server MUST send the ChecksumAlgorithm in the response headers
	// as specified in section 2.2.2.3.
	ChecksumAlgorithm string
}

// 3.8: The GetDscAction request SHOULD<12> get the action, as specified in
// section 3.8.5.1.1.2, from the server.
type GetDscActionRequest struct {
	AgentID string
	Body    GetDscActionRequestBody
}

// 3.8.5.1.1.2: The ActionContent packet is used by the server to transfer the
// following data fields:
//   - NodeStatus: MUST be either GetConfiguration, UpdateMetaConfiguration, Retry, or OK.
//   - Details: An array of the following fields:
//     - ConfigurationName: An opaque name as specified in section 2.2.2.4.
//     - Status: MUST be either GetConfiguration, UpdateMetaConfiguration, Retry, or OK.
type GetDscActionResponse struct {
	Body GetDscActionResponseBody
}

// 3.9: The RegisterDscAgent request SHOULD<13> register a client with a
// server, as specified in section 3.9.5.1.1.1.
type RegisterDscAgentRequest struct {
	AgentID string
	Body    RegisterDscAgentRequestBody
}

// 3.9.5.1.1.2: RegisterDscAgentContent represents a BLOB.
type RegisterDscAgentResponse struct{}

// 3.10: The SendReport request SHOULD<15> send the status report, as specified
// in section 3.10.5.1.1.1, to the server.
type SendReportRequest struct {
	AgentID string
	Body    SendReportRequestBody
}

// 3.10.5.1.1.2: The ReportContent packet does not contain any data
type SendReportResponse struct{}

// 3.11: The GetReports request SHOULD<16> get the status report from the
// server, as specified in section 3.11.5.1.1.1.
type GetReportsRequest struct {
	AgentID string
	JobID   string
}

// 3.11.5.1.1.2: ReportContent represents a BLOB.
type GetReportsResponse struct {
	Response []byte
}

// 3.12: With the CertificateRotation request, the server MUST rotate the
// client certificate with the new certificate as specified in section
// 3.12.5.1.1.1.<17>
type CertificateRotationRequest struct {
	AgentID string
	Body    RotateCertificateRequestBody
}

// 3.12.5.1.1.2 ResponseBody: None.
type CertificateRotationResponse struct{}
