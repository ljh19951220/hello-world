package executor

import (
	log "github.com/33cn/chain33/common/log/log15"
	drivers "github.com/33cn/chain33/system/dapp"
	"github.com/33cn/chain33/types"
	financetypes "github.com/33cn/plugin/plugin/dapp/finance/types"
)

/*
 * 执行器相关定义
 * 重载基类相关接口
 */

var (
	//日志
	elog = log.New("module", "finance.executor")
)

var driverName = financetypes.FinanceX

// Init register dapp
func Init(name string, cfg *types.Chain33Config, sub []byte) {
	drivers.Register(cfg, GetName(), newFinance, cfg.GetDappFork(driverName, "Enable"))
	InitExecType()
}

// InitExecType Init Exec Type
func InitExecType() {
	ety := types.LoadExecutorType(driverName)
	ety.InitFuncList(types.ListMethod(&finance{}))
}

type finance struct {
	drivers.DriverBase
}

func newFinance() drivers.Driver {
	t := &finance{}
	t.SetChild(t)
	t.SetExecutorType(types.LoadExecutorType(driverName))
	return t
}

// GetName get driver name
func GetName() string {
	return newFinance().GetName()
}

func (f *finance) GetDriverName() string {
	return driverName
}

// CheckTx 实现自定义检验交易接口，供框架调用
func (f *finance) CheckTx(tx *types.Transaction, index int) error {
	// implement code
	return nil
}
