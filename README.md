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
    token:="aaaabbbccxcdasda"
    client,err := gohetznerdns.NewClient(token)
    if err != nil {
        fmt.Errorf("invalid token %s", token)
    }
}
```

## Examples

### List all domains

```go
package main

import (
	"fmt"
	"os"

	"github.com/opsheaven/gohetznerdns"
)

func main() {
	var err error
	var client *gohetznerdns.HetznerDNS
	var zones []*gohetznerdns.Zone

	if client, err = gohetznerdns.NewClient(os.Getenv("HETZNER_DNS_TOKEN")); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if zones, err = client.ZoneService.GetAllZones(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("%-24s  |  %s\n", "ID", "NAME")
	for _, zone := range zones {
		fmt.Printf("%-24s  |  %s\n", zone.Id, zone.Name)
	}
}
```

### Query Domains By Name

```go
package main

import (
	"fmt"
	"os"

	"github.com/opsheaven/gohetznerdns"
)

func main() {
	var err error
	var client *gohetznerdns.HetznerDNS
	var zones []*gohetznerdns.Zone

	if client, err = gohetznerdns.NewClient(os.Getenv("HETZNER_DNS_TOKEN")); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if zones, err = client.ZoneService.GetAllZonesByName("opsheaven"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("%-24s  |  %s\n", "ID", "NAME")
	for _, zone := range zones {
		fmt.Printf("%-24s  |  %s\n", zone.Id, zone.Name)
	}
}
```


### Query Records of Domain

```go
package main

import (
	"fmt"
	"os"

	"github.com/opsheaven/gohetznerdns"
)

func main() {
	var err error
	var client *gohetznerdns.HetznerDNS
	var zones []*gohetznerdns.Zone
	var records []*gohetznerdns.Record

	if client, err = gohetznerdns.NewClient(os.Getenv("HETZNER_DNS_TOKEN")); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if zones, err = client.ZoneService.GetAllZonesByName("opsheaven.space"); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if records, err = client.RecordService.GetAllRecords(zones[0].Id); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("%-32s | %-32s | %-5s | %s\n", "ID", "NAME", "TYPE", "VALUE")
	for _, record := range records {
		fmt.Printf("%-32s | %-32s | %-5s | %s\n", record.Id, record.Name, record.Type, record.Value)
	}
}
```
## Versioning

Each version of the client is tagged and the version is updated accordingly.

To see the list of past versions, run `git tag` or check [releases](https://github.com/opsheaven/gohetznerdns/releases)
