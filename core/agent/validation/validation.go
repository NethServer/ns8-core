/*
 * Copyright (C) 2021 Nethesis S.r.l.
 * http://www.nethesis.it - nethserver@nethesis.it
 *
 * This script is part of NethServer.
 *
 * NethServer is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License,
 * or any later version.
 *
 * NethServer is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with NethServer.  If not, see COPYING.
 */

package validation

import (
	"encoding/json"
	"github.com/xeipuuv/gojsonschema"
	"strings"
	"io/ioutil"
	"os"
	"fmt"
)

type ValidationError struct {
	Parameter string      `json:"parameter"`
	Error     string      `json:"error"`
	Value     interface{} `json:"value"`
	Field     string      `json:"field"`
}

func ToJSON(errorList []gojsonschema.ResultError) ([]byte, error) {
	errors := make([]ValidationError, len(errorList))

	// Convert gojsonschema errors to our format
	for idx, resError := range errorList {
		field := resError.Field()
		parameter := field
		if dotPos := strings.IndexByte(field, '.'); dotPos >= 0 {
			parameter = field[:dotPos]
		}
		errors[idx] = ValidationError{
			Parameter: parameter,
			Error:     parameter + "_" + resError.Type(),
			Value:     resError.Value(),
			Field:     field,
		}
	}
	return json.Marshal(errors)
}

func ValidateGoStruct(schemaPath string, data interface{}) ([]gojsonschema.ResultError, error) {


	documentLoader := gojsonschema.NewGoLoader(data)
	schemaLoader := gojsonschema.NewSchemaLoader()

	// Read JSON Schema definitions from a well-known path. Ignore read errors.
	libraryPath := os.Getenv("AGENT_INSTALL_DIR") + "/validator-definitions.json"
	if schemaData, err := ioutil.ReadFile(libraryPath) ; err == nil {
		loader := gojsonschema.NewStringLoader(string(schemaData))
		if err := schemaLoader.AddSchemas(loader) ; err != nil {
			return nil, fmt.Errorf("Schema loader error while loding %s: %w", libraryPath, err)
		}
	}

	// Read the JSON Schema from the given path. This must succeed.
	schemaData, readFileError := ioutil.ReadFile(schemaPath)
	if readFileError != nil {
		return nil, fmt.Errorf("Failed to read schema file %s: %w", schemaPath, readFileError)
	}

	loader := gojsonschema.NewStringLoader(string(schemaData))
	schema, compileError := schemaLoader.Compile(loader)
	if compileError != nil {
		return nil, fmt.Errorf("Schema compile error for file %s: %w", schemaPath, compileError)
	}

	// Validate "data" against the JSON Schema documents
	result, validateError := schema.Validate(documentLoader)
	if validateError != nil {
		return nil, fmt.Errorf("Schema validate error: %w", validateError)
	}

	// Validation has completed, return the errors
	return result.Errors(), nil
}

func ValidateJsonString(schemaPath string, data []byte) ([]gojsonschema.ResultError, error) {
	var ddata interface{}
	err := json.Unmarshal(data, &ddata)
	if err != nil {
		return nil, fmt.Errorf("JSON unmarshal error: %w. Input data: %v", err, data)
	}
	return ValidateGoStruct(schemaPath, ddata)
}
