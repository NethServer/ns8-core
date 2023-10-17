/*
 * Copyright (C) 2023 Nethesis S.r.l.
 * SPDX-License-Identifier: GPL-3.0-or-later
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
 
 func toValidationErrorArray(errorList []gojsonschema.ResultError) []ValidationError {
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
	 return errors
 }
 
 func validateGoStruct(schemaPath string, data interface{}) ([]gojsonschema.ResultError, error) {
 
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
 
 func ValidatePayload(schemaPath string, data []byte) ([]ValidationError, error) {
	 var ddata interface{}
	 err := json.Unmarshal(data, &ddata)
	 if err != nil {
		 return nil, fmt.Errorf("JSON unmarshal error: %w. Input data: %v", err, data)
	 }
	 errorList, errInfo := validateGoStruct(schemaPath, ddata)
	 if errInfo != nil {
		return nil, errInfo
	 }

	 return toValidationErrorArray(errorList), nil
 }