package watcher

import (
	"github.com/fsnotify/fsnotify"
	"github.com/handofgod94/dhcpwatch/config"
	"github.com/handofgod94/dhcpwatch/dhcp"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sirupsen/logrus"
	"net/http"
)

func Start() {
	logrus.SetLevel(config.LogLevel())
	logrus.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logrus.Fatal(err)
	}
	defer watcher.Close()

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					logrus.WithField("Ok", ok).Errorf("failed while watching event")
					return
				}

				if event.Op&fsnotify.Write == fsnotify.Write {
					logrus.Info("received file modified event")
					leaseDb, err := dhcp.ReadDatabase(event.Name)
					if err != nil {
						logrus.WithError(err).Error("failed to read")
					}
					logrus.WithField("db", leaseDb).Debugf("lease db loaded successfully")
				}
			case err := <-watcher.Errors:
				logrus.WithError(err).Error("failed to watch file for changes")
			}
		}
	}()

	err = watcher.Add(config.DhcpDbFilePath())
	if err != nil {
		logrus.WithError(err).Fatal("failed to start watcher")
	}

	logrus.WithField("db", config.DhcpDbFilePath()).Info("started watcher")

	logrus.WithField("Address", config.Addr()).Info("Starting server")
	http.Handle("/metrics", promhttp.Handler())
	logrus.Fatal(http.ListenAndServe(config.Addr(), nil))
}
