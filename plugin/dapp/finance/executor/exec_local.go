package executor

import (
	"encoding/hex"
	"github.com/33cn/chain33/common/db/table"
	"github.com/33cn/chain33/types"
	financetypes "github.com/33cn/plugin/plugin/dapp/finance/types"
	pty "github.com/33cn/plugin/plugin/dapp/finance/types"
)

/*
 * 实现交易相关数据本地执行，数据不上链
 * 非关键数据，本地存储(localDB), 用于辅助查询，效率高
 */

func (f *finance) ExecLocal_AddCreditToken(payload *financetypes.AddCreditToken, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	dbSet := &types.LocalDBSet{}
	if receiptData.GetTy() != types.ExecOk {
		financelog.Error("ExecLocal_AddDpdtToken: exec not OK")
		return dbSet, nil
	}

	i := pty.NegtiveNum
	for k, log := range receiptData.Logs {
		if log.GetTy() != financetypes.TyAddCreditTokenLog {
			continue
		}
		i = k
		break
	}
	if i == pty.NegtiveNum {
		return nil, types.ErrLogType
	}

	log := receiptData.Logs[i].GetLog()
	var receipt pty.ReceiptLogAddCreditToken
	err := types.Decode(log, &receipt)
	if err != nil {
		financelog.Error("ExecLocal_AddCreditToken", "DecodeErr", err)
		return nil, types.ErrDecode
	}

	var kvPairs []*types.KeyValue
	for _, v := range receipt.TransfItems.Items {
		//添加用户没有的资产symbol,授信代币
		// 这个方法是将相应的代币标识symbol存到用户资产下面，为了查询账户资产时使用
		kv := AddTokenToAssets(v.ToAddr, f.GetLocalDB(), v.Symbol)
		if kv != nil {
			kvPairs = append(kvPairs, kv...)
		}
	}

	//在LocalDB中添加授信记录
	creditItem := &pty.CreditRecordLocalDB{
		CreditAddr:  payload.CreditAddr,
		GranteeAddr: payload.GranteeAddr,
		Symbol:      payload.Symbol,
		Amount:      payload.Amount,
		Rate:        payload.Rate,
		Expire:      receipt.Expire,
		Duration:    payload.Duration,
		Remark:      payload.Remark,
		BlockTime:   f.GetBlockTime(), //授信时间在查询信息的时候已经添加了这里可以省略
		TxHash:      "0x" + hex.EncodeToString(tx.Hash()),
		Height:      f.GetHeight(),
	}

	tab, err := table.NewTable(NewCreditRecordRow(), f.GetLocalDB(), optCreditRecordLocalDB)
	if err != nil {
		financelog.Error("ExecLocal_AddCreditToken", "AddCreditToken NewTable err ", err)
		return nil, err
	}
	//将授信记录添加进LocalDB中的表中
	tab.Add(creditItem)
	tabkv, err := tab.Save()
	if err != nil {
		return nil, err
	}

	kvPairs = append(kvPairs, tabkv...)
	dbSet = &types.LocalDBSet{KV: kvPairs}

	return f.addAutoRollBack(tx, dbSet.KV), nil
}

func (f *finance) ExecLocal_DepositAsset(payload *financetypes.DepositAsset, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	dbSet := &types.LocalDBSet{}
	if receiptData.GetTy() != types.ExecOk {
		return dbSet, nil
	}

	//将存入的symbol代币存到addr资产列表中
	kv := AddTokenToAssets(payload.DepositAddr, f.GetLocalDB(), payload.Symbol)
	if kv != nil {
		dbSet.KV = append(dbSet.KV, kv...)
	}

	//将入金操作记录到DWRecord_Addr表中
	record := pty.DWRecordLocalDB{
		Addr:      payload.DepositAddr,
		Action:    pty.DEPOSIT,
		Amount:    payload.Amount,
		Symbol:    payload.Symbol,
		Timestamp: f.GetBlockTime(),
		TxHash:    "0x" + hex.EncodeToString(tx.Hash()),
	}

	tab, err := table.NewTable(NewDWRecordRow(), f.GetLocalDB(), getDWRecordOpt(payload.DepositAddr))
	if err != nil {
		financelog.Error("ExecLocal_DepositAsset", "DepositAsset NewTable err ", err)
		return nil, err
	}
	tab.Add(&record)
	tabkv, err := tab.Save()
	if err != nil {
		return nil, err
	}
	dbSet.KV = append(dbSet.KV, tabkv...)

	return f.addAutoRollBack(tx, dbSet.KV), nil
}

