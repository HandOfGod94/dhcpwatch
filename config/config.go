package config

import (
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	DefaultBindAddr       = ":8080"
	DefaultDhcpDbFilePath = "/var/lib/dhcp/dhcpd.leases"
)

func Addr() string {
	key := "BIND"
	if viper.IsSet(key) {
		val := viper.GetString(key)
		log.Debugf("Config %v=%v", key, val)
		return val
	}

	log.Debugf("Using Default Config %v=%v", key, DefaultBindAddr)
	return DefaultBindAddr
}

func DhcpDbFilePath() string {
	key := "DHCP_DB_FILE_PATH"
	if viper.IsSet(key) {
		val := viper.GetString(key)
		log.Debugf("Config %v=%v", key, val)
		return val
	}

	log.Debugf("Using Default Config %v=%v", key, DefaultDhcpDbFilePath)
	return DefaultDhcpDbFilePath
}
