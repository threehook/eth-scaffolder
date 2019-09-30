package scaffold

import (
	"fmt"
	"github.com/threehook/eth-scaffolder/cmdline"
	"github.com/threehook/eth-scaffolder/config"
	"github.com/threehook/eth-scaffolder/tmpl"
	"github.com/threehook/eth-scaffolder/util"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func ScaffoldNetwork() {
	var err error

	removeNodeDirs()
	createNodeDirs()
	err = createGenesisAccounts()
	if err != nil {
		log.Fatalf("Could not create the genesis account: %v", err)
	}
	err = createOtherNodesAccounts()
	if err != nil {
		log.Fatalf("Could not create nodes accounts: %v", err)
	}

	err = tmpl.CreateGenesis()
	if err != nil {
		log.Fatalf("Could not create the genesis file (genesis.json): %v", err)
	}
	copyGenesisToOtherNodes()
	err = initializeNodes()
	if err != nil {
		log.Fatalf("Could not initialize node(s): %v", err)
	}
	enodes, err := collectEnodes()
	if err != nil {
		log.Fatalf("Could not collect enodes: %v", err)
	}

	err = tmpl.CreateConfigTomls(enodes)
	if err != nil {
		log.Fatalf("Could not create config.toml file(s): %v", err)
	}

	showPublicKeys()
}

func removeNodeDirs() {
	installroot := cmdline.CmdlineArgs().GetArg("installroot")
	genDir := installroot.(string) + "/" + config.GetNetwork().GenesisNode.Dir

	_ = os.RemoveAll(genDir)
	for _, node := range *config.GetNetwork().Nodes {
		nodeDir := installroot.(string) + "/" + node.Dir
		_ = os.RemoveAll(nodeDir)
	}
}

func createNodeDirs() {
	installroot := cmdline.CmdlineArgs().GetArg("installroot")
	genDir := installroot.(string) + "/" + config.GetNetwork().GenesisNode.Dir
	_ = exec.Command("mkdir", "-p", genDir).Run()

	for _, node := range *config.GetNetwork().Nodes {
		nodeDir := node.Dir
		_ = exec.Command("mkdir", "-p", installroot.(string)+"/"+nodeDir).Run()
	}
}

func createGenesisAccounts() error {
	installroot := cmdline.CmdlineArgs().GetArg("installroot")
	dir := config.GetNetwork().GenesisNode.Dir
	genDir := installroot.(string) + "/" + dir
	accs := *config.GetNetwork().GenesisNode.Accounts
	for i, acc := range accs {
		b, err := exec.Command("geth", "--datadir", genDir, "--password", acc.PasswordFile, "account", "new").Output()
		if err != nil {
			return err
		}

		accs[i].PublicKey = pubKeyFromAccountOutput(string(b))
		log.Printf("Created public key %v for account under node %v", accs[i].PublicKey, dir)
	}

	return nil
}

func pubKeyFromAccountOutput(output string) string {
	re := regexp.MustCompile("0x([0-9a-fA-F]+)")
	return fmt.Sprintf("%s", re.Find([]byte(output)))
}

func copyGenesisToOtherNodes() {
	// get the genesis node dir
	installroot := cmdline.CmdlineArgs().GetArg("installroot")
	genDir := installroot.(string) + "/" + config.GetNetwork().GenesisNode.Dir

	for _, node := range *config.GetNetwork().Nodes {
		nodeDir := node.Dir
		srcFile := genDir + "/genesis.json"
		targetDir := installroot.(string) + "/" + nodeDir
		if err := util.CopyToDir(srcFile, targetDir); err != nil {
			log.Fatal(err)
		}
	}
}

func createOtherNodesAccounts() error {
	installroot := cmdline.CmdlineArgs().GetArg("installroot")
	for _, node := range *config.GetNetwork().Nodes {
		nodeDir := installroot.(string) + "/" + node.Dir
		accs := *node.Accounts
		for i, acc := range accs {
			b, err := exec.Command("geth", "--datadir", nodeDir, "--password", acc.PasswordFile, "account", "new").Output()
			if err != nil {
				return err
			}
			accs[i].PublicKey = pubKeyFromAccountOutput(string(b))
			log.Printf("Created public key %v for account under node %v", accs[i].PublicKey, node.Dir)
		}
	}
	return nil
}

func initializeNodes() error {
	var err error
	genNode := config.GetNetwork().GenesisNode
	installroot := cmdline.CmdlineArgs().GetArg("installroot")

	genDir := filepath.FromSlash(installroot.(string) + "/" + genNode.Dir)
	err = exec.Command("geth", "--nousb", "--datadir", genDir, "init", genDir+"/genesis.json").Run()
	if err != nil {
		return err
	}

	for _, node := range *config.GetNetwork().Nodes {
		nodeDir := installroot.(string) + "/" + node.Dir
		err = exec.Command("geth", "--nousb", "--datadir", nodeDir, "init", nodeDir+"/genesis.json").Run()
		if err != nil {
			return err
		}
	}

	return nil
}

func collectEnodes() (map[string]string, error) {
	var err error
	var b []byte

	enodeMap := make(map[string]string)
	chainId := config.GetNetwork().ChainId
	genNode := config.GetNetwork().GenesisNode
	installroot := cmdline.CmdlineArgs().GetArg("installroot")

	genDir := installroot.(string) + "/" + genNode.Dir
	b, err = exec.Command("geth", "--datadir", genDir, "--nousb", "console", "--networkid", fmt.Sprint(chainId)).CombinedOutput()
	if err != nil {
		return nil, err
	}
	enode := enodeFromNodeInitOutput(string(b)) + "@127.0.0.1:" + strconv.Itoa(int(genNode.ListenAddr))
	enodeMap[genNode.Dir] = enode

	for _, node := range *config.GetNetwork().Nodes {
		nodeDir := installroot.(string) + "/" + node.Dir
		b, err = exec.Command("geth", "--datadir", nodeDir, "--nousb", "console", "--networkid", fmt.Sprint(chainId)).CombinedOutput()
		if err != nil {
			return nil, err
		}
		port := node.ListenAddr
		enode := enodeFromNodeInitOutput(string(b)) + "@127.0.0.1:" + strconv.Itoa(int(port))
		enodeMap[node.Dir] = enode
	}

	return enodeMap, nil
}

func showPublicKeys() {
	fmt.Println("\nThe following public keys were generated:")
	fmt.Println(config.GetNetwork().GenesisNode.Dir + ":")
	for _, acc := range *config.GetNetwork().GenesisNode.Accounts {
		fmt.Println("\t" + acc.PublicKey)
	}
	for _, node := range *config.GetNetwork().Nodes {
		fmt.Println(node.Dir + ":")
		for _, acc := range *node.Accounts {
			fmt.Println("\t" + acc.PublicKey)
		}
	}
}

func enodeFromNodeInitOutput(output string) string {
	re := regexp.MustCompile("self=enode://([0-9a-fA-F]+)")
	s := fmt.Sprintf("%s", re.Find([]byte(output)))
	i := strings.LastIndexByte(s, '=')
	s = s[i+1:]

	return s
}
