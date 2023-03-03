// Copyright 2016-2023, Pulumi Corporation.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"bytes"
	"encoding/json"
	"log"
	"os"

	codegenrpc "github.com/pulumi/pulumi/sdk/v3/proto/go/codegen"
	"google.golang.org/protobuf/encoding/protojson"
)

func main() {
	coreSchema := getCoreSchema()

	var rawJSON []byte
	var err error
	rawJSON, err = protojson.Marshal(coreSchema)
	if err != nil {
		log.Fatalf("cannot marshal proto message to json: %v", err)
	}

	var formattedJSON bytes.Buffer
	err = json.Indent(&formattedJSON, rawJSON, "", "  ")
	if err != nil {
		log.Fatalf("failed to format core JSON: %v", err)
	}

	err = os.WriteFile("../core.json", formattedJSON.Bytes(), 0o755)
	if err != nil {
		log.Fatalf("failed to write core.json: %v", err)
	}
}

func getCoreSchema() *codegenrpc.Core {
	unitType := &codegenrpc.TypeReference{
		Element: &codegenrpc.TypeReference_Primitive{Primitive: codegenrpc.PrimitiveType_TYPE_UNIT},
	}
	boolType := &codegenrpc.TypeReference{
		Element: &codegenrpc.TypeReference_Primitive{Primitive: codegenrpc.PrimitiveType_TYPE_BOOL},
	}
	byteType := &codegenrpc.TypeReference{
		Element: &codegenrpc.TypeReference_Primitive{Primitive: codegenrpc.PrimitiveType_TYPE_BYTE},
	}
	stringType := &codegenrpc.TypeReference{
		Element: &codegenrpc.TypeReference_Primitive{Primitive: codegenrpc.PrimitiveType_TYPE_STRING},
	}
	propertyValueType := &codegenrpc.TypeReference{
		Element: &codegenrpc.TypeReference_Primitive{Primitive: codegenrpc.PrimitiveType_TYPE_PROPERTY_VALUE},
	}
	durationType := &codegenrpc.TypeReference{
		Element: &codegenrpc.TypeReference_Primitive{Primitive: codegenrpc.PrimitiveType_TYPE_DURATION},
	}
	propertyMapType := &codegenrpc.TypeReference{
		Element: &codegenrpc.TypeReference_Map{Map: propertyValueType},
	}

	makeRef := func(ref string) *codegenrpc.TypeReference {
		return &codegenrpc.TypeReference{
			Element: &codegenrpc.TypeReference_Ref{Ref: ref},
		}
	}

	return &codegenrpc.Core{
		Sdk: &codegenrpc.SDK{
			TypeDeclarations: []*codegenrpc.TypeDeclaration{
				{
					Element: &codegenrpc.TypeDeclaration_Enumeration{
						Enumeration: &codegenrpc.Enumeration{
							Name:        "pulumi.experimental.providers.log_severity",
							Description: "The severity level of a log message. Errors are fatal; all others are informational.",
							Values: []*codegenrpc.EnumerationValue{
								{
									Name:        "debug",
									Description: "A debug-level message not displayed to end-users (the default).",
									Value:       0,
								},
								{
									Name:        "info",
									Description: "An informational message printed to output during resource operations.",
									Value:       1,
								},
								{
									Name:        "warning",
									Description: "A warning to indicate that something went wrong.",
									Value:       2,
								},
								{
									Name:        "error",
									Description: "A fatal error indicating that the tool should stop processing subsequent resource operations.",
									Value:       3,
								},
							},
						},
					},
				},
				{
					Element: &codegenrpc.TypeDeclaration_Record{
						Record: &codegenrpc.Record{
							Name:        "pulumi.experimental.providers.log_message",
							Description: "A log message to be sent to the Pulumi engine.",
							Properties: []*codegenrpc.Property{
								{
									Name:        "severity",
									Description: "The logging level of this message.",
									Type:        makeRef("pulumi.experimental.providers.log_severity"),
								},
								{
									Name:        "message",
									Description: "The contents of the logged message.",
									Type:        stringType,
								},
								{
									Name:        "URN",
									Description: "The (optional) resource urn this log is associated with.",
									Type:        stringType,
								},
								{
									Name: "stream_id",
									Description: "The (optional) stream id that a stream of log messages can be associated with. This allows" +
										" clients to not have to buffer a large set of log messages that they all want to be" +
										" conceptually connected.  Instead the messages can be sent as chunks (with the same stream id)" +
										" and the end display can show the messages as they arrive, while still stitching them together" +
										" into one total log message. 0 means do not associate with any stream.",
									Type: stringType,
								},
								{
									Name:        "ephemeral",
									Description: "Optional value indicating whether this is a status message.",
									Type:        boolType,
								},
							},
						},
					},
				},
				{
					Element: &codegenrpc.TypeDeclaration_Interface{
						Interface: &codegenrpc.Interface{
							Name:        "pulumi.experimental.providers.host",
							Description: "An interface to the engine host running this plugin.",
							Methods: []*codegenrpc.Method{
								{
									Name:        "log",
									Description: "Send a log message to the host.",
									ReturnType:  unitType,
									Parameters: []*codegenrpc.Parameter{
										{
											Name: "message",
											Type: makeRef("pulumi.experimental.providers.log_message"),
										},
									},
								},
							},
						},
					},
				},
				{
					Element: &codegenrpc.TypeDeclaration_Record{
						Record: &codegenrpc.Record{
							Name:        "pulumi.experimental.providers.check_request",
							Description: "A request to validate the inputs for a resource.",
							Properties: []*codegenrpc.Property{
								{
									Name:        "URN",
									Description: "The Pulumi URN for this resource.",
									Type:        stringType,
								},
								{
									Name:        "type",
									Description: "The Pulumi type for this resource.",
									Type:        stringType,
								},
								{
									Name:        "name",
									Description: "The Pulumi name for this resource.",
									Type:        stringType,
								},
								{
									Name:        "olds",
									Description: "The old Pulumi inputs for this resource, if any.",
									Type:        propertyMapType,
								},
								{
									Name:        "news",
									Description: "The new Pulumi inputs for this resource.",
									Type:        propertyMapType,
								},
								{
									Name:        "random_seed",
									Description: "A deterministically random hash, primarily intended for global unique naming.",
									Type: &codegenrpc.TypeReference{
										Element: &codegenrpc.TypeReference_Array{Array: byteType},
									},
								},
							},
						},
					},
				},
				{
					Element: &codegenrpc.TypeDeclaration_Record{
						Record: &codegenrpc.Record{
							Name:        "pulumi.experimental.providers.check_failure",
							Description: "Indicates that a call to check failed; it contains the property and reason for the failure.",
							Properties: []*codegenrpc.Property{
								{
									Name:        "property",
									Description: "The property that failed validation.",
									Type:        stringType,
								},
								{
									Name:        "reason",
									Description: "The reason that the property failed validation.",
									Type:        stringType,
								},
							},
						},
					},
				},
				{
					Element: &codegenrpc.TypeDeclaration_Record{
						Record: &codegenrpc.Record{
							Name:        "pulumi.experimental.providers.check_response",
							Description: "The response from checking a resource.",
							Properties: []*codegenrpc.Property{
								{
									Name:        "inputs",
									Description: "The provider inputs for this resource.",
									Type:        propertyMapType,
								},
								{
									Name:        "failures",
									Description: "Any validation failures that occurred.",
									Type: &codegenrpc.TypeReference{
										Element: &codegenrpc.TypeReference_Array{Array: makeRef("pulumi.experimental.providers.check_failure")},
									},
								},
							},
						},
					},
				},
				{
					Element: &codegenrpc.TypeDeclaration_Record{
						Record: &codegenrpc.Record{
							Name:        "pulumi.experimental.providers.diff_request",
							Description: "A request to diff the inputs for a resource.",
							Properties: []*codegenrpc.Property{
								{
									Name:        "id",
									Description: "The ID of the resource to diff.",
									Type:        stringType,
								},
								{
									Name:        "URN",
									Description: "The Pulumi URN for this resource.",
									Type:        stringType,
								},
								{
									Name:        "type",
									Description: "The Pulumi type for this resource.",
									Type:        stringType,
								},
								{
									Name:        "name",
									Description: "The Pulumi name for this resource.",
									Type:        stringType,
								},
								{
									Name:        "olds",
									Description: "The old values of provider inputs to diff.",
									Type:        propertyMapType,
								},
								{
									Name:        "news",
									Description: "The new values of provider inputs to diff.",
									Type:        propertyMapType,
								},
								{
									Name:        "ignore_changes",
									Description: "A set of property paths that should be treated as unchanged.",
									Type: &codegenrpc.TypeReference{
										Element: &codegenrpc.TypeReference_Array{Array: stringType},
									},
								},
							},
						},
					},
				},
				{
					Element: &codegenrpc.TypeDeclaration_Record{
						Record: &codegenrpc.Record{
							Name:        "pulumi.experimental.providers.diff_response",
							Description: "The response from diffing a resource.",
							Properties: []*codegenrpc.Property{
								{
									Name:        "id",
									Description: "The ID of the resource to diff.",
									Type:        stringType,
								},
							},
						},
					},
				},
				{
					Element: &codegenrpc.TypeDeclaration_Record{
						Record: &codegenrpc.Record{
							Name:        "pulumi.experimental.providers.configure_request",
							Description: "A request to configure a provider.",
							Properties:  []*codegenrpc.Property{},
						},
					},
				},
				{
					Element: &codegenrpc.TypeDeclaration_Record{
						Record: &codegenrpc.Record{
							Name:        "pulumi.experimental.providers.configure_response",
							Description: "",
							Properties:  []*codegenrpc.Property{},
						},
					},
				},
				{
					Element: &codegenrpc.TypeDeclaration_Record{
						Record: &codegenrpc.Record{
							Name:        "pulumi.experimental.providers.create_request",
							Description: "A request to create a resource.",
							Properties: []*codegenrpc.Property{
								{
									Name:        "URN",
									Description: "The Pulumi URN for this resource.",
									Type:        stringType,
								},
								{
									Name:        "type",
									Description: "The Pulumi type for this resource.",
									Type:        stringType,
								},
								{
									Name:        "name",
									Description: "The Pulumi name for this resource.",
									Type:        stringType,
								},
								{
									Name:        "properties",
									Description: "The provider inputs to set during creation.",
									Type:        propertyMapType,
								},
								{
									Name:        "timeout",
									Description: "The create request timeout.",
									Type:        durationType,
								},
								{
									Name:        "preview",
									Description: "true if this is a preview and the provider should not actually create the resource.",
									Type:        boolType,
								},
							},
						},
					},
				},
				{
					Element: &codegenrpc.TypeDeclaration_Record{
						Record: &codegenrpc.Record{
							Name:        "pulumi.experimental.providers.create_response",
							Description: "The response from creating a resource.",
							Properties: []*codegenrpc.Property{
								{
									Name:        "id",
									Description: "The ID of the created resource.",
									Type:        stringType,
								},
								{
									Name:        "properties",
									Description: "Any any properties that were computed during creation.",
									Type:        propertyMapType,
								},
							},
						},
					},
				},
				{
					Element: &codegenrpc.TypeDeclaration_Record{
						Record: &codegenrpc.Record{
							Name:        "pulumi.experimental.providers.read_request",
							Description: "A request to read a resource.",
							Properties: []*codegenrpc.Property{
								{
									Name:        "id",
									Description: "The ID of the resource to read.",
									Type:        stringType,
								},
								{
									Name:        "URN",
									Description: "The Pulumi URN for this resource.",
									Type:        stringType,
								},
								{
									Name:        "type",
									Description: "The Pulumi type for this resource.",
									Type:        stringType,
								},
								{
									Name:        "name",
									Description: "The Pulumi name for this resource.",
									Type:        stringType,
								},
								{
									Name:        "properties",
									Description: "The current state (sufficiently complete to identify the resource).",
									Type:        propertyMapType,
								},
								{
									Name:        "inputs",
									Description: "The current inputs, if any (only populated during refresh).",
									Type:        propertyMapType,
								},
							},
						},
					},
				},
				{
					Element: &codegenrpc.TypeDeclaration_Record{
						Record: &codegenrpc.Record{
							Name:        "pulumi.experimental.providers.read_response",
							Description: "The response from reading a resource.",
							Properties: []*codegenrpc.Property{
								{
									Name:        "id",
									Description: "The ID of the resource read back (or empty if missing).",
									Type:        stringType,
								},
								{
									Name:        "properties",
									Description: "The state of the resource read from the live environment.",
									Type:        propertyMapType,
								},
								{
									Name:        "inputs",
									Description: "The inputs for this resource that would be returned from Check.",
									Type:        propertyMapType,
								},
							},
						},
					},
				},
				{
					Element: &codegenrpc.TypeDeclaration_Record{
						Record: &codegenrpc.Record{
							Name:        "pulumi.experimental.providers.update_request",
							Description: "A request to update a resource.",
							Properties: []*codegenrpc.Property{
								{
									Name:        "id",
									Description: "The ID of the resource to update.",
									Type:        stringType,
								},
								{
									Name:        "URN",
									Description: "The Pulumi URN for this resource.",
									Type:        stringType,
								},
								{
									Name:        "type",
									Description: "The Pulumi type for this resource.",
									Type:        stringType,
								},
								{
									Name:        "name",
									Description: "The Pulumi name for this resource.",
									Type:        stringType,
								},
								{
									Name:        "old",
									Description: "The old values of provider inputs for the resource to update.",
									Type:        propertyMapType,
								},
								{
									Name:        "new",
									Description: "The new values of provider inputs for the resource to update.",
									Type:        propertyMapType,
								},
								{
									Name:        "timeout",
									Description: "The create request timeout.",
									Type:        durationType,
								},
								{
									Name:        "ignore_changes",
									Description: "A set of property paths that should be treated as unchanged.",
									Type: &codegenrpc.TypeReference{
										Element: &codegenrpc.TypeReference_Array{Array: stringType},
									},
								},
								{
									Name:        "preview",
									Description: "true if this is a preview and the provider should not actually update the resource.",
									Type:        boolType,
								},
							},
						},
					},
				},
				{
					Element: &codegenrpc.TypeDeclaration_Record{
						Record: &codegenrpc.Record{
							Name:        "pulumi.experimental.providers.update_response",
							Description: "",
							Properties: []*codegenrpc.Property{
								{
									Name:        "properties",
									Description: "Any properties that were computed during updating.",
									Type:        propertyMapType,
								},
							},
						},
					},
				},
				{
					Element: &codegenrpc.TypeDeclaration_Record{
						Record: &codegenrpc.Record{
							Name:        "pulumi.experimental.providers.delete_request",
							Description: "A request to delete a resource.",
							Properties: []*codegenrpc.Property{
								{
									Name:        "id",
									Description: "The ID of the resource to delete.",
									Type:        stringType,
								},
								{
									Name:        "URN",
									Description: "The Pulumi URN for this resource.",
									Type:        stringType,
								},
								{
									Name:        "type",
									Description: "The Pulumi type for this resource.",
									Type:        stringType,
								},
								{
									Name:        "name",
									Description: "The Pulumi name for this resource.",
									Type:        stringType,
								},
							},
						},
					},
				},
				{
					Element: &codegenrpc.TypeDeclaration_Record{
						Record: &codegenrpc.Record{
							Name:        "pulumi.experimental.providers.delete_response",
							Description: "",
							Properties:  []*codegenrpc.Property{},
						},
					},
				},
				{
					Element: &codegenrpc.TypeDeclaration_Interface{
						Interface: &codegenrpc.Interface{
							Name: "pulumi.experimental.providers.provider",
							Description: "Provider presents a simple interface for orchestrating resource create, read, update, and delete operations. Each" +
								" provider understands how to handle all of the resource types within a single package.\n" +
								"\n" +
								"It is important to note that provider operations are not transactional (Some providers might decide to offer" +
								" transactional semantics, but such a provider is a rare treat). As a result, failures in the operations below can" +
								" range from benign to catastrophic (possibly leaving behind a corrupt resource). It is up to the provider to make a" +
								" best effort to ensure catastrophes do not occur. The errors returned from mutating operations indicate both the" +
								" underlying error condition in addition to a bit indicating whether the operation was successfully rolled back.",
							Methods: []*codegenrpc.Method{
								{
									Name:        "check_config",
									Description: "Validates the configuration for this resource provider.",
									ReturnType:  makeRef("pulumi.experimental.providers.check_response"),
									Parameters: []*codegenrpc.Parameter{
										{
											Name: "request",
											Type: makeRef("pulumi.experimental.providers.check_request"),
										},
									},
								},
								{
									Name:        "diff_config",
									Description: "Checks what impacts a hypothetical change to this provider's configuration will have on the provider.",
									ReturnType:  makeRef("pulumi.experimental.providers.diff_response"),
									Parameters: []*codegenrpc.Parameter{
										{
											Name: "request",
											Type: makeRef("pulumi.experimental.providers.diff_request"),
										},
									},
								},
								{
									Name:        "configure",
									Description: "Configures the resource provider with \"globals\" that control its behavior.",
									ReturnType:  makeRef("pulumi.experimental.providers.configure_response"),
									Parameters: []*codegenrpc.Parameter{
										{
											Name: "request",
											Type: makeRef("pulumi.experimental.providers.configure_request"),
										},
									},
								},
								{
									Name: "check",
									Description: "Validates that the given property bag is valid for a resource of the given type and returns the inputs" +
										" that should be passed to successive calls to Diff, Create, or Update for this resource. As a rule, the provider" +
										" inputs returned by a call to Check should preserve the original representation of the properties as present in" +
										" the program inputs. Though this rule is not required for correctness, violations thereof can negatively impact" +
										" the end-user experience, as the provider inputs are using for detecting and rendering diffs.",
									ReturnType: makeRef("pulumi.experimental.providers.check_response"),
									Parameters: []*codegenrpc.Parameter{
										{
											Name: "request",
											Type: makeRef("pulumi.experimental.providers.check_request"),
										},
									},
								},
								{
									Name:        "diff",
									Description: "Checks what impacts a hypothetical update will have on the resource's properties.",
									ReturnType:  makeRef("pulumi.experimental.providers.diff_response"),
									Parameters: []*codegenrpc.Parameter{
										{
											Name: "request",
											Type: makeRef("pulumi.experimental.providers.diff_request"),
										},
									},
								},
								{
									Name: "create",
									Description: "Allocates a new instance of the provided resource and returns its unique ID afterwards. (The input ID" +
										" must be blank.)  If this call fails, the resource must not have been created (i.e., it is \"transactional\").",
									ReturnType: makeRef("pulumi.experimental.providers.create_response"),
									Parameters: []*codegenrpc.Parameter{
										{
											Name: "request",
											Type: makeRef("pulumi.experimental.providers.create_request"),
										},
									},
								},
								{
									Name:        "update",
									Description: "Updates an existing resource with new values.",
									ReturnType:  makeRef("pulumi.experimental.providers.update_response"),
									Parameters: []*codegenrpc.Parameter{
										{
											Name: "request",
											Type: makeRef("pulumi.experimental.providers.update_request"),
										},
									},
								},
								{
									Name:        "delete",
									Description: "Tears down an existing resource with the given ID. If it fails, the resource is assumed to still exist.",
									ReturnType:  makeRef("pulumi.experimental.providers.delete_response"),
									Parameters: []*codegenrpc.Parameter{
										{
											Name: "request",
											Type: makeRef("pulumi.experimental.providers.delete_request"),
										},
									},
								},
								{
									Name: "read",
									Description: "Reads the current live state associated with a resource. Enough state must be include in the inputs to uniquely" +
										" identify the resource; this is typically just the resource ID, but may also include some properties.",
									ReturnType: makeRef("pulumi.experimental.providers.read_response"),
									Parameters: []*codegenrpc.Parameter{
										{
											Name: "request",
											Type: makeRef("pulumi.experimental.providers.read_request"),
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}
