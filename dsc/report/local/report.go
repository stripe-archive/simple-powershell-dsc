package local

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/stripe-archive/simple-powershell-dsc/dsc/internal/fsmock"
	"github.com/stripe-archive/simple-powershell-dsc/dsc/types"
)

type ReportServer struct {
	root string
	fs   fsmock.FileSystem
}

func New(root string) *ReportServer {
	return &ReportServer{root, fsmock.OSFS{}}
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
	// Make the directory if it doesn't exist.
	dirPath := filepath.Join(
		c.root,
		strings.ToLower(req.AgentID),
	)
	if err := c.fs.MkdirAll(dirPath, 0750); err != nil && !os.IsExist(err) {
		return nil, err
	}

	// Save into a file in the given directory
	path := filepath.Join(dirPath, strings.ToLower(req.Body.JobID)+".json")
	f, err := c.fs.Create(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	err = json.NewEncoder(f).Encode(&req.Body)
	if err != nil {
		return nil, err
	}

	return &types.SendReportResponse{}, nil
}

func (c *ReportServer) GetReports(
	ctx context.Context,
	req types.GetReportsRequest,
) (*types.GetReportsResponse, error) {
	// Open from the given directory
	path := filepath.Join(
		c.root,
		strings.ToLower(req.AgentID),
		strings.ToLower(req.JobID)+".json",
	)
	f, err := c.fs.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, types.ReportNotFoundError{
				AgentID: req.AgentID,
				JobID:   req.JobID,
			}
		}
		return nil, err
	}
	defer f.Close()

	report, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	ret := &types.GetReportsResponse{
		Response: report,
	}
	return ret, nil
}
