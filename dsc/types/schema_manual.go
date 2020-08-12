package types

import (
	"fmt"
	"net"
	"strings"
	"time"
)

// This file contains types that don't have JSON schema definitions and are
// generated manually.

// 3.10.5.1.1.1: The ReportRequest packet is used by the client to transfer the
// following data fields:
type SendReportRequestBody struct {
	// JobId: The JobId parameter is a universally unique identifier (UUID)
	// as specified in [RFC4122] section 3.
	JobID string `json:"JobId"`

	// OperationType: A value that identifies the type for the operation.
	OperationType string `json:"OperationType"`

	// RefreshMode: A value that identifies whether the client is in PUSH
	// or PULL mode.
	RefreshMode string `json:"RefreshMode,omitempty"`

	// Status: A value that identifies the status of the current operation.
	Status string `json:"Status,omitempty"`

	// LCMVersion: A value that identifies the report generator on the
	// client.
	LCMVersion string `json:"LCMVersion,omitempty"`

	// ReportFormatVersion: A value that finds the identifier for the
	// report.
	ReportFormatVersion string `json:"ReportFormatVersion"`

	// ConfigurationVersion: A value that contains a string of two to four
	// groups of digits where the groups are separated by a period.
	ConfigurationVersion string `json:"ConfigurationVersion,omitempty"`

	// NodeName: A value that is used to identify the name of the client.
	NodeName string `json:"NodeName,omitempty"`

	// IpAddress: A value that identifies the IP addresses of the client
	// separated by a semicolon (;).
	IPAddress string `json:"IpAddress,omitempty"`

	// StartTime: A value that identifies the start time of an operation on
	// the client.
	StartTime string `json:"StartTime,omitempty"`

	// EndTime: A value that identifies the end time of an operation on the
	// client.
	EndTime string `json:"EndTime,omitempty"`

	// RebootRequested: A value that identifies whether the client
	// requested a reboot.
	RebootRequested string `json:"RebootRequested,omitempty"`

	// Errors: A value that represents the errors for an operation on the
	// client.
	Errors []string `json:"Errors,omitempty"`

	// StatusData: A value that represents the status of an operation on
	// the client.
	StatusData []string `json:"StatusData,omitempty"`

	// TODO: AdditionalData
}

// SendReportRequestBodyNorm is a version of the SendReport request body that
// is a "nicer" structure with more Go-like types. This isn't the default,
// because we need to be able to return the original report to a caller of
// GetReports.
type SendReportRequestBodyNorm struct {
	JobID                string       `json:"JobId"`
	OperationType        string       `json:"OperationType"`
	RefreshMode          string       `json:"RefreshMode,omitempty"`
	Status               string       `json:"Status,omitempty"`
	LCMVersion           string       `json:"LCMVersion,omitempty"`
	ReportFormatVersion  string       `json:"ReportFormatVersion"`
	ConfigurationVersion string       `json:"ConfigurationVersion,omitempty"`
	NodeName             string       `json:"NodeName,omitempty"`
	IPAddress            []net.IPAddr `json:"IpAddress,omitempty"`
	StartTime            *time.Time   `json:"StartTime,omitempty"`
	EndTime              *time.Time   `json:"EndTime,omitempty"`
	RebootRequested      *bool        `json:"RebootRequested,omitempty"`
	Errors               []string     `json:"Errors,omitempty"`
	StatusData           []string     `json:"StatusData,omitempty"`
}

// Convert the raw SendReport request body into a "nicer" structure with more
// Go-like types. This isn't the default, because we need to be able to return
// the original report to a caller of GetReports.
func (r SendReportRequestBody) Normalized() (*SendReportRequestBodyNorm, error) {
	ret := &SendReportRequestBodyNorm{
		// Fields we pass through as-is
		JobID:                r.JobID,
		OperationType:        r.OperationType,
		RefreshMode:          r.RefreshMode,
		Status:               r.Status,
		LCMVersion:           r.LCMVersion,
		ReportFormatVersion:  r.ReportFormatVersion,
		ConfigurationVersion: r.ConfigurationVersion,
		NodeName:             r.NodeName,
		Errors:               r.Errors,
		StatusData:           r.StatusData,
	}

	// Convert fields
	if r.StartTime != "" {
		t, err := time.Parse(time.RFC3339Nano, r.StartTime)
		if err != nil {
			return nil, err
		}
		ret.StartTime = &t
	}

	if r.EndTime != "" {
		t, err := time.Parse(time.RFC3339Nano, r.EndTime)
		if err != nil {
			return nil, err
		}
		ret.EndTime = &t
	}
	switch r.RebootRequested {
	case "":
		ret.RebootRequested = nil
	case "True", "true":
		ret.RebootRequested = new(bool)
		*ret.RebootRequested = true
	case "False", "false":
		ret.RebootRequested = new(bool)
		*ret.RebootRequested = false
	default:
		return nil, fmt.Errorf("unrecognized boolean value: %s", r.RebootRequested)
	}

	if r.IPAddress != "" {
		for _, val := range strings.Split(r.IPAddress, ";") {
			var zone, ip string

			// Remove any IPv6 Zone
			j := strings.IndexByte(val, '%')
			if j >= 0 {
				ip = val[0:j]
				zone = val[j+1:]
			} else {
				ip = val
			}

			// Parse the IP component
			if pip := net.ParseIP(ip); pip != nil {
				ret.IPAddress = append(ret.IPAddress, net.IPAddr{
					IP:   pip,
					Zone: zone,
				})
				continue
			}

			return nil, fmt.Errorf("could not parse IP address: %s", ip)
		}
	}

	// TODO: any validations?

	return ret, nil
}
