package webhook_docs

import (
	"encoding/json"
	"fmt"
	"os"

	server_model "github.com/Rfluid/whatsapp-webhook-server/src/server/model"
	webhook_service "github.com/Rfluid/whatsapp-webhook-server/src/webhook/service"
)

// Swagger represents the structure of the Swagger JSON file.
type Swagger struct {
	Paths map[string]interface{} `json:"paths"`
}

// Updates the generated docs.
func UpdateDocs(server server_model.Config, hook webhook_service.Config, docsPath *string) error {
	if docsPath == nil {
		defaultPath := "./docs/swagger.json"
		docsPath = &defaultPath
	}
	//
	// Load the generated Swagger file
	file, err := os.ReadFile(*docsPath)
	if err != nil {
		return err
	}

	// Parse the JSON data
	var swaggerData Swagger
	err = json.Unmarshal(file, &swaggerData)
	if err != nil {
		return err
	}

	// Define the old and new paths
	oldPath := "/{webhook_path}"
	newPath := fmt.Sprintf("%s%s", server.Path, hook.Path)

	// Update the path
	if pathData, exists := swaggerData.Paths[oldPath]; exists {
		swaggerData.Paths[newPath] = pathData
		delete(swaggerData.Paths, oldPath)
	} else {
		return fmt.Errorf("old path %s not found in Swagger JSON", oldPath)
	}

	// Convert the updated data back to JSON
	updatedFile, err := json.MarshalIndent(swaggerData, "", "  ")
	if err != nil {
		return err
	}

	// Save the updated JSON back to the file
	err = os.WriteFile(*docsPath, updatedFile, os.ModePerm)
	if err != nil {
		return err
	}

	fmt.Println("Swagger file updated successfully.")

	return nil
}
