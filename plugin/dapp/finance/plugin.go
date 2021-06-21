package types

import (
	"github.com/33cn/chain33/pluginmgr"
	"github.com/33cn/plugin/plugin/dapp/finance/commands"
	"github.com/33cn/plugin/plugin/dapp/finance/executor"
	"github.com/33cn/plugin/plugin/dapp/finance/rpc"
	financetypes "github.com/33cn/plugin/plugin/dapp/finance/types"
)

/*
 * 初始化dapp相关的组件
 */

func init() {
	pluginmgr.Register(&pluginmgr.PluginBase{
		Name:     financetypes.FinanceX,
		ExecName: executor.GetName(),
		Exec:     executor.Init,
		Cmd:      commands.Cmd,
		RPC:      rpc.Init,
	})
}
