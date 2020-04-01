# ip
*This is a IP filter, used in filter ip belongs to some whitelist or blacklist defined by CIDR.
It supports both IPV4 and IPV6, and use the block[min, max] to describe CIDR, so do not worry about the countless IPV6 ips. And it use binary search to check if some ip in the blocks, the query speed very quickly.*



#### How to use

>**imports ("github.com/koketama/ip")**



#### Example

```go
zoneUS4, err := MkZone("us", "http://ipverse.net/ipblocks/data/countries/us.zone")

zoneUS16, err := MkZone("us", "http://ipverse.net/ipblocks/data/countries/us-ipv6.zone")

filter, err := NewFilter(zoneUS4, zoneUS16)

ok, zone, err := filter.Bingo("2001:430::")
```

