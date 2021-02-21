package dhcp_test

import (
	"fmt"
	"github.com/handofgod94/dhcpwatch/dhcp"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	leaseDatabase        = fmt.Sprintf("%s\n%s", validLease, validLease)
	invalidLeaseDatabase = fmt.Sprintf("%s\n%s", invalidLease, invalidLease)
)

func TestLeaseDatabase_UnmarshalText_WithValidLease(t *testing.T) {
	lease := dhcp.Lease{
		Hostname:   "MyLocalClient",
		Ip:         "192.168.0.1",
		MacAddress: "12:ab:CD:78:90:91",
		IsActive:   true,
	}
	expected := dhcp.LeaseDatabase{Leases: []dhcp.Lease{lease, lease}}

	actual := dhcp.LeaseDatabase{}
	err := actual.UnmarshalText([]byte(leaseDatabase))
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestLeaseDatabase_UnmarshalText_WithInvalidLease(t *testing.T) {
	ld := dhcp.LeaseDatabase{}
	err := ld.UnmarshalText([]byte(invalidLeaseDatabase))
	assert.Error(t, err)
}
