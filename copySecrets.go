package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
	vault "github.com/hashicorp/vault/api"
)

func copySecretsFunction(keyVaultName string, vaultAddr string, vaultNamespace string, jsonFile string) {
	// Read the JSON file
	data, err := os.ReadFile(jsonFile)
	if err != nil {
		log.Fatalf("Error reading file: %v. Use --file flag if you would like to read from a different file besides secrets.json", err)
	}

	// Unmarshal the JSON data into the Secrets struct
	var secretsData Secrets
	err = json.Unmarshal(data, &secretsData)
	if err != nil {
		log.Fatalf("Error unmarshalling JSON: %v", err)
	}

	for _, individualSecret := range secretsData.Secrets {
		if !individualSecret.Copy {
			fmt.Printf("Skipping %v as copy: false\n", individualSecret.KeyVaultSecretName)
		} else {
			fmt.Printf("Retrieving secret from keyvault: %v\n", individualSecret.KeyVaultSecretName)
			secretValue, err := getSecretFromKeyVault(keyVaultName, individualSecret.KeyVaultSecretName)
			if err != nil {
				log.Fatalf("failed to get secret from keyvault: %v\n", err)
			}
			fmt.Printf("Writing secret to vault: %v/%v %v=REDACTED\n", individualSecret.VaultSecretPath, individualSecret.VaultSecretName, individualSecret.VaultSecretKey)
			err = writeSecretToVault(vaultAddr, vaultNamespace, individualSecret.VaultSecretPath, individualSecret.VaultSecretName, individualSecret.VaultSecretKey, secretValue)
			if err != nil {
				log.Fatalf("failed to write secret to Hashicorp Vault: %v\n", err)
			}
		}

	}
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

func writeSecretToVault(addr string, namespace string, mount string, name string, key string, value string) error {
	ctx := context.Background()
	config := vault.DefaultConfig()
	config.Address = addr
	client, err := vault.NewClient(config)
	if err != nil {
		log.Fatalf("unable to initialize Vault client: %v", err)
	}

	client.SetToken(os.Getenv("VAULT_TOKEN"))
	client.SetNamespace(namespace)

	secretData := map[string]interface{}{
		key: value,
	}

	// Check for existing secret

	_, err = client.KVv2(mount).Get(ctx, name)
	if err != nil {
		// Use put method to create secret if it doesn't exist
		_, err = client.KVv2(mount).Put(ctx, name, secretData)
		if err != nil {
			log.Fatalf("unable to write secret: %v", err)
		}
	} else {
		// Use patch method to update a secret if it already exists
		_, err = client.KVv2(mount).Patch(ctx, name, secretData)
		if err != nil {
			log.Fatalf("unable to write secret: %v", err)
		}
	}
	return err
}
