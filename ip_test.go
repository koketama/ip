package ip

import (
	"encoding/binary"
	"math"
	"net"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	zoneUS4 *Zone
	zoneCA4 *Zone
	zoneAU4 *Zone
)

func Test_Init(t *testing.T) {
	assert := assert.New(t)
	var err error

	zoneUS4, err = MkZone("us", "http://ipverse.net/ipblocks/data/countries/us.zone")
	assert.Nil(err)

	zoneCA4, err = MkZone("ca", "http://ipverse.net/ipblocks/data/countries/ca.zone")
	assert.Nil(err)

	zoneAU4, err = MkZone("ca", "http://ipverse.net/ipblocks/data/countries/au.zone")
	assert.Nil(err)
}

func Test_IP4(t *testing.T) {
	assert := assert.New(t)
	filter, err := NewFilter(zoneUS4, zoneCA4)
	assert.Nil(err)

	// us bingo
	// us: 6.0.0.0/7  6.0.0.0 - 7.255.255.255
	first := binary.BigEndian.Uint32([]byte{6, 0, 0, 0})
	last := binary.BigEndian.Uint32([]byte{7, 255, 255, 255})
	for k := first; k <= last; k++ {
		ip := make([]byte, 4)
		binary.BigEndian.PutUint32(ip, k)
		ok, name, err := filter.Bingo(net.IP(ip).String())
		assert.Nil(err)
		assert.True(ok)
		assert.Equal(name, "us")
	}

	maxLoop := 100

	// ca bingo
	for index, cidr := range zoneCA4.CIDR {
		if index == maxLoop {
			return
		}

		_, ipnet, _ := net.ParseCIDR(cidr)
		ones, _ := ipnet.Mask.Size()

		raw := binary.BigEndian.Uint32(ipnet.IP)
		first := raw >> (32 - ones) << (32 - ones)
		last := first | math.MaxUint32>>ones

		for k := first; k <= last; k++ {
			ip := make([]byte, 4)
			binary.BigEndian.PutUint32(ip, k)
			ok, name, err := filter.Bingo(net.IP(ip).String())
			assert.Nil(err)
			assert.True(ok)
			assert.Equal(name, "ca")
		}
	}

	// out
	for index, cidr := range zoneAU4.CIDR {
		if index == maxLoop {
			return
		}

		_, ipnet, _ := net.ParseCIDR(cidr)
		ones, _ := ipnet.Mask.Size()

		raw := binary.BigEndian.Uint32(ipnet.IP)
		first := raw >> (32 - ones) << (32 - ones)
		last := first | math.MaxUint32>>ones

		for k := first; k <= last; k++ {
			ip := make([]byte, 4)
			binary.BigEndian.PutUint32(ip, k)
			ok, name, err := filter.Bingo(net.IP(ip).String())
			assert.Nil(err)
			assert.False(ok)
			assert.Empty(name)
		}
	}
}
