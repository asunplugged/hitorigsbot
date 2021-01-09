package options

import (
	"flag"
)

// Options - flag options
type Options struct {
	Help   bool
	Config string
	DBTool bool
}

var opts *Options

// RegFlag - register flag
func RegFlag() {
	opts = &Options{}
	flag.StringVar(&opts.Config, "config", "", "config file path (defualt {PWD}/config.yml")
	flag.StringVar(&opts.Config, "f", "", "config file path (short) (defualt {PWD}/config.yml")
	flag.BoolVar(&opts.DBTool, "dbtool", false, "run dbtool deploy schema")
	flag.BoolVar(&opts.Help, "help", false, "show help")
}

// GetFlag -
func GetFlag() *Options {
	return opts
}
