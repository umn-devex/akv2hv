# akv2hv
Go app for migrating secrets from Azure KeyVault to Hashicorp Vault.

## Prerequisites

1. Download the latest akv2hv binary for your OS from https://github.com/umn-secm/akv2hv/releases or [build locally](./README.md#building-locally) and open a command line window in the directory that you downloaded the binary to.

2. Install azure cli <https://learn.microsoft.com/en-us/cli/azure/install-azure-cli>

3. Login to azure cli `az login`

4. Get a vault token with permissions to write secrets (if you want to use your own token, login to the vault GUI and go to the Person Icon>Copy token). Tokens are good for a limited amount of time and may expire.

5. Export your vault token `export VAULT_TOKEN=TOKEN_FROM_STEP_3`

## Running

1. The first step is to generate a json file with the list of all secrets you have in KeyVault. This will only retrieve their names, not their values.

    - Generate with default kvv2 mount `secret`

        ```bash
        # Linux
        ./akv2hv --kv=INSERT_AZ_KV_NAME --gen

        # Windows
        akv2hv.exe --kv=INSERT_AZ_KV_NAME --gen
        ```

    - Generate with alternative kvv2 mount location

        ```bash
        # Linux
        ./akv2hv --kv=INSERT_AZ_KV_NAME --mount=secret2 --gen

        # Windows
        akv2hv.exe --kv=INSERT_AZ_KV_NAME --mount=secret2 --gen
        ```

2. The second step is to edit the secrets.json file that was generated in step 1. The fields that you will want to edit include:

    - `vault_kvv2_mount`      - the kvv2 mount location (defaults to secret)
    - `vault_secret_name` 	  - the name of the secret including path that you would like to store the value to i.e. github/super_secret
    - `vault_secret_key`      - the key of the field within the secret (each secret can contain multiple key/value pairs)
    - `copy`                  - true if you would like it copied to vault (defaults to false so nothing will be copied)

3. The final step is to run the copy function to retrieve the secret data from Azure KeyVault and write the secrets to Hashicorp Vault.

    ```bash
    # Linux
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
