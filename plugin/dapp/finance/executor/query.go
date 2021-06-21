package executor

import (
	"fmt"
	"github.com/33cn/chain33/account"
	"github.com/33cn/chain33/common/db/table"
	"github.com/33cn/chain33/system/dapp"
	"github.com/33cn/chain33/types"
	financetypes "github.com/33cn/plugin/plugin/dapp/finance/types"
	pty "github.com/33cn/plugin/plugin/dapp/finance/types"
)

/***************************query credit 查询授信相关操作********************************************************************/
//Query_CreditDetail query credit detail
func (f *finance) Query_CreditDetail(in *financetypes.ReqCreditDetail) (types.Message, error) {
	var reply financetypes.ReplyCreditDetail

	credit, err := getCreditTokenFromDB(f.GetStateDB(), in.CreditSymbol)
	if err != nil {
		financelog.Error("Query_CreditDetail", "getCreditTokenFromDB", err)
		return nil, err
	}
	reply.Credit = credit

	return &reply, nil
}

//Query_CreditForFunder query credit list for funder
func (f *finance) Query_CreditForFunder(in *financetypes.ReqCreditFunder) (types.Message, error) {
	var reply financetypes.ReplyCreditFunder
	tab, err := table.NewTable(NewCreditRecordRow(), f.GetLocalDB(), optCreditRecordLocalDB)
	if err != nil {
		return nil, err
	}
	query := tab.GetQuery(f.GetLocalDB())

	//列出指定数量的行数
	var primaryKey []byte
	if len(in.PrimaryKey) != 0 {
		primaryKey = []byte(in.PrimaryKey)
	}
	rows, err := query.ListIndex("creditAddr", []byte(in.FunderAddr), primaryKey, in.EntrieCount, 0)
	if err != nil {
		financelog.Error("Query_CreditForFunder", "QueryErr", err)
		return nil, err
	}

	for _, v := range rows {
		local := v.Data.(*financetypes.CreditRecordLocalDB)
		reply.Items = append(reply.Items, local)
		reply.PrimaryKey = string(v.Primary) //主键是授信Token的symbol
	}
	reply.BlockTime = f.GetBlockTime()

	financelog.Debug("Fetch list from credit record for funder", "items", reply.Items, "Last primary key", reply.PrimaryKey, "blockTime", reply.BlockTime)

	return &reply, nil
}

//Query_CreditForCoreFirm query credit list for core firm
func (f *finance) Query_CreditForCoreFirm(in *financetypes.ReqCreditCoreFirm) (types.Message, error) {
	var reply financetypes.ReplyCreditCoreFirm
	tab, err := table.NewTable(NewCreditRecordRow(), f.GetLocalDB(), optCreditRecordLocalDB)
	if err != nil {
		return nil, err
	}
	query := tab.GetQuery(f.GetLocalDB())

	//列出指定数量的行数
	var primaryKey []byte
	if len(in.PrimaryKey) != 0 {
		primaryKey = []byte(in.PrimaryKey)
	}
	rows, err := query.ListIndex("granteeAddr", []byte(in.CoreAddr), primaryKey, in.EntrieCount, 0)
	if err != nil {
		financelog.Error("Query_CreditForCoreFirm", "QueryErr", err)
		return nil, err
	}

	execAddr := dapp.ExecAddress(f.GetAPI().GetConfig().GetTitle() + financetypes.FinanceX)

	for _, v := range rows {
		local := v.Data.(*financetypes.CreditRecordLocalDB)
		//通过StateDB查询授信代币信息
		acct := f.GetAcctInfo(local.Symbol, execAddr, in.CoreAddr)
		if acct == nil || (acct.Balance+acct.Frozen) == 0 {
			continue
		}
		asset := &pty.AssetItem{Balance: acct.Balance, Frozen: acct.Frozen, Symbol: local.Symbol}
		item := &pty.CreditAsset{AssetItem: asset, CreditItem: local}
		reply.Items = append(reply.Items, item)
		reply.PrimaryKey = string(v.Primary)
	}
	reply.BlockTime = f.GetBlockTime()

	financelog.Debug("Fetch list from credit record for core firm", "items", reply.Items, "Last primary key", reply.PrimaryKey, "blockTime", reply.BlockTime)

	return &reply, nil
}

//GetAcctInfo get account info
func (f *finance) GetAcctInfo(symbol string, execAddr string, addr string) *types.Account {
	cfg := f.GetAPI().GetConfig()
	exer := pty.TokenExecer
	if symbol == pty.CCNY {
		exer = pty.CoinsExecer
	}

	accdb, err := account.NewAccountDB(cfg, exer, symbol, f.GetStateDB())
	if err != nil {
		financelog.Error("GetAcctInfo", "NewAccountDB", err)
		return nil
	}
	acc := accdb.LoadExecAccount(addr, execAddr)

	return acc
}

/************************查询出入金(DWRecord)相关操作信息存款操作*******************************************************/
//Query_DWRecord query cash list
func (f *finance) Query_DWRecord(in *financetypes.ReqDWRecord) (types.Message, error) {
	var reply financetypes.ReplyDWRecord
	var items []*pty.DWRecordLocalDB

	tab, err := table.NewTable(NewDWRecordRow(), f.GetLocalDB(), getDWRecordOpt(in.Addr))
	if err != nil {
		return nil, err
	}
	//列出指定数量的行数
	var primaryKey []byte
	if len(in.PrimaryKey) != 0 {
		primaryKey = []byte(in.PrimaryKey)
	}

	query := tab.GetQuery(f.GetLocalDB())
	rows, err := query.ListIndex("primary", nil, primaryKey, in.EntrieCount, 0)
	if err != nil {
		financelog.Error("Query_DWRecord", "QueryErr", err)
		return nil, err
	}

	for _, v := range rows {
		recordItem, ok := v.Data.(*pty.DWRecordLocalDB)
		if !ok {
			financelog.Error("Get DWRecord item failed")
			return nil, err
		}
		items = append(items, recordItem)
	}

	reply.Items = append(reply.Items, items...)
	financelog.Debug("Fetch deposits and withdrawals records", "Items", reply)

	return &reply, nil
}

