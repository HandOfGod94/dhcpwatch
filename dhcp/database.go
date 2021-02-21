package dhcp

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"io/ioutil"
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
				logrus.
					WithError(err).
					Warnf("skipping lease parsing between lines %d:%d", leaseStart+1, leaseEnd+1)
			} else {
				ld.Leases = append(ld.Leases, lease)
			}

			fragmentStarted = false
			fragmentEnded = false
		}
	}

	if len(ld.Leases) == 0 {
		return &InvalidLeaseFormatError{Arg: "invalid lease database"}
	}

	return nil
}

func ReadDatabase(filename string) (*LeaseDatabase, error) {
	logrus.WithField("db file", filename)

	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to read dhcp database file %v", filename)
	}

	leaseDb := &LeaseDatabase{}
	if err := leaseDb.UnmarshalText(content); err != nil {
		return nil, errors.Wrapf(err, "failed to parse database file %v", filename)
	}

	logrus.Info("successfully parsed dhcp db file")
	return leaseDb, err
}
