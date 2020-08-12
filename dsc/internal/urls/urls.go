package urls

// Constants for URLs
const (
	// 3.1 GetConfiguration Versions 1.0 and 1.1
	GetConfigurationV1URL = `Action\(ConfigurationId='(?P<configuration_id>` + ConfigurationId + `)'\)\/ConfigurationContent`

	// 3.2 GetModule Versions 1.0 and 1.1
	GetModuleV1URL = `Module\(` +
		`ConfigurationId='(?P<configuration_id>` + ConfigurationId + `)',` +
		`ModuleName='(?P<module_name>` + ModuleName + `)',` +
		`ModuleVersion='(?P<module_name>` + ModuleVersion + `)'` +
		`\)\/ModuleContent`

	// 3.3 GetAction Versions 1.0 and 1.1
	GetActionV1URL = `Action\(ConfigurationId='(?P<configuration_id>` + ConfigurationId + `)'\)\/GetAction`

	// 3.4 SendStatusReport Versions 1.0 and 1.1
	SendStatusReportURL = `Node\(ConfigurationId='(?P<configuration_id>` + ConfigurationId + `)'\)\/SendStatusReport`

	// 3.5 GetStatusReport Versions 1.0 and 1.1
	GetStatusReportURL = `Node\(ConfigurationId='(?P<configuration_id>` + ConfigurationId + `)'\)\/` +
		`Reports\(JobId='(?P<job_id>` + JobId + `)'\)`

	// 3.6 GetConfiguration Version 2.0
	GetConfigurationV2URL = `Nodes\(AgentId='(?P<agent_id>` + AgentId + `)'\)\/` +
		`Configurations\(ConfigurationName='(?P<configuration_name>` + ConfigurationName + `)'\)\/` +
		`ConfigurationContent`

	// 3.7 GetModule Version 2.0
	GetModuleV2URL = `Modules\(ModuleName='(?P<module_name>` + ModuleName + `)',` +
		`ModuleVersion='(?P<module_version>` + ModuleVersion + `)'\)\/ModuleContent`

	// 3.8 GetDscAction Version 2.0
	GetDscActionV2URL = `Nodes\(AgentId='(?P<agent_id>` + AgentId + `)'\)\/GetDscAction`

	// 3.9 RegisterDscAgent Version 2
	RegisterDscAgentV2URL = `Nodes\(AgentId='(?P<agent_id>` + AgentId + `)'\)`

	// 3.10 SendReport Version 2.0
	SendReportV2URL = `Nodes\(AgentId='(?P<agent_id>` + AgentId + `)'\)\/SendReport`

	// 3.11 GetReports Version 2.0
	GetReportsV2URL = `Nodes\(AgentId='(?P<agent_id>` + AgentId + `)'\)\/` +
		`Reports\(JobId='(?P<job_id>` + JobId + `)'\)`

	// 3.12 CertificateRotation
	CertificateRotationURL = `Nodes\(AgentId='(?P<agent_id>` + AgentId + `)'\)\/CertificateRotation`
)