//Query_DWRecordBySymbol query cash list
func (f *finance) Query_DWRecordBySymbol(in *financetypes.ReqDWRecord) (types.Message, error) {
	var reply financetypes.ReplyDWRecord
	var items []*pty.DWRecordLocalDB

	tab, err := table.NewTable(NewDWRecordRow(), f.GetLocalDB(), getDWRecordOpt(in.Addr))
	if err != nil {
		return nil, err
	}
	//列出指定数量的行数
	var primaryKey []byte
	if len(in.PrimaryKey) != 0 {
		primaryKey = []byte(in.PrimaryKey)
	}
	query := tab.GetQuery(f.GetLocalDB())
	rows, err := query.ListIndex("symbol", []byte(in.AssetSymbol), primaryKey, in.EntrieCount, 0)
	if err != nil {
		financelog.Error("Query_DWRecordBySymbol", "QueryErr", err)
		return nil, err
	}

	for _, v := range rows {
		recordItem, ok := v.Data.(*pty.DWRecordLocalDB)
		if !ok {
			financelog.Error("Get DWRecord item failed")
			return nil, err
		}
		items = append(items, recordItem)
	}

	reply.Items = append(reply.Items, items...)
	financelog.Debug("Fetch deposits and withdrawals records", "Items", reply)

	return &reply, nil
}

//Query_DWRecordByAddr query cash list 查的就是这个人的出入金金额
func (f *finance) Query_DWRecordByAddr(in *financetypes.ReqDWRecord) (types.Message, error) {
	var reply financetypes.ReplyDWRecord
	var items []*pty.DWRecordLocalDB

	tab, err := table.NewTable(NewDWRecordRow(), f.GetLocalDB(), getDWRecordOpt(in.Addr))
	if err != nil {
		return nil, err
	}
	//列出指定数量的行数
	var primaryKey []byte
	if len(in.PrimaryKey) != 0 {
		primaryKey = []byte(in.PrimaryKey)
	}
	query := tab.GetQuery(f.GetLocalDB())
	rows, err := query.ListIndex("addr", []byte(in.Addr), primaryKey, in.EntrieCount, 0)
	if err != nil {
		financelog.Error("Query_DWRecordByAddr", "QueryErr", err)
		return nil, err
	}

	for _, v := range rows {
		recordItem, ok := v.Data.(*pty.DWRecordLocalDB)
		if !ok {
			financelog.Error("Get DWRecord item failed")
			return nil, err
		}
		items = append(items, recordItem)
	}

	reply.Items = append(reply.Items, items...)
	financelog.Debug("Fetch deposits and withdrawals records", "Items", reply)

	return &reply, nil
}

/***********************************QueryAccountAsset查询合约子账户资产信息*******************************************/
//Query_AccountAsset query account asset
func (f *finance) Query_AccountAsset(in *financetypes.ReqAccountAsset) (types.Message, error) {
	var reply financetypes.ReplyAccountAsset

	cfg := f.GetAPI().GetConfig()
	//通过StateDB查询用户在合约finance下的资产信息
	accdb, err := account.NewAccountDB(cfg, in.Execer, in.Symbol, f.GetStateDB())
	if err != nil {
		financelog.Error("Query_AccountAsset", "NewAccountDB", err)
		return nil, err
	}

	//获取执行器地址
	execAddr := dapp.ExecAddress(cfg.GetTitle() + financetypes.FinanceX)
	acc := accdb.LoadExecAccount(in.Addr, execAddr)

	reply.Addr = acc.Addr
	reply.Balance = acc.Balance
	reply.Frozen = acc.Frozen

	return &reply, nil
}

//Query_UserAssets query user assets
func (f *finance) Query_UserAssets(in *financetypes.ReqFinanceUserAssets) (types.Message, error) {
	var reply financetypes.ReplyFinanceUserAssets
	reply.Addr = in.Addr

	//assets是字符串形式的，所以这个查询是将用户所有的资产信息都查出来
	//AddTokenToAssets(v.ToAddr, f.GetLocalDB(), v.Symbol)
	assets, err := getTokenAssetsSymbol(in.Addr, f.GetLocalDB())
	if err != nil {
		return nil, err
	}
	if assets == nil {
		return &reply, nil
	}

	execAddr := dapp.ExecAddress(f.GetAPI().GetConfig().GetTitle() + financetypes.FinanceX)

	for _, symbol := range assets.Datas {
		exer := pty.TokenExecer
		if symbol == pty.CCNY {
			exer = pty.CoinsExecer
		}
		cfg := f.GetAPI().GetConfig()
		accdb, err := account.NewAccountDB(cfg, exer, symbol, f.GetStateDB())
		if err != nil {
			financelog.Error("Query_UserAssets", "NewAccountDB", err)
			return nil, err
		}
		acc := accdb.LoadExecAccount(in.Addr, execAddr)
		item := &pty.AssetItem{}
		item.Balance = acc.Balance
		item.Frozen = acc.Frozen
		item.Symbol = symbol
		reply.AssetItems = append(reply.AssetItems, item)
	}

	return &reply, nil
}

