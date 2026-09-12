
<div align="center">
  <a href="https://github.com/descope/descopecli">
    <img src=".github/images/descope-logo.png" alt="Descope Logo" width="160" height="160">
  </a>

  <h3 align="center">Descope CLI</h3>

  <p align="center">
    A command line tool for managing your Descope project 
  </p>
</div>

<br />

## About

The `descope` command line tool provides a convenient way to perform common tasks on your Descope project by leveraging Descope's management APIs.

* Create and modify project entities such as users, tenants and access keys.
* Manage project settings and configurations using snapshots that can be exported, validated and imported into other projects.
* Search and display audit logs for projects.
* Supports JSON output for easy integration into scripts and CI/CD workflows.

<br/>

## Installation

### All Platforms

The `descope` tool is available as a downloadable binary from the [releases page](https://github.com/descope/descopecli/releases/latest).

### Debian/Ubuntu

Install `descope` using APT:

```bash
sudo apt-key adv --keyserver keyserver.ubuntu.com --recv-keys e8365d8513142909
echo "deb https://descope.github.io/packages stable main" | sudo tee /etc/apt/sources.list.d/descope.list
sudo apt-get update
sudo apt-get install descope
```

### Fedora/CentOS

You can install `descope` using DNF by adding the Descope repository:

```bash
sudo dnf config-manager --add-repo https://descope.github.io/packages/descope.repo
sudo dnf install descope
```

### Build from Source

You can build the `descope` command line tool directly with the `go` compiler:

1.  Verify that you have Go 1.21 or newer installed, and if not, follow the instructions on the [Go website](https://go.dev/dl):

    ```bash
    go version
    ```

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

<br/>

## Getting Started

### Requirements

-   The Descope project's `Project ID` is required by `descope` to know which project
    to work with. You can find it in the [Project section](https://app.descope.com/settings/project)
    in the Descope console.
-   You'll also need credentials for the above project, either a management key or a
    workload identity token. You can create a management key in the
    [Company section](https://app.descope.com/settings/company) in the Descope console.

### Usage

All `descope` commands expect credentials to be provided in one of two environment
variables:

| Variable | Value |
|----------|-------|
| `DESCOPE_MANAGEMENT_KEY` | A management key created in the Descope console. |
| `DESCOPE_WORKLOAD_TOKEN` | A workload identity token, for CI jobs that authenticate with the OIDC token they request for themselves rather than a stored secret. See [Using a workload identity token](#using-a-workload-identity-token). |

Set whichever suits how the command is being run. The two credentials are exclusive: a
management key that a workload identity token acts as does not accept its own secret, so
setting both is rejected rather than resolved in favour of one. You'll also have to provide
your Descope project's unique id either in the `DESCOPE_PROJECT_ID` environment variable or
as a command argument, depending on the command.

```bash
export DESCOPE_PROJECT_ID='P...'

# with a management key
export DESCOPE_MANAGEMENT_KEY='K...'

# or with a workload identity token
export DESCOPE_WORKLOAD_TOKEN='eyJ...'

descope --help
```

```
A command line utility for working with the Descope management APIs

Usage:
  descope [command]

Entity Commands:
  access-key  Commands for creating and managing access keys
  apps        Commands for creating and managing applications and integrations
  tenant      Commands for creating and managing tenants
  user        Commands for creating and managing users

Project Commands:
  audit       Commands for working with audit logs
  project     Commands for managing projects

Additional Commands:
  completion  Generate the autocompletion script for the specified shell
  help        Help about any command
```

#### Using a workload identity token

A CI job can authenticate with the short-lived OIDC token it requests for itself instead of
a stored management key, so there is no secret to leak or rotate.

The token is not a credential on its own. The Management API authorizes it as one specific
management key, the way `sts:AssumeRoleWithWebIdentity` requires a role to assume. That key
is named by the token's own `aud` claim, as `<apiBaseUrl>/<keyId>` — so the identifier is
covered by the issuer's signature rather than asserted next to it. The key's permissions
decide what the token can do, and that key holds no usable secret: it authenticates only
through a token from its trusted issuer, and presenting a static secret for it is refused.

The issuer must first be registered as a trusted issuer for your company, which creates the
management key at the same time and returns its id. That configuration decides which
subjects are accepted.

Because the audience has to match exactly, the `import` and `export` actions request the
token themselves rather than taking one as an input. Give them the audience and grant the
job `id-token: write`:

```yaml
permissions:
  id-token: write
  contents: read

steps:
  # the export action writes into the working tree and diffs it, so it needs a checkout
  - uses: actions/checkout@v4

  - name: Export Snapshot
    uses: descope/descopecli/.github/actions/export@main
    with:
      project_id: ${{ vars.PRODUCTION_PROJECT_ID }}
      token_audience: ${{ vars.DESCOPE_TOKEN_AUDIENCE }}
      files_path: ./descope_export
```

The audience is not a secret, so it belongs in a CI variable. Read it back from the
management key rather than composing it by hand: Descope chooses it, and it is reported
read-only on the key.

When running the `descope` binary directly, request the token yourself and put it in
`DESCOPE_WORKLOAD_TOKEN`, using the audience the key reports:

```yaml
  - name: Request an OIDC token
    id: token
    uses: actions/github-script@v7
    with:
      script: |
        // the audience names the management key being assumed
        const token = await core.getIDToken(process.env.AUDIENCE)
        core.setSecret(token)
        core.setOutput('token', token)
    env:
      AUDIENCE: ${{ vars.DESCOPE_TOKEN_AUDIENCE }}

  - name: Run descope
    run: descope project snapshot export "$DESCOPE_PROJECT_ID" --path ./descope_export
    env:
      DESCOPE_PROJECT_ID: P...
      DESCOPE_WORKLOAD_TOKEN: ${{ steps.token.outputs.token }}
```

<br/>

## Examples

### Tenants

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

### Users

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

### Project

#### Manage project settings

```bash
# to prevent mistakes some project commands require the projectId as
# an argument, rather than as an environment variable

# export a snapshot of all the project's settings and configurations
descope project snapshot export 'P2abc...' --path ./descope_export

# import the exported snapshot from the first project into another project
descope project snapshot import 'P2xyz...' --path ./descope_export
```

#### Search audit records

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

<br/>

## Support

#### Contributing

If anything is missing or not working correctly please open an issue or pull request.

#### Learn more

To learn more please see the [Descope documentation](https://docs.descope.com).

#### Contact us

If you need help you can hop on our [Slack community](https://www.descope.com/community) or send an email to [Descope support](mailto:support@descope.com).
