package sf

import flag "github.com/spf13/pflag"

// ParseFlags returns the supplied flags
func ParseFlags() (string, bool, bool) {
	configFilePath := flag.StringP("config-path", "c", "config.yml", "Path to config file")
	printVersion := flag.BoolP("version", "v", false, "Print the current version and exit")
	debug := flag.BoolP("debug", "d", false, "Enable debug mode")
	flag.Parse()

	return *configFilePath, *printVersion, *debug
}
