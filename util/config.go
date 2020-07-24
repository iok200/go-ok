package util

import "github.com/magiconair/properties"

func GetConfig(path ...string) *properties.Properties {
	var pf []string
	if path == nil || len(path) == 0 {
		pf = []string{"ok.properties"}
	} else {
		pf = path
	}
	return properties.MustLoadFiles(pf, properties.UTF8, true)
}