func (f *finance) ExecLocal_WithdrawlAsset(payload *financetypes.WithdrawlAsset, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	dbSet := &types.LocalDBSet{}

	record := pty.DWRecordLocalDB{
		Addr:      payload.WithdrawAddr,
		Action:    pty.WITHDRAWAL,
		Amount:    payload.Amount,
		Symbol:    payload.Symbol,
		Timestamp: f.GetBlockTime(),
		TxHash:    "0x" + hex.EncodeToString(tx.Hash()),
	}

	//转出操作记录到DWRecord_Addr表中
	tab, err := table.NewTable(NewDWRecordRow(), f.GetLocalDB(), getDWRecordOpt(payload.WithdrawAddr))
	if err != nil {
		financelog.Error("ExecLocal_WithdrawlAsset", "WithdrawlAsset NewTable err ", err)
		return nil, err
	}
	tab.Add(&record)
	tabkv, err := tab.Save()
	if err != nil {
		return nil, err
	}
	dbSet.KV = append(dbSet.KV, tabkv...)
	return f.addAutoRollBack(tx, dbSet.KV), nil
}

func (f *finance) ExecLocal_TransferCoins(payload *financetypes.TransferCoins, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	dbSet := &types.LocalDBSet{}
	if receiptData.GetTy() != types.ExecOk {
		return dbSet, nil
	}
	i := pty.NegtiveNum
	for k, log := range receiptData.Logs {
		if log.GetTy() != financetypes.TyTransferCoinsLog {
			continue
		}
		i = k
		break
	}
	if i == pty.NegtiveNum {
		return nil, types.ErrLogType
	}
	log := receiptData.Logs[i].GetLog()
	var receipt financetypes.ReceiptLogTransferCoins
	err := types.Decode(log, &receipt)
	if err != nil {
		financelog.Error("ExecLocal_TransferCoins", "DecodeErr", err)
		return nil, types.ErrDecode
	}

	var kvPairs []*types.KeyValue
	for _, v := range receipt.TransfItems.Items {
		//将转到ToAddr地址中的symbol代币标识添加到ToAddr资产列表中
		kv := AddTokenToAssets(v.ToAddr, f.GetLocalDB(), v.Symbol)
		if kv != nil {
			kvPairs = append(kvPairs, kv...)
		}
	}

	//记录一下TransferCoinRecord
	transferCoin := pty.TransferCoinsRecordLocalDB{
		FromAddr:  payload.FromAddr,
		ToAddr:    payload.ToAddr,
		Exec:      payload.Exec,
		Amount:    payload.Amount,
		Symbol:    payload.Symbol,
		Timestamp: f.GetBlockTime(),
		TxHash:    "0x" + hex.EncodeToString(tx.Hash()),
		Remark:    payload.Remark,
	}
	//转出操作记录到TransferCoinsRecordLocalDB表中
	tab, err := table.NewTable(NewTransferCoinsRecordRow(), f.GetLocalDB(), getTransferCoinsRecordOpt(payload.FromAddr))
	if err != nil {
		financelog.Error("ExecLocal_TransferCoins", "TransferCoins NewTable err ", err)
		return nil, err
	}
	tab.Add(&transferCoin)
	tabkv, err := tab.Save()
	if err != nil {
		return nil, err
	}
	kvPairs = append(kvPairs, tabkv...)
	dbSet = &types.LocalDBSet{KV: kvPairs}
	return f.addAutoRollBack(tx, dbSet.KV), nil
}

