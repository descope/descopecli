
# Descope CLI

The `descopecli` tool provides a convenient way to perform common tasks on your Descope project.

## Getting Started

### Requirements

- The Descope project's `Project ID` is required by `descopecli` to know which project
  to work with. You can find it in the [project page](https://app.descope.com/settings/project)
  in the Descope console.
- You'll also need a valid Descope management key for the above project. You can create
  a management key in the [Company section](https://app.descope.com/settings/company) in
  the Descope console.

### Installing

For the moment, the `descopecli` tool requires the `go` compiler to be installed.

1.  Verify that you have Go 1.21 or newer installed:

    ```bash
    go version
    ```

    If `go` is not installed follow the instructions on the [Go website](https://go.dev/dl).

2.  Install `descopecli` with `go install`:

    ```bash
    # installs to $GOPATH/bin by default
    go install github.com/descope/descopecli
    ```

### Usage

All `descopecli` commands expect the Descope management key to be provided in
the `DESCOPE_MANAGEMENT_KEY` environment variable. You'll have to provide your
Descope project's unique id either in the `DESCOPE_PROJECT_ID` environment
variable or as a command argument, depending on the command.

```bash
export DESCOPE_MANAGEMENT_KEY=...
export DESCOPE_PROJECT_ID=...
descopecli help
```
```
A command line utility for working with the Descope management APIs

Usage:
  descopecli [command]

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
descopecli tenant create 'AcmeCorp' --id 'acmecorp'
```
```
* Created new tenant with id: acmecorp
```

#### Create a user in a tenant

```bash
# creates a user and sends them an invitation if configured in the Descope console
descopecli user create andyr@example.com --name 'Andy Rhoads' -t 'acmecorp'
```
```
* Created user:
  {
    "name": "Andy Rhoads",
    "userId": "U2cm9Iy",
    "loginIds": ["andyr@example.com"],
    "email": "andyr@example.com",
    "verifiedEmail": true,
    "status": "invited",
    "userTenants": [{"tenantId":"acmecorp","tenantName":"AcmeCorp"}],
    "createdTime": 1708700000
  }
```

#### List all users

```bash
# returns a page of user results
descopecli user load-all --limit 10 --page 0
```
```
* Found 3 users
  - User 0: { "name": ... }
  - User 1: { "name": ... }
  - User 2: { "name": ... }
```

### Manage project settings

```bash
# to prevent mistakes these command require the projectId as an argument,
# rather than as an environment variable
descopecli project snapshot export P2abc... --path ./descope_export

# import the exported snapshot from the first project into another project
descopecli project snapshot import P2xyz... --path ./descope_export
```

### Search audit records

```bash
# searches for any audit records about the user we created above
descopecli audit search 'andyr'
```
```
* Found 1 record
  - Record 0:
    {
      "action": "UserCreated",
      "loginIds": [
        "andyr@example.com"
      ],
      ...
    }
```

## Support

#### Contributing

If anything is missing or not working correctly please open an issue or pull request.

#### Learn more

To learn more please see the [Descope documentation](https://docs.descope.com).

#### Contact us

If you need help you can hop on our [Slack community](https://www.descope.com/community) or send an email to [Descope support](mailto:support@descope.com).
