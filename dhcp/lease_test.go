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
  hardware ethernet 12:ab:CD:78:90:91;
  uid "\001\204\330\033E\023=";
  set vendor-class-identifier = "MSFT 5.0";
  client-hostname "MyLocalClient";
}
`
	invalidLease = `foo {bar;}`
)

func TestLease_UnmarshalText(t *testing.T) {
	tests := []struct {
		name    string
		fields  dhcp.Lease
		args    []byte
		wantErr bool
	}{
		{
			name: "should parse to lease for valid data",
			fields: dhcp.Lease{
				Hostname:   "MyLocalClient",
				Ip:         "192.168.0.1",
				MacAddress: "12:ab:CD:78:90:91",
				IsActive:   true,
			},
			args:    []byte(validLease),
			wantErr: false,
		},
		{
			name:    "should return error for invalid data",
			args:    []byte(invalidLease),
			wantErr: true,
		},
		{
			name: "should return error for invalid mac address",
			args: []byte(`
					lease 192.168.0.1 {
						hardware ethernet fo:ba:rf:iz;
					}`),
			wantErr: true,
		},
		{
			name: "should return error when client-hostname is absent",
			args: []byte(`
					lease 192.168.0.1 {
  						binding state active;
  						hardware ethernet 12:ab:CD:78:90:91;
					}`),
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expected := dhcp.Lease{
				Hostname:   tt.fields.Hostname,
				Ip:         tt.fields.Ip,
				MacAddress: tt.fields.MacAddress,
				IsActive:   tt.fields.IsActive,
			}
			actual := dhcp.Lease{}
			err := actual.UnmarshalText(tt.args)

			if tt.wantErr {
				assert.Error(t, err, tt.name)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, expected, actual, tt.name)
			}
		})
	}
}