func (f *finance) ExecLocal_CreateBill(payload *financetypes.CreateBill, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	dbSet := &types.LocalDBSet{}
	if receiptData.GetTy() != types.ExecOk {
		return dbSet, nil
	}

	var kvPairs []*types.KeyValue

	//白条状态为持有中
	borrowerIou := &financetypes.BorrowerBillLocalDB{Id: payload.Id, Status: financetypes.IouHolded}
	tab, err := table.NewTable(NewBorrowerBillRow(), f.GetLocalDB(), getBorrowerBillOpt(payload.Addr))
	if err != nil {
		financelog.Error("ExecLocal_CreateBill", "CreateBill NewTable err ", err)
		return nil, err
	}
	tab.Add(borrowerIou)
	tabkv, err := tab.Save()
	if err != nil {
		return nil, err
	}
	kvPairs = append(kvPairs, tabkv...)
	dbSet = &types.LocalDBSet{KV: kvPairs}
	return f.addAutoRollBack(tx, dbSet.KV), nil
}

func (f *finance) ExecLocal_ReleaseBill(payload *financetypes.ReleaseBill, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	dbSet := &types.LocalDBSet{}
	if receiptData.GetTy() != types.ExecOk {
		financelog.Error("ExecLocal_ReleaseBill: exec not OK")
		return dbSet, nil
	}

	lenderIou := &financetypes.BorrowerBillLocalDB{Id: payload.Id, Status: financetypes.IouFulfilling}
	tab, err := table.NewTable(NewBorrowerBillRow(), f.GetLocalDB(), getBorrowerBillOpt(payload.ReleaseAddr))
	if err != nil {
		financelog.Error("ExecLocal_ReleaseBill", "ReleaseBill NewTable err ", err)
		return nil, err
	}
	err = tab.Update([]byte(payload.Id), lenderIou)
	if err != nil {
		financelog.Error("ExecLocal_ReleaseBill: Update table failed")
		return nil, err
	}
	tabkv, err := tab.Save()
	if err != nil {
		financelog.Error("ExecLocal_ReleaseIou: save table failed")
		return nil, err
	}
	dbSet = &types.LocalDBSet{KV: tabkv}
	financelog.Debug("ExecLocal_ReleaseIou: save to table successfully", "borrowerIou", lenderIou)
	return f.addAutoRollBack(tx, dbSet.KV), nil
}

func (f *finance) ExecLocal_UnReleaseBill(payload *financetypes.UnReleaseBill, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	dbSet := &types.LocalDBSet{}
	if receiptData.GetTy() != types.ExecOk {
		return dbSet, nil
	}

	lenderIou := &financetypes.BorrowerBillLocalDB{Id: payload.Id, Status: financetypes.IouHolded}
	tab, err := table.NewTable(NewBorrowerBillRow(), f.GetLocalDB(), getBorrowerBillOpt(payload.UnReleaseAddr))
	if err != nil {
		financelog.Error("ExecLocal_UnReleaseIou", "UnReleaseBill NewTable err ", err)
		return nil, err
	}
	err = tab.Update([]byte(payload.Id), lenderIou)
	if err != nil {
		financelog.Error("ExecLocal_UnReleaseIou: Update table failed")
		return nil, err
	}
	tabkv, err := tab.Save()
	if err != nil {
		return nil, err
	}
	dbSet = &types.LocalDBSet{KV: tabkv}
	return f.addAutoRollBack(tx, dbSet.KV), nil
}

func (f *finance) ExecLocal_DeliverBill(payload *financetypes.DeliverBill, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	dbSet := &types.LocalDBSet{}
	if receiptData.GetTy() != types.ExecOk {
		financelog.Error("ExecLocal_DeliverBill: exec not OK")
		return dbSet, nil
	}

	//存入到LocalDB中的数据
	deliverItem := &pty.DeliverItemLocalDB{DeliverAddr: payload.DeliverAddr, ToAddr: payload.ToAddr, BillID: payload.BillID, Remark: payload.Remark, Amount: payload.Amount}
	deliverItem.Timestamp = f.GetBlockTime()
	//交付id就是txHash
	deliverItem.DeliverID = "0x" + hex.EncodeToString(tx.Hash())

	tab, err := table.NewTable(NewDeliverBillListRow(), f.GetLocalDB(), optDeliverBillList)
	if err != nil {
		financelog.Error("ExecLocal_DeliverBill", "DeliverBill NewTable err ", err)
		return nil, err
	}

	tab.Add(deliverItem)
	tabkv, err := tab.Save()
	if err != nil {
		return nil, err
	}
	var kvPairs []*types.KeyValue
	kvPairs = append(kvPairs, tabkv...)
	dbSet = &types.LocalDBSet{KV: kvPairs}

	return f.addAutoRollBack(tx, dbSet.KV), nil
}

