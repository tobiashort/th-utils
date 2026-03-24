package main

import "github.com/tobiashort/th-utils/lib/clog"

func main() {
	clog.Level = clog.LevelDebug
	clog.Debug("Hey", "world")
	clog.Debugf("Hey %s", "you")
	clog.Debugs("What", "foo", true, "bar", true)
	clog.Info("Hey", "world")
	clog.Infof("Hey %s", "you")
	clog.Infos("What", "foo", true, "bar", true)
	clog.Warn("Hey", "world")
	clog.Warnf("Hey %s", "you")
	clog.Warns("What", "foo", true, "bar", true)
	clog.Error("Hey", "world")
	clog.Errorf("Hey %s", "you")
	clog.Errors("What", "foo", true, "bar", true)
}
