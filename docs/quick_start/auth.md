# Authenticating in Janus

To start using Janus you need to get a [Json Web Token](https://jwt.io) and provide it in every single request
using the `Authorization` header.

To get a token you must execute:

```sh
http -v --json POST localhost:8081/login username=admin password=admin
```

The username and password are defined in an environmental variable called `ADMIN_USERNAME` and `ADMIN_PASSWORD`. It defaults to *admin*/*admin*.

<p align="center">
  <a href="http://g.recordit.co/dDjkyDKobL.gif">
    <img src="http://g.recordit.co/dDjkyDKobL.gif">
  </a>
</p>

With this token you can now request the administration endpoints of Janus
