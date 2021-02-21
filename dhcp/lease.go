package dhcp

import (
	"fmt"
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

	isActive, err := l.extractIsActive(lines[4])
	if err != nil {
		return err
	}
	macAddress, err := l.extractMacAddress(lines[7])
	if err != nil {
		return err
	}
	hostname, err := l.extractHostname(lines[10])
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

func (l *Lease) extractIsActive(line string) (bool, error) {
	// binding state active;
	bindingStateError := &InvalidLeaseFormatError{Arg: "binding state"}
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

func (l *Lease) extractMacAddress(line string) (string, error) {
	// hardware ethernet <mac-address>
	macAddressError := &InvalidLeaseFormatError{Arg: "mac address"}
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

func (l *Lease) extractHostname(line string) (string, error) {
	// client-hostname "MyLocalClient";
	hostnameError := &InvalidLeaseFormatError{Arg: "client-hostname"}
	line = strings.TrimSpace(line)
	line = strings.TrimSuffix(line, ";")
	s := strings.Split(line, " ")
	if len(s) < 2 {
		return "", hostnameError
	}

	hostname := strings.ReplaceAll(s[1], `"`, "")
	return hostname, nil
}
