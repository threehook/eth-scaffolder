package tmpl

import (
	"fmt"
	"github.com/threehook/eth-scaffolder/cmdline"
	"github.com/threehook/eth-scaffolder/config"
	"github.com/threehook/eth-scaffolder/util"
	"github.com/tidwall/pretty"
	"log"
	"os"
	"strings"
)

// Genesis data
type GenesisData struct {
	ChainId       uint32
	AllocAccounts []config.Account
}

type ConfigData struct {
	NetworkId   uint32
	InstallRoot string
	DataDir     string
	HttpPort    uint16
	StaticNodes string
	ListenAddr  uint16
}

// CreateGenesis creates the 'genesis.json' file based on the 'genesis.tmpl' template
func CreateGenesis() error {
	vars := make(map[string]interface{})

	accs := *config.GetNetwork().GenesisNode.Accounts
	genesisDta := GenesisData{
		ChainId:       config.GetNetwork().ChainId,
		AllocAccounts: accs,
	}
	vars["GenesisDta"] = genesisDta

	// process template file
	result := ProcessFile("tmpl/genesis.tmpl", vars)
	genesisDir := config.GetNetwork().GenesisNode.Dir
	installRoot := cmdline.CmdlineArgs().GetArg("installroot")
	prettyJson := pretty.Pretty([]byte(result))
	return writeResultToFile(string(prettyJson), installRoot.(string)+"/"+genesisDir+"/genesis.json")
}

func CreateConfigTomls(enodes map[string]string) error {
	err := createConfigTomlGenesisNode(enodes)
	if err != nil {
		return err
	}
	err = createConfigTomlOtherNodes(enodes)
	if err != nil {
		return err
	}
	return nil
}

func createConfigTomlGenesisNode(enodes map[string]string) error {
	vars := make(map[string]interface{})
	installRoot := cmdline.CmdlineArgs().GetArg("installroot").(string)
	genNode := config.GetNetwork().GenesisNode

	configDta := ConfigData{
		NetworkId:   config.GetNetwork().ChainId,
		InstallRoot: util.OsAwareFilePath(installRoot + "/"),
		DataDir:     util.OsAwareFilePath(installRoot + "/" + genNode.Dir),
		HttpPort:    genNode.HttpPort,
		StaticNodes: createStaticNodes(enodes, genNode.Dir),
		ListenAddr:  genNode.ListenAddr,
	}

	vars["ConfigDta"] = configDta

	// process template file
	result := ProcessFile("tmpl/config.tmpl", vars)
	err := writeResultToFile(result, installRoot+"/"+genNode.Dir+"/config.toml")
	if err != nil {
		return err
	}

	return nil
}

func createConfigTomlOtherNodes(enodes map[string]string) error {
	vars := make(map[string]interface{})
	installRoot := cmdline.CmdlineArgs().GetArg("installroot").(string)

	for _, node := range *config.GetNetwork().Nodes {
		configDta := ConfigData{
			NetworkId:   config.GetNetwork().ChainId,
			InstallRoot: installRoot,
			DataDir:     util.OsAwareFilePath(installRoot + "/" + node.Dir),
			HttpPort:    node.HttpPort,
			StaticNodes: createStaticNodes(enodes, node.Dir),
			ListenAddr:  node.ListenAddr,
		}

		vars["ConfigDta"] = configDta

		// process template file
		result := ProcessFile("tmpl/config.tmpl", vars)
		err := writeResultToFile(result, installRoot+"/"+node.Dir+"/config.toml")
		if err != nil {
			return err
		}
	}

	return nil
}

func createStaticNodes(enodes map[string]string, dir string) string {
	enodesFound := make([]string, 0)
	// Lookup the enodes that are not equal to dir
	for key, enode := range enodes {
		if key != dir {
			enodesFound = append(enodesFound, "\""+enode+"\"")
		}
	}
	return fmt.Sprint(strings.Join(enodesFound[:], ","))
}

func writeResultToFile(result string, filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(result)
	if err != nil {
		return err
	}

	log.Printf("Generated %v from template and wrote to filesystem", filename)
	return nil
}
