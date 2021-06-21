package executor

import (
	drivers "github.com/33cn/chain33/system/dapp"
	"github.com/33cn/chain33/types"
	financetypes "github.com/33cn/plugin/plugin/dapp/finance/types"
)

/*
 * 实现交易的链上执行接口
 * 关键数据上链（statedb）并生成交易回执（log）
 */
func (f *finance) Exec_AddCreditToken(payload *financetypes.AddCreditToken, tx *types.Transaction, index int) (*types.Receipt, error) {
	financelog.Debug("AddCreditToken action")
	actiondb := NewFinanceAction(f, tx)
	return actiondb.addCreditToken(payload)
}

func (f *finance) Exec_DepositAsset(payload *financetypes.DepositAsset, tx *types.Transaction, index int) (*types.Receipt, error) {
	financelog.Debug("DepositAsset action")
	actiondb := NewFinanceAction(f, tx)
	return actiondb.depositAsset(payload)
}

func (f *finance) Exec_WithdrawlAsset(payload *financetypes.WithdrawlAsset, tx *types.Transaction, index int) (*types.Receipt, error) {
	financelog.Debug("WithdrawlAsset action")
	actiondb := NewFinanceAction(f, tx)
	return actiondb.withdrawlAsset(payload)
}

func (f *finance) Exec_TransferCoins(payload *financetypes.TransferCoins, tx *types.Transaction, index int) (*types.Receipt, error) {
	financelog.Debug("TransferCoins action")
	actiondb := NewFinanceAction(f, tx)
	return actiondb.transferCoins(payload)
}

func (f *finance) Exec_CreateBill(payload *financetypes.CreateBill, tx *types.Transaction, index int) (*types.Receipt, error) {
	financelog.Debug("CreateBill action")
	actiondb := NewFinanceAction(f, tx)
	return actiondb.createBill(payload)
}

func (f *finance) Exec_ReleaseBill(payload *financetypes.ReleaseBill, tx *types.Transaction, index int) (*types.Receipt, error) {
	financelog.Debug("ReleaseBill action")
	actiondb := NewFinanceAction(f, tx)
	return actiondb.releaseBill(payload)
}

func (f *finance) Exec_UnReleaseBill(payload *financetypes.UnReleaseBill, tx *types.Transaction, index int) (*types.Receipt, error) {
	financelog.Debug("UnReleaseBill action")
	actiondb := NewFinanceAction(f, tx)
	return actiondb.unReleaseBill(payload)
}

func (f *finance) Exec_DeliverBill(payload *financetypes.DeliverBill, tx *types.Transaction, index int) (*types.Receipt, error) {
	financelog.Debug("DeliverBill action")
	actiondb := NewFinanceAction(f, tx)
	return actiondb.deliverBill(payload)
}

func (f *finance) Exec_UnDeliverBill(payload *financetypes.UnDeliverBill, tx *types.Transaction, index int) (*types.Receipt, error) {
	financelog.Debug("UnDeliverBill action")
	actiondb := NewFinanceAction(f, tx)
	return actiondb.unDeliverBill(payload)
}

func (f *finance) Exec_ConfirmDeliverBill(payload *financetypes.ConfirmDeliverBill, tx *types.Transaction, index int) (*types.Receipt, error) {
	financelog.Debug("ConfirmDeliverBill action")
	actiondb := NewFinanceAction(f, tx)
	return actiondb.confirmDeliverBill(payload)
}

func (f *finance) Exec_SplitBill(payload *financetypes.SplitBill, tx *types.Transaction, index int) (*types.Receipt, error) {
	financelog.Debug("SplitBill action")
	actiondb := NewFinanceAction(f, tx)
	return actiondb.splitBill(payload)
}

func (f *finance) Exec_UnSplitBill(payload *financetypes.UnSplitBill, tx *types.Transaction, index int) (*types.Receipt, error) {
	financelog.Debug("UnSplitBill action")
	actiondb := NewFinanceAction(f, tx)
	return actiondb.unSplitBill(payload)
}

func (f *finance) Exec_ConfirmSplitBill(payload *financetypes.ConfirmSplitBill, tx *types.Transaction, index int) (*types.Receipt, error) {
	financelog.Debug("ConfirmSplitBill action")
	actiondb := NewFinanceAction(f, tx)
	return actiondb.confirmSplitBill(payload)
}

func (f *finance) Exec_ApplyForFinancing(payload *financetypes.ApplyForFinancing, tx *types.Transaction, index int) (*types.Receipt, error) {
	financelog.Debug("ApplyForFinancing action")
	actiondb := NewFinanceAction(f, tx)
	return actiondb.applyForFinancing(payload)
}

func (f *finance) Exec_UnApplyForFinancing(payload *financetypes.UnApplyForFinancing, tx *types.Transaction, index int) (*types.Receipt, error) {
	financelog.Debug("UnApplyForFinancing action")
	actiondb := NewFinanceAction(f, tx)
	return actiondb.unApplyForFinancing(payload)
}

func (f *finance) Exec_ConfirmFinancing(payload *financetypes.ConfirmFinancing, tx *types.Transaction, index int) (*types.Receipt, error) {
	financelog.Debug("ConfirmFinancing action")
	actiondb := NewFinanceAction(f, tx)
	return actiondb.confirmFinancing(payload)
}

func (f *finance) Exec_CashBill(payload *financetypes.CashBill, tx *types.Transaction, index int) (*types.Receipt, error) {
	financelog.Debug("CashBill action")
	actiondb := NewFinanceAction(f, tx)
	return actiondb.CashBill(payload)
}

func (f *finance) Exec_RepayBill(payload *financetypes.RepayBill, tx *types.Transaction, index int) (*types.Receipt, error) {
	financelog.Debug("RepayBill action")
	actiondb := NewFinanceAction(f, tx)
	return actiondb.RepayBill(payload)
}

func (f *finance) Exec_ReportBroken(payload *financetypes.ReportBroken, tx *types.Transaction, index int) (*types.Receipt, error) {
	financelog.Debug("ReportBroken action")
	actiondb := NewFinanceAction(f, tx)
	return actiondb.ReportBroken(payload)
}

//ExecutorOrder Exec 的时候 同时执行 ExecLocal
func (f *finance) ExecutorOrder() int64 {
	return drivers.ExecLocalSameTime
}
