# Descope CLI

The `descope` command line tool provides a convenient way to perform common tasks on your Descope project.

## Getting Started

### Requirements

-   The Descope project's `Project ID` is required by `descope` to know which project
    to work with. You can find it in the [project page](https://app.descope.com/settings/project)
    in the Descope console.
-   You'll also need a valid Descope management key for the above project. You can create
    a management key in the [Company section](https://app.descope.com/settings/company) in
    the Descope console.

### Installing

The `descope` tool is available as a downloadable binary from the [releases page](https://github.com/descope/descopecli/releases/latest).

Alternatively, you can build it directly from source with the `go` compiler.

1.  Verify that you have Go 1.21 or newer installed:

    ```bash
    go version
    ```

    If `go` is not installed follow the instructions on the [Go website](https://go.dev/dl).

2.  Clone or download the repository:

    ```bash
    git clone https://github.com/descope/descopecli
    cd descopecli
    ```

3.  Install `descope` with `make install`:

    ```bash
    # installs to $GOPATH/bin by default
    make install
    ```

#### Using Debian/Ubuntu or Fedora/CentOS

```bash
# Debian/Ubuntu
sudo apt-key adv --keyserver keyserver.ubuntu.com --recv-keys e8365d8513142909
echo "deb https://descope.github.io/packages stable main" | sudo tee /etc/apt/sources.list.d/descope.list
sudo apt-get update
sudo apt-get install descope

# Fedora/CentOS
sudo dnf config-manager --add-repo https://descope.github.io/packages/descope.repo
sudo dnf install descope
```

### Usage

All `descope` commands expect the Descope management key to be provided in
the `DESCOPE_MANAGEMENT_KEY` environment variable. You'll have to provide your
Descope project's unique id either in the `DESCOPE_PROJECT_ID` environment
variable or as a command argument, depending on the command.

```bash
export DESCOPE_MANAGEMENT_KEY=...
export DESCOPE_PROJECT_ID=...
descope help
```

```
A command line utility for working with the Descope management APIs

Usage:
  descope [command]

Entity Commands:
  access-key  Commands for creating and managing access keys
  tenant      Commands for creating and managing tenants
  user        Commands for creating and managing users

Project Commands:
  audit       Commands for working with audit logs
  project     Commands for managing projects

Additional Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
```

### Examples

#### Create a tenant

```bash
# creates a new tenant with a predefined tenantId
descope tenant create 'AcmeCorp' --id 'acmecorp'
```

```
* Created new tenant with id: acmecorp
```

#### List all tenants

```bash
# use the --json option to get structured JSON output from any command
descope tenant load-all --json
```

```json
{
    "count": 1,
    "ok": true,
    "tenants": [
        {
            "id": "acmecorp",
            "name": "AcmeCorp",
            "selfProvisioningDomains": [],
            "authType": "none"
        }
    ]
}
```

#### Create a user in a tenant

```bash
# creates a user and sends them an invitation if configured in the Descope console
descope user create 'andyr@example.com' --name 'Andy Rhoads' -t 'acmecorp' --json
```

```json
{
    "ok": true,
    "user": {
        "name": "Andy Rhoads",
        "email": "andyr@example.com",
        "userId": "U2eY8ZRNUlC9IKqLGzmAww7qgK0T",
        "loginIds": ["andyr@example.com"],
        "verifiedEmail": true,
        "userTenants": [
            {
                "tenantId": "acmecorp",
                "tenantName": "AcmeCorp"
            }
        ],
        "status": "invited",
        "createdTime": 1712070205
    }
}
```

#### List all users

```bash
# returns a page of user results
descope user load-all --limit 10 --page 0
```

```
* Loaded 3 users
  - User 0: { "name": ... }
  - User 1: { "name": ... }
  - User 2: { "name": ... }
```

### Manage project settings

```bash
# to prevent mistakes some project commands require the projectId as
# an argument, rather than as an environment variable

# export a snapshot of all the project's settings and configurations
descope project snapshot export P2abc... --path ./descope_export

# import the exported snapshot from the first project into another project
descope project snapshot import P2xyz... --path ./descope_export
```

### Search audit records

```bash
# searches for any audit records about the user we created above
descope audit search 'andyr' --json
```

```json
{
    "count": 1,
    "ok": true,
    "records": [
        {
            "action": "UserCreated",
            "loginIds": ["andyr@example.com"]
        }
    ]
}
```

## Support

#### Contributing

If anything is missing or not working correctly please open an issue or pull request.

#### Learn more

To learn more please see the [Descope documentation](https://docs.descope.com).

#### Contact us

If you need help you can hop on our [Slack community](https://www.descope.com/community) or send an email to [Descope support](mailto:support@descope.com).