/***********************************TransferCoinsRecord查询合约内转账记录*******************************************/
//Query_DWRecord query cash list
func (f *finance) Query_TransferCoinsRecord(in *financetypes.ReqTransferCoinsRecord) (types.Message, error) {
	var reply financetypes.ReplyTransferCoinsRecord
	var items []*pty.TransferCoinsRecordLocalDB

	tab, err := table.NewTable(NewTransferCoinsRecordRow(), f.GetLocalDB(), getTransferCoinsRecordOpt(in.Addr))
	if err != nil {
		return nil, err
	}
	//列出指定数量的行数
	var primaryKey []byte
	if len(in.PrimaryKey) != 0 {
		primaryKey = []byte(in.PrimaryKey)
	}

	query := tab.GetQuery(f.GetLocalDB())
	rows, err := query.ListIndex("primary", nil, primaryKey, in.EntrieCount, 0)
	if err != nil {
		financelog.Error("Query_TransferCoinsRecord", "QueryErr", err)
		return nil, err
	}

	for _, v := range rows {
		recordItem, ok := v.Data.(*pty.TransferCoinsRecordLocalDB)
		if !ok {
			financelog.Error("Get TransferCoinsRecord item failed")
			return nil, err
		}
		items = append(items, recordItem)
	}

	reply.Items = append(reply.Items, items...)
	financelog.Debug("Fetch TransferCoinsRecord records", "Items", reply)

	return &reply, nil
}

//Query_DWRecordBySymbol query cash list
func (f *finance) Query_TransferCoinsRecordBySymbol(in *financetypes.ReqTransferCoinsRecord) (types.Message, error) {
	var reply financetypes.ReplyTransferCoinsRecord
	var items []*pty.TransferCoinsRecordLocalDB

	tab, err := table.NewTable(NewTransferCoinsRecordRow(), f.GetLocalDB(), getTransferCoinsRecordOpt(in.Addr))
	if err != nil {
		return nil, err
	}
	//列出指定数量的行数
	var primaryKey []byte
	if len(in.PrimaryKey) != 0 {
		primaryKey = []byte(in.PrimaryKey)
	}

	query := tab.GetQuery(f.GetLocalDB())
	rows, err := query.ListIndex("symbol", []byte(in.AssetSymbol), primaryKey, in.EntrieCount, 0)
	if err != nil {
		financelog.Error("Query_TransferCoinsRecord", "QueryErr", err)
		return nil, err
	}

	for _, v := range rows {
		recordItem, ok := v.Data.(*pty.TransferCoinsRecordLocalDB)
		if !ok {
			financelog.Error("Get TransferCoinsRecord item failed")
			return nil, err
		}
		items = append(items, recordItem)
	}

	reply.Items = append(reply.Items, items...)
	financelog.Debug("Fetch TransferCoinsRecord records", "Items", reply)

	return &reply, nil
}

//Query_DWRecordByAddr query cash list 查的就是这个人的出入金金额
func (f *finance) Query_TransferCoinsRecordByAddr(in *financetypes.ReqTransferCoinsRecord) (types.Message, error) {
	var reply financetypes.ReplyTransferCoinsRecord
	var items []*pty.TransferCoinsRecordLocalDB

	tab, err := table.NewTable(NewTransferCoinsRecordRow(), f.GetLocalDB(), getTransferCoinsRecordOpt(in.Addr))
	if err != nil {
		return nil, err
	}
	//列出指定数量的行数
	var primaryKey []byte
	if len(in.PrimaryKey) != 0 {
		primaryKey = []byte(in.PrimaryKey)
	}

	query := tab.GetQuery(f.GetLocalDB())
	rows, err := query.ListIndex("fromAddr", []byte(in.Addr), primaryKey, in.EntrieCount, 0)
	if err != nil {
		financelog.Error("Query_TransferCoinsRecord", "QueryErr", err)
		return nil, err
	}

	for _, v := range rows {
		recordItem, ok := v.Data.(*pty.TransferCoinsRecordLocalDB)
		if !ok {
			financelog.Error("Get TransferCoinsRecord item failed")
			return nil, err
		}
		items = append(items, recordItem)
	}

	reply.Items = append(reply.Items, items...)
	financelog.Debug("Fetch TransferCoinsRecord records", "Items", reply)

	return &reply, nil
}

/***********************************Query_BillInfo查询白条信息*******************************************/
//Query_BillInfo query Bill information
func (f *finance) Query_BillInfo(in *financetypes.ReqBillInfo) (types.Message, error) {
	var reply financetypes.ReplyBillInfo
	var bill financetypes.Bill

	//从StateDB数据库中查询白条信息
	key := calcBillNewKey(in.GetId())
	oldVal, err := f.GetStateDB().Get(key)
	if err != nil && err != types.ErrNotFound {
		return nil, err
	}
	err = types.Decode(oldVal, &bill)
	if err != nil {
		elog.Error("Query_BillInfo", "DecodeErr", err)
		return nil, types.ErrDecode
	}
	reply.BillInfo = &bill

	return &reply, nil
}

