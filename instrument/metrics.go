package instrument

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	DhcpTable = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "router",
		Name:      "dhcp_table",
		Help:      "dhcp table containing ip assignment",
	}, []string{"hostname", "ip_address", "mac_address", "is_active", "lease_end"})
)