func (f *finance) ExecLocal_UnDeliverBill(payload *financetypes.UnDeliverBill, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	dbSet := &types.LocalDBSet{}
	if receiptData.GetTy() != types.ExecOk {
		return dbSet, nil
	}

	tab, err := table.NewTable(NewDeliverBillListRow(), f.GetLocalDB(), optDeliverBillList)
	if err != nil {
		financelog.Error("ExecLocal_UnDeliverBill", "new deliverBill List init err ", err)
		return nil, err
	}
	//删除相应的交付信息在LocalDB中
	err = tab.Del([]byte(payload.DeliverID))
	if err != nil {
		financelog.Error("Delete row from deliverBill List failed")
		return nil, err
	}
	tabkv, err := tab.Save()
	if err != nil {
		return nil, err
	}

	dbSet = &types.LocalDBSet{KV: tabkv}

	return f.addAutoRollBack(tx, dbSet.KV), nil
}

func (f *finance) ExecLocal_ConfirmDeliverBill(payload *financetypes.ConfirmDeliverBill, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	dbSet := &types.LocalDBSet{}
	if receiptData.GetTy() != types.ExecOk {
		return dbSet, nil
	}

	i := pty.NegtiveNum
	for k, log := range receiptData.Logs {
		if log.GetTy() != financetypes.TyConfirmDeliverBillLog {
			continue
		}
		i = k
		break
	}
	if i == pty.NegtiveNum {
		financelog.Error("ExecLocal_ConfirmDeliverBill: log type error", "logType", i)
		return nil, types.ErrLogType
	}
	log := receiptData.Logs[i].GetLog()
	var receipt financetypes.ReceiptLogConfirmDeliver
	err := types.Decode(log, &receipt)
	if err != nil {
		financelog.Error("ExecLocal_ConfirmDeliverBill", "DecodeErr", err)
		return nil, types.ErrDecode
	}
	var kvPairs []*types.KeyValue
	for _, v := range receipt.TransfItems.Items {
		//将白条币symbol添加到供应商资产列表中
		kv := AddTokenToAssets(v.ToAddr, f.GetLocalDB(), v.Symbol)
		if kv != nil {
			kvPairs = append(kvPairs, kv...)
		}
	}

	tab, err := table.NewTable(NewDeliverBillListRow(), f.GetLocalDB(), optDeliverBillList)
	if err != nil {
		financelog.Error("ExecLocal_ConfirmDeliverBill", "new deliver List init err ", err)
		return nil, err
	}
	//白条资产交付成功，从交付列表中删除交付信息
	err = tab.Del([]byte(payload.DeliverID))
	if err != nil {
		financelog.Error("Delete row from Deliver List failed")
		return nil, err
	}
	tabkv, err := tab.Save()
	if err != nil {
		return nil, err
	}

	kvPairs = append(kvPairs, tabkv...)
	dbSet = &types.LocalDBSet{KV: kvPairs}
	return f.addAutoRollBack(tx, dbSet.KV), nil
}

