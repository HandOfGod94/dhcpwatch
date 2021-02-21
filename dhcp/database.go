package dhcp

import (
	"github.com/pkg/errors"
	"strings"
)

type LeaseDatabase struct {
	Leases []Lease
}

func (ld *LeaseDatabase) UnmarshalText(text []byte) error {
	var leaseStart, leaseEnd int
	var fragmentStarted, fragmentEnded bool
	lines := strings.Split(strings.TrimSpace(string(text)), "\n")

	for idx, line := range lines {
		if strings.HasPrefix(line, "lease") {
			fragmentStarted = true
			leaseStart = idx
		}

		if strings.Contains(line, "}") {
			fragmentEnded = true
			leaseEnd = idx
		}

		if fragmentStarted && fragmentEnded {
			rawLease := strings.Join(lines[leaseStart:leaseEnd+1], "\n")
			lease := Lease{}
			err := lease.UnmarshalText([]byte(rawLease))

			if err != nil {
				return errors.Wrapf(err, "invalid lease content between %d:%d", leaseStart, leaseEnd)
			}

			ld.Leases = append(ld.Leases, lease)

			fragmentStarted = false
			fragmentEnded = false
		}
	}

	if len(ld.Leases) == 0 {
		return &InvalidLeaseFormatError{Arg: "invalid lease database"}
	}

	return nil
}
