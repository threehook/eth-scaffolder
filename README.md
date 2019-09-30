# eth-scaffolder
Setup your private Ethereum network with a single config file.

Sometimes you want to test your Ethereum smart contract before you deploy it to Ethereum mainnet.<br/>
Now you can use **eth-scaffolder** to setup this private network on your local machine and test any process.

**eth-scaffolder** is a commandline tool with the following usage:<br/>

| Flags                       | Description    
| ------------------------    | ------------- 
| -c, --configfile `string`   | path to yaml ethereum private network config file w/o extension (default "config/config")
| -i, --installroot `string`  | root location of installation (default "/home/`<user>`")
| -h, --help                  | prints this message


The only prerequisite is that you have a recent *`geth`* running on your local machine. It must be on your PATH.<br/>
You can download *`geth`* here: https://geth.ethereum.org/downloads/index.html<br/>

You clone the **eth-scaffolder** sources by `git clone git@github.com:threehook/eth-scaffolder.git` and build it by running `make` from the root directory.

<H4>Configuration</H4>
Configuration of the private network is defined in a yaml file.<br/>
Default the builtin config file is used (config/config.yaml). It is placed in the dist directory when building.<br/> 
It looks like this:<br/>

```
Network:
  ChainId: 17
  Difficulty: "1"
  GasLimit: "0"
  GenesisNode:
    Dir: "testdata1"
    ListenAddr: 30301
    HttpPort: 8081
      # Only one node can have the genesis accounts
    accounts:
      - PasswordFile: "./password.txt"
        # Balance in Wei
        Balance: 300000000000000000000
      - PasswordFile: "./password.txt"
        Balance: 400000000000000000000
      - PasswordFile: "./password.txt"
        Balance: 500000000000000000000
  otherNodes:
    - Dir: "testdata2"
      ListenAddr: 30302
      HttpPort: 8082
      accounts:
        - PasswordFile: "./password.txt"
    - Dir: "testdata3"
      ListenAddr: 30303
      HttpPort: 8083
      accounts:
        - PasswordFile: "./password.txt"
```

If you want to define a different network you place it anywhere on your filesystem and use the -c flag to point to it.

<H4>Password files</H4>
Every account defined in the config file needs a password file.<br/>
A password file is a text file with on the first line your chosen password.<br/>
During testing you need this password to send ether or to deploy a smart contract.<br/>
A default password.txt file is placed in the dist direcory.<br/>
The password in this default password file is: *`secret`*<br/>

If you want to use your own password file(s) you can place them anywhere and point to them from the config file:

```
- PasswordFile: "./password.txt"
```

<H4>Running a node</H4>
You can run a node like this:<br/>
`geth --config /path/to/your_config.toml console 2>> ~/node1.log`<br/>
If you want to read more about using geth in a private ethereum network try googling with similar searches like `ethereum private network`<br/>


<H4>Note</H4>
Scaffolding a new network can **not** be done when a config file's `ListenAddr` is in use during the scaffolding.
It will show an error message when you do!
 