package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	// Define flags
	keyVaultName := flag.String("kv", "", "The name of the Azure Key Vault.")
	vaultAddr := flag.String("vault_addr", "", "The url for vault (i.e. https://examplevault.com).")
	vaultNamespace := flag.String("vault_namespace", "", "The namespace for vault (i.e. https://examplevault.com).")
	defaultMount := flag.String("mount", "", "The path of the kvv2 mount (will default to secret).")
	generateJSON := flag.Bool("gen", false, "Generate json file secrets.json with a list of secrets from KeyVault as keys.")
	copySecrets := flag.Bool("copy", false, "Run the function to copy the secrets from KeyVault to HashiCorp Vault based on the secrets.json locations.")
	jsonFile := flag.String("file", "", "json file to write or read list of secrets from/to. Defaults to secrets.json in the current directory")
	flag.Parse()

	// Retrieve token from vault cli login (if available)
	if *copySecrets {
		vaultToken := os.Getenv("VAULT_TOKEN")
		if vaultToken == "" {
			userHomeDir, err := os.UserHomeDir()
			if err != nil {
				log.Fatal(err)
			}
			vaultTokenPath := filepath.Join(userHomeDir, ".vault-token")
			_, err = os.Stat(vaultTokenPath)
			if err == nil {
				content, err := os.ReadFile(vaultTokenPath)
				if err != nil {
					log.Fatal(err)
				}
				os.Setenv("VAULT_TOKEN", string(content))
				fmt.Printf("VAULT_TOKEN environment variable not set so using token found at %v.\n", vaultTokenPath)
			} else {
				log.Fatalf("VAULT_TOKEN environment variable not set and no login token found at %v so aborting. You will either need to set VAULT_TOKEN or run the vault login cli command.\n", vaultTokenPath)
			}
		} else {
			fmt.Println("Using VAULT_TOKEN environment variable to connect to vault. If you just ran vault login, you may want to unset VAULT_TOKEN before running this to use your login token from your home directory instead.")
		}
	}

	// Validate flags

	if *keyVaultName == "" {
		log.Fatalf("--kv flag is required")
	}

	if *copySecrets && *vaultAddr == "" && *vaultNamespace == "" {
		log.Fatalf("--vault_addr & --vault_namespace required to copy secrets")
	}

	// Set default variables

	if *defaultMount == "" {
		*defaultMount = "secret"
	}

	if *vaultNamespace == "" {
		*vaultNamespace = ""
	}

	if *jsonFile == "" {
		*jsonFile = "secrets.json"
	}

	// Call functions

	if *generateJSON {
		generateJSONFunction(*keyVaultName, *defaultMount, *jsonFile)
	} else if *copySecrets {
		copySecretsFunction(*keyVaultName, *vaultAddr, *vaultNamespace, *jsonFile)
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
