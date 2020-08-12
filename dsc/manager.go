package dsc

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"

	"github.com/sirupsen/logrus"
	"goji.io"

	"github.com/stripe-archive/simple-powershell-dsc/dsc/internal/middleware"
	"github.com/stripe-archive/simple-powershell-dsc/dsc/internal/regexpat"
	"github.com/stripe-archive/simple-powershell-dsc/dsc/internal/urls"
	"github.com/stripe-archive/simple-powershell-dsc/dsc/types"
)

// Manager is the core interface to a DSC server; it serves the appropriate
// routes over HTTP, converts incoming requests into request structs, calls the
// interfaces that implement actual functionality, and responds to the request
// with the approriate HTTP response.
type Manager struct {
	config ConfigurationRepository
	report ReportServer
	status NodeStatus

	mux  *goji.Mux
	log  logrus.FieldLogger
	keys []string // TODO: don't keep around?
}

// NewManager creates a new Manager with the given ConfigurationRepository and
// ReportServer implementations.
func NewManager(config ConfigurationRepository, report ReportServer, status NodeStatus, opts ...Option) *Manager {
	ret := &Manager{
		config: config,
		report: report,
		status: status,
	}

	for _, opt := range opts {
		opt(ret)
	}

	// Make mux + middleware
	ret.mux = goji.NewMux()
	ret.mux.Use(protocolVersion)
	ret.mux.Use(middleware.MaxBodySize(1 * 1024 * 1024))
	ret.mux.Use(middleware.LogrusLogger(ret.log))

	// Optionally add middlware if we have the right config
	if len(ret.keys) > 0 {
		ret.mux.Use(checkRegistration(ret.keys))
	}

	// Register routes on mux
	routes := []struct {
		Method  string
		Regexp  string
		Handler func(http.ResponseWriter, *http.Request)
	}{
		// API Versions 1.0 and 1.1 (not supported)
		{"GET", urls.GetConfigurationV1URL, ret.methodNotSupported},
		{"GET", urls.GetModuleV1URL, ret.methodNotSupported},
		{"POST", urls.GetActionV1URL, ret.methodNotSupported},
		{"POST", urls.SendStatusReportURL, ret.methodNotSupported},
		{"GET", urls.GetStatusReportURL, ret.methodNotSupported},

		// API Version 2.0
		{"GET", urls.GetConfigurationV2URL, ret.getConfiguration},
		{"GET", urls.GetModuleV2URL, ret.getModule},
		{"POST", urls.GetDscActionV2URL, ret.getDscAction},
		{"PUT", urls.RegisterDscAgentV2URL, ret.registerDscAgent},
		{"POST", urls.SendReportV2URL, ret.sendReport},
		{"GET", urls.GetReportsV2URL, ret.getReports},

		// TODO: support this, or at least do nothing
		{"POST", urls.CertificateRotationURL, ret.methodNotSupported},
	}

	for _, route := range routes {
		// Anchor the route regex at the root of the mux
		re := `^\/` + route.Regexp + `$`

		// Create new pattern that matches the given regex/method
		pat := regexpat.NewWithMethods(re, route.Method)

		// Attach it to the mux
		ret.mux.HandleFunc(pat, route.Handler)
	}

	return ret
}

// ServeHTTP implements the http.Handler interface
func (m *Manager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.mux.ServeHTTP(w, r)
}

func (m *Manager) registerDscAgent(w http.ResponseWriter, r *http.Request) {
	var err error

	var body types.RegisterDscAgentRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "error decoding body: %s", err)
		return
	}

	var regType string
	if s := body.RegistrationInformation.RegistrationMessageType; s != nil {
		regType = *s
	}

	// Get other stuff
	agentId := regexpat.Param(r, "agent_id")
	m.log.WithFields(logrus.Fields{
		"registration_type":   regType,
		"agent_id":            agentId,
		"configuration_names": body.ConfigurationNames,
	}).Debug("registering new agent")

	// Dispatch to the correct implementation depending on what we're
	// registering. This method will be called multiple times, usually once
	// per type.
	if regType == "ConfigurationRepository" {
		// Register with both config repository and node status handler
		req := types.RegisterDscAgentRequest{
			AgentID: regexpat.Param(r, "agent_id"),
			Body:    body,
		}

		_, err = m.config.RegisterDscAgent(r.Context(), req)
		if err == nil {
			_, err = m.status.RegisterDscAgent(r.Context(), req)
		}

	} else if regType == "ReportServer" {
		_, err = m.report.RegisterDscAgent(
			r.Context(),
			types.RegisterDscAgentRequest{
				AgentID: regexpat.Param(r, "agent_id"),
				Body:    body,
			},
		)
	} else {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "unknown registration type: %s", regType)
		return
	}

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error registering agent: %s", err)
		m.log.WithError(err).Errorf("error registering agent")
		return
	}

	// No response supported right now
	w.WriteHeader(http.StatusNoContent)
}