func (f *finance) ExecLocal_SplitBill(payload *financetypes.SplitBill, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	dbSet := &types.LocalDBSet{}
	if receiptData.GetTy() != types.ExecOk {
		financelog.Error("ExecLocal_SplitBill: exec not OK")
		return dbSet, nil
	}

	//存入到LocalDB中的数据
	splitItem := &pty.SplitBillRecordLocalDB{
		SplitAddr: payload.SplitAddr,
		ToAddr:    payload.ToAddr,
		BillID:    payload.BillID,
		Remark:    payload.Remark,
		Amount:    payload.Amount}
	splitItem.BlockTime = f.GetBlockTime()
	splitItem.TxHash = "0x" + hex.EncodeToString(tx.Hash())

	tab, err := table.NewTable(NewSplitBillListRow(), f.GetLocalDB(), optSplitBillList)
	if err != nil {
		financelog.Error("ExecLocal_SplitBill", "SplitBill NewTable err ", err)
		return nil, err
	}

	tab.Add(splitItem)
	tabkv, err := tab.Save()
	if err != nil {
		return nil, err
	}
	var kvPairs []*types.KeyValue
	kvPairs = append(kvPairs, tabkv...)
	dbSet = &types.LocalDBSet{KV: kvPairs}

	return f.addAutoRollBack(tx, dbSet.KV), nil
}

func (f *finance) ExecLocal_UnSplitBill(payload *financetypes.UnSplitBill, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	dbSet := &types.LocalDBSet{}
	if receiptData.GetTy() != types.ExecOk {
		financelog.Error("ExecLocal_UnSplitBill: exec not OK")
		return dbSet, nil
	}

	tab, err := table.NewTable(NewSplitBillListRow(), f.GetLocalDB(), optSplitBillList)
	if err != nil {
		financelog.Error("ExecLocal_UnSplitBill", "SplitBill NewTable err ", err)
		return nil, err
	}

	tab.Del([]byte(payload.TxHash))
	if err != nil {
		financelog.Error("Delete row from SplitBill List failed")
		return nil, err
	}

	tabkv, err := tab.Save()
	if err != nil {
		return nil, err
	}
	var kvPairs []*types.KeyValue
	kvPairs = append(kvPairs, tabkv...)
	dbSet = &types.LocalDBSet{KV: kvPairs}

	return f.addAutoRollBack(tx, dbSet.KV), nil
}

func (f *finance) ExecLocal_ConfirmSplitBill(payload *financetypes.ConfirmSplitBill, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	dbSet := &types.LocalDBSet{}
	if receiptData.GetTy() != types.ExecOk {
		return dbSet, nil
	}

	i := pty.NegtiveNum
	for k, log := range receiptData.Logs {
		if log.GetTy() != financetypes.TyConfirmSplitBillLog {
			continue
		}
		i = k
		break
	}
	if i == pty.NegtiveNum {
		financelog.Error("ExecLocal_ConfirmSplitBill: log type error", "logType", i)
		return nil, types.ErrLogType
	}
	log := receiptData.Logs[i].GetLog()
	var receipt financetypes.ReceiptLogConfirmSplitBill
	err := types.Decode(log, &receipt)
	if err != nil {
		financelog.Error("ExecLocal_ConfirmSplitBill", "DecodeErr", err)
		return nil, types.ErrDecode
	}
	var kvPairs []*types.KeyValue
	for _, v := range receipt.TransfItems.Items {
		//将白条币symbol添加到二级供应商资产列表中
		kv := AddTokenToAssets(v.ToAddr, f.GetLocalDB(), v.Symbol)
		if kv != nil {
			kvPairs = append(kvPairs, kv...)
		}
	}

	tab, err := table.NewTable(NewSplitBillListRow(), f.GetLocalDB(), optSplitBillList)
	if err != nil {
		financelog.Error("ExecLocal_ConfirmSplitBill", "new SplitBill List init err ", err)
		return nil, err
	}
	//白条资产拆分成功，从拆分列表中删除交付信息
	err = tab.Del([]byte(payload.TxHash))
	if err != nil {
		financelog.Error("Delete row from Split List failed")
		return nil, err
	}
	tabkv, err := tab.Save()
	if err != nil {
		return nil, err
	}

	kvPairs = append(kvPairs, tabkv...)
	dbSet = &types.LocalDBSet{KV: kvPairs}

	return f.addAutoRollBack(tx, dbSet.KV), nil
}

