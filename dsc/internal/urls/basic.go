package urls

// This file contains regex definitions for basic types.

const (
	// 2.2.2.4 ConfigurationName
	// The ConfigurationName header field SHOULD<1> be used in the request
	// message sent to the server as part of a GET request for the
	// configuration.
	//
	//     ConfigurationName = "ConfigurationName:" DQUOTE Configuration-Namevalue DQUOTE CRLF
	//     Configuration-Namevalue = Element *(Element)
	//     Element = DIGIT / ALPHA
	//
	// Example: "ConfigurationName":"SubPart1"
	ConfigurationName = `[a-zA-Z0-9]+`

	// 2.2.3.1 ConfigurationId
	// The ConfigurationId parameter is a universally unique identifier
	// (UUID) as specified in [RFC4122] section 3.
	ConfigurationId = `[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}`

	// 2.2.3.2 ModuleName
	// The ModuleName parameter is a string that is used by the server to
	// identify a specific module.
	//
	//     MODULENAME = Element *(Element)
	//     Element = DIGIT / ALPHA / "_"
	ModuleName = `[a-zA-Z0-9_]+`

	// 2.2.3.3 ModuleVersion
	// The ModuleVersion parameter identifies the version of a module. It
	// can be either an empty string or a string containing two to four
	// groups of digits where the groups are separated by a period.
	//
	//     MODULEVERSION = SQUOTE MULTIDIGIT "." MULTIDIGIT SQUOTE
	//     / SQUOTE MULTIDIGIT "." MULTIDIGIT "." MULTIDIGIT SQUOTE
	//
	//     / SQUOTE MULTIDIGIT "." MULTIDIGIT "." MULTIDIGIT "." MULTIDIGIT SQUOTE
	//     / SQUOTE SQUOTE (NULL character)
	//     MULTIDIGIT = DIGIT *[DIGIT]
	ModuleVersion = `([0-9]+(\.[0-9]+){1,3}|)`

	// 3.5.5.1.1 (GetStatusReport)
	//     JOBID = UUID; as specified in [RFC4122]
	JobId = `[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}`

	// 3.6.5.2 (GetConfiguration v2)
	//     AGENTID = UUID; as specified in [RFC4122]
	AgentId = `[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{4}-[a-fA-F0-9]{12}`
)
