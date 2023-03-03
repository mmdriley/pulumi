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
	"fmt"
	"html/template"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/text/cases"
	"golang.org/x/text/language"

	codegenrpc "github.com/pulumi/pulumi/sdk/v3/proto/go/codegen"
	"google.golang.org/protobuf/encoding/protojson"
)

func format(fullname string) {
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get current working directory: %v", err)
	}

	abs, err := filepath.Abs(fullname)
	if err != nil {
		log.Fatalf("failed to get absolute path for %q: %v", fullname, err)
	}

	parent := filepath.Dir(cwd)
	rel, err := filepath.Rel(parent, abs)
	if err != nil {
		log.Fatalf("failed to get relative path for %q from %q: %v", abs, parent, err)
	}

	gofmt := exec.Command("gofumpt", "-w", rel)
	gofmt.Dir = parent

	stderr, err := gofmt.StderrPipe()
	if err != nil {
		log.Fatalf("failed to pipe stderr from gofmt: %v", err)
	}
	go func() {
		_, err := io.Copy(os.Stderr, stderr)
		if err != nil {
			panic(fmt.Sprintf("unexpected error running gofmt: %v", err))
		}
	}()
	if err := gofmt.Run(); err != nil {
		log.Fatalf("failed to gofmt %v: %v", fullname, err)
	}
}

func allUpper(name string) bool {
	for _, c := range name {
		if c >= 'a' && c <= 'z' {
			return false
		}
	}
	return true
}

func pulumiNameToGoName(name string) string {
	// A name looks like "separate_word_AKA", that is each part is separated by "_" and acronyms are
	// uppercased.

	nameParts := strings.Split(name, "_")
	goName := ""
	titleCaser := cases.Title(language.English)
	for _, part := range nameParts {
		if allUpper(part) {
			goName += part
			continue
		}
		goName += titleCaser.String(part)
	}
	return goName
}

func pulumiTypeToGoType(typ *codegenrpc.TypeReference) string {
	switch e := typ.Element.(type) {
	case *codegenrpc.TypeReference_Primitive:
		switch e.Primitive {
		case codegenrpc.PrimitiveType_TYPE_UNIT:
			return ""
		case codegenrpc.PrimitiveType_TYPE_BOOL:
			return "bool"
		case codegenrpc.PrimitiveType_TYPE_BYTE:
			return "byte"
		case codegenrpc.PrimitiveType_TYPE_INT:
			return "int"
		case codegenrpc.PrimitiveType_TYPE_STRING:
			return "string"
		case codegenrpc.PrimitiveType_TYPE_PROPERTY_VALUE:
			return "resource.PropertyValue"
		case codegenrpc.PrimitiveType_TYPE_DURATION:
			return "time.Duration"
		}
	case *codegenrpc.TypeReference_Ref:
		parts := strings.Split(e.Ref, ".")
		return pulumiNameToGoName(parts[len(parts)-1])
	case *codegenrpc.TypeReference_Map:
		// Special case for PropertyMap.
		if e.Map.GetPrimitive() == codegenrpc.PrimitiveType_TYPE_PROPERTY_VALUE {
			return "resource.PropertyMap"
		}
		return "map[string]" + pulumiTypeToGoType(e.Map)
	case *codegenrpc.TypeReference_Array:
		return "[]" + pulumiTypeToGoType(e.Array)
	}

	log.Fatalf("unhandled type: %v", typ)
	return ""
}

// Generate the code for a Pulumi package.
func main() {
	coreSchemaJson, err := os.ReadFile("../../../proto/core.json")
	if err != nil {
		log.Fatalf("read core.json: %v", err)
	}

	var core codegenrpc.Core
	if err := protojson.Unmarshal(coreSchemaJson, &core); err != nil {
		log.Fatalf("parse core schema: %v", err)
	}

	templates, err := template.New("templates").ParseGlob("./templates/*")
	if err != nil {
		log.Fatalf("failed to parse templates: %v", err)
	}

	for _, typ := range core.Sdk.TypeDeclarations {
		log.Printf("Generating %v", typ)

		var name string
		var templateName string
		var data map[string]interface{}

		switch e := typ.Element.(type) {
		case *codegenrpc.TypeDeclaration_Record:
			name = e.Record.Name
			templateName = "record.go.template"

			properties := make([]interface{}, 0)
			for _, prop := range e.Record.Properties {
				templateProp := map[string]interface{}{
					"Name":        pulumiNameToGoName(prop.Name),
					"Description": strings.Split(prop.Description, "\n"),
					"Type":        pulumiTypeToGoType(prop.Type),
				}

				properties = append(properties, templateProp)
			}

			data = map[string]interface{}{
				"Description": strings.Split(e.Record.Description, "\n"),
				"Properties":  properties,
			}
		case *codegenrpc.TypeDeclaration_Enumeration:
			name = e.Enumeration.Name
			templateName = "enum.go.template"

			values := make([]interface{}, 0)
			for _, value := range e.Enumeration.Values {
				templateValue := map[string]interface{}{
					"Name":        pulumiNameToGoName(value.Name),
					"Description": strings.Split(value.Description, "\n"),
					"Value":       value.Value,
				}

				values = append(values, templateValue)
			}

			data = map[string]interface{}{
				"Description": strings.Split(e.Enumeration.Description, "\n"),
				"Values":      values,
			}
		case *codegenrpc.TypeDeclaration_Interface:
			name = e.Interface.Name
			templateName = "interface.go.template"

			methods := make([]interface{}, 0)
			for _, method := range e.Interface.Methods {
				params := make([]interface{}, 0)
				for _, param := range method.Parameters {
					templateParam := map[string]interface{}{
						"Name":        param.Name,
						"Description": strings.Split(param.Description, "\n"),
						"Type":        pulumiTypeToGoType(param.Type),
					}

					params = append(params, templateParam)
				}

				templateMethod := map[string]interface{}{
					"Name":        pulumiNameToGoName(method.Name),
					"Description": strings.Split(method.Description, "\n"),
					"ReturnType":  pulumiTypeToGoType(method.ReturnType),
					"Parameters":  params,
				}

				methods = append(methods, templateMethod)
			}

			data = map[string]interface{}{
				"Description": strings.Split(e.Interface.Description, "\n"),
				"Methods":     methods,
			}

		default:
			log.Fatalf("unexpected type declaration: %v", typ.Element)
		}

		parts := strings.Split(name, ".")

		data["Name"] = pulumiNameToGoName(parts[len(parts)-1])
		if len(parts) > 1 {
			data["Package"] = parts[len(parts)-2]
		} else {
			data["Package"] = "pulumi"
		}

		path := filepath.Join(parts[:len(parts)-1]...)
		fullname := filepath.Join("..", path, parts[len(parts)-1]+".go")
		log.Printf("Writing %v", fullname)
		if err := os.MkdirAll(filepath.Dir(fullname), 0755); err != nil {
			log.Fatalf("create directory %q: %v", filepath.Dir(fullname), err)
		}

		f, err := os.Create(fullname)
		if err != nil {
			log.Fatalf("create %q: %v", fullname, err)
		}
		if err := templates.ExecuteTemplate(f, templateName, data); err != nil {
			log.Fatalf("execute %q: %v", templateName, err)
		}
		f.Close()

		format(fullname)
	}
}