func (f *finance) ExecLocal_ApplyForFinancing(payload *financetypes.ApplyForFinancing, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	dbSet := &types.LocalDBSet{}
	if receiptData.GetTy() != types.ExecOk {
		financelog.Error("ExecLocal_ApplyForFinancing: exec not OK")
		return dbSet, nil
	}

	applyFinancingItem := &pty.ApplyFinancingItemLocalDB{
		ApplyAddr:  payload.ApplyAddr,
		Id:         payload.Id,
		Amount:     payload.Amount,
		DpdtSymbol: payload.DpdtSymbol,
		Remark:     payload.Remark,
		Timestamp:  f.GetBlockTime(),
		TxHash:     "0x" + hex.EncodeToString(tx.Hash()),
	}

	//获取授信记录
	credit, err := getCreditTokenFromDB(f.GetStateDB(), payload.DpdtSymbol)
	if err != nil && err != types.ErrNotFound {
		financelog.Error("ExecLocal_ApplyForFinancing", "getCreditTokenFromDB", err)
		return nil, err
	} else if err == types.ErrNotFound {
		return nil, pty.ErrCreditNotFound
	}
	//申请融资的资金方等于授信人
	applyFinancingItem.CreditorAddr = credit.CreditAddr

	tab, err := table.NewTable(NewApplyForFinancingListRow(), f.GetLocalDB(), optApplyForFinancingList)
	if err != nil {
		financelog.Error("ExecLocal_ApplyForFinancing", "new ApplyForFinancing list init err ", err)
		return nil, err
	}
	tab.Add(applyFinancingItem)
	tabkv, err := tab.Save()
	if err != nil {
		return nil, err
	}
	var kvPairs []*types.KeyValue
	kvPairs = append(kvPairs, tabkv...)
	dbSet = &types.LocalDBSet{KV: kvPairs}
	return f.addAutoRollBack(tx, dbSet.KV), nil
}

func (f *finance) ExecLocal_UnApplyForFinancing(payload *financetypes.UnApplyForFinancing, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	dbSet := &types.LocalDBSet{}
	if receiptData.GetTy() != types.ExecOk {
		return dbSet, nil
	}

	tab, err := table.NewTable(NewApplyForFinancingListRow(), f.GetLocalDB(), optApplyForFinancingList)
	if err != nil {
		financelog.Error("ExecLocal_UnApplyForFinancing", "new ApplyForFinancing list init err ", err)
		return nil, err
	}
	err = tab.Del([]byte(payload.TxHash))
	if err != nil {
		financelog.Error("Delete row from ApplyForFinancing List failed")
		return nil, err
	}
	tabkv, err := tab.Save()
	if err != nil {
		return nil, err
	}

	dbSet = &types.LocalDBSet{KV: tabkv}
	return f.addAutoRollBack(tx, dbSet.KV), nil
}

func (f *finance) ExecLocal_ConfirmFinancing(payload *financetypes.ConfirmFinancing, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	dbSet := &types.LocalDBSet{}
	if receiptData.GetTy() != types.ExecOk {
		return dbSet, nil
	}

	i := pty.NegtiveNum
	for k, log := range receiptData.Logs {
		if log.GetTy() != financetypes.TyConfirmFinancingLog {
			continue
		}
		i = k
		break
	}
	if i == pty.NegtiveNum {
		return nil, types.ErrLogType
	}
	log := receiptData.Logs[i].GetLog()
	var receipt financetypes.ReceiptLogConfirmFinancing
	err := types.Decode(log, &receipt)
	if err != nil {
		financelog.Error("ExecLocal_ConfirmFinancing", "DecodeErr", err)
		return nil, types.ErrDecode
	}
	var kvPairs []*types.KeyValue
	for _, v := range receipt.TransfItems.Items {
		//将白条币资产symbol添加到资金方地址中
		kv := AddTokenToAssets(v.ToAddr, f.GetLocalDB(), v.Symbol)
		if kv != nil {
			kvPairs = append(kvPairs, kv...)
		}
	}

	tab, err := table.NewTable(NewApplyForFinancingListRow(), f.GetLocalDB(), optApplyForFinancingList)
	if err != nil {
		financelog.Error("ExecLocal_ConfirmFinancing", "new ApplyFinancing List init err ", err)
		return nil, err
	}
	err = tab.Del([]byte(payload.TxHash))
	if err != nil {
		financelog.Error("Delete row from ApplyFinancing List failed")
		return nil, err
	}
	tabkv, err := tab.Save()
	if err != nil {
		return nil, err
	}

	kvPairs = append(kvPairs, tabkv...)
	dbSet = &types.LocalDBSet{KV: kvPairs}
	return f.addAutoRollBack(tx, dbSet.KV), nil
}

