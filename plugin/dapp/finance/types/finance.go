package types

import (
	log "github.com/33cn/chain33/common/log/log15"
	"github.com/33cn/chain33/types"
	"reflect"
)

/*
 * 交易相关类型定义
 * 交易action通常有对应的log结构，用于交易回执日志记录
 * 每一种action和log需要用id数值和name名称加以区分
 */

// action类型id和name，这些常量可以自定义修改
const (
	TyUnknowAction = iota + 100
	TyAddCreditTokenAction
	TyDepositAssetAction
	TyWithdrawlAssetAction
	TyTransferCoinsAction
	TyCreateBillAction
	TyReleaseBillAction
	TyUnReleaseBillAction
	TyApplyForFinancingAction
	TyUnApplyForFinancingAction
	TyConfirmFinancingAction
	TyCashBillAction
	TyRepayBillAction
	TyReportBrokenAction
	TyDeliverBillAction
	TyUnDeliverBillAction
	TyConfirmDeliverBillAction
	TySplitBillAction
	TyUnSplitBillAction
	TyConfirmSplitBillAction

	NameAddCreditTokenAction      = "AddCreditToken"
	NameDepositAssetAction        = "DepositAsset"
	NameWithdrawlAssetAction      = "WithdrawlAsset"
	NameTransferCoinsAction       = "TransferCoins"
	NameCreateBillAction          = "CreateBill"
	NameReleaseBillAction         = "ReleaseBill"
	NameUnReleaseBillAction       = "UnReleaseBill"
	NameApplyForFinancingAction   = "ApplyForFinancing"
	NameUnApplyForFinancingAction = "UnApplyForFinancing"
	NameConfirmFinancingAction    = "ConfirmFinancing"
	NameCashBillAction            = "CashBill"
	NameRepayBillAction           = "RepayBill"
	NameReportBrokenAction        = "ReportBroken"
	NameDeliverBillAction         = "DeliverBill"
	NameUnDeliverBillAction       = "UnDeliverBill"
	NameConfirmDeliverBillAction  = "ConfirmDeliverBill"
	NameSplitBillAction           = "SplitBill"
	NameUnSplitBillAction         = "UnSplitBill"
	NameConfirmSplitBillAction    = "ConfirmSplitBill"
)

// log类型id值
const (
	TyUnknownLog = iota + 100
	TyAddCreditTokenLog
	TyDepositAssetLog
	TyWithdrawlAssetLog
	TyTransferCoinsLog
	TyCreateBillLog
	TyReleaseBillLog
	TyUnReleaseBillLog
	TyApplyForFinancingLog
	TyUnApplyForFinancingLog
	TyConfirmFinancingLog
	TyCashBillLog
	TyRepayBillLog
	TyReportBrokenLog
	TyDeliverBillLog
	TyUnDeliverBillLog
	TyConfirmDeliverBillLog
	TySplitBillLog
	TyUnSplitBillLog
	TyConfirmSplitBillLog
)