/***********************************Query_BillInfo查询借款人相关的白条信息*******************************************/
//Query_BorrowerBill query borrower's relate Bill
func (f *finance) Query_BorrowerBill(in *financetypes.ReqBorrowerBill) (types.Message, error) {
	var reply financetypes.ReplyBorrowerBill

	tab, err := table.NewTable(NewBorrowerBillRow(), f.GetLocalDB(), getBorrowerBillOpt(in.Addr))
	if err != nil {
		return nil, err
	}
	query := tab.GetQuery(f.GetLocalDB())

	//列出指定数量的行数
	var primaryKey []byte
	if len(in.PrimaryKey) != 0 {
		primaryKey = []byte(in.PrimaryKey)
	}
	rows, err := query.ListIndex("primary", nil, primaryKey, in.EntrieCount, 0)
	if err != nil {
		financelog.Error("Query_BorrowerBill", "QueryErr", err)
		return nil, err
	}

	for _, v := range rows {
		local := v.Data.(*financetypes.BorrowerBillLocalDB)
		iou, err := getBillFromDB(f.GetStateDB(), local.Id)
		if err != nil && err != types.ErrNotFound {
			financelog.Error("BorrowerBill", "getBillFromDB", err)
			return nil, err
		} else if err == types.ErrNotFound {
			return nil, pty.ErrBillNotFound
		}
		item := &financetypes.BorrowerBillInfo{Id: local.Id, Status: local.Status, CollectedAmount: iou.TotalCollectedAmount}
		item.LoanAmount = iou.LoanAmount
		item.CreateDate = iou.CreateTime
		item.WaitForGuaranteePeriod = iou.WaitForGuaranteePeriod
		reply.Items = append(reply.Items, item)
		reply.PrimaryKey = string(v.Primary)
	}
	financelog.Debug("Fetch Borrower Bills", "items", reply.Items)

	return &reply, nil
}

//Query_BorrowerBillByStatus query borrower's relate Bill by status
func (f *finance) Query_BorrowerBillByStatus(in *financetypes.ReqBorrowerBill) (types.Message, error) {
	var reply financetypes.ReplyBorrowerBill

	tab, err := table.NewTable(NewBorrowerBillRow(), f.GetLocalDB(), getBorrowerBillOpt(in.Addr))
	if err != nil {
		return nil, err
	}
	query := tab.GetQuery(f.GetLocalDB())

	//列出指定数量的行数
	var primaryKey []byte
	if len(in.PrimaryKey) != 0 {
		primaryKey = []byte(in.PrimaryKey)
	}
	rows, err := query.ListIndex("status", []byte(fmt.Sprintf("%d", in.Status)), primaryKey, in.EntrieCount, 0)
	if err != nil {
		financelog.Error("Query_BorrowerBillByStatus", "QueryErr", err)
		return nil, err
	}

	if in.Status != financetypes.IouFulfilling {
		for _, v := range rows {
			local := v.Data.(*financetypes.BorrowerBillLocalDB)
			iou, err := getBillFromDB(f.GetStateDB(), local.Id)
			if err != nil && err != types.ErrNotFound {
				financelog.Error("BorrowerBillByStatus", "getBillFromDB", err)
				return nil, err
			} else if err == types.ErrNotFound {
				return nil, pty.ErrBillNotFound
			}
			item := &financetypes.BorrowerBillInfo{Id: local.Id, Status: local.Status}
			item.LoanAmount = iou.LoanAmount
			item.CreateDate = iou.CreateTime
			item.WaitForGuaranteePeriod = iou.WaitForGuaranteePeriod
			reply.Items = append(reply.Items, item)
			reply.PrimaryKey = string(v.Primary)
		}
	} else {
		//正在履行的白条IouFulfilling，ReleaseIou发行白条之后能查看到相应的还款信息
		for _, v := range rows {
			local := v.Data.(*financetypes.BorrowerBillLocalDB)
			bill, err := getBillFromDB(f.GetStateDB(), local.Id)
			if err != nil && err != types.ErrNotFound {
				financelog.Error("BorrowerBillByStatus", "getBillFromDB", err)
				return nil, err
			} else if err == types.ErrNotFound {
				return nil, pty.ErrBillNotFound
			}
			// //没有筹集到款项，表示不需要履行还款职责
			// if iou.TotalCollectedAmount == 0 {
			// 	continue
			// }
			item := &financetypes.BorrowerBillInfo{Id: local.Id, Status: local.Status, CollectedAmount: bill.TotalCollectedAmount}
			item.LoanAmount = bill.LoanAmount
			item.CreateDate = bill.CreateTime
			item.WaitForGuaranteePeriod = bill.WaitForGuaranteePeriod
			method := &BillRepayMethod{}
			//amount=0,随着时间的缩进，还款金额在变化
			detail := method.CalcIouRepayInfo(bill, f.GetBlockTime(), bill.LoanAmount)
			item.Detail = detail

			reply.Items = append(reply.Items, item)
			reply.PrimaryKey = string(v.Primary)
		}
	}
	financelog.Debug("Fetch Borrower Bills by status", "items", reply.Items)

	return &reply, nil
}

/***********************************Query_Deliver查询交付白条相关的白条信息*******************************************/
//Query_DeliverListByDeliverID query deliver item by deliverID
func (f *finance) Query_DeliverListByDeliverID(in *financetypes.ReqDeliverList) (types.Message, error) {
	var reply financetypes.ReplyDeliverList
	tab, err := table.NewTable(NewDeliverBillListRow(), f.GetLocalDB(), optDeliverBillList)
	if err != nil {
		return nil, err
	}
	query := tab.GetQuery(f.GetLocalDB())

	rows, err := query.ListIndex("deliverID", []byte(in.DeliverID), nil, 1, 0)
	if err != nil {
		financelog.Error("Fetch deliver item by deliverID", "QueryErr", err)
		return nil, err
	}

	for _, v := range rows {
		local := v.Data.(*financetypes.DeliverItemLocalDB)
		reply.Items = append(reply.Items, local)
		reply.PrimaryKey = string(v.Primary)
	}
	financelog.Debug("Fetch deliver item from deliver list", "deliver pending order", reply.Items, "Last primary key", reply.PrimaryKey)

	return &reply, nil
}

