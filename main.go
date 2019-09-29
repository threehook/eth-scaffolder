package main

import (
	"fmt"
	flag "github.com/spf13/pflag"
	"github.com/threehook/eth-scaffolder/cmdline"
	"github.com/threehook/eth-scaffolder/config"
	"github.com/threehook/eth-scaffolder/scaffold"
	"log"
	"os"
)

var (
	user = userHome()

	configfile  = flag.StringP("configfile", "c", "config/config", "path to yaml ethereum private network config file w/o extension")
	installroot = flag.StringP("installroot", "i", user, "root location of installation")
	help        = flag.BoolP("help", "h", false, "prints this message")
)

func userHome() string {
	user, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Could not determine users' home dir")
	}
	return user
}

func main() {
	_ = processArgs()
	// Read the Viper config
	config.ReadConfig()
	scaffold.ScaffoldNetwork()
}

func processArgs() map[string]interface{} {
	initMap := make(map[string]interface{})

	os.Args = os.Args[1:]
	flag.CommandLine.SortFlags = false
	_ = flag.CommandLine.Parse(os.Args)

	if *help {
		showUsage()
		os.Exit(0)
	}

	initMap["configfile"] = *configfile
	initMap["installroot"] = *installroot

	cmdlineArgs := cmdline.CmdlineArgs()
	cmdlineArgs.AddArgs(initMap)

	return cmdlineArgs.GetArgs()
}

func showUsage() {
	flag.CommandLine.SortFlags = false
	_, _ = fmt.Fprintln(os.Stderr, "Usage:")
	_, _ = fmt.Fprint(os.Stderr, flag.CommandLine.FlagUsagesWrapped(120))
}