var (
	//FinanceX 执行器名称定义
	FinanceX = "finance"
	//定义actionMap
	actionMap = map[string]int32{
		NameAddCreditTokenAction:      TyAddCreditTokenAction,
		NameDepositAssetAction:        TyDepositAssetAction,
		NameWithdrawlAssetAction:      TyWithdrawlAssetAction,
		NameTransferCoinsAction:       TyTransferCoinsAction,
		NameCreateBillAction:          TyCreateBillAction,
		NameReleaseBillAction:         TyReleaseBillAction,
		NameUnReleaseBillAction:       TyUnReleaseBillAction,
		NameApplyForFinancingAction:   TyApplyForFinancingAction,
		NameUnApplyForFinancingAction: TyUnApplyForFinancingAction,
		NameConfirmFinancingAction:    TyConfirmFinancingAction,
		NameCashBillAction:            TyCashBillAction,
		NameRepayBillAction:           TyRepayBillAction,
		NameReportBrokenAction:        TyReportBrokenAction,
		NameDeliverBillAction:         TyDeliverBillAction,
		NameUnDeliverBillAction:       TyUnDeliverBillAction,
		NameConfirmDeliverBillAction:  TyConfirmDeliverBillAction,
		NameSplitBillAction:           TySplitBillAction,
		NameUnSplitBillAction:         TyUnSplitBillAction,
		NameConfirmSplitBillAction:    TyConfirmSplitBillAction,
	}
	//定义log的id和具体log类型及名称，填入具体自定义log类型
	logMap = map[int64]*types.LogInfo{
		TyRepayBillLog:          {Ty: reflect.TypeOf(ReceiptLogRepayBill{}), Name: "RepayBillLog"},
		TyCashBillLog:           {Ty: reflect.TypeOf(ReceiptLogCashBill{}), Name: "CashBillLog"},
		TyReportBrokenLog:       {Ty: reflect.TypeOf(BrokenRecordStateDB{}), Name: "ReportBrokenLog"},
		TyTransferCoinsLog:      {Ty: reflect.TypeOf(ReceiptLogTransferCoins{}), Name: "TransferCoinsLog"},
		TyConfirmFinancingLog:   {Ty: reflect.TypeOf(ReceiptLogConfirmFinancing{}), Name: "ConfirmFinancingLog"},
		TyDepositAssetLog:       {Ty: reflect.TypeOf(ReceiptLogDepositAsset{}), Name: "DepositAssetLog"},
		TyWithdrawlAssetLog:     {Ty: reflect.TypeOf(ReceiptLogWithdrawlAsset{}), Name: "WithdrawlAssetLog"},
		TyAddCreditTokenLog:     {Ty: reflect.TypeOf(ReceiptLogAddCreditToken{}), Name: "AddCreditTokenLog"},
		TyCreateBillLog:         {Ty: reflect.TypeOf(ReceiptLogCreateBill{}), Name: "CreateBillLog"},
		TyReleaseBillLog:        {Ty: reflect.TypeOf(ReceiptLogReleaseBill{}), Name: "ReleaseBillLog"},
		TyUnReleaseBillLog:      {Ty: reflect.TypeOf(ReceiptLogUnReleaseBill{}), Name: "UnReleaseBillLog"},
		TyConfirmDeliverBillLog: {Ty: reflect.TypeOf(ReceiptLogConfirmDeliver{}), Name: "ConfirmDeliverBillLog"},
		TyConfirmSplitBillLog:   {Ty: reflect.TypeOf(ReceiptLogConfirmSplitBill{}), Name: "ConfirmSplitBillLog"},
		//LogID:	{Ty: reflect.TypeOf(LogStruct), Name: LogName},
	}
	tlog = log.New("module", "finance.types")
)

// init defines a register function
func init() {
	types.AllowUserExec = append(types.AllowUserExec, []byte(FinanceX))
	//注册合约启用高度
	types.RegFork(FinanceX, InitFork)
	types.RegExec(FinanceX, InitExecutor)
}

// InitFork defines register fork
func InitFork(cfg *types.Chain33Config) {
	cfg.RegisterDappFork(FinanceX, "Enable", 0)
}

// InitExecutor defines register executor
func InitExecutor(cfg *types.Chain33Config) {
	types.RegistorExecutor(FinanceX, NewType(cfg))
}

type financeType struct {
	types.ExecTypeBase
}

func NewType(cfg *types.Chain33Config) *financeType {
	c := &financeType{}
	c.SetChild(c)
	c.SetConfig(cfg)
	return c
}

// GetPayload 获取合约action结构
func (f *financeType) GetPayload() types.Message {
	return &FinanceAction{}
}

// GeTypeMap 获取合约action的id和name信息
func (f *financeType) GetTypeMap() map[string]int32 {
	return actionMap
}

// GetLogMap 获取合约log相关信息
func (f *financeType) GetLogMap() map[int64]*types.LogInfo {
	return logMap
}