//Query_DeliverListByDeliverAddr query deliver item by deliver address
func (f *finance) Query_DeliverListByDeliverAddr(in *financetypes.ReqDeliverList) (types.Message, error) {
	var reply financetypes.ReplyDeliverList
	tab, err := table.NewTable(NewDeliverBillListRow(), f.GetLocalDB(), optDeliverBillList)
	if err != nil {
		return nil, err
	}
	//列出指定数量的行数
	var primaryKey []byte
	if len(in.PrimaryKey) != 0 {
		primaryKey = []byte(in.PrimaryKey)
	}
	query := tab.GetQuery(f.GetLocalDB())

	rows, err := query.ListIndex("deliverAddr", []byte(in.DeliverAddr), primaryKey, in.EntrieCount, 0)
	if err != nil {
		financelog.Error("Fetch deliver items by deliverAddr", "QueryErr", err)
		return nil, err
	}

	for _, v := range rows {
		local := v.Data.(*financetypes.DeliverItemLocalDB)
		reply.Items = append(reply.Items, local)
		reply.PrimaryKey = string(v.Primary)
	}
	financelog.Debug("Fetch deliver items from deliver list", "deliver pending order", reply.Items, "Last primary key", reply.PrimaryKey)

	return &reply, nil
}

//Query_DeliverListByPeerAddr query deliver item by peer address
func (f *finance) Query_DeliverListByPeerAddr(in *financetypes.ReqDeliverList) (types.Message, error) {
	var reply financetypes.ReplyDeliverList
	tab, err := table.NewTable(NewDeliverBillListRow(), f.GetLocalDB(), optDeliverBillList)
	if err != nil {
		return nil, err
	}
	//列出指定数量的行数
	var primaryKey []byte
	if len(in.PrimaryKey) != 0 {
		primaryKey = []byte(in.PrimaryKey)
	}
	query := tab.GetQuery(f.GetLocalDB())

	rows, err := query.ListIndex("toAddr", []byte(in.ToAddr), primaryKey, in.EntrieCount, 0)
	if err != nil {
		financelog.Error("Fetch deliver items by peerAddr", "QueryErr", err)
		return nil, err
	}

	for _, v := range rows {
		local := v.Data.(*financetypes.DeliverItemLocalDB)
		reply.Items = append(reply.Items, local)
		reply.PrimaryKey = string(v.Primary)
	}
	financelog.Debug("Fetch deliver items from deliver list", "deliver pending order", reply.Items, "Last primary key", reply.PrimaryKey)

	return &reply, nil
}

/***********************************Query_SplitBill查询拆分白条相关的白条信息*******************************************/
//Query_SplitBillByTxHash query deliver item by txHash
func (f *finance) Query_SplitBillByTxHash(in *financetypes.ReqSplitBillList) (types.Message, error) {
	var reply financetypes.ReplySplitBillList
	tab, err := table.NewTable(NewSplitBillListRow(), f.GetLocalDB(), optSplitBillList)
	if err != nil {
		return nil, err
	}
	query := tab.GetQuery(f.GetLocalDB())

	rows, err := query.ListIndex("txHash", []byte(in.TxHash), nil, 1, 0)
	if err != nil {
		financelog.Error("Fetch SplitBill item by txHash", "QueryErr", err)
		return nil, err
	}

	for _, v := range rows {
		local := v.Data.(*financetypes.SplitBillRecordLocalDB)
		reply.Items = append(reply.Items, local)
		reply.PrimaryKey = string(v.Primary)
	}
	financelog.Debug("Fetch SplitBill item from Split list", "SplitBill pending order", reply.Items, "Last primary key", reply.PrimaryKey)

	return &reply, nil
}

//Query_SplitBillByTxHash query deliver item by txHash
func (f *finance) Query_SplitBillBySplitAddr(in *financetypes.ReqSplitBillList) (types.Message, error) {
	var reply financetypes.ReplySplitBillList
	tab, err := table.NewTable(NewSplitBillListRow(), f.GetLocalDB(), optSplitBillList)
	if err != nil {
		return nil, err
	}
	//列出指定数量的行数
	var primaryKey []byte
	if len(in.PrimaryKey) != 0 {
		primaryKey = []byte(in.PrimaryKey)
	}
	query := tab.GetQuery(f.GetLocalDB())

	rows, err := query.ListIndex("splitAddr", []byte(in.SplitAddr), primaryKey, in.EntrieCount, 0)
	if err != nil {
		financelog.Error("Fetch SplitBill item by splitAddr", "QueryErr", err)
		return nil, err
	}

	for _, v := range rows {
		local := v.Data.(*financetypes.SplitBillRecordLocalDB)
		reply.Items = append(reply.Items, local)
		reply.PrimaryKey = string(v.Primary)
	}
	financelog.Debug("Fetch SplitBill item from Split list", "SplitBill pending order", reply.Items, "Last primary key", reply.PrimaryKey)

	return &reply, nil
}

//Query_SplitBillByTxHash query deliver item by txHash
func (f *finance) Query_SplitBillByToAddr(in *financetypes.ReqSplitBillList) (types.Message, error) {
	var reply financetypes.ReplySplitBillList
	tab, err := table.NewTable(NewSplitBillListRow(), f.GetLocalDB(), optSplitBillList)
	if err != nil {
		return nil, err
	}
	//列出指定数量的行数
	var primaryKey []byte
	if len(in.PrimaryKey) != 0 {
		primaryKey = []byte(in.PrimaryKey)
	}
	query := tab.GetQuery(f.GetLocalDB())

	rows, err := query.ListIndex("toAddr", []byte(in.ToAddr), primaryKey, in.EntrieCount, 0)
	if err != nil {
		financelog.Error("Fetch SplitBill item by toAddr", "QueryErr", err)
		return nil, err
	}

	for _, v := range rows {
		local := v.Data.(*financetypes.SplitBillRecordLocalDB)
		reply.Items = append(reply.Items, local)
		reply.PrimaryKey = string(v.Primary)
	}
	financelog.Debug("Fetch SplitBill item from Split list", "SplitBill pending order", reply.Items, "Last primary key", reply.PrimaryKey)

	return &reply, nil
}

