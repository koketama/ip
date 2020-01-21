package ip

import (
	"encoding/binary"
	"fmt"
	"math"
	"net"
	"sort"

	"github.com/pkg/errors"
)

var _ Filter = (*filter)(nil)

type Filter interface {
	Bingo(ip string) (ok bool, name string, err error)
	info()
}

type interval4 struct {
	zone string
	min  uint32
	max  uint32
}

type ip4 [math.MaxUint16 + 1][]*interval4

type filter struct {
	ip4 ip4
}

type Zone struct {
	Name string
	CIDR []string
}

func (f *filter) info() {
	var ip4Bclass int
	var ip4MaxInterval int
	var ip4MaxIntervalIndex int

	for i := range f.ip4 {
		if f.ip4[i] != nil {
			ip4Bclass++
			intervals := len(f.ip4[i])
			if intervals > ip4MaxInterval {
				ip4MaxInterval = intervals
				ip4MaxIntervalIndex = i
			}
		}
	}

	fmt.Println(fmt.Sprintf("ip4Bclass: %d, ip4MaxInterval:%d", ip4Bclass, ip4MaxInterval))

	bClass := make([]byte, 2)
	binary.BigEndian.PutUint16(bClass, uint16(ip4MaxIntervalIndex))
	fmt.Println(bClass)
	for i, v := range f.ip4[ip4MaxIntervalIndex] {
		fmt.Println(i, v)
	}
}

// NewFilter return new instance
// filter no support dynamic change zones reason for performance;
// if zones changed, just new another filter.
func NewFilter(zones ...*Zone) (Filter, error) {
	f := new(filter)

	if len(zones) == 0 {
		return nil, errors.New("zones required")
	}

	var ip4 ip4

	do4 := func(zone string, ip net.IP, ones int) {
		raw := binary.BigEndian.Uint32(ip)
		min := raw >> (32 - ones) << (32 - ones)
		max := min | math.MaxUint32>>ones

		if ones >= 16 { // c or d class
			bClass := binary.BigEndian.Uint16(ip)
			ip4[bClass] = append(ip4[bClass],
				&interval4{zone: zone, min: min, max: max})
			return
		}

		// a or b class
		firstBclass := binary.BigEndian.Uint16(ip)
		ip4[firstBclass] = append(ip4[firstBclass],
			&interval4{zone: zone, min: min, max: min | math.MaxUint32>>16})

		secondBclass := min>>16 + 1
		lastBclass := max >> 16
		for bClass := secondBclass; bClass <= lastBclass; bClass++ {
			ip4[uint16(bClass)] = append(ip4[uint16(bClass)],
				&interval4{zone: zone, min: bClass << 16, max: bClass<<16 | math.MaxUint32>>16})
		}
	}

	for _, zone := range zones {
		for _, cidr := range zone.CIDR {
			_, netip, err := net.ParseCIDR(cidr)
			if err != nil {
				return nil, errors.WithMessagef(err, "parse cidr %s err", cidr)
			}

			ones, bits := netip.Mask.Size()
			switch bits {
			case 32:
				do4(zone.Name, netip.IP, ones)

			case 128:
				// do16(zone.Name, netip.IP, ones)
			}
		}
	}

	if len(ip4) == 0 { // && len(ip16) == 0
		return nil, errors.New("both ip4 and ip16 are empty")
	}

	for k := range ip4 {
		if ip4[k] != nil {
			sort.Slice(ip4[k], func(i, j int) bool {
				return ip4[k][i].min < ip4[k][j].min
			})
		}
	}

	f.ip4 = ip4
	// f.ip16 = ip16

	return f, nil
}

func (f *filter) Bingo(ip string) (ok bool, zone string, err error) {
	netIP := net.ParseIP(ip)
	if netIP == nil {
		err = errors.Errorf("%s is not ip4 or ip16", ip)
		return
	}

	do4 := func(ip net.IP) {
		raw := binary.BigEndian.Uint32(ip)
		bClass := binary.BigEndian.Uint16(ip)

		if intervals := f.ip4[bClass]; intervals != nil {
			index := sort.Search(len(intervals), func(i int) bool {
				return raw <= intervals[i].max
			})
			if index != -1 && index < len(intervals) && intervals[index].min <= raw {
				ok, zone = true, intervals[index].zone
			}
		}
	}

	if ip := []byte(netIP.To4()); ip != nil {
		do4(ip)
	}

	return
}
