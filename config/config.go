package config

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	DefaultBindAddr       = ":8080"
	DefaultDhcpDbFilePath = "/var/lib/dhcp/dhcpd.leases"
	DefaultLogLevel       = logrus.InfoLevel
)

func Addr() string {
	key := "BIND"
	if viper.IsSet(key) {
		val := viper.GetString(key)
		logrus.Debugf("Config %v=%v", key, val)
		return val
	}

	logrus.Debugf("Using Default Config %v=%v", key, DefaultBindAddr)
	return DefaultBindAddr
}

func DhcpDbFilePath() string {
	key := "DHCP_DB_FILE_PATH"
	if viper.IsSet(key) {
		val := viper.GetString(key)
		logrus.Debugf("Config %v=%v", key, val)
		return val
	}

	logrus.Debugf("Using Default Config %v=%v", key, DefaultDhcpDbFilePath)
	return DefaultDhcpDbFilePath
}

func LogLevel() logrus.Level {
	key := "LOG_LEVEL"
	if viper.IsSet(key) {
		val := viper.GetString(key)
		logrus.Debugf("Config %v=%v", key, val)
		if level, err := logrus.ParseLevel(val); err == nil {
			return level
		}
		logrus.Fatalf("invalid %v config: %v", key, val)
	}

	logrus.Debugf("Using Default Config %v=%v", key, DefaultLogLevel)
	return DefaultLogLevel
}