/***********************************Query_ApplyFinancing查询申请融资相关的信息*******************************************/
//Query_ApplyFinancingListByCreditor query ApplyFinancing list by Creditor
func (f *finance) Query_ApplyFinancingListByCreditor(in *financetypes.ReqApplyList) (types.Message, error) {
	var reply financetypes.ReplyApplyList
	tab, err := table.NewTable(NewApplyForFinancingListRow(), f.GetLocalDB(), optApplyForFinancingList)
	if err != nil {
		return nil, err
	}
	query := tab.GetQuery(f.GetLocalDB())

	//列出指定数量的行数
	var primaryKey []byte
	if len(in.PrimaryKey) != 0 {
		primaryKey = []byte(in.PrimaryKey)
	}

	rows, err := query.ListIndex("creditorAddr", []byte(in.CreditorAddr), primaryKey, in.EntrieCount, 0)
	if err != nil {
		financelog.Error("Query_ApplyFinancingListByCreditor", "QueryErr", err)
		return nil, err
	}

	for _, v := range rows {
		local := v.Data.(*financetypes.ApplyFinancingItemLocalDB)
		reply.Items = append(reply.Items, local)
		reply.PrimaryKey = string(v.Primary)
	}

	return &reply, nil
}

//Query_ApplyFinancingListByApplicant query ApplyFinancing list by Applicant
func (f *finance) Query_ApplyFinancingListByApplicant(in *financetypes.ReqApplyList) (types.Message, error) {
	var reply financetypes.ReplyApplyList
	tab, err := table.NewTable(NewApplyForFinancingListRow(), f.GetLocalDB(), optApplyForFinancingList)
	if err != nil {
		return nil, err
	}
	query := tab.GetQuery(f.GetLocalDB())

	//列出指定数量的行数
	var primaryKey []byte
	if len(in.PrimaryKey) != 0 {
		primaryKey = []byte(in.PrimaryKey)
	}

	rows, err := query.ListIndex("applyAddr", []byte(in.ApplyAddr), primaryKey, in.EntrieCount, 0)
	if err != nil {
		financelog.Error("Query_ApplyFinancingListByApplicant", "QueryErr", err)
		return nil, err
	}

	for _, v := range rows {
		local := v.Data.(*financetypes.ApplyFinancingItemLocalDB)
		reply.Items = append(reply.Items, local)
		reply.PrimaryKey = string(v.Primary)
	}

	return &reply, nil
}

//Query_ApplyFinancingListByBillID query ApplyFinancing list by iou-id
func (f *finance) Query_ApplyFinancingListByBillID(in *financetypes.ReqApplyList) (types.Message, error) {
	var reply financetypes.ReplyApplyList
	tab, err := table.NewTable(NewApplyForFinancingListRow(), f.GetLocalDB(), optApplyForFinancingList)
	if err != nil {
		return nil, err
	}
	query := tab.GetQuery(f.GetLocalDB())

	//列出指定数量的行数
	var primaryKey []byte
	if len(in.PrimaryKey) != 0 {
		primaryKey = []byte(in.PrimaryKey)
	}

	rows, err := query.ListIndex("id", []byte(in.BillID), primaryKey, in.EntrieCount, 0)
	if err != nil {
		financelog.Error("Query_ApplyFinancingListByBillID", "QueryErr", err)
		return nil, err
	}

	for _, v := range rows {
		local := v.Data.(*financetypes.ApplyFinancingItemLocalDB)
		reply.Items = append(reply.Items, local)
		reply.PrimaryKey = string(v.Primary)
	}

	return &reply, nil
}

//Query_ApplyFinancingListByDpdtSmybol query Financing list by dpdtSymbol
func (f *finance) Query_ApplyFinancingListByDpdtSmybol(in *financetypes.ReqApplyList) (types.Message, error) {
	var reply financetypes.ReplyApplyList
	tab, err := table.NewTable(NewApplyForFinancingListRow(), f.GetLocalDB(), optApplyForFinancingList)
	if err != nil {
		return nil, err
	}
	query := tab.GetQuery(f.GetLocalDB())

	//列出指定数量的行数
	var primaryKey []byte
	if len(in.PrimaryKey) != 0 {
		primaryKey = []byte(in.PrimaryKey)
	}

	rows, err := query.ListIndex("dpdtSymbol", []byte(in.DpdtSymbol), primaryKey, in.EntrieCount, 0)
	if err != nil {
		financelog.Error("Query_ApplyFinancingListByDpdtSmybol", "QueryErr", err)
		return nil, err
	}

	for _, v := range rows {
		local := v.Data.(*financetypes.ApplyFinancingItemLocalDB)
		reply.Items = append(reply.Items, local)
		reply.PrimaryKey = string(v.Primary)
	}

	return &reply, nil
}

//Query_ApplyFinancingListByTxHash query ApplyFinancing list by txHash
func (f *finance) Query_ApplyFinancingListByTxHash(in *financetypes.ReqApplyList) (types.Message, error) {
	var reply financetypes.ReplyApplyList
	tab, err := table.NewTable(NewApplyForFinancingListRow(), f.GetLocalDB(), optApplyForFinancingList)
	if err != nil {
		return nil, err
	}
	query := tab.GetQuery(f.GetLocalDB())

	rows, err := query.ListIndex("txHash", []byte(in.TxHash), nil, 1, 0)
	if err != nil {
		financelog.Error("Query_ApplyFinancingListByTxHash", "QueryErr", err)
		return nil, err
	}

	for _, v := range rows {
		local := v.Data.(*financetypes.ApplyFinancingItemLocalDB)
		reply.Items = append(reply.Items, local)
		reply.PrimaryKey = string(v.Primary)
	}

	return &reply, nil
}