func (f *finance) ExecLocal_CashBill(payload *financetypes.CashBill, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	dbSet := &types.LocalDBSet{}

	if receiptData.GetTy() != types.ExecOk {
		return dbSet, nil
	}

	i := pty.NegtiveNum
	for k, log := range receiptData.Logs {
		if log.GetTy() != financetypes.TyCashBillLog {
			continue
		}
		i = k
		break
	}
	if i == pty.NegtiveNum {
		return nil, types.ErrLogType
	}

	log := receiptData.Logs[i].GetLog()
	var receipt financetypes.ReceiptLogCashBill
	err := types.Decode(log, &receipt)
	if err != nil {
		financelog.Error("ExecLocal_CashBill", "DecodeErr", err)
		return nil, types.ErrDecode
	}

	cashRecordItem := receipt.Record
	//创建一张兑现记录的表
	tab, err := table.NewTable(NewCashRecordLocalDBRow(), f.GetLocalDB(), optCashRecordLocalDB)
	if err != nil {
		financelog.Error("ExecLocal_CashBill", "new cash record init err ", err)
		return nil, err
	}
	tab.Add(cashRecordItem)
	tabkv, err := tab.Save()
	if err != nil {
		return nil, err
	}
	var pairs []*types.KeyValue
	pairs = append(pairs, tabkv...)

	if cashRecordItem.BillType == pty.BILL {
		cashItem := &financetypes.CashList{
			BillId:     cashRecordItem.BillId,
			Amount:     cashRecordItem.Amount,
			TxHash:     cashRecordItem.TxHash,
			TokenOwner: cashRecordItem.CashAddr,
			Timestamp:  cashRecordItem.Timestamp,
		}
		//创建兑现记录list 根据billID创建兑现记录
		tab, err := table.NewTable(NewCashListRow(), f.GetLocalDB(), getCashListOpt(cashRecordItem.BillId))
		if err != nil {
			financelog.Error("ExecLocal_CashBill", "new bill market init err ", err)
			return nil, err
		}
		tab.Add(cashItem)
		tabkv, err := tab.Save()
		if err != nil {
			return nil, err
		}
		pairs = append(pairs, tabkv...)
	}

	dbSet = &types.LocalDBSet{KV: pairs}
	return f.addAutoRollBack(tx, dbSet.KV), nil
}

