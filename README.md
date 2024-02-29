# gohetznerdns

gohetznerdns is a Go client library for accessing the Hetzner DNS API.

You can view Hetzner DNS API docs here: [https://dns.hetzner.com/api-docs](https://dns.hetzner.com/api-docs)

## Install
```sh
go get github.com/opsheaven/gohetznerdns@vX.Y.Z
```

where X.Y.Z is the [version](https://github.com/opsheaven/gohetznerdns/releases) you need.

or
```sh
go get github.com/opsheaven/gohetznerdns
```
for non Go modules usage or latest version.

## Usage

```go
import "github.com/opsheaven/gohetznerdns"
```

Create a new client, then use the exposed services to access different parts of the Hetzner DNS API.

### Authentication

API Access token is needed to access Public API. You can create the api token by following the [manual](https://docs.hetzner.com/dns-console/dns/general/api-access-token).

You can then use your token in the client:

```go
package main

import (
    "github.com/opsheaven/gohetznerdns"
)

func main() {
    client := gohetznerdns.NewClient()
    err := client.SetToken(token)
}
```

## Examples

// TODO
## Versioning

Each version of the client is tagged and the version is updated accordingly.

To see the list of past versions, run `git tag` or check [releases](https://github.com/opsheaven/gohetznerdns/releases)
