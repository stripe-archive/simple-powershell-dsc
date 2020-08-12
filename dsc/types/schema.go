package types

// This file contains types that were generated from the JSON schema provided
// in the specification document, with some hand edits to make the types nicer.
//
// Schemas can be found in the ../schemas folder

type GetDscActionRequestBody struct {
	ClientStatus []ClientStatusItem `json:"ClientStatus,omitempty"`
}

type ClientStatusItem struct {
	Checksum          string `json:"Checksum,omitempty"`
	ChecksumAlgorithm string `json:"ChecksumAlgorithm,omitempty"`
	ConfigurationName string `json:"ConfigurationName,omitempty"`
}

type GetDscActionResponseBody struct {
	Details    []GetDscActionResponseBodyDetail `json:"Details,omitempty"`
	NodeStatus string                           `json:"NodeStatus,omitempty"`
}

type GetDscActionResponseBodyDetail struct {
	ConfigurationName string `json:"ConfigurationName,omitempty"`
	Status            string `json:"Status,omitempty"`
}

type CertificateInformation struct {
	FriendlyName *string `json:"FriendlyName,omitempty"`
	Issuer       *string `json:"Issuer,omitempty"`
	NotAfter     *string `json:"NotAfter,omitempty"`
	NotBefore    *string `json:"NotBefore,omitempty"`
	PublicKey    *string `json:"PublicKey,omitempty"`
	Subject      *string `json:"Subject,omitempty"`
	Thumbprint   *string `json:"Thumbprint,omitempty"`
	Version      *int    `json:"Version,omitempty"`
}

type RegisterDscAgentRequestBody struct {
	AgentInformation        RegisterAgentInformation `json:"AgentInformation,omitempty"`
	ConfigurationNames      []string                 `json:"ConfigurationNames,omitempty"`
	RegistrationInformation RegistrationInformation  `json:"RegistrationInformation,omitempty"`
}

type RegisterAgentInformation struct {
	IPAddress  *string `json:"IPAddress,omitempty"`
	LCMVersion *string `json:"LCMVersion,omitempty"`
	NodeName   *string `json:"NodeName,omitempty"`
}

type RegistrationInformation struct {
	RegistrationMessageType *string                `json:"RegistrationMessageType,omitempty"`
	CertificateInformation  CertificateInformation `json:"CertificateInformation,omitempty"`
}

type RotateCertificateRequestBody struct {
	RotationInformation RotationInformation `json:"RotationInformation,omitempty"`
}

type RotationInformation struct {
	CertificateInformation CertificateInformation `json:"CertificateInformation,omitempty"`
}