func (m *Manager) getConfiguration(w http.ResponseWriter, r *http.Request) {
	agentId := regexpat.Param(r, "agent_id")
	configName := regexpat.Param(r, "configuration_name")

	m.log.WithFields(logrus.Fields{
		"agent_id":           agentId,
		"configuration_name": configName,
	}).Debug("getting configuration")

	resp, err := m.config.GetConfiguration(
		r.Context(),
		types.GetConfigurationRequest{
			AgentID:           agentId,
			ConfigurationName: configName,
		},
	)
	if err != nil {
		switch v := err.(type) {
		case types.ConfigurationNotFoundError:
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "%s", v)

		default:
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "error getting configuration: %s", err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Checksum", resp.Checksum)

	// Case-sensitive
	w.Header()["ChecksumAlgorithm"] = []string{resp.ChecksumAlgorithm}

	_, err = io.Copy(w, resp.Content)
	if cl, ok := resp.Content.(io.Closer); ok {
		cl.Close()
	}
	if err != nil {
		m.log.WithError(err).Errorf("error copying response body")
		return
	}
}

var agentIdRegexp = regexp.MustCompile(urls.AgentId)

func (m *Manager) getModule(w http.ResponseWriter, r *http.Request) {
	agentId := r.Header.Get("AgentId")
	matched := agentIdRegexp.Match([]byte(agentId))
	if !matched {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "invalid agent id")
		return
	}

	moduleName := regexpat.Param(r, "module_name")
	moduleVersion := regexpat.Param(r, "module_version")

	m.log.WithFields(logrus.Fields{
		"agent_id":       agentId,
		"module_name":    moduleName,
		"module_version": moduleVersion,
	}).Debug("getting module")

	resp, err := m.config.GetModule(
		r.Context(),
		types.GetModuleRequest{
			AgentID: agentId,
			Name:    moduleName,
			Version: moduleVersion,
		},
	)
	if err != nil {
		switch v := err.(type) {
		case types.ModuleNotFoundError:
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "%s", v)

		default:
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "error getting module: %s", err)
			m.log.WithError(err).Errorf("error getting module")
		}
		return
	}

	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Checksum", resp.Checksum)

	// Case-sensitive
	w.Header()["ChecksumAlgorithm"] = []string{resp.ChecksumAlgorithm}

	_, err = io.Copy(w, resp.Content)
	if cl, ok := resp.Content.(io.Closer); ok {
		cl.Close()
	}
	if err != nil {
		m.log.WithError(err).Errorf("error copying response body")
		return
	}
}

func (m *Manager) getDscAction(w http.ResponseWriter, r *http.Request) {
	agentId := regexpat.Param(r, "agent_id")

	var body types.GetDscActionRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "error decoding body: %s", err)
		return
	}

	m.log.WithFields(logrus.Fields{
		"agent_id": agentId,
		"body":     body,
	}).Infof("getting action for client")

	resp, err := m.status.GetDscAction(
		r.Context(),
		types.GetDscActionRequest{
			AgentID: agentId,
			Body:    body,
		},
	)
	if err != nil {
		switch v := err.(type) {
		case types.AgentNotRegisteredError:
			// If the agent isn't registered, we send a 400 response to trigger a re-registration.
			// TODO: figure out why this doesn't actually re-register
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "agent not registered: %s", v.AgentID)

		default:
			m.log.WithError(err).Errorf("error getting status")

			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "error getting status: %s", err)
		}

		return
	}

	// Log information about about the response
	for _, d := range resp.Body.Details {
		m.log.WithFields(logrus.Fields{
			"agent_id":           agentId,
			"configuration_name": d.ConfigurationName,
			"status":             d.Status,
			"node_status":        resp.Body.NodeStatus,
		}).Debug("replying with DSC action")
	}

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	if err := json.NewEncoder(w).Encode(&resp.Body); err != nil {
		m.log.WithError(err).Error("error encoding response")
		return
	}
}

func (m *Manager) sendReport(w http.ResponseWriter, r *http.Request) {
	agentId := regexpat.Param(r, "agent_id")
	var body types.SendReportRequestBody
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "error decoding body: %s", err)
		return
	}

	_, err := m.report.SendReport(
		r.Context(),
		types.SendReportRequest{
			AgentID: agentId,
			Body:    body,
		},
	)
	if err != nil {
		m.log.WithError(err).Error("error saving report")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "error saving report: %s", err)
		return
	}

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	fmt.Fprintf(w, `{"value":"SavedReport"}`)
}

func (m *Manager) getReports(w http.ResponseWriter, r *http.Request) {
	agentId := regexpat.Param(r, "agent_id")
	jobId := regexpat.Param(r, "job_id")

	resp, err := m.report.GetReports(
		r.Context(),
		types.GetReportsRequest{
			AgentID: agentId,
			JobID:   jobId,
		},
	)
	if err != nil {
		switch v := err.(type) {
		case types.ReportNotFoundError:
			w.WriteHeader(http.StatusNotFound)
			fmt.Fprintf(w, "%s", v)

		default:
			m.log.WithError(err).Error("error getting report")

			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "error getting report: %s", err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")

	// The returned value is already JSON-encoded, so we can use
	// json.RawMessage to ensure that we don't need to decode and re-encode
	// it, and instead return it as-is.
	//
	// TODO: validate that this is the right response format?
	var jsonResponse struct {
		Value []json.RawMessage `json:"value"`
	}
	jsonResponse.Value = []json.RawMessage{json.RawMessage(resp.Response)}

	if err := json.NewEncoder(w).Encode(&jsonResponse); err != nil {
		m.log.WithError(err).Errorf("error copying response body")
		return
	}
}

func (m *Manager) methodNotSupported(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotImplemented)

	// TODO: real error
	w.Header().Set("Content-Type", "text/plain")
	fmt.Fprint(w, "method not supported")
}
