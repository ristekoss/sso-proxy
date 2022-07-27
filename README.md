# SSO Proxy

Proxy for SSO UI to use REST API call.

This project is meant to run on Google Cloud Function, but you can run it manually with `make run`

to use this, send a post request with username and password in JSON format

```json
{
  "username": "username",
  "password": "password"
}
```

The API call can take up to 4 seconds. This is purely because of many synchronous back to back call between the proxy and the SSO server.
