package dhcp

import (
	"fmt"
	"github.com/pkg/errors"
	"regexp"
	"strings"
)

const macAddressRegexp = `([0-9A-Fa-f]{2}[:-]){5}([0-9A-Fa-f]{2})`

type Lease struct {
	Hostname   string
	Ip         string
	MacAddress string
	IsActive   bool
}

type InvalidLeaseFormatError struct {
	Arg   string
	Cause error
}

func (i *InvalidLeaseFormatError) Error() string {
	err := fmt.Sprint("failed to convert text to lease")
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

	isActive, err := l.extractIsActive(lines)
	if err != nil {
		return err
	}

	macAddress, err := l.extractMacAddress(lines)
	if err != nil {
		return err
	}

	hostname, err := l.extractHostname(lines)
	if err != nil {
		return err
	}

	l.Ip = ip
	l.IsActive = isActive
	l.MacAddress = macAddress
	l.Hostname = hostname
	return nil
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
	return false, bindingStateError
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

func findLineWithPrefix(lines []string, prefix string) int {
	for idx, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), prefix) {
			return idx
		}
	}

	return -1
}
