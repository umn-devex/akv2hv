package main

import (
	"flag"
	"log"
)

func main() {
	// Define flags
	keyVaultName := flag.String("kv", "", "The name of the Azure Key Vault.")
	vaultAddr := flag.String("vault_addr", "", "The url for vault (i.e. https://examplevault.com).")
	vaultNamespace := flag.String("vault_namespace", "", "The namespace for vault (i.e. https://examplevault.com).")
	defaultMount := flag.String("mount", "", "The path of the kvv2 mount (will default to secret).")
	generateJSON := flag.Bool("gen", false, "Generate json file secrets.json with a list of secrets from KeyVault as keys.")
	copySecrets := flag.Bool("copy", false, "Run the function to copy the secrets from KeyVault to HashiCorp Vault based on the secrets.json locations.")
	flag.Parse()

	// Validate flags

	if *keyVaultName == "" {
		log.Fatalf("--kv flag is required")
	}

	if *copySecrets == true && *vaultAddr == "" && *vaultNamespace == "" {
		log.Fatalf("--vault_addr & --vault_namespace required to copy secrets")
	}

	// Set default variables

	if *defaultMount == "" {
		*defaultMount = "secret"
	}

	if *vaultNamespace == "" {
		*vaultNamespace = ""
	}

	// Call functions

	if *generateJSON {
		generateJSONFunction(*keyVaultName, *defaultMount)
	} else if *copySecrets {
		copySecretsFunction(*keyVaultName, *vaultAddr, *vaultNamespace)
	} else {
		log.Fatalf("Either --gen or --copy flag must be specified")
	}
}

// Structure of each secret
type Secret struct {
	KeyVaultSecretName string `json:"key_vault_secret_name"`
	VaultSecretPath    string `json:"vault_secret_mount"`
	VaultSecretName    string `json:"vault_secret_name"`
	VaultSecretKey     string `json:"vault_secret_key"`
	Copy               bool   `json:"copy"`
}

// Overall structutre of secrets
type Secrets struct {
	Secrets []Secret `json:"secrets"`
}
