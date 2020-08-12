package types

// This file is responsible for generating Go types from the
// Microsoft-published JSON schemas for requests and responses. These JSON
// schema files are extracted from "Appendix A: Full JSON Schema" in the
// published specification.
//
// The 'schematyper' utility can be found at:
//     https://github.com/idubinskiy/schematyper

//go:generate schematyper -o schema_GetDscActionRequest.go --package=types --root-type=GetDscActionRequestBody --prefix=GetDscActionRequestBody_ ../schemas/GetDscAction_request.json
//go:generate schematyper -o schema_GetDscActionResponse.go --package=types --root-type=GetDscActionResponseBody --prefix=GetDscActionResponseBody_ ../schemas/GetDscAction_response.json
//go:generate schematyper -o schema_RegisterDscAgentRequest.go --package=types --root-type=RegisterDscAgentRequestBody --prefix=RegisterDscAgentRequestBody_ ../schemas/RegisterDscAgent_request.json
//go:generate schematyper -o schema_RotateCertificateRequest.go --package=types --root-type=RotateCertificateRequestBody --prefix=RotateCertificateRequestBody_ ../schemas/RotateCertificate_request.json
