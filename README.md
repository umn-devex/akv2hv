# akv2hv

Go app for migrating secrets from Azure KeyVault to Hashicorp Vault.

## Prerequisites

1. Download the latest akv2hv binary for your OS from https://github.com/umn-secm/akv2hv/releases or [build locally](./README.md#building-locally)

      - Windows (probably want windows_amd64)
      - Linux (probably want linux_amd64)
      - MacOS (darwin_amd64 for intel based macs, darwin_arm64 for apple silicon based macs). You may get a warning that the the app is from an unidentified developer. You will need to be an administrator on your mac and follow [these instructions](https://support.apple.com/en-us/102445)

2. Open a command line window in the directory that you downloaded the binary to.

3. If you downloaded the binary, unzip it i.e. `tar -xzvf akv2hv_0.0.7_linux_amd64.tar.gz`

4. Install azure cli <https://learn.microsoft.com/en-us/cli/azure/install-azure-cli>

5. Login to azure cli `az login`

6. Set your azure cli subscription to the subscription that contains the keyvault `az account set --subscription <subID>`

7. Get a vault token with permissions to write secrets

    - If you have the [vault enterprise](https://www.hashicorp.com/en/resources/getting-vault-enterprise-installed-running) cli installed, run `vault login --method=saml --namespace=admin`. If you need to install the vault enterprise binary, make sure to uninstall the non-enterprise binary first. To verify that you have the enterprise binary, run `vault --version` and make sure that `+ent` is at the end of the version number (i.e. `Vault v1.19.1+ent`)

    - If you do not have the vault enterprise cli installed login to the vault GUI and go to the `Person Icon>Copy token` and use the `--vault_token` flag.

## Copy Secrets

1. The first step is to generate a json file with the list of all secrets you have in KeyVault. This will only retrieve their names, not their values. CLI flags for this step: 
      
      - **Required:** `-gen` - Generate json file secrets.json with a list of secrets from KeyVault as keys.
      - **Required:** `-kv` *string* - The name of the Azure Key Vault.
      - **Optional:** `-file` *string* - json file to write or read list of secrets from/to. Defaults to secrets.json in the current directory
      - **Optional:** `-default_copy` - sets the copy attribute as true for all secrets in the json file if you want to default move them all
      - **Optional:** `-default_mount` *string* - sets the default kvv2 mount location in the json file instead of `secret/`
      - **Optional:** `-default_path` *string* - sets the default path for all secrets in the json file

    ```bash
    # Linux & MacOS
    ./akv2hv --kv=INSERT_AZ_KV_NAME --gen --default_path aap/

    # Windows
    akv2hv.exe --kv=INSERT_AZ_KV_NAME --gen -default_path aap/
    ```

2. The second step is to edit the secrets.json file that was generated in step 1. The fields that you will want to edit include:

    - `vault_secret_mount`    - the kvv2 mount location (defaults to secret)
    - `vault_secret_path`     - the path to place the secret at (i.e. app1/dev/) if blank, will place the secret at the root of your mount
    - `vault_secret_name` 	- the name of the secret that you would like to store the value to (i.e. super_secret)
    - `vault_secret_key`      - the key of the field within the secret (each secret can contain multiple key/value pairs)
    - `copy`                  - true if you would like it copied to vault (defaults to false so nothing will be copied)

3. The final step is to run the copy function to retrieve the secret data from Azure KeyVault and write the secrets to Hashicorp Vault. CLI flags for this step: 

      - **Required:** `-copy` - Run the function to copy the secrets from KeyVault to HashiCorp Vault based on the secrets.json locations.
      - **Required:** `-kv` *string* - The name of the Azure Key Vault.
      - **Required:** `-vault_addr` *string* - The url for vault (i.e. https://examplevault.com).
      - **Required:** `-vault_namespace` *string* - The namespace for vault (i.e. admin/abc).
      - **Optional:** `-file` *string* - json file to write or read list of secrets from/to. Defaults to secrets.json in the current directory
      - **Optional:** `-vault_token` *string* - Vault token, will override VAULT_TOKEN environment variable and vault login token file.

    ```bash
    # Linux & MacOS
    ./akv2hv --kv=INSERT_AZ_KV_NAME --vault_addr=https://EXAMPLE.z1.hashicorp.cloud:8200/ --vault_namespace=admin/namespace --copy

    # Windows
    akv2hv.exe --kv=INSERT_AZ_KV_NAME --vault_addr=https://EXAMPLE.z1.hashicorp.cloud:8200/ --vault_namespace=admin/namespace --copy
    ```


## Building Locally

```bash
git clone git@github.com:umn-secm/akv2hv.git
cd akv2hv
go build -o .
```
