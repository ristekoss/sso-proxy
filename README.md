# SSO Proxy

Proxy for SSO UI to use REST API call.

This project is meant to run on Google Cloud Function, but you can run it manually with `make run` on linux. On windows you need to export variable  `FUNCTION_TARGET` = `Proxy`.

to use this, send a post request with username and password in JSON format

```json
{
  "username": "username",
  "password": "password"
}
```

The API call can take up to 4 seconds. This is purely because the default implementation use dev server that is slow. If you have access to production serviceUrl cas, you can add that to the request.

```json
{
  "username": "username",
  "password": "password",
  "casUrl": "https://sso.ui.ac.id/cas/",
  "serviceUrl": "http%3A%2F%2Flocalhost%3A8080%2F"
}
```

note that `casUrl` need defined protocol and ended with backslash.
