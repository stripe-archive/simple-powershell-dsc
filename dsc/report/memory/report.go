package memory

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/stripe-archive/simple-powershell-dsc/dsc/types"
)

type ReportServer struct {
	reports map[types.GetReportsRequest][]byte
}

func New() *ReportServer {
	return &ReportServer{
		reports: make(map[types.GetReportsRequest][]byte),
	}
}

func (c *ReportServer) RegisterDscAgent(
	ctx context.Context,
	req types.RegisterDscAgentRequest,
) (*types.RegisterDscAgentResponse, error) {
	// No registration required
	return nil, nil
}

func (c *ReportServer) SendReport(
	ctx context.Context,
	req types.SendReportRequest,
) (*types.SendReportResponse, error) {
	body, err := json.Marshal(&req.Body)
	if err != nil {
		return nil, err
	}

	key := types.GetReportsRequest{
		AgentID: strings.ToLower(req.AgentID),
		JobID:   strings.ToLower(req.Body.JobID),
	}
	c.reports[key] = body

	return &types.SendReportResponse{}, nil
}

func (c *ReportServer) GetReports(
	ctx context.Context,
	req types.GetReportsRequest,
) (*types.GetReportsResponse, error) {
	key := types.GetReportsRequest{
		AgentID: strings.ToLower(req.AgentID),
		JobID:   strings.ToLower(req.JobID),
	}

	report, ok := c.reports[key]
	if !ok {
		return nil, types.ReportNotFoundError{
			AgentID: req.AgentID,
			JobID:   req.JobID,
		}
	}

	return &types.GetReportsResponse{Response: report}, nil
}
