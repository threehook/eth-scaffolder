package cmdline

import "sync"

type Args interface {
	GetArgs() map[string]interface{}
	AddArgs(map[string]interface{})
	AddArg(string, interface{})
	GetArg(key string) interface{}
}

type ethscafArgs struct {
	argsMap map[string]interface{}
}

var instance *ethscafArgs
var once sync.Once

func CmdlineArgs() Args {
	once.Do(func() {
		newMap := make(map[string]interface{})
		instance = &ethscafArgs{argsMap: newMap}
	})

	return instance
}

func (a *ethscafArgs) GetArgs() map[string]interface{} {
	return instance.argsMap
}

func (a *ethscafArgs) AddArgs(aMap map[string]interface{}) {
	for k, v := range aMap {
		instance.argsMap[k] = v
	}
}

func (a *ethscafArgs) AddArg(key string, value interface{}) {
	instance.argsMap[key] = value
}

func (a *ethscafArgs) GetArg(key string) interface{} {
	return instance.argsMap[key]
}
