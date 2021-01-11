package sf

import "flag"

// ParseFlags returns the supplied flags
func ParseFlags() (string, bool, bool) {
	configFilePath := flag.String("config-path", "config.yml", "Path to config file")
	printVersion := flag.Bool("version", false, "Print the current version and exit")
	debug := flag.Bool("debug", false, "Enable debug mode")
	flag.Parse()

	return *configFilePath, *printVersion, *debug
}
