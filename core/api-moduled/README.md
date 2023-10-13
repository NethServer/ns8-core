# api-moduled

Generic HTTP API server with JSON schema validation and JWT authentication

## Configuration

Configuration parameters can be passed as environment variables.

- `AMLD_PUBLIC_DIR`, the directory from which static files are served --
  default `./public/`

- `AMLD_HANDLER_DIR`, the root directory where the handler commands are --
  default `./handlers/`

- `AMLD_BIND_ADDRESS`, the IP address and port number where the HTTP
  server listens for connections -- default `:9313`

- `AMLD_ID_KEY`, the name of the key in the JWT claims payload which holds
  the authenticated identity value -- default `uid`

- `AMLD_JWT_SECRET`, a secret random string used for JWT signature --
  default `` (empty)

- `AMLD_JWT_TIMEOUT`, duration in seconds of a JWT token -- default
  `14400` (4 hours)

- `AMLD_JWT_TOKEN_LOOKUP`, configure GIN-JWT token lookup. Refer to its
  README for details -- default `header: Authorization`

- `AMLD_JWT_REALM`, set the Realm for GIN-JWT cookie-based authorization
  -- default `api-moduled`

- `AMLD_EXPORT_ENV`, space-separated list of environment variable names that
  are exported to handler commands -- default `` (empty)

## How to implement the API

Each directory under `AMLD_HANDLER_DIR` is considered an API HTTP route.

For instance, if a `list/` directory exists under `AMLD_HANDLER_DIR` the
corresponding `/api/list` URL path is bound to it. Under the `list/`
directory the following files can exist:

- `post` executable command, which implements the POST HTTP verb.
- `validate-input.json` file, with the JSON schema of a valid input
- `validate-output.json` file, with the JSON schema of a valid output

Note that the `login/` directory has a special meaning.

## Command execution

When a API command is executed, the following environment variables are
passed:

- `PATH`, a copy of the existing `PATH` variable from the server
  environment
- `JWT_ID`, the value of the identity key in the JWT claims payload
- `JWT_CLAIMS`, the JWT claims payload, in JSON-encoded form
- other variables listed in `AMLD_EXPORT_ENV`

The command stderr is sent to the server log.

Furthermore, in a POST request handler:

- The command process can read from standard input the request payload, if
  it passes the JSON schema input validation step
- The command output is returned as response to the API caller, if it
  passes the JSON schema output validation step

Note that the output of the `login/post` command is not returned to the
client: it has a special meaning decribed in the *Authentication* section.

## Authentication

The `login` handler, bound to `/api/login` URL path is the only one that
can be accessed without a valid JWT token. It must provide a `post`
command, implementing the following protocol:

- The command stdin receives the request payload as usual
- If the login is successful the command prints to stdout the JWT claims
  in JSON format.
  * Print out as little data as possible!
  * Ensure the claims include the JWT identifier attribute, matching
    `JWT_ID` configuration.
  * The command exit code must be 0.
- If the login fails, exit code must be non-zero; write errors to stderr
  as usual.

Note that for the `login` handler, the `validate-output.json` file is
ignored. This is an example of `/api/login` response payload:

```json
{
  "claims": {
    "groups": [],
    "uid": "first.user"
  },
  "expire": "2023-10-13T14:11:27Z",
  "message": "login succeeded",
  "code": 200,
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTcyMDYyODcsImdyb3VwcyI6W10sIm9yaWdfaWF0IjoxNjk3MTkxODg3LCJ1aWQiOiJmaXJzdC51c2VyIn0.ItWqHn94-vLWB3sIS5ontvKqJhcIrnxoYn-yG4hY9xw"
}
```

Format description:

- The `claims` attribute is a copy of JWT claims
- JWT itself is returned as a string in the `token` attribute
- `expire` timestamp corresponds to the JWT `exp` date

In case of login failure, HTTP status 401 is returned, with a similar payload:

```json
{
  "code": 401,
  "message": "incorrect Username or Password"
}
```
