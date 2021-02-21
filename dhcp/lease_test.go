package dhcp_test

import (
	"github.com/handofgod94/dhcpwatch/dhcp"
	"github.com/stretchr/testify/assert"
	"testing"
)

const (
	validLease = `
lease 192.168.0.1 {
  starts 6 2021/02/20 14:11:36;
  ends 6 2021/02/20 14:21:36;
  cltt 6 2021/02/20 14:11:36;
  binding state active;
  next binding state free;
  rewind binding state free;
  hardware ethernet 12:ab:CD:78:91;
  uid "\001\204\330\033E\023=";
  set vendor-class-identifier = "MSFT 5.0";
  client-hostname "MyLocalClient";
}
`
	invalidLease = `foo {bar;}`
)

func TestLease_UnmarshalText_ForValidConfig(t *testing.T) {
	expected := dhcp.Lease{
		Hostname:   "MyLocalClient",
		Ip:         "192.168.0.1",
		MacAddress: "12:ab:CD:78:90:91",
		IsActive:   true,
	}
	actual := dhcp.Lease{}
	err := actual.UnmarshalText([]byte(validLease))
	assert.NoError(t, err)
	assert.Equal(t, expected, actual)
}

func TestLease_UnmarshalText_ErroredForInvalidConfig(t *testing.T) {
	var invalidLeaseError *dhcp.InvalidLeaseFormatError
	l := dhcp.Lease{}
	err := l.UnmarshalText([]byte(invalidLease))
	assert.ErrorAs(t, err, &invalidLeaseError)
}
