package executor

import (
	dbm "github.com/33cn/chain33/common/db"
	"github.com/33cn/chain33/common/db/table"
	"github.com/33cn/chain33/types"
	pty "github.com/33cn/plugin/plugin/dapp/finance/types"
)

func getCreditTokenFromStateDB(db dbm.KV, symbol string) (*pty.CreditTokenStateDB, error) {
	//mavl-finance-credit-symbol
	key := calcCreditNewKey(symbol)
	value, err := db.Get(key)
	if err != nil {
		return nil, err
	}

	var credit pty.CreditTokenStateDB
	if err = types.Decode(value, &credit); err != nil {
		financelog.Error("getCreditTokenFromDB", "Fail to decode types.CreditToken for key", string(key), "err info is", err)
		return nil, err
	}

	return &credit, nil
}

// AddTokenToAssets 添加个人资产列表
func AddTokenToAssets(addr string, db dbm.KVDB, symbol string) []*types.KeyValue {
	tokenAssets, err := getTokenAssetsSymbol(addr, db)
	if err != nil {
		return nil
	}
	if tokenAssets == nil {
		tokenAssets = &types.ReplyStrings{}
	}

	var found = false
	for _, sym := range tokenAssets.Datas {
		if sym == symbol {
			found = true
			break
		}
	}
	if !found {
		//如果addr中没有这种的代币就将这个symbol代币添加到addr资产中
		tokenAssets.Datas = append(tokenAssets.Datas, symbol)
		var kv []*types.KeyValue
		kv = append(kv, &types.KeyValue{Key: calcTokenAssetsKey(addr), Value: types.Encode(tokenAssets)})
		return kv
	}

	return nil
}

//getTokenAssetsSymbol 获取个人账户资产列表 tokenSymbol
func getTokenAssetsSymbol(addr string, db dbm.KVDB) (*types.ReplyStrings, error) {
	//LODB-finance-assets:addr
	key := calcTokenAssetsKey(addr)
	value, err := db.Get(key)
	if err != nil && err != types.ErrNotFound {
		financelog.Error("financeDB", "GetTokenAssetsSymbol", err)
		return nil, err
	}
	var assets types.ReplyStrings
	if err == types.ErrNotFound {
		return &assets, nil
	}
	err = types.Decode(value, &assets)
	if err != nil {
		financelog.Error("financeDB", "GetTokenAssetsSymbol", err)
		return nil, err
	}
	return &assets, nil
}

func getCreditTokenFromDB(db dbm.KV, symbol string) (*pty.CreditTokenStateDB, error) {
	key := calcCreditNewKey(symbol)
	value, err := db.Get(key)
	if err != nil {
		return nil, err
	}

	var credit pty.CreditTokenStateDB
	if err = types.Decode(value, &credit); err != nil {
		financelog.Error("getCreditTokenFromDB", "Fail to decode types.CreditToken for key", string(key), "err info is", err)
		return nil, err
	}

	return &credit, nil
}

func getBillFromDB(db dbm.KV, id string) (*pty.Bill, error) {
	key := calcBillNewKey(id)
	value, err := db.Get(key)
	if err != nil {
		return nil, err
	}
	var bill pty.Bill
	if err = types.Decode(value, &bill); err != nil {
		financelog.Error("getBillFromDB", "Fail to decode types.Bill for key", string(key), "err info is", err)
		return nil, err
	}
	return &bill, nil
}

//getCashListFromDB 获取指定IOUID的Cash列表
func getCashListFromDB(billID string, db dbm.KVDB, primary string, count int32) ([]*pty.CashList, error) {
	var cashList []*pty.CashList
	tab, err := table.NewTable(NewCashListRow(), db, getCashListOpt(billID))
	if err != nil {
		return nil, err
	}
	//列出指定数量的行数
	var primaryKey []byte
	if len(primary) != 0 {
		primaryKey = []byte(primary)
	}

	query := tab.GetQuery(db)
	rows, err := query.ListIndex("primary", nil, primaryKey, count, 0)
	if err != nil {
		financelog.Error("getCashListFromDB", "QueryErr", err)
		return nil, err
	}

	for _, v := range rows {
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
