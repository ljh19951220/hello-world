package executor

import (
	"github.com/33cn/chain33/account"
	"github.com/33cn/chain33/common/db/table"
	"github.com/33cn/chain33/types"
	pty "github.com/33cn/plugin/plugin/dapp/finance/types"
)

func (action *Action) getBalance(execer string, symbol string, addr string) (int64, error) {
	cfg := action.api.GetConfig()
	accdb, err := account.NewAccountDB(cfg, execer, symbol, action.db)
	if err != nil {
		financelog.Error("getBalance", "NewAccountDB", err)
		return 0, err
	}
	acc := accdb.LoadExecAccount(addr, action.execaddr)

	return acc.Balance, nil
}

func (action *Action) transferAssets(execer string, from string, to string, amount int64, symbol string) (*types.Receipt, error) {
	cfg := action.api.GetConfig()
	accdb, err := account.NewAccountDB(cfg, execer, symbol, action.db)
	if err != nil {
		financelog.Error("transferAssets", "NewAccountDB", err)
		return nil, err
	}

	receipt, err := accdb.ExecTransfer(from, to, action.execaddr, amount)
	if err != nil {
		financelog.Debug("transferAssets", "ExecTransfer", err)
		return nil, err
	}

	return receipt, nil
}

func (action *Action) frozenAssets(execer string, symbol string, addr string, amount int64) (*types.Receipt, error) {
	cfg := action.api.GetConfig()
	accdb, err := account.NewAccountDB(cfg, execer, symbol, action.db)
	if err != nil {
		financelog.Error("frozenAssets", "NewAccountDB", err)
		return nil, err
	}

	receipt, err := accdb.ExecFrozen(addr, action.execaddr, amount)
	if err != nil {
		financelog.Error("Fail to frozen assets", "assets owner", addr, "assets symbol", symbol)
		return nil, err
	}

	return receipt, nil
}

func (action *Action) getFrozenAsset(execer string, symbol string, addr string) (int64, error) {
	cfg := action.api.GetConfig()
	accdb, err := account.NewAccountDB(cfg, execer, symbol, action.db)
	if err != nil {
		financelog.Error("getFrozenAsset", "NewAccountDB", err)
		return 0, err
	}
	acc := accdb.LoadExecAccount(addr, action.execaddr)

	return acc.Frozen, nil
}

//activeFrozenAssets 激活冻结的资金
func (action *Action) activeFrozenAssets(execer string, addr string, amount int64, symbol string) (*types.Receipt, error) {
	cfg := action.api.GetConfig()
	accdb, err := account.NewAccountDB(cfg, execer, symbol, action.db)
	if err != nil {
		financelog.Error("frozenAssets", "NewAccountDB", err)
		return nil, err
	}

	receipt, err := accdb.ExecActive(addr, action.execaddr, amount)
	if err != nil {
		financelog.Error("activeFrozenAssets", "ExecActive", err, "Active frozen amount", amount)
		return nil, err
	}

	return receipt, nil
}

func (action *Action) getDeliverItem(deliverID string) (*pty.DeliverItemLocalDB, error) {
	tab, err := table.NewTable(NewDeliverBillListRow(), action.localdb, optDeliverBillList)
	if err != nil {
		financelog.Error("getDeliverItem", "new deliver list init err ", err)
		return nil, err
	}
	query := tab.GetQuery(action.localdb)
	rows, err := query.ListIndex("deliverID", []byte(deliverID), nil, 1, 0)
	if err != nil {
		financelog.Error("getDeliverItem", "QueryErr", err)
		return nil, err
	}
	if len(rows) == 0 {
		financelog.Debug("Current deliver item not exist in deliver list", "deliverID", deliverID)
		return nil, pty.ErrDeliverItemNotFound
	}

	local := rows[0].Data.(*pty.DeliverItemLocalDB)

	return local, nil
}

func (action *Action) transferFrozenAssets(execer string, from string, to string, amount int64, symbol string) (*types.Receipt, error) {
	cfg := action.api.GetConfig()
	accdb, err := account.NewAccountDB(cfg, execer, symbol, action.db)
	if err != nil {
		financelog.Error("frozenAssets", "NewAccountDB", err)
		return nil, err
	}

	receipt, err := accdb.ExecTransferFrozen(from, to, action.execaddr, amount)
	if err != nil {
		financelog.Debug("transferFrozenAssets", "ExecTransferFrozen", err)
		return nil, err
	}

	return receipt, nil
}

func (action *Action) getSplitBillItem(txHash string) (*pty.SplitBillRecordLocalDB, error) {
	tab, err := table.NewTable(NewSplitBillListRow(), action.localdb, optSplitBillList)
	if err != nil {
		financelog.Error("getSplitBillItem", "new Split list init err ", err)
		return nil, err
	}
	query := tab.GetQuery(action.localdb)
	rows, err := query.ListIndex("txHash", []byte(txHash), nil, 1, 0)
	if err != nil {
		financelog.Error("getSplitBillItem", "QueryErr", err)
		return nil, err
	}
	if len(rows) == 0 {
		financelog.Debug("Current deliver item not exist in deliver list", "txHash", txHash)
		return nil, pty.ErrSplitItemNotFound
	}

	local := rows[0].Data.(*pty.SplitBillRecordLocalDB)

	return local, nil
}

//getApplyFinancingItem 获取指定申请哈息的融资申请条目
func (action *Action) getApplyFinancingItem(txHash string) (*pty.ApplyFinancingItemLocalDB, error) {
	tab, err := table.NewTable(NewApplyForFinancingListRow(), action.localdb, optApplyForFinancingList)
	if err != nil {
		financelog.Error("getApplyFinancingItem", "new applyFinancing list init err ", err)
		return nil, err
	}
	query := tab.GetQuery(action.localdb)
	rows, err := query.ListIndex("txHash", []byte(txHash), nil, 1, 0)
	//rows, err := query.List("txHash", []byte(txHash), nil, 1, 0)
	if err != nil {
		financelog.Error("getApplyFinancingItem", "QueryErr", err)
		return nil, err
	}
	if len(rows) == 0 {
		financelog.Debug("Current applyFinancing item not exist in applyFinancing list", "TxHash", txHash)
		return nil, pty.ErrApplyFinancingNotFound
	}

	local := rows[0].Data.(*pty.ApplyFinancingItemLocalDB)

	return local, nil
}

//getCashList 获取指定BillID的Cash列表
func (action *Action) getCashList(billID string) ([]*pty.CashList, error) {
	var cashList []*pty.CashList
	tab, err := table.NewTable(NewCashListRow(), action.localdb, getCashListOpt(billID))
	if err != nil {
		return nil, err
	}
	query := tab.GetQuery(action.localdb)
	rows, err := query.ListIndex("primary", nil, nil, 0, 0)
	if err != nil {
		financelog.Error("Query_CashList", "QueryErr", err)
		return nil, err
	}

	for _, v := range rows {
		//根据白条ID可以获取所有的关于指定白条id的申请兑现列表信息
		cashItem, ok := v.Data.(*pty.CashList)
		if !ok {
			financelog.Error("Get cash item failed")
			return nil, err
		}
		cashList = append(cashList, cashItem)
	}
	financelog.Debug("Fetch Cash List", "items", cashList)

	return cashList, nil
}
