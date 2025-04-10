package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func main() {
	// Define flags
	keyVaultName := flag.String("kv", "", "The name of the Azure Key Vault.")
	vaultAddr := flag.String("vault_addr", "", "The url for vault (i.e. https://examplevault.com).")
	vaultNamespace := flag.String("vault_namespace", "", "The namespace for vault (i.e. admin/abc).")
	vaultToken := flag.String("vault_token", "", "Vault token, will override VAULT_TOKEN environment variable and vault login token file.")
	defaultMount := flag.String("default_mount", "", "Generate the json file with a default kvv2 mount.")
	defaultPath := flag.String("default_path", "", "Generate the json file with a default path of the secret including trailing slash.")
	defaultCopy := flag.Bool("default_copy", false, "Generate json file with copy: true as the default")
	generateJSON := flag.Bool("gen", false, "Generate json file secrets.json with a list of secrets from KeyVault as keys.")
	copySecrets := flag.Bool("copy", false, "Run the function to copy the secrets from KeyVault to HashiCorp Vault based on the secrets.json locations.")
	jsonFile := flag.String("file", "", "json file to write or read list of secrets from/to. Defaults to secrets.json in the current directory")
	flag.Parse()

	// Retrieve token from vault cli login (if available)
	var token string
	if *copySecrets {
		token_env := os.Getenv("VAULT_TOKEN")
		userHomeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatal(err)
		}
		vaultTokenPath := filepath.Join(userHomeDir, ".vault-token")
		_, pathErr := os.Stat(vaultTokenPath)

		switch {

		// First check for a vault token flag and use that token
		case *vaultToken != "":
			token = *vaultToken

		// Next check for a VAULT_TOKEN environment variable and use that token
		case token_env != "":
			token = token_env
			fmt.Println("Using VAULT_TOKEN environment variable to connect to vault. If you just ran vault login, you may want to unset VAULT_TOKEN before running this to use your login token from your home directory instead.")

		// Finally, check for the vault login token from the home directory and use that token
		case pathErr == nil:
			content, err := os.ReadFile(vaultTokenPath)
			if err != nil {
				log.Fatal(err)
			}
			token = string(content)
			fmt.Printf("Using vault login token found at %v.\n", vaultTokenPath)

		default:
			log.Fatalf("VAULT_TOKEN environment variable not set and no login token found at %v so aborting. You will either need to run the vault login cli command, set the VAULT_TOKEN environment variable, or pass it in with the --vault_token flag.\n", vaultTokenPath)
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

	// Append trailing / to path if it doesn't already include it

	if *defaultPath != "" && !strings.HasSuffix(*defaultPath, "/") {
		*defaultPath = *defaultPath + "/"
	}

	// Call functions

	if *generateJSON {
		generateJSONFunction(*keyVaultName, *defaultMount, *defaultPath, *defaultCopy, *jsonFile)
	} else if *copySecrets {
		copySecretsFunction(*keyVaultName, *vaultAddr, token, *vaultNamespace, *jsonFile)
	} else {
		log.Fatalf("Either --gen or --copy flag must be specified")
	}
}

// Structure of each secret
type Secret struct {
	KeyVaultSecretName string `json:"key_vault_secret_name"`
	VaultSecretMount   string `json:"vault_secret_mount"`
	VaultSecretPath    string `json:"vault_secret_path"`
	VaultSecretName    string `json:"vault_secret_name"`
	VaultSecretKey     string `json:"vault_secret_key"`
	Copy               bool   `json:"copy"`
}

// Overall structutre of secrets
type Secrets struct {
	Secrets []Secret `json:"secrets"`
}
