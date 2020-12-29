package sf

import "flag"

// ParseFlags returns the supplied flags
func ParseFlags() (string, bool) {
	configFilePath := flag.String("config-path", "config.yml", "Path to config file")
	debug := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	return *configFilePath, *debug
}