func (f *finance) ExecLocal_RepayBill(payload *financetypes.RepayBill, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	dbSet := &types.LocalDBSet{}
	if receiptData.GetTy() != types.ExecOk {
		return dbSet, nil
	}

	i := pty.NegtiveNum
	for k, log := range receiptData.Logs {
		if log.GetTy() != financetypes.TyRepayBillLog {
			continue
		}
		i = k
		break
	}
	if i == pty.NegtiveNum {
		return nil, types.ErrLogType
	}

	log := receiptData.Logs[i].GetLog()
	var receipt financetypes.ReceiptLogRepayBill
	err := types.Decode(log, &receipt)
	if err != nil {
		financelog.Error("ExecLocal_RepayBill", "DecodeErr", err)
		return nil, types.ErrDecode
	}

	var kvPairs []*types.KeyValue
	repayRecordItem := receipt.Record
	if repayRecordItem.BillType == pty.BILL {
		//Remove item from cash list
		for _, v := range repayRecordItem.CashHashes {
			cashTab, err := table.NewTable(NewCashListRow(), f.GetLocalDB(), getCashListOpt(payload.Id))
			if err != nil {
				financelog.Error("ExecLocal_RepayBill", "new Cash list init err ", err)
				return nil, err
			}
			//兑现完成之后从申请兑现记录中的删除申请兑现信息
			cashTab.Del([]byte(v))
			tabkv, err := cashTab.Save()
			if err != nil {
				return nil, err
			}
			kvPairs = append(kvPairs, tabkv...)
		}
	}

	//Add repay record
	tab, err := table.NewTable(NewRepayRecordLocalDBRow(), f.GetLocalDB(), optRepayRecordLocalDB)
	if err != nil {
		financelog.Error("ExecLocal_RepayBill", "new repay record init err ", err)
		return nil, err
	}
	tab.Add(repayRecordItem)
	tabkv, err := tab.Save()
	if err != nil {
		return nil, err
	}
	kvPairs = append(kvPairs, tabkv...)

	//如果白条已经还款完成
	if repayRecordItem.IsArrearsCleared == true {
		//Finished 已完成履行的白条(还款结束),修改白条状态
		lenderBill := &financetypes.BorrowerBillLocalDB{Id: payload.Id, Status: financetypes.IouFinished}
		tab, err := table.NewTable(NewBorrowerBillRow(), f.GetLocalDB(), getBorrowerBillOpt(repayRecordItem.BillCreator))
		if err != nil {
			financelog.Error("ExecLocal_RepayBill", "new iou market init err ", err)
			return nil, err
		}
		err = tab.Update([]byte(payload.Id), lenderBill)
		if err != nil {
			financelog.Error("ExecLocal_RepayBill: Update table failed")
			return nil, err
		}
		tabkv, err := tab.Save()
		if err != nil {
			return nil, err
		}
		kvPairs = append(kvPairs, tabkv...)

	}

	//Remove from broken list once repay successfully
	tab, err = table.NewTable(NewBrokenRecordStateDBRow(), f.GetLocalDB(), optBrokenRecord)
	if err != nil {
		financelog.Error("ExecLocal_ReportBroken", "new broken record init err ", err)
		return nil, err
	}
	tab.Del([]byte(payload.Id))
	tabkv, err = tab.Save()
	if err != nil {
		return nil, err
	}
	kvPairs = append(kvPairs, tabkv...)

	for _, v := range receipt.TransfItems.Items {
		//将所有转账操作的symbol资产全都加入到转入地址中
		kv := AddTokenToAssets(v.ToAddr, f.GetLocalDB(), v.Symbol)
		if kv != nil {
			kvPairs = append(kvPairs, kv...)
		}
	}

	dbSet = &types.LocalDBSet{KV: kvPairs}
	return f.addAutoRollBack(tx, dbSet.KV), nil
}

func (f *finance) ExecLocal_ReportBroken(payload *financetypes.ReportBroken, tx *types.Transaction, receiptData *types.ReceiptData, index int) (*types.LocalDBSet, error) {
	dbSet := &types.LocalDBSet{}
	if receiptData.GetTy() != types.ExecOk {
		return dbSet, nil
	}

	i := pty.NegtiveNum
	for k, log := range receiptData.Logs {
		if log.GetTy() != financetypes.TyReportBrokenLog {
			continue
		}
		i = k
		break
	}
	if i == pty.NegtiveNum {
		return nil, types.ErrLogType
	}

	log := receiptData.Logs[i].GetLog()
	var brokenRecordItem financetypes.BrokenRecordStateDB
	err := types.Decode(log, &brokenRecordItem)
	if err != nil {
		financelog.Error("ExecLocal_ReportBroken", "DecodeErr", err)
		return nil, types.ErrDecode
	}

	tab, err := table.NewTable(NewBrokenRecordStateDBRow(), f.GetLocalDB(), optBrokenRecord)
	if err != nil {
		financelog.Error("ExecLocal_ReportBroken", "new broken record init err ", err)
		return nil, err
	}
	tab.Add(&brokenRecordItem)
	tabkv, err := tab.Save()
	if err != nil {
		return nil, err
	}

	dbSet = &types.LocalDBSet{KV: tabkv}
	return f.addAutoRollBack(tx, dbSet.KV), nil
}

//设置自动回滚
func (f *finance) addAutoRollBack(tx *types.Transaction, kv []*types.KeyValue) *types.LocalDBSet {

	dbSet := &types.LocalDBSet{}
	dbSet.KV = f.AddRollbackKV(tx, tx.Execer, kv)
	return dbSet
}
