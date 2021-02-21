package watcher

import (
	"github.com/fsnotify/fsnotify"
	"github.com/handofgod94/dhcpwatch/config"
	"github.com/handofgod94/dhcpwatch/dhcp"
	"github.com/sirupsen/logrus"
)

func Start() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		logrus.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					logrus.WithField("Ok", ok).Errorf("failed while watching event")
				}

				if event.Op&fsnotify.Write == fsnotify.Write {
					logrus.Info("received file modified event")
					leaseDb, err := dhcp.ReadDatabase(event.Name)
					if err != nil {
						logrus.WithError(err).Error("failed to read")
					}
					logrus.WithField("db", leaseDb).Debugf("lease db loaded successfully")
				}
			}
		}
	}()

	err = watcher.Add(config.DhcpDbFilePath())
	if err != nil {
		logrus.WithError(err).Fatal("failed to start watcher")
	}

	<-done
}
