package config

import (
	"github.com/spf13/viper"
	"log"
	"sync"
)

var instance *Network
var once sync.Once

// GetNetworkSetup returns the single instance of NetworkSetup
func GetNetwork() *Network {
	once.Do(func() {
		instance = unmarshalNetwork()
	})
	return instance
}

// Network
type Network struct {
	ChainId     uint32  `mapstructure:"ChainId"`
	Difficulty  string  `mapstructure:"Difficulty"`
	GasLimit    string  `mapstructure:"GasLimit"`
	GenesisNode Node    `mapstructure:"GenesisNode"`
	Nodes       *[]Node `mapstructure:"otherNodes"`
}

type Node struct {
	Dir        string     `mapstructure:"Dir"`
	ListenAddr uint16     `mapstructure:"ListenAddr"`
	HttpPort   uint16     `mapstructure:"HttpPort"`
	Enode      string     `mapstructure:"Enode"`
	Accounts   *[]Account `mapstructure:"accounts"`
}

type Account struct {
	PublicKey    string `mapstructure:"PublicKey"`
	PasswordFile string `mapstructure:"PasswordFile"`
	Balance      string `mapstructure:"Balance"`
}

func unmarshalNetwork() *Network {
	var network Network
	err := viper.Sub("Network").Unmarshal(&network)
	if err != nil {
		log.Fatalf("Yaml struct 'Network' in config.yaml does not have the right format, %v", err)
	}
	return &network
}
