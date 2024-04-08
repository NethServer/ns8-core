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

### JWT token
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
- If the login fails, exit code must be between 2 and 7 for bad
  credentials. A value of 1 or 8 and above is considered a generic error.
  Different log lines are printed to stderr accordingly. For example:

      Bad login attempt! Exit code 2; remote address: 1.2.3.4

  In case of internal error

      Error from /home/openldap1/.config/api-moduled/handlers/login/post: exit code 1

Note that for the `login` handler, the `validate-output.json` file is
ignored. This is an example of `/api/login` response payload:

### Basic auth
You can check whether a user's credentials are valid via the `GET /api/auth` API which
does not return an authentication token but responds:
- `200 OK` if credentials are valid
- `401 Unauthorized` if credentials are invalid

```json
{
  "claims": {
    "groups": [],
    "scope" : ["change-password"],
    "uid": "first.user"
  },
  "expire": "2023-10-13T14:11:27Z",
  "message": "login succeeded",
  "code": 200,
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTcyMDYyODcsImdyb3VwcyI6W10sIm9yaWdfaWF0IjoxNjk3MTkxODg3LCJ1aWQiOiJmaXJzdC51c2VyIn0.ItWqHn94-vLWB3sIS5ontvKqJhcIrnxoYn-yG4hY9xw"
}
```

Format description:

- The `claims` attribute is a copy of JWT claims. The `claims.scope`
  attribute is used for authorization checks, as explained in the next
  section.
- JWT itself is returned as a string in the `token` attribute
- `expire` timestamp corresponds to the JWT `exp` date

In case of login failure, HTTP status 401 is returned, with a similar payload:

```json
{
  "code": 401,
  "message": "incorrect Username or Password"
}
```

## Authorization

The `login` script can optionally return a `scope` claim, so the token can
be used only for a restricted list of handlers.

- If the `scope` claim is not set, the token is valid for any handler.
- Otherwise, if the `scope` claim is set and its value is an array of
  strings, on each request the handler name is seeked in the array. If a
  match is found, the request is authorized, otherwise a `403` status is
  returned to the client.
