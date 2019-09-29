package config

import (
	"github.com/spf13/viper"
	"github.com/threehook/eth-scaffolder/cmdline"
	"log"
)

func ReadConfig() {
	config := cmdline.CmdlineArgs().GetArg("configfile").(string)
	viper.SetConfigName(config)
	viper.AddConfigPath(".")
	viper.AddConfigPath("/")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file, %s", err)
	}
}
