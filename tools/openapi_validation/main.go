package main

import (
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/pb33f/libopenapi"
	validator "github.com/pb33f/libopenapi-validator"
)

func main() {
	// 1. Load an OpenAPI Spec into bytes
	myAPI, err := os.ReadFile("../docs/openapi.yaml")

	if err != nil {
		log.Fatal(err.Error())
	}

	// 2. Create a new OpenAPI document using libopenapi
	document, docErrs := libopenapi.NewDocument(myAPI)

	if docErrs != nil {
		log.Fatal(docErrs.Error())
	}

	highLevelValidator, validatorErrs := validator.NewValidator(document)
	if validatorErrs != nil {
		for _, e := range validatorErrs {
			log.Printf("Failure: %s\n", e.Error())
		}
		log.Fatal("Validation failed")
	}

	// 4. Validate!
	valid, validationErrs := highLevelValidator.ValidateDocument()

	if !valid {
		for _, e := range validationErrs {
			// 5. Handle the error
			log.Printf("Type: %s, Failure: %s\n", e.ValidationType, e.Message)
			log.Printf("Fix: %s\n\n", e.HowToFix)
		}
		log.Fatal("Validation failed")
	}
	v3, err1 := document.BuildV3Model()
	if err1 != nil {
		for _, e := range err1 {
			log.Printf("Failure: %s\n", e.Error())
		}
		log.Fatal("Validation failed")
	}

	for _, component := range v3.Model.Components.Schemas.FromNewest() {
		schema, err := component.BuildSchema()

		if err != nil {
			log.Fatalf("Error: %s\n", err.Error())
		}
		if len(schema.Enum) > 0 {
			uniqueEnums := make(map[string]bool)
			for _, v := range schema.Enum {
				val := strings.ToLower(v.Value)
				if _, ok := uniqueEnums[val]; ok {
					log.Fatalf("Duplicate enum value: %s in %v\n", v.Value, component.GetSchemaKeyNode().Value)
				} else {
					uniqueEnums[val] = true
				}
			}
		}
	}
	for uri, path := range v3.Model.Paths.PathItems.FromNewest() {
		for operation := range path.GetOperations().ValuesFromNewest() {
			method := operation.GoLow().KeyNode.Value
			if operation.Responses == nil {
				log.Fatalf("No responses for operation %s\n", operation.OperationId)
			}
			for code, response := range operation.Responses.Codes.FromNewest() {
				codeInt, err := strconv.ParseInt(code, 10, 64)
				if err != nil {
					log.Fatalf("Error: %s\n", err.Error())
				}
				if codeInt < 200 || codeInt >= 300 {
					continue
				}
				if response.Content != nil {
					for contentType, mediaType := range response.Content.FromNewest() {
						if contentType == "application/json" {
							mediaTypeSchema, err := mediaType.Schema.BuildSchema()
							if err != nil {
								log.Fatalf("Error: %s\n", err.Error())
							}
							if len(mediaTypeSchema.Type) > 0 && mediaTypeSchema.Type[0] == "object" {
								if mediaType.Schema.GoLow().GetReference() == "" {
									log.Println(method, uri, code, operation.OperationId)
									log.Fatal("No reference for object item response")
								}
							}
							if len(mediaTypeSchema.Type) > 0 && mediaTypeSchema.Type[0] == "array" {
								arrayItemRef := mediaType.Schema.Schema().Items.A.GetReference()
								if arrayItemRef == "" {
									log.Println(method, uri, code, operation.OperationId)
									log.Fatal("No reference for array item response")
								}
							}

						}
					}
				}
			}
		}
	}
}