/***********************************Query_CashList查询指定白条的兑现列表信息*******************************************/
//Query_CashList query cash list
func (f *finance) Query_CashList(in *financetypes.ReqCashList) (types.Message, error) {
	var reply financetypes.ReplyCashList

	cashList, err := getCashListFromDB(in.Id, f.GetLocalDB(), in.PrimaryKey, in.EntrieCount)
	if err != nil {
		return nil, err
	}

	reply.CashItems = append(reply.CashItems, cashList...)
	financelog.Debug("Fetch cash list", "Cash Items", cashList)

	return &reply, nil
}

//Query_CashListByTokenOwner query cash list
func (f *finance) Query_CashListByTokenOwner(in *financetypes.ReqCashList) (types.Message, error) {
	var reply financetypes.ReplyCashList
	var cashList []*pty.CashList

	tab, err := table.NewTable(NewCashListRow(), f.GetLocalDB(), getCashListOpt(in.Id))
	if err != nil {
		return nil, err
	}
	//列出指定数量的行数
	var primaryKey []byte
	if len(in.PrimaryKey) != 0 {
		primaryKey = []byte(in.PrimaryKey)
	}

	query := tab.GetQuery(f.GetLocalDB())
	rows, err := query.ListIndex("tokenOwner", []byte(in.TokenOwner), primaryKey, in.EntrieCount, 0)
	if err != nil {
		financelog.Error("Query_CashListByTokenOwner", "QueryErr", err)
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

	reply.CashItems = append(reply.CashItems, cashList...)
	financelog.Debug("Fetch cash list", "Cash Items", cashList)

	return &reply, nil
}

/***********************************Query_BillCashRecord查询指定白条的兑现记录信息*******************************************/
//Query_BillCashRecord query cash record for Bill
func (f *finance) Query_BillCashRecord(in *financetypes.ReqCashRecord) (types.Message, error) {
	var reply financetypes.ReplyCashRecord
	tab, err := table.NewTable(NewCashRecordLocalDBRow(), f.GetLocalDB(), optCashRecordLocalDB)
	if err != nil {
		return nil, err
	}
	query := tab.GetQuery(f.GetLocalDB())

	rows, err := query.ListIndex("billId", []byte(in.Id), nil, 0, 0)
	if err != nil {
		financelog.Error("Query_BillCashRecord", "QueryErr", err)
		return nil, err
	}

	for _, v := range rows {
		local := v.Data.(*financetypes.CashRecordLocalDB)
		reply.Items = append(reply.Items, local)
	}
	financelog.Debug("Fetch cash record", "items", reply.Items, "Bill ID", in.Id)

	return &reply, nil
}

/***********************************Query_BillRepayInfo查询bill的还款信息*******************************************/
//Query_BillRepayInfo query bill repay info
func (f *finance) Query_BillRepayInfo(in *financetypes.ReqBillRepayInfo) (types.Message, error) {
	var reply financetypes.ReplyBillRepayInfo

	cashList, err := getCashListFromDB(in.Id, f.GetLocalDB(), "", 0)
	if err != nil {
		return nil, err
	}

	bill, err := getBillFromDB(f.GetStateDB(), in.Id)
	if err != nil && err != types.ErrNotFound {
		financelog.Error("Query_BillRepayInfo", "getBillFromDB", err)
		return nil, err
	} else if err == types.ErrNotFound {
		return nil, pty.ErrBillNotFound
	}

	var totalNormalAmount int64
	var totalActualAmount int64
	for _, item := range cashList {
		cashTime := f.GetBlockTime()
		if item.Timestamp > bill.IssueDate+bill.CirculationTime {
			cashTime = bill.IssueDate + bill.CirculationTime + f.GetBlockTime() - item.Timestamp
		}
		m := &BillRepayMethod{}
		repayInfo := m.CalcIouRepayInfo(bill, cashTime, item.Amount)
		if repayInfo == nil {
			continue
		}
		if repayInfo.ActualRepayAmount == 0 {
			continue
		}
		totalNormalAmount += repayInfo.NormalRepayAmount
		totalActualAmount += repayInfo.ActualRepayAmount
		if repayInfo.IncludeOverdueRepay == true {
			reply.IncludeOverdueRepay = true
		}
		reply.OverdueRate = repayInfo.OverdueRate
		reply.RepayDueDate = repayInfo.RepayDueDate
	}
	reply.ActualRepayAmount = totalActualAmount
	reply.NormalRepayAmount = totalNormalAmount
	//在逾期宽限日之内不算逾期，在逾期宽限日之后算逾期
	if totalActualAmount > totalNormalAmount {
		reply.IncludeOverdueRepay = true
	} else {
		reply.IncludeOverdueRepay = false
	}

	return &reply, nil
}

//Query_RepayInfo query repay info for Bill
func (f *finance) Query_RepayInfo(in *financetypes.ReqBillRepayInfo) (types.Message, error) {
	var reply *financetypes.ReplyBillRepayInfo

	bill, err := getBillFromDB(f.GetStateDB(), in.Id)
	if err != nil && err != types.ErrNotFound {
		financelog.Error("ReqGuaranteeRepay", "getIouFromDB", err)
		return nil, err
	} else if err == types.ErrNotFound {
		return nil, pty.ErrBillNotFound
	}
	//method := &ACAPRepayMethod{}
	//reply = method.CalcIouRepayInfo(bill, f.GetBlockTime(), 0)
	method := &BillRepayMethod{}
	reply = method.CalcIouRepayInfo(bill, f.GetBlockTime(), bill.LoanAmount)
	financelog.Debug("Fetch Bill repay info", "repay info", reply, "Bill ID", in.Id)
	if reply == nil {
		return nil, pty.ErrNeedNotRepay
	}

	return reply, nil
}

//Query_BillRepayRecord query repay record for BIll
func (f *finance) Query_BillRepayRecord(in *financetypes.ReqRepayRecord) (types.Message, error) {
	var reply financetypes.ReplyRepayRecord
	tab, err := table.NewTable(NewRepayRecordLocalDBRow(), f.GetLocalDB(), optRepayRecordLocalDB)
	if err != nil {
		return nil, err
	}
	query := tab.GetQuery(f.GetLocalDB())

	rows, err := query.ListIndex("billId", []byte(in.Id), nil, 0, 0)
	if err != nil {
		financelog.Error("Query_BillRepayRecord", "QueryErr", err)
		return nil, err
	}

	for _, v := range rows {
		local := v.Data.(*financetypes.RepayRecordLocalDB)
		reply.Items = append(reply.Items, local)
	}
	financelog.Debug("Fetch Repay record", "items", reply.Items, "Bill ID", in.Id)

	return &reply, nil
}

/***********************************Query_BillBrokenList查询失信人信息*******************************************/
//Query_IouBrokenList query broken list
func (f *finance) Query_BillBrokenList(in *financetypes.ReqBrokenList) (types.Message, error) {
	var reply financetypes.ReplyBrokenList
	tab, err := table.NewTable(NewBrokenRecordStateDBRow(), f.GetLocalDB(), optBrokenRecord)
	if err != nil {
		return nil, err
	}
	query := tab.GetQuery(f.GetLocalDB())

	//列出指定数量的行数
	var primaryKey []byte
	if len(in.PrimaryKey) != 0 {
		primaryKey = []byte(in.PrimaryKey)
	}
	rows, err := query.ListIndex("primary", nil, primaryKey, in.EntrieCount, 0)
	if err != nil {
		financelog.Error("Query_BillBrokenList", "QueryErr", err)
		return nil, err
	}

	for _, v := range rows {
		local := v.Data.(*financetypes.BrokenRecordStateDB)
		reply.Items = append(reply.Items, local)
		reply.PrimaryKey = string(v.Primary)
	}
	financelog.Debug("Fetch list from broken list", "items", reply.Items, "Last primary key", reply.PrimaryKey)

	return &reply, nil
}

//Query_BillBrokenListByPhone query broken list by phone
func (f *finance) Query_BillBrokenListByPhone(in *financetypes.ReqBrokenList) (types.Message, error) {
	var reply financetypes.ReplyBrokenList
	tab, err := table.NewTable(NewBrokenRecordStateDBRow(), f.GetLocalDB(), optBrokenRecord)
	if err != nil {
		return nil, err
	}
	query := tab.GetQuery(f.GetLocalDB())

	//列出指定数量的行数
	var primaryKey []byte
	if len(in.PrimaryKey) != 0 {
		primaryKey = []byte(in.PrimaryKey)
	}
	rows, err := query.ListIndex("phone", []byte(in.Phone), primaryKey, in.EntrieCount, 0)
	if err != nil {
		financelog.Error("Query_BillBrokenListByPhone", "QueryErr", err)
		return nil, err
	}

	for _, v := range rows {
		local := v.Data.(*financetypes.BrokenRecordStateDB)
		reply.Items = append(reply.Items, local)
		reply.PrimaryKey = string(v.Primary)
	}
	financelog.Debug("Fetch list from broken list", "items", reply.Items, "Last primary key", reply.PrimaryKey)

	return &reply, nil
}

//Query_BillBrokenListByAddr query broken list by address
func (f *finance) Query_BillBrokenListByAddr(in *financetypes.ReqBrokenList) (types.Message, error) {
	var reply financetypes.ReplyBrokenList
	tab, err := table.NewTable(NewBrokenRecordStateDBRow(), f.GetLocalDB(), optBrokenRecord)
	if err != nil {
		return nil, err
	}
	query := tab.GetQuery(f.GetLocalDB())

	//列出指定数量的行数
	var primaryKey []byte
	if len(in.PrimaryKey) != 0 {
		primaryKey = []byte(in.PrimaryKey)
	}
	rows, err := query.ListIndex("borrowerAddr", []byte(in.Addr), primaryKey, in.EntrieCount, 0)
	if err != nil {
		financelog.Error("Query_BillBrokenListByAddr", "QueryErr", err)
		return nil, err
	}

	for _, v := range rows {
		local := v.Data.(*financetypes.BrokenRecordStateDB)
		reply.Items = append(reply.Items, local)
		reply.PrimaryKey = string(v.Primary)
	}
	financelog.Debug("Fetch list from broken list", "items", reply.Items, "Last primary key", reply.PrimaryKey)

	return &reply, nil
}

/***********************************Query_TokenValue查询token价值*******************************************/
//Query_TokenValue query token value
func (f *finance) Query_TokenValue(in *financetypes.ReqTokenValue) (types.Message, error) {
	var reply financetypes.ReplyTokenValue

	bill, err := getBillFromDB(f.GetStateDB(), in.Id)
	if err != nil && err != types.ErrNotFound {
		financelog.Error("ReqTokenValue", "getBillFromDB", err)
		return nil, err
	} else if err == types.ErrNotFound {
		return nil, pty.ErrBillNotFound
	}
	method := &BillRepayMethod{}
	tokenValue, err := method.GetTokenValue(in.Amount, bill, f.GetBlockTime())
	if err != nil {
		financelog.Error("Query_TokenValue error happen", "error", err)
		return nil, err
	}
	reply.Value = tokenValue
	financelog.Debug("Fetch token value info", "Token amount", in.Amount, "Token value", tokenValue)

	return &reply, nil
}
