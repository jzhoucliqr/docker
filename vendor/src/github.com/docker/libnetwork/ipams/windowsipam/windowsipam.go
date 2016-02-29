package windowsipam

import (
	"net"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/libnetwork/discoverapi"
	"github.com/docker/libnetwork/ipamapi"
	"github.com/docker/libnetwork/types"
)

const (
	localAddressSpace  = "LocalDefault"
	globalAddressSpace = "GlobalDefault"
)

var (
	defaultPool, _ = types.ParseCIDR("0.0.0.0/0")
)

type allocator struct {
}

// GetInit registers the built-in ipam service with libnetwork
func GetInit(ipamName string) func(ic ipamapi.Callback, l, g interface{}) error {
	return func(ic ipamapi.Callback, l, g interface{}) error {
		return ic.RegisterIpamDriver(ipamName, &allocator{})
	}
}

func (a *allocator) GetDefaultAddressSpaces() (string, string, error) {
	return localAddressSpace, globalAddressSpace, nil
}

// RequestPool returns an address pool along with its unique id. This is a null ipam driver. It allocates the
// subnet user asked and does not validate anything. Doesnt support subpool allocation
func (a *allocator) RequestPool(addressSpace, pool, subPool string, options map[string]string, v6 bool) (string, *net.IPNet, map[string]string, error) {
	log.Debugf("RequestPool(%s, %s, %s, %v, %t)", addressSpace, pool, subPool, options, v6)
	if subPool != "" || v6 {
		return "", nil, nil, types.InternalErrorf("This request is not supported by null ipam driver")
	}

	var ipNet *net.IPNet
	var err error

	if pool != "" {
		_, ipNet, err = net.ParseCIDR(pool)
		if err != nil {
			return "", nil, nil, err
		}
	} else {
		ipNet = defaultPool
	}

	return ipNet.String(), ipNet, nil, nil
}

// ReleasePool releases the address pool - always succeeds
func (a *allocator) ReleasePool(poolID string) error {
	log.Debugf("ReleasePool(%s)", poolID)
	return nil
}

// RequestAddress returns an address from the specified pool ID.
// Always allocate the 0.0.0.0/32 ip if no preferred address was specified
func (a *allocator) RequestAddress(poolID string, prefAddress net.IP, opts map[string]string) (*net.IPNet, map[string]string, error) {
	log.Debugf("RequestAddress(%s, %v, %v) %s", poolID, prefAddress, opts, opts["RequestAddressType"])
	_, ipNet, err := net.ParseCIDR(poolID)

	if err != nil {
		return nil, nil, err
	}
	if prefAddress == nil {
		return ipNet, nil, nil
	}
	return &net.IPNet{IP: prefAddress, Mask: ipNet.Mask}, nil, nil
}

// ReleaseAddress releases the address - always succeeds
func (a *allocator) ReleaseAddress(poolID string, address net.IP) error {
	log.Debugf("ReleaseAddress(%s, %v)", poolID, address)
	return nil
}

// DiscoverNew informs the allocator about a new global scope datastore
func (a *allocator) DiscoverNew(dType discoverapi.DiscoveryType, data interface{}) error {
	return nil
}

// DiscoverDelete is a notification of no interest for the allocator
func (a *allocator) DiscoverDelete(dType discoverapi.DiscoveryType, data interface{}) error {
	return nil
}
