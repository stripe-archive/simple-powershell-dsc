package status

import (
	"context"
	"fmt"

	"github.com/stripe-archive/simple-powershell-dsc/dsc"
	"github.com/stripe-archive/simple-powershell-dsc/dsc/types"
	"github.com/stripe-archive/simple-powershell-dsc/dsc/util"
)

// ReconcileDscStatus contains common logic to check hashes for and return the
// status of a GetDscAction request.
func ReconcileDscStatus(
	ctx context.Context,
	repo dsc.ConfigurationRepository,
	registered []string,
	req *types.GetDscActionRequest,
) (*types.GetDscActionResponse, error) {

	var resp []types.GetDscActionResponseBodyDetail

	// If we have no configurations, return nothing.
	if len(req.Body.ClientStatus) == 0 {
		ret := &types.GetDscActionResponse{
			Body: types.GetDscActionResponseBody{
				Details:    []types.GetDscActionResponseBodyDetail{},
				NodeStatus: "OK",
			},
		}
		return ret, nil
	}

	// If we get a single query...
	if len(req.Body.ClientStatus) == 1 {
		// If it has an empty checksum, it's because this
		// client hasn't run a consistency check yet; just return
		// "GetConfiguration" for every configuration that was registered,
		// rather than trying to check hashes, match names, etc.
		checksum := req.Body.ClientStatus[0].Checksum
		if checksum == "" {
			for _, reg := range registered {
				resp = append(resp, types.GetDscActionResponseBodyDetail{
					ConfigurationName: reg,
					Status:            "GetConfiguration",
				})
			}

			ret := &types.GetDscActionResponse{
				Body: types.GetDscActionResponseBody{
					Details:    resp,
					NodeStatus: "GetConfiguration",
				},
			}
			return ret, nil
		}

		// Otherwise, we should only have a single registered
		// configuration; check the hash and return.
		if len(registered) != 1 {
			return nil, fmt.Errorf("dsc/status: have %d registered configurations but client is only requesting 1", len(registered))
		}

		expected, _, err := util.GetConfigHash(ctx, repo, req.AgentID, registered[0])
		if err != nil {
			return nil, err
		}

		var statusStr string
		if expected != checksum {
			statusStr = "GetConfiguration"
		} else {
			statusStr = "OK"
		}

		ret := &types.GetDscActionResponse{
			Body: types.GetDscActionResponseBody{
				NodeStatus: statusStr,
				Details: []types.GetDscActionResponseBodyDetail{{
					ConfigurationName: registered[0],
					Status:            statusStr,
				}},
			},
		}
		return ret, nil
	}

	// As a heuristic, if we have a ConfigurationName in the request, it is
	// because the client is using PartialConfigurations and is requesting
	// the status of configurations with those names. Look up each of the
	// configurations with those names and return the appropriate status
	// depending on whether the hash matches.
	if req.Body.ClientStatus[0].ConfigurationName != "" {
		var updates int

		for _, status := range req.Body.ClientStatus {
			expected, _, err := util.GetConfigHash(ctx, repo, req.AgentID, status.ConfigurationName)
			if err != nil {
				// TODO: log/error/etc.
				continue
			}

			var statusStr string
			if expected != status.Checksum {
				updates++
				statusStr = "GetConfiguration"
			} else {
				statusStr = "OK"
			}

			resp = append(resp, types.GetDscActionResponseBodyDetail{
				ConfigurationName: status.ConfigurationName,
				Status:            statusStr,
			})
		}

		// The top-level NodeStatus depends on whether any of the
		// configurations need to be updated
		var status string
		if updates > 0 {
			status = "GetConfiguration"
		} else {
			status = "OK"
		}

		ret := &types.GetDscActionResponse{
			Body: types.GetDscActionResponseBody{
				Details:    resp,
				NodeStatus: status,
			},
		}
		return ret, nil
	}

	// If we get here, we have a request with more than one configuration,
	// but no names for them. We could assume that the order matches that
	// of the registration request, but while that appears to be true in my
	// tests, it isn't documented and feels fragile. Let's return an error.
	return nil, fmt.Errorf("dsc/status: don't support multiple (%d) non-partial configurations", len(req.Body.ClientStatus))
}
