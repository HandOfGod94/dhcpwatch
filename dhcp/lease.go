package dhcp

import (
	"fmt"
	"github.com/handofgod94/dhcpwatch/instrument"
	"github.com/pkg/errors"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const (
	macAddressRegexp = `([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})`
	timeLayout       = "2006/01/02 15:04:05"
)

type Lease struct {
	Hostname   string
	Ip         string
	MacAddress string
	IsActive   bool
	LeaseStart time.Time
	LeaseEnd   time.Time
}

type InvalidLeaseFormatError struct {
	Arg   string
	Cause error
}

func (i *InvalidLeaseFormatError) Error() string {
	err := "failed to convert text to lease"
	if i.Arg != "" {
		err = fmt.Sprintf("%s - %v is invalid", err, i.Arg)
	}
	if i.Cause != nil {
		err = fmt.Sprintf("%s - error: %v", err, i.Cause)
	}
	return err
}

func (l *Lease) UnmarshalText(text []byte) error {
	defaultError := &InvalidLeaseFormatError{Cause: fmt.Errorf("invalid lease")}

	lines := strings.Split(strings.TrimSpace(string(text)), "\n")
	if !strings.HasPrefix(lines[0], "lease") {
		return defaultError
	}

	ip, err := l.extractIp(lines[0])
	if err != nil {
		return err
	}

	leaseStart, err := l.extractLeaseStart(lines)
	if err != nil {
		return err
	}

	leaseEnd, err := l.extractLeaseEnd(lines)
	if err != nil {
		return err
	}

	isActive, _ := l.extractIsActive(lines)
	hostname, _ := l.extractHostname(lines)
	macAddress, _ := l.extractMacAddress(lines)

	l.Ip = ip
	l.IsActive = isActive
	l.MacAddress = macAddress
	l.Hostname = hostname
	l.LeaseStart = leaseStart
	l.LeaseEnd = leaseEnd
	l.publish()
	return nil
}

func (l *Lease) publish() {
	instrument.DhcpTable.WithLabelValues(
		l.Hostname,
		l.Ip,
		l.MacAddress,
		strconv.FormatBool(l.IsActive),
	).Set(float64(l.LeaseEnd.Unix()))
}

func (l *Lease) extractIp(line string) (string, error) {
	// lease 192.168.0.1 {
	s := strings.Split(strings.TrimSpace(line), " ")
	if len(s) < 3 {
		return "", &InvalidLeaseFormatError{Arg: "ip"}
	}
	return s[1], nil
}

func (l *Lease) extractIsActive(lines []string) (bool, error) {
	// binding state active;
	var lineIdx int
	bindingStateError := &InvalidLeaseFormatError{Arg: "binding state"}

	if lineIdx = findLineWithPrefix(lines, "binding"); lineIdx < 0 {
		return false, errors.Wrap(bindingStateError, "binding state absent")
	}

	line := lines[lineIdx]
	line = strings.TrimSpace(line)
	line = strings.TrimSuffix(line, ";")
	s := strings.Split(line, " ")
	if len(s) < 3 {
		return false, bindingStateError
	}

	if s[2] == "active" {
		return true, nil
	}
	return false, nil
}

func (l *Lease) extractMacAddress(lines []string) (string, error) {
	// hardware ethernet <mac-address>
	var lineIdx int
	macAddressError := &InvalidLeaseFormatError{Arg: "mac address"}
	if lineIdx = findLineWithPrefix(lines, "hardware"); lineIdx < 0 {
		return "", errors.Wrap(macAddressError, "mac address absent")
	}

	line := lines[lineIdx]
	line = strings.TrimSpace(line)
	line = strings.TrimSuffix(line, ";")
	s := strings.Split(line, " ")
	if len(s) < 3 {
		return "", macAddressError
	}

	if matched, _ := regexp.Match(macAddressRegexp, []byte(s[2])); !matched {
		return "", macAddressError
	}

	return s[2], nil
}

func (l *Lease) extractHostname(lines []string) (string, error) {
	// client-hostname "MyLocalClient";
	var lineIdx int
	hostnameError := &InvalidLeaseFormatError{Arg: "client-hostname"}
	if lineIdx = findLineWithPrefix(lines, "client-hostname"); lineIdx < 0 {
		return "", errors.Wrap(hostnameError, "client-hostname absent")
	}

	line := lines[lineIdx]
	line = strings.TrimSpace(line)
	line = strings.TrimSuffix(line, ";")
	s := strings.Split(line, " ")
	if len(s) < 2 {
		return "", hostnameError
	}

	hostname := strings.ReplaceAll(s[1], `"`, "")
	return hostname, nil
}

func (l *Lease) extractLeaseStart(lines []string) (time.Time, error) {
	// starts 6 2021/02/20 14:11:36;
	var lineIdx int
	leaseStartError := &InvalidLeaseFormatError{Arg: "lease starts"}
	if lineIdx = findLineWithPrefix(lines, "starts"); lineIdx < 0 {
		return time.Time{}, errors.Wrap(leaseStartError, "'starts' is absent")
	}

	line := lines[lineIdx]
	line = strings.TrimSpace(line)
	line = strings.TrimSuffix(line, ";")
	s := strings.Split(line, " ")
	if len(s) < 4 {
		return time.Time{}, leaseStartError
	}

	rawStartLeaseTime := fmt.Sprintf("%s %s", s[2], s[3])
	leaseStart, err := time.Parse(timeLayout, rawStartLeaseTime)
	if err != nil {
		return time.Time{}, errors.Wrap(err, leaseStartError.Error())
	}
	return leaseStart, nil
}

func (l *Lease) extractLeaseEnd(lines []string) (time.Time, error) {
	// ends 6 2021/02/20 14:21:36;
	var lineIdx int
	leaseEndError := &InvalidLeaseFormatError{Arg: "lease ends"}
	if lineIdx = findLineWithPrefix(lines, "ends"); lineIdx < 0 {
		return time.Time{}, errors.Wrap(leaseEndError, "'ends' is absent")
	}

	line := lines[lineIdx]
	line = strings.TrimSpace(line)
	line = strings.TrimSuffix(line, ";")
	s := strings.Split(line, " ")
	if len(s) < 4 {
		return time.Time{}, leaseEndError
	}

	rawEndLeaseTime := fmt.Sprintf("%s %s", s[2], s[3])
	leaseEnd, err := time.Parse(timeLayout, rawEndLeaseTime)
	if err != nil {
		return time.Time{}, errors.Wrap(err, leaseEndError.Error())
	}
	return leaseEnd, nil
}

func findLineWithPrefix(lines []string, prefix string) int {
	for idx, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), prefix) {
			return idx
		}
	}

	return -1
}
