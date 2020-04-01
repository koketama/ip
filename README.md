# ip
*This is a IP filter, used in filter ip belongs to some whitelist or blacklist defined by CIDR.
It supports both IPV4 and IPV6, and use the block[min, max] to describe CIDR, so do not worry about the countless IPV6 ips. And it use binary search to check if some ip in the blocks, the query speed very quickly.*



#### Example

```go
package main

import (
	"github.com/koketama/ip"
)

func main() {
	zoneUS4, err := ip.MkZone("us", "http://ipverse.net/ipblocks/data/countries/us.zone")
	if err != nil {
		// TODO
	}

	zoneUS16, err := ip.MkZone("us", "http://ipverse.net/ipblocks/data/countries/us-ipv6.zone")
	if err != nil {
		// TODO
	}

	filter, err := ip.NewFilter(zoneUS4, zoneUS16)
	if err != nil {
		// TODO
	}

	ok, zone, err := filter.Bingo("2001:430::")
	if err != nil {
		// TODO
	}
}
```

