// Run the gen command to generate a json file secrets.json listing all secrets
//
// go run main.go --kv=jst-awx --gen
//
// Update that secrets.json file:
// `vault_secret_location` - where you would like the secret to go in vault
// `vault_secret_key` - the key of the field within the secret (each secret can contain multiple key/value pairs)
// `copy` true if you would like it copied to vault
//
// Run the copy command to copy the secrets to vault:
//
// go run main.go --kv=jst-axw --copy
//

package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
)

func main() {
	// Define the flags
	keyVaultName := flag.String("kv", "", "The name of the Azure Key Vault.")
	defaultVaultLocation := flag.String("location", "", "The default location to place secrets in Hashicorp Vault (will default to secret/).")
	generateJSON := flag.Bool("gen", false, "Generate json file secrets.json with a list of secrets from KeyVault as keys.")
	copySecrets := flag.Bool("copy", false, "Run the function to copy the secrets from KeyVault to HashiCorp Vault based on the secrets.json locations.")
	flag.Parse()

	if *keyVaultName == "" {
		log.Fatalf("kv flag is required")
	}

	if *defaultVaultLocation == "" {
		*defaultVaultLocation = "secret/"
	}

	if *generateJSON {
		generateJSONFunction(*keyVaultName, *defaultVaultLocation)
	} else if *copySecrets {
		copySecretsFunction(*keyVaultName)
	} else {
		log.Fatalf("Either --gen or --copy flag must be specified")
	}
}

// Secret represents the structure of each secret in the JSON
type Secret struct {
	KeyVaultSecretName  string `json:"key_vault_secret_name"`
	KeyVaultSecretValue string `json:"key_vault_secret_value"`
	VaultSecretLocation string `json:"vault_secret_location"`
	VaultSecretKey      string `json:"vault_secret_key"`
	Copy                bool   `json:"copy"`
}

// Secrets represents the overall JSON structure
type Secrets struct {
	Secrets []Secret `json:"Secrets"`
}

func generateJSONFunction(keyVaultName string, defaultVaultLocation string) {

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
			newSecret := Secret{KeyVaultSecretName: secret.ID.Name(), KeyVaultSecretValue: "", VaultSecretLocation: defaultVaultLocation, VaultSecretKey: "secret", Copy: false}
			retrivedList.Secrets = append(retrivedList.Secrets, newSecret)
		}
	}

	// Convert to JSON
	jsonData, err := json.MarshalIndent(retrivedList, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal secrets to JSON: %v", err)
	}

	// Write to file
	file, err := os.Create("secrets.json")
	if err != nil {
		log.Fatalf("Failed to create JSON file: %v", err)
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		log.Fatalf("Failed to write JSON data to file: %v", err)
	}

	fmt.Println("List of secrets in Keyvault has been written to secrets.json")
}

func copySecretsFunction(keyVaultName string) {
	// Read the JSON file
	data, err := os.ReadFile("secrets.json")
	if err != nil {
		log.Fatalf("Error reading file: %v", err)
	}

	// Unmarshal the JSON data into the Secrets struct
	var secretsData Secrets
	err = json.Unmarshal(data, &secretsData)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}

	for i, individualSecret := range secretsData.Secrets {
		if !individualSecret.Copy {
			fmt.Printf("Skipping %v as copy: false\n", individualSecret.KeyVaultSecretName)
		} else {
			fmt.Printf("Retrieving %v\n", individualSecret.KeyVaultSecretName)
			secretsData.Secrets[i].KeyVaultSecretValue, err = getSecretFromKeyVault(keyVaultName, individualSecret.KeyVaultSecretName)
			if err != nil {
				log.Fatalf("failed to get secret: %v", err)
			}
		}

	}

	// // Initialize Vault client
	// config := vault.DefaultConfig()
	// client, err := vault.NewClient(config)
	// if err != nil {
	// 	log.Fatalf("Error creating Vault client: %v", err)
	// }

	// // Set Vault token (replace with your actual token)
	// client.SetToken(os.Getenv("VAULT_TOKEN"))

	// // Write secrets to Vault
	// for key, value := range secretsData.Secrets {
	// 	err := writeSecretToVault(client, "secret/data/"+key, value)
	// 	if err != nil {
	// 		log.Fatalf("Error writing secret to Vault: %v", err)
	// 	}
	// }

	// fmt.Println("Secrets written to Vault successfully!")

	// REMOVE LATER

	// Convert to JSON
	jsonData, err := json.MarshalIndent(secretsData, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal secrets to JSON: %v", err)
	}

	// Write to file
	file, err := os.Create("secrets2.json")
	if err != nil {
		log.Fatalf("Failed to create JSON file: %v", err)
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		log.Fatalf("Failed to write JSON data to file: %v", err)
	}

	fmt.Println("List of secrets in Keyvault has been written to secrets2.json")

	// REMOVE LATER
}

func getSecretFromKeyVault(keyVaultName string, secretName string) (string, error) {
	vaultURL := fmt.Sprintf("https://%s.vault.azure.net/", keyVaultName)

	// Create a credential using DefaultAzureCredential
	cred, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		return "", fmt.Errorf("failed to obtain a credential: %v", err)
	}

	// Create a SecretClient using the vault URL and credential
	client, err := azsecrets.NewClient(vaultURL, cred, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create secret client: %v", err)
	}

	// Retrieve the secret
	secretResp, err := client.GetSecret(context.Background(), secretName, "", nil)
	if err != nil {
		return "", fmt.Errorf("failed to get secret: %v", err)
	}

	return *secretResp.Value, nil
}

// func writeSecretToVault(client *vault.Client, path, value string) error {
// 	// Prepare the data to be written
// 	data := map[string]interface{}{
// 		"data": map[string]string{
// 			"value": value,
// 		},
// 	}

// 	// Write the data to the specified path
// 	_, err := client.Logical().Write(path, data)
// 	return err
// }
