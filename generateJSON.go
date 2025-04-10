package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
)

func generateJSONFunction(keyVaultName string, defaultMount string, defaultPath string, jsonFile string) {

	var retrivedList Secrets

	vaultURL := fmt.Sprintf("https://%s.vault.azure.net/", keyVaultName)

	// Create a new DefaultAzureCredential
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		log.Fatalf("Failed to obtain a credential: %v", err)
	}

	// Create a new client
	client, err := azsecrets.NewClient(vaultURL, cred, nil)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}

	// List secrets
	pager := client.NewListSecretPropertiesPager(nil)

	for pager.More() {
		page, err := pager.NextPage(context.Background())
		if err != nil {
			log.Fatalf("Failed to get next page of secrets: %v", err)
		}

		for _, secret := range page.Value {
			newSecret := Secret{KeyVaultSecretName: secret.ID.Name(), VaultSecretMount: defaultMount, VaultSecretPath: defaultPath, VaultSecretName: secret.ID.Name(), VaultSecretKey: "secret", Copy: false}
			retrivedList.Secrets = append(retrivedList.Secrets, newSecret)
		}
	}

	// Convert to JSON
	jsonData, err := json.MarshalIndent(retrivedList, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal secrets to json: %v", err)
	}

	// Write to file
	file, err := os.Create(jsonFile)
	if err != nil {
		log.Fatalf("Failed to create json file: %v. Use --file flag if you would like to write to a different file besides secrets.json", err)
	}

	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println("Error closing file:", err)
		}
	}()

	_, err = file.Write(jsonData)
	if err != nil {
		log.Fatalf("Failed to write json data to file: %v.", err)
	}

	fmt.Printf("List of secrets in Keyvault has been written to %v\n", jsonFile)
}
