# join-cluster input Schema

```txt
http://schema.nethserver.org/cluster/join-cluster-input.json
```

Discard current Redis DB and installed modules, then join an existing cluster

| Abstract            | Extensible | Status         | Identifiable | Custom Properties | Additional Properties | Access Restrictions | Defined In                                                                        |
| :------------------ | :--------- | :------------- | :----------- | :---------------- | :-------------------- | :------------------ | :-------------------------------------------------------------------------------- |
| Can be instantiated | No         | Unknown status | No           | Forbidden         | Allowed               | none                | [join-cluster-input.json](cluster/join-cluster-input.json "open original schema") |

## join-cluster input Type

unknown ([join-cluster input](join-cluster-input.md))

## join-cluster input Examples

```json
{
  "tls_verify": true,
  "url": "https://my.example.com",
  "jwt": "eyJ0eXAiOiJKV1QiLCJhbGciOiJSUzI1NiJ9.eyJmb28iOiJiYXIiLCJpYXQiOjE0MzQ0Nzk4ODN9.HQyx15jWm1upqsrKSf89X_iP0sg7N46a9pqBVGPMYdiqZeuU_ZZOdU-zizHJoIHMIJxtEWzpSMaVubJW0AJsTqjqQf6GoJ4cmFAfmfUFXmMC4Xv5oc4UqvGizpoLjfZedd834PcwbS-WskZcL4pVNmBIGRtDXkoU1j2X1P5M_sNJ9lYZ5vITyqe4MYJovQzNdQziUNhcMI5wkXncV7XzGInBeQsPquASWVG4gb3Y--k1P3xWA4Df3rKeEQBbInDKXczvDpfIlTojx4Ch8OM8vXWWNxW-mIQrV31wRrS9XtNoig7irx8N0MzokiYKrQ8WP_ezPicHvVPIHhz-InOw"
}
```

# join-cluster input Properties

| Property                   | Type      | Required | Nullable       | Defined by                                                                                                                                                    |
| :------------------------- | :-------- | :------- | :------------- | :------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| [tls\_verify](#tls_verify) | `boolean` | Required | cannot be null | [join-cluster input](join-cluster-input-properties-tls-verification.md "http://schema.nethserver.org/cluster/join-cluster-input.json#/properties/tls_verify") |
| [url](#url)                | `string`  | Required | cannot be null | [join-cluster input](join-cluster-input-properties-leader-node-url.md "http://schema.nethserver.org/cluster/join-cluster-input.json#/properties/url")         |
| [jwt](#jwt)                | `string`  | Required | cannot be null | [join-cluster input](join-cluster-input-properties-jwt-authorization-token.md "http://schema.nethserver.org/cluster/join-cluster-input.json#/properties/jwt") |
| [endpoint](#endpoint)      | `string`  | Optional | cannot be null | [join-cluster input](join-cluster-input-properties-endpoint.md "http://schema.nethserver.org/cluster/join-cluster-input.json#/properties/endpoint")           |

## tls\_verify

If set to `false` the given @url is not subject to TLS certificate validation

`tls_verify`

* is required

* Type: `boolean` ([TLS verification](join-cluster-input-properties-tls-verification.md))

* cannot be null

* defined in: [join-cluster input](join-cluster-input-properties-tls-verification.md "http://schema.nethserver.org/cluster/join-cluster-input.json#/properties/tls_verify")

### tls\_verify Type

`boolean` ([TLS verification](join-cluster-input-properties-tls-verification.md))

### tls\_verify Examples

```json
true
```

```json
false
```

## url

Public URL of the leader API endpoint

`url`

* is required

* Type: `string` ([Leader node URL](join-cluster-input-properties-leader-node-url.md))

* cannot be null

* defined in: [join-cluster input](join-cluster-input-properties-leader-node-url.md "http://schema.nethserver.org/cluster/join-cluster-input.json#/properties/url")

### url Type

`string` ([Leader node URL](join-cluster-input-properties-leader-node-url.md))

### url Constraints

**pattern**: the string must match the following regular expression:&#x20;

```regexp
^https?://.+$
```

[try pattern](https://regexr.com/?expression=%5Ehttps%3F%3A%2F%2F.%2B%24 "try regular expression with regexr.com")

## jwt



`jwt`

* is required

* Type: `string` ([JWT authorization token](join-cluster-input-properties-jwt-authorization-token.md))

* cannot be null

* defined in: [join-cluster input](join-cluster-input-properties-jwt-authorization-token.md "http://schema.nethserver.org/cluster/join-cluster-input.json#/properties/jwt")

### jwt Type

`string` ([JWT authorization token](join-cluster-input-properties-jwt-authorization-token.md))

### jwt Constraints

**minimum length**: the minimum number of characters for this string is: `1`

## endpoint

The public WireGuard VPN endpoint in `hostname:port` form. The given port may differ from 55820 as it depends on external configurations (i.e. port-forwards)

`endpoint`

* is optional

* Type: `string`

* cannot be null

* defined in: [join-cluster input](join-cluster-input-properties-endpoint.md "http://schema.nethserver.org/cluster/join-cluster-input.json#/properties/endpoint")

### endpoint Type

`string`

### endpoint Examples

```json
"my.example.com:55820"
```

```json
"1.2.3.4:60000"
```
