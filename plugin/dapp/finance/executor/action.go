package executor

import (
	"encoding/hex"
	"errors"
	"github.com/33cn/chain33/account"
	"github.com/33cn/chain33/client"
	dbm "github.com/33cn/chain33/common/db"
	log "github.com/33cn/chain33/common/log/log15"
	"github.com/33cn/chain33/system/dapp"
	"github.com/33cn/chain33/types"
	pty "github.com/33cn/plugin/plugin/dapp/finance/types"
)

var financelog = log.New("module", "execs.financeOnChain")

// Action def
type Action struct {
	coinsAccount *account.DB
	db           dbm.KV
	txhash       []byte
	fromaddr     string
	blocktime    int64
	height       int64
	execaddr     string
	localdb      dbm.KVDB
	api          client.QueueProtocolAPI
}

// NewLoanAcction gen instance
func NewFinanceAction(r *finance, tx *types.Transaction) *Action {
	hash := tx.Hash()
	fromaddr := tx.From()
	return &Action{
		r.GetCoinsAccount(),
		r.GetStateDB(),
		hash,
		fromaddr,
		r.GetBlockTime(),
		r.GetHeight(),
		dapp.ExecAddress(string(tx.Execer)),
		r.GetLocalDB(),
		r.GetAPI(),
	}
}

//addCreditToken 向核心企业添加授信代币
func (action *Action) addCreditToken(c *pty.AddCreditToken) (*types.Receipt, error) {
	financelog.Debug("addCreditToken")
	if c.CreditAddr != action.fromaddr {
		return nil, pty.ErrOperateAddrMismatch
	}
	if c.CreditIndentity != pty.Funder {
		return nil, pty.ErrOperateAddrIsNotFunder
	}
	if c.GranteeIndentity != pty.Core {
		return nil, pty.ErrOperateAddrIsNotCore
	}
	if c == nil {
		return nil, types.ErrInvalidParam
	}

	//从StateDB中获取授信credit Token的记录,第一次查找授信记录的时候是查不到的
	credit, err := getCreditTokenFromStateDB(action.db, c.Symbol)
	if err != nil && err != types.ErrNotFound {
		financelog.Error("addCreditToken", "getCreditTokenFromStateDB", err)
		return nil, err
	} else if err == types.ErrNotFound {
		//一开始是找不到授信Token的因为还没有进行授信的操作，所以这里是types.ErrNotFound
		//此处需要注意，目前db接口，获取key未找到记录，返回空时候也带一个notFound错误，需要特殊处理，而不是直接返回错误
		financelog.Debug("addCreditToken", "credit token", c.Symbol)
	}

	//添加CreditToken存储到StateDB中
	credit = &pty.CreditTokenStateDB{
		CreditAddr:  c.CreditAddr,
		GranteeAddr: c.GranteeAddr,
		Symbol:      c.Symbol,
		Amount:      c.Amount,
		Rate:        c.Rate,
		Remark:      c.Remark,
	}
	credit.AddTime = action.blocktime
	credit.Duration = c.Duration

	//记录日志Log
	receiptLogAddCreditToken := &pty.ReceiptLogAddCreditToken{}
	receiptLogAddCreditToken.TransfItems = &pty.AssetTransferTags{Items: []*pty.AssetTransferTag{}}
	receiptLogAddCreditToken.Expire = action.blocktime + c.Duration

	//Transfer credit token to peerAddr
	//查询授信地址的授信代币余额
	//在执行这个操作之前，在finance外部已经创建好相应的授信Token,并且已经转入到finance合约中，现在可以查看到授信Token的余额
	tokenAmount, err := action.getBalance(pty.TokenExecer, c.Symbol, c.CreditAddr)
	if err != nil {
		financelog.Error("Get token amount failed")
		return nil, err
	}
	if tokenAmount < c.Amount {
		financelog.Error("Credit token amount not enough", "symbol", c.Symbol, "real amount", tokenAmount)
		return nil, pty.ErrAssetsNotEnough
	}
	//把授信代币转账到被授信人地址
	receiptTransferToken, err := action.transferAssets(pty.TokenExecer, c.CreditAddr, c.GranteeAddr, c.Amount, c.Symbol)
	if err != nil {
		financelog.Error("addCreditToken: transferAssets failed", "err", err)
		return nil, err
	}

	//Add transfer tag to log 将转账记录写到Log中
	transferItem := &pty.AssetTransferTag{FromAddr: c.CreditAddr, ToAddr: c.GranteeAddr, Symbol: c.Symbol, Amount: c.Amount, Remark: c.Remark}
	receiptLogAddCreditToken.TransfItems.Items = append(receiptLogAddCreditToken.TransfItems.Items, transferItem)

	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	//fill with kv and logs
	kv = append(kv, receiptTransferToken.KV...)
	logs = append(logs, receiptTransferToken.Logs...)

	//将log进行编码
	log := types.Encode(receiptLogAddCreditToken)
	logs = append(logs, &types.ReceiptLog{Ty: pty.TyAddCreditTokenLog, Log: log})

	//将授信Token写到StateDB中，创建key:mavl-finance-credit-token就能查到value：CreditToken的信息
	key := calcCreditNewKey(c.Symbol)
	value := types.Encode(credit)
	kv = append(kv, &types.KeyValue{Key: key, Value: value})
	receipt := &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}
	return receipt, nil

}

//depositAsset 向合约对应子账户存款
func (action *Action) depositAsset(c *pty.DepositAsset) (*types.Receipt, error) {
	financelog.Debug("depositAsset")
	if c == nil {
		return nil, types.ErrInvalidParam
	}
	if c.DepositAddr != action.fromaddr {
		return nil, pty.ErrOperateAddrMismatch
	}

	execer := pty.TokenExecer
	if c.Symbol == "token.CCNY" {
		execer = pty.CoinsExecer
	}
	//查询合约子账户在合约中的代币余额
	// (在进行depositAsset操作之前，代币已经从token/paracross合约中转进合约，这个操作只是记录一下转入金额和一些转入信息)
	balance, err := action.getBalance(execer, c.Symbol, c.DepositAddr)
	if err != nil {
		financelog.Error("Get asset balance failed", "Asset symbol", c.Symbol)
		return nil, err
	}

	if balance < c.Amount {
		financelog.Error("Deposit asset to loan execer failed", "asset symbol", c.Symbol, "balance", balance)
		return nil, pty.ErrDepositAsset
	}

	//Add deposit record log
	receiptLogDepositAsset := &pty.ReceiptLogDepositAsset{}
	receiptLogDepositAsset.TransfItems = &pty.AssetTransferTags{Items: []*pty.AssetTransferTag{}}
	transferItem := &pty.AssetTransferTag{ToAddr: c.DepositAddr, Symbol: c.Symbol, Amount: c.Amount, Remark: c.Remark}
	receiptLogDepositAsset.TransfItems.Items = append(receiptLogDepositAsset.TransfItems.Items, transferItem)

	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	log := types.Encode(receiptLogDepositAsset)
	logs = append(logs, &types.ReceiptLog{Ty: pty.TyDepositAssetLog, Log: log})
	receipt := &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}

	return receipt, nil
}

//withdrawlAsset 从合约中取款
func (action *Action) withdrawlAsset(c *pty.WithdrawlAsset) (*types.Receipt, error) {
	financelog.Debug("withdrawlAsset")
	if c == nil {
		return nil, types.ErrInvalidParam
	}
	if c.WithdrawAddr != action.fromaddr {
		return nil, pty.ErrOperateAddrMismatch
	}

	//Add withdrawl record logic
	receiptLogWithdrawlAsset := &pty.ReceiptLogWithdrawlAsset{}
	receiptLogWithdrawlAsset.TransfItems = &pty.AssetTransferTags{Items: []*pty.AssetTransferTag{}}
	transferItem := &pty.AssetTransferTag{FromAddr: c.WithdrawAddr, Symbol: c.Symbol, Amount: c.Amount, Remark: c.Remark}
	receiptLogWithdrawlAsset.TransfItems.Items = append(receiptLogWithdrawlAsset.TransfItems.Items, transferItem)

	//在进行withdrawlAsset操作之前就已经将相应的代币转出到合约了，所以不需要查询余额，此操作只是为了记录转出信息
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	log := types.Encode(receiptLogWithdrawlAsset)
	logs = append(logs, &types.ReceiptLog{Ty: pty.TyWithdrawlAssetLog, Log: log})
	receipt := &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}

	return receipt, nil
}

// TransferCoins Action
func (action *Action) transferCoins(c *pty.TransferCoins) (*types.Receipt, error) {
	financelog.Debug("TransferCoins")
	if c == nil {
		return nil, types.ErrInvalidParam
	}
	if c.FromAddr != action.fromaddr {
		return nil, pty.ErrOperateAddrMismatch
	}

	balance, err := action.getBalance(c.Exec, c.Symbol, c.FromAddr)
	if err != nil {
		financelog.Error("Get asset balance failed", "Asset symbol", c.Symbol)
		return nil, err
	}
	if balance < c.Amount {
		financelog.Error("Transfer Coins failed, balance not enough", "asset symbol", c.Symbol, "balance", balance)
		return nil, pty.ErrAssetsNotEnough
	}

	var logs []*types.ReceiptLog
	var kv []*types.KeyValue

	receiptTransfer, err := action.transferAssets(c.Exec, c.FromAddr, c.ToAddr, c.Amount, c.Symbol)
	if err != nil {
		return nil, err
	}
	kv = append(kv, receiptTransfer.KV...)
	logs = append(logs, receiptTransfer.Logs...)

	//Tag asset symbol
	receiptLogTransferCoins := &pty.ReceiptLogTransferCoins{}
	receiptLogTransferCoins.TransfItems = &pty.AssetTransferTags{Items: []*pty.AssetTransferTag{}}
	transferItem := &pty.AssetTransferTag{ToAddr: c.ToAddr, Symbol: c.Symbol, FromAddr: c.FromAddr, Amount: c.Amount, Remark: c.Remark}
	receiptLogTransferCoins.TransfItems.Items = append(receiptLogTransferCoins.TransfItems.Items, transferItem)

	log := types.Encode(receiptLogTransferCoins)
	logs = append(logs, &types.ReceiptLog{Ty: pty.TyTransferCoinsAction, Log: log})

	receipt := &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}
	return receipt, nil
}

//createBill create Bill
func (action *Action) createBill(c *pty.CreateBill) (*types.Receipt, error) {
	financelog.Debug("createBill")
	if c.Addr != action.fromaddr {
		return nil, pty.ErrOperateAddrMismatch
	}
	if c.CreatorIndentity != pty.Core {
		return nil, pty.ErrOperateAddrIsNotCore
	}
	if c == nil {
		return nil, types.ErrInvalidParam
	}

	bill, err := getBillFromDB(action.db, c.Id)
	if err != nil && err != types.ErrNotFound {
		financelog.Error("CreateBill", "getBillFromDB", err)
		return nil, err
	} else if err == types.ErrNotFound {
		financelog.Debug("CreateBill", "id", c.Id)
	}

	if c.BillType == pty.BILL {
		bill = &pty.Bill{
			Id:                     c.Id,
			TokenSymbol:            c.TokenSymbol,
			CoinsSymbol:            c.CoinsSymbol,
			Borrower:               c.Borrower,
			LoanAmount:             c.LoanAmount,
			Rate:                   c.Rate,
			OverdueRate:            c.OverdueRate,
			CirculationTime:        c.CirculationTime,
			Anonymous:              c.Anonymous,
			Phone:                  c.Phone,
			Identifier:             c.Identifier,
			Name:                   c.Name,
			Remark:                 c.Remark,
			BillType:               pty.BILL,
			WaitForGuaranteePeriod: c.WaitForGuaranteePeriod,
			Split:                  c.Split, //是否可以拆分白条
		}

		if c.OverdueLimit != 0 {
			bill.OverdueLimit = c.OverdueLimit
		} else {
			bill.OverdueLimit = pty.DefaultOverdueLimit
		}
		if c.OverdueGracePeriod != 0 {
			bill.OverdueGracePeriod = c.OverdueGracePeriod
		} else {
			bill.OverdueGracePeriod = pty.DefaultOverdueGracePeriod
		}
		if c.ExLimit != 0 {
			//If the value is zero, mean no limit
			bill.ExLimit = c.ExLimit
		}
		//For Bill, the repayCount will be only one
		bill.RepayCount = 1

		bill.OriginPrice = calcOriginPriceForBill(float64(c.Rate)/1e8, c.CirculationTime/pty.OneDaySeconds)
	}
	bill.Status = pty.StatusCreated
	bill.CreateTime = action.blocktime
	financelog.Debug("createBill", "blocktime", action.blocktime)

	bill.NeedDpdtToken = c.NeedDpdtToken
	if c.NeedDpdtToken {
		var total int64
		for _, v := range c.DpdtTokens {
			//检查授信的有效性,查看授信是否在授信列表里面
			credit, err := getCreditTokenFromDB(action.db, v.Symbol)
			if err != nil {
				financelog.Error("CreateBill", "getCreditTokenFromDB", err)
				continue
			}
			//授信代币有效性检查; 为了确保用户在白条兑付日前能够融到资金，授信到期日必须大于白条兑付日;
			if action.blocktime > (credit.AddTime+credit.Duration) || credit.GranteeAddr != c.Borrower {
				err := errors.New("credit token is invalid")
				financelog.Error("CreateBill", "getCreditTokenFromDB", err)
				continue
			}

			item := &pty.DpdtToken{Symbol: v.Symbol, Amount: v.Amount}
			bill.DpdtTokens = append(bill.DpdtTokens, item)
			total += v.Amount
		}
		if total < c.LoanAmount {
			return nil, errors.New("the whole amount of credit token is no enough")
		}
	}

	//Add CreateBill tag to log
	receiptLogCreateBill := &pty.ReceiptLogCreateBill{}
	receiptLogCreateBill.CreateBill = c

	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	key := calcBillNewKey(c.Id)
	value := types.Encode(bill)
	kv = append(kv, &types.KeyValue{Key: key, Value: value})

	//将log进行编码
	log := types.Encode(receiptLogCreateBill)
	//返回CreateBill的信息
	logs = append(logs, &types.ReceiptLog{Ty: pty.TyCreateBillLog, Log: log})
	receipt := &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}
	return receipt, nil
}

//releaseBill
func (action *Action) releaseBill(c *pty.ReleaseBill) (*types.Receipt, error) {
	financelog.Debug("releaseBill")
	if c.ReleaseAddr != action.fromaddr {
		return nil, pty.ErrOperateAddrMismatch
	}
	if c.ReleaseIndentity != pty.Core {
		return nil, pty.ErrOperateAddrIsNotCore
	}
	if c == nil {
		return nil, types.ErrInvalidParam
	}

	bill, err := getBillFromDB(action.db, c.Id)
	if err != nil && err != types.ErrNotFound {
		financelog.Error("releaseBill", "getBillFromDB", err)
		return nil, err
	} else if err == types.ErrNotFound {
		return nil, pty.ErrBillNotFound
	}

	//操作人必须就是借款人
	if c.ReleaseAddr != bill.Borrower {
		financelog.Error("Actual operator not the Bill creator", "creator", bill.Borrower)
		return nil, pty.ErrNoPermissionAction
	}

	//只有状态为created的白条，才能被发行;
	if bill.GetStatus() != pty.StatusCreated {
		financelog.Error("releaseBill error", "Bill status", bill.GetStatus())
		return nil, pty.ErrBillStatus
	}

	//查看借款人的白条币余额是否充足
	tokenAmount, err := action.getBalance(pty.TokenExecer, bill.TokenSymbol, bill.Borrower)
	if err != nil {
		financelog.Error("Get token amount failed")
		return nil, err
	}
	if tokenAmount < bill.LoanAmount {
		financelog.Error("Token amount not enough", "symbol", bill.TokenSymbol, "real amount", tokenAmount)
		return nil, pty.ErrAssetsNotEnough
	}

	//Update Bill
	bill.IssueDate = action.blocktime
	bill.Status = pty.StatusReleased
	bill.RepayDate = bill.IssueDate + bill.CirculationTime

	var logs []*types.ReceiptLog
	var kv []*types.KeyValue

	if bill.NeedDpdtToken {
		for _, v := range bill.DpdtTokens {
			//查看借款人的授信代币是否足够
			balance, err := action.getBalance(pty.TokenExecer, v.Symbol, c.ReleaseAddr)
			if err != nil {
				financelog.Error("Get credit token amount failed")
				return nil, err
			}
			if balance < v.Amount {
				financelog.Error("Credit token amount not enough", "balance", balance, "symbol", v.Symbol)
				return nil, pty.ErrAssetsNotEnough
			}
			//从当前账户中冻结授信授信token用于发行白条
			receiptFrozen, err := action.frozenAssets(pty.TokenExecer, v.Symbol, c.ReleaseAddr, v.Amount)
			if err != nil {
				return nil, err
			}
			kv = append(kv, receiptFrozen.KV...)
			logs = append(logs, receiptFrozen.Logs...)
		}
	}
	receiptLogReleaseBill := &pty.ReceiptLogReleaseBill{Id: bill.Id, State: bill.Status}
	key := calcBillNewKey(c.Id)
	value := types.Encode(bill)
	kv = append(kv, &types.KeyValue{Key: key, Value: value})
	//将log进行编码
	log := types.Encode(receiptLogReleaseBill)
	//返回ReleaseBill的信息
	logs = append(logs, &types.ReceiptLog{Ty: pty.TyReleaseBillLog, Log: log})
	receipt := &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}
	return receipt, nil
}

//unReleaseBill
func (action *Action) unReleaseBill(c *pty.UnReleaseBill) (*types.Receipt, error) {
	financelog.Debug("unReleaseBill")
	if c.UnReleaseAddr != action.fromaddr {
		return nil, pty.ErrOperateAddrMismatch
	}
	if c.UnReleaseIndentity != pty.Core {
		return nil, pty.ErrOperateAddrIsNotCore
	}
	if c == nil {
		return nil, types.ErrInvalidParam
	}

	bill, err := getBillFromDB(action.db, c.Id)
	if err != nil && err != types.ErrNotFound {
		financelog.Error("unReleaseBIll", "getBillFromDB", err)
		return nil, err
	} else if err == types.ErrNotFound {
		return nil, pty.ErrBillNotFound
	}

	//操作人必须就是借款人
	if c.UnReleaseAddr != bill.Borrower {
		financelog.Error("Actual operator not the Bill creator", "creator", bill.Borrower)
		return nil, pty.ErrNoPermissionAction
	}

	//只有状态为Released的白条，才能被撤发行;
	if bill.GetStatus() != pty.StatusReleased {
		financelog.Error("UnReleaseBill error", "Bill status", bill.GetStatus())
		return nil, pty.ErrBillStatus
	}

	//查询借款人的白条币余额，在白条币转给供应商之前才能进行撤发行白条
	balanceAmount, err := action.getBalance(pty.TokenExecer, bill.TokenSymbol, bill.Borrower)
	if err != nil {
		financelog.Error("Get token balance amount failed")
		return nil, err
	}
	//确保发行的白条没有被其他人持有; -todo, should be not equal?
	if balanceAmount < bill.LoanAmount {
		financelog.Error("Current Bill has been published to market, UnRelease Bill failed", "balance token amount", balanceAmount, "symbol", bill.TokenSymbol)
		return nil, pty.ErrUnReleaseIou
	}

	//update Bill
	bill.Status = pty.StatusCreated

	var logs []*types.ReceiptLog
	var kv []*types.KeyValue

	if bill.NeedDpdtToken {
		for _, v := range bill.DpdtTokens {
			//解冻依赖代币(授信Token)的资产
			frozen, err := action.getFrozenAsset(pty.TokenExecer, v.Symbol, c.UnReleaseAddr)
			if err != nil {
				financelog.Error("Get credit token amount failed")
				return nil, err
			}
			if frozen < v.Amount {
				financelog.Error("Credit token amount not enough", "frozen", frozen, "symbol", v.Symbol)
				return nil, pty.ErrAssetsNotEnough
			}
			//从当前账户中解冻授信token
			receiptFrozen, err := action.activeFrozenAssets(pty.TokenExecer, c.UnReleaseAddr, v.Amount, v.Symbol)
			if err != nil {
				return nil, err
			}
			kv = append(kv, receiptFrozen.KV...)
			logs = append(logs, receiptFrozen.Logs...)
		}
	}

	key := calcBillNewKey(c.Id)
	//修改bill的内容
	value := types.Encode(bill)
	kv = append(kv, &types.KeyValue{Key: key, Value: value})

	//记录UnreleaseBill
	receiptLogUnReleaseBill := &pty.ReceiptLogUnReleaseBill{Id: bill.Id, State: bill.Status}
	//将log进行编码
	log := types.Encode(receiptLogUnReleaseBill)
	logs = append(logs, &types.ReceiptLog{Ty: pty.TyUnReleaseBillLog, Log: log})
	receipt := &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}
	return receipt, nil
}

//deliverBill action 核心企业将白条交付给供应商
func (action *Action) deliverBill(c *pty.DeliverBill) (*types.Receipt, error) {
	financelog.Debug("deliverBill")
	if c == nil {
		return nil, types.ErrInvalidParam
	}
	if c.DeliverAddr != action.fromaddr {
		return nil, pty.ErrOperateAddrMismatch
	}
	//兑付白条的人必须是核心企业
	if c.DeliverIndentity != pty.Core {
		return nil, pty.ErrOperateAddrIsNotCore
	}
	//转到地址必须是供应商
	if c.ToIndentity != pty.Supplier {
		return nil, pty.ErrOperateAddrIsNotSupplier
	}

	bill, err := getBillFromDB(action.db, c.BillID)
	if err != nil {
		financelog.Error("deliverBill", "getBillFromDB", err)
		return nil, err
	}

	//判断白条是否过期，过期则不能交付
	cashDate := bill.IssueDate + bill.CirculationTime
	if cashDate <= action.blocktime {
		financelog.Error("deliverBill timeout")
		return nil, pty.ErrDeliverTimeout
	}
	//只有StatusReleased的白条才允许做deliverBill操作
	if bill.Status != pty.StatusReleased {
		financelog.Error("deliverBill: wrong status", "status", bill.Status)
		return nil, pty.ErrBillStatus
	}

	//Frozen Iou token
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	//从当前账户中冻结token,为后续交付白条做准备,冻结白条币
	receiptFrozen, err := action.frozenAssets(pty.TokenExecer, bill.TokenSymbol, c.DeliverAddr, c.Amount)
	if err != nil {
		return nil, err
	}
	kv = append(kv, receiptFrozen.KV...)
	logs = append(logs, receiptFrozen.Logs...)

	receipt := &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}
	return receipt, nil
}

//unDeliverBill action
func (action *Action) unDeliverBill(c *pty.UnDeliverBill) (*types.Receipt, error) {
	financelog.Debug("unDeliverBill")
	if c == nil {
		return nil, types.ErrInvalidParam
	}
	if c.UnDeliverAddr != action.fromaddr {
		return nil, pty.ErrOperateAddrMismatch
	}
	//撤销兑付白条的设情人必须是核心企业
	if c.UnDeliverIndentity != pty.Core {
		return nil, pty.ErrOperateAddrIsNotCore
	}
	//根据DeliverID(txHash)查询交付信息
	local, err := action.getDeliverItem(c.DeliverID)
	if err != nil {
		financelog.Error("unDeliverIou: getDeliverItem failed", "err", err)
		return nil, err
	}
	iou, err := getBillFromDB(action.db, local.BillID)
	if err != nil {
		financelog.Error("unDeliverBill: getBillFromDB failed", "err", err)
		return nil, err
	}

	//操作人必须是白条持有人
	if local.DeliverAddr != c.UnDeliverAddr {
		financelog.Error("unDeliverBill failed", "deliverAddr", local.DeliverAddr, "action address", c.UnDeliverAddr)
		return nil, pty.ErrOperateAddrMismatch
	}

	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	//Unfrozen token
	//Active relate token assets
	receiptActive, err := action.activeFrozenAssets(pty.TokenExecer, c.UnDeliverAddr, local.Amount, iou.TokenSymbol)
	if err != nil {
		financelog.Error("Active frozen assets failed")
		return nil, err
	}
	kv = append(kv, receiptActive.KV...)
	logs = append(logs, receiptActive.Logs...)

	receipt := &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}
	return receipt, nil
}

//confirmDeliverBill action 供应商时确认交付的地址
func (action *Action) confirmDeliverBill(c *pty.ConfirmDeliverBill) (*types.Receipt, error) {
	financelog.Debug("confirmDeliverBill")
	if c == nil {
		return nil, types.ErrInvalidParam
	}
	if c.ConfirmAddr != action.fromaddr {
		return nil, pty.ErrOperateAddrMismatch
	}
	//确认人的身份必须是供应商
	if c.ConfirmIndentity != pty.Supplier {
		return nil, pty.ErrOperateAddrIsNotSupplier
	}

	local, err := action.getDeliverItem(c.DeliverID)
	if err != nil {
		financelog.Error("confirmDeliver: getDeliverItem failed", "err", err)
		return nil, err
	}
	bill, err := getBillFromDB(action.db, local.BillID)
	if err != nil {
		financelog.Error("confirmDeliver: getBillFromDB failed", "err", err)
		return nil, err
	}

	//被交付对象必须是toAddr所指定的对象，供应商决定是否想要白条币
	if local.ToAddr != c.ConfirmAddr {
		financelog.Error("confirmDeliver: toAddr mismatch", "ToAddr", local.ToAddr, "operate addr", c.ConfirmAddr)
		return nil, pty.ErrOperateAddrMismatch
	}

	if action.blocktime > (bill.IssueDate + bill.CirculationTime) {
		//白条已逾期, 白条逾期之后则不允许融资
		return nil, pty.ErrBillHasExpired
	} else if bill.Status != pty.StatusReleased {
		//白条无效, 不能为其融资
		return nil, pty.ErrBillStatus
	}

	var logs []*types.ReceiptLog
	var kv []*types.KeyValue

	receiptConfirmD := &pty.ReceiptLogConfirmDeliver{}
	receiptConfirmD.TransfItems = &pty.AssetTransferTags{Items: []*pty.AssetTransferTag{}}
	if c.Ack {
		//同意白条交付操作，核心企业将冻结的白条币转账到供应商地址
		receiptTransferTokens, err := action.transferFrozenAssets(pty.TokenExecer, local.DeliverAddr, local.ToAddr, local.Amount, bill.TokenSymbol)
		if err != nil {
			financelog.Error("confirmDeliver: transferFrozenAssets failed", "err", err)
			return nil, err
		}
		//Add transfer tag to log
		transferItem := &pty.AssetTransferTag{ToAddr: c.ConfirmAddr, Symbol: bill.TokenSymbol, FromAddr: local.DeliverAddr, Amount: local.Amount, Remark: c.DeliverID}
		receiptConfirmD.TransfItems.Items = append(receiptConfirmD.TransfItems.Items, transferItem)

		//fill with kv and logs
		kv = append(kv, receiptTransferTokens.KV...)
		logs = append(logs, receiptTransferTokens.Logs...)
	} else {
		//不同意白条交付操作
		//解冻冻结资产
		receipt, err := action.activeFrozenAssets(pty.TokenExecer, local.DeliverAddr, local.Amount, bill.TokenSymbol)
		if err != nil {
			//解冻失败
			financelog.Error("confirmDeliver: Active frozen assets failed", "deliver address", local.DeliverAddr, "deliver token symbol", bill.TokenSymbol, "deliver amount", local.Amount)
			return nil, err
		}
		kv = append(kv, receipt.KV...)
		logs = append(logs, receipt.Logs...)
	}
	log := types.Encode(receiptConfirmD)
	logs = append(logs, &types.ReceiptLog{Ty: pty.TyConfirmDeliverBillLog, Log: log})

	receipt := &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}
	return receipt, nil
}

//splitBill action 一级供应商将白条币拆分给下一级供应商
func (action *Action) splitBill(c *pty.SplitBill) (*types.Receipt, error) {
	financelog.Debug("splitBill")
	if c == nil {
		return nil, types.ErrInvalidParam
	}
	if c.SplitAddr != action.fromaddr {
		return nil, pty.ErrOperateAddrMismatch
	}
	//拆分白条的人必须是供应商
	if c.SplitIndentity != pty.Supplier {
		return nil, pty.ErrOperateAddrIsNotSupplier
	}
	//接收拆分白条的人必须是供应商
	if c.ToIndentity != pty.Supplier {
		return nil, pty.ErrOperateAddrIsNotSupplier
	}

	bill, err := getBillFromDB(action.db, c.BillID)
	if err != nil {
		financelog.Error("splitBill", "getBillFromDB", err)
		return nil, err
	}

	//判断白条币是否能拆分
	if !bill.Split {
		//白条币不能拆分
		return nil, pty.ErrBillCanNotSplit
	}

	//判断白条是否过期，过期则不能拆分
	cashDate := bill.IssueDate + bill.CirculationTime
	if cashDate <= action.blocktime {
		financelog.Error("splitBill timeout")
		return nil, pty.ErrSplitTimeout
	}
	//只有StatusReleased的白条才允许做splitBill操作
	if bill.Status != pty.StatusReleased {
		financelog.Error("splitBill: wrong status", "status", bill.Status)
		return nil, pty.ErrBillStatus
	}

	//Frozen Iou token
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	//从当前(供应商)账户中冻结相应的白条token,为后续拆分白条做准备,冻结白条币
	//冻结一级供应商相应的白条币
	receiptFrozen, err := action.frozenAssets(pty.TokenExecer, bill.TokenSymbol, c.SplitAddr, c.Amount)
	if err != nil {
		return nil, err
	}
	kv = append(kv, receiptFrozen.KV...)
	logs = append(logs, receiptFrozen.Logs...)

	receipt := &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}
	return receipt, nil
}

//unSplitBill action 撤销拆分白条
func (action *Action) unSplitBill(c *pty.UnSplitBill) (*types.Receipt, error) {
	financelog.Debug("unSplitBill")
	if c == nil {
		return nil, types.ErrInvalidParam
	}
	if c.UnSplitAddr != action.fromaddr {
		return nil, pty.ErrOperateAddrMismatch
	}
	//拆分白条的人必须是供应商
	if c.UnSplitIndentity != pty.Supplier {
		return nil, pty.ErrOperateAddrIsNotSupplier
	}

	//根据(txHash)查询拆分信息
	local, err := action.getSplitBillItem(c.TxHash)
	if err != nil {
		financelog.Error("unSplitBill: getSplitBillItem failed", "err", err)
		return nil, err
	}
	bill, err := getBillFromDB(action.db, local.BillID)
	if err != nil {
		financelog.Error("unSplitBill: getBillFromDB failed", "err", err)
		return nil, err
	}

	//操作人必须是白条持有人
	if local.SplitAddr != c.UnSplitAddr {
		financelog.Error("unSplitBill failed", "SplitAddr", local.SplitAddr, "action address", c.UnSplitAddr)
		return nil, pty.ErrOperateAddrMismatch
	}

	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	//Unfrozen token
	//Active relate token assets,接触冻结供应商的白条币资产
	receiptActive, err := action.activeFrozenAssets(pty.TokenExecer, c.UnSplitAddr, local.Amount, bill.TokenSymbol)
	if err != nil {
		financelog.Error("Active frozen assets failed")
		return nil, err
	}
	kv = append(kv, receiptActive.KV...)
	logs = append(logs, receiptActive.Logs...)

	receipt := &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}
	return receipt, nil
}

//confirmSplitBill action 供应商时确认接收白条的地址
func (action *Action) confirmSplitBill(c *pty.ConfirmSplitBill) (*types.Receipt, error) {
	financelog.Debug("confirmSplitBill")
	if c == nil {
		return nil, types.ErrInvalidParam
	}
	if c.ConfirmAddr != action.fromaddr {
		return nil, pty.ErrOperateAddrMismatch
	}
	//确认人的身份必须是供应商
	if c.ConfirmIndentity != pty.Supplier {
		return nil, pty.ErrOperateAddrIsNotSupplier
	}

	local, err := action.getSplitBillItem(c.TxHash)
	if err != nil {
		financelog.Error("confirmSplitBill: getSplitBillItem failed", "err", err)
		return nil, err
	}
	bill, err := getBillFromDB(action.db, local.BillID)
	if err != nil {
		financelog.Error("confirmSplitBill: getSplitBillItem failed", "err", err)
		return nil, err
	}

	//接收拆分白条的对象必须是toAddr所指定的对象，二级供应商决定是否想要白条币
	if local.ToAddr != c.ConfirmAddr {
		financelog.Error("confirmSplitBill: toAddr mismatch", "ToAddr", local.ToAddr, "operate addr", c.ConfirmAddr)
		return nil, pty.ErrOperateAddrMismatch
	}

	if action.blocktime > (bill.IssueDate + bill.CirculationTime) {
		//白条已逾期, 白条逾期之后则不允许融资
		return nil, pty.ErrBillHasExpired
	} else if bill.Status != pty.StatusReleased {
		//白条无效, 不能为其融资
		return nil, pty.ErrBillStatus
	}

	var logs []*types.ReceiptLog
	var kv []*types.KeyValue

	receiptConfirmSplitBill := &pty.ReceiptLogConfirmSplitBill{}
	receiptConfirmSplitBill.TransfItems = &pty.AssetTransferTags{Items: []*pty.AssetTransferTag{}}
	if c.Ack {
		//同意白条拆分操作，供应商将冻结的白条币转账到二级供应商地址
		receiptTransferTokens, err := action.transferFrozenAssets(pty.TokenExecer, local.SplitAddr, local.ToAddr, local.Amount, bill.TokenSymbol)
		if err != nil {
			financelog.Error("confirmSplitBill: transferFrozenAssets failed", "err", err)
			return nil, err
		}
		//Add transfer tag to log
		transferItem := &pty.AssetTransferTag{ToAddr: c.ConfirmAddr, Symbol: bill.TokenSymbol, FromAddr: local.SplitAddr, Amount: local.Amount, Remark: c.TxHash}
		receiptConfirmSplitBill.TransfItems.Items = append(receiptConfirmSplitBill.TransfItems.Items, transferItem)

		//fill with kv and logs
		kv = append(kv, receiptTransferTokens.KV...)
		logs = append(logs, receiptTransferTokens.Logs...)
	} else {
		//不同意白条拆分操作
		//解冻冻结资产
		receipt, err := action.activeFrozenAssets(pty.TokenExecer, local.SplitAddr, local.Amount, bill.TokenSymbol)
		if err != nil {
			//解冻失败
			financelog.Error("confirmSplitBill: Active frozen assets failed", "Split address", local.SplitAddr, "Split token symbol", bill.TokenSymbol, "Split amount", local.Amount)
			return nil, err
		}
		kv = append(kv, receipt.KV...)
		logs = append(logs, receipt.Logs...)
	}
	log := types.Encode(receiptConfirmSplitBill)
	logs = append(logs, &types.ReceiptLog{Ty: pty.TyConfirmSplitBillLog, Log: log})

	receipt := &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}
	return receipt, nil
}

//applyForFinancing Action 多个供应商可以拿着白条id进行融資，在这里融資只是冻结相应的借款人的授信Token
func (action *Action) applyForFinancing(c *pty.ApplyForFinancing) (*types.Receipt, error) {
	financelog.Debug("applyForFinancing")
	if c == nil {
		return nil, types.ErrInvalidParam
	}
	if c.ApplyAddr != action.fromaddr {
		return nil, pty.ErrOperateAddrMismatch
	}
	//申请融资的人必须是供应商
	if c.ApplyIndentity != pty.Supplier {
		return nil, pty.ErrOperateAddrIsNotSupplier
	}

	bill, err := getBillFromDB(action.db, c.Id)
	if err != nil && err != types.ErrNotFound {
		financelog.Error("applyForFinancing", "getBillFromDB", err)
		return nil, err
	} else if err == types.ErrNotFound {
		return nil, pty.ErrBillNotFound
	}

	var dpdtToken *pty.DpdtToken
	for _, v := range bill.DpdtTokens {
		if v.Symbol == c.DpdtSymbol {
			dpdtToken = v
		}
	}
	if dpdtToken == nil {
		return nil, pty.ErrCreditTokenDismatch
	}
	//Check if the amount of unused credit token is enough
	if c.Amount > (dpdtToken.Amount - dpdtToken.Used) {
		financelog.Error("applyForFinancing: unused not enough", "unused", dpdtToken.Amount-dpdtToken.Used)
		return nil, pty.ErrUnusedCreditTNotEnough
	}

	//Frozen Bill token, plus amount of used credit token
	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	//从当前账户中冻结token用于向资金方融资
	receiptFrozen, err := action.frozenAssets(pty.TokenExecer, bill.TokenSymbol, c.ApplyAddr, c.Amount)
	if err != nil {
		return nil, err
	}
	kv = append(kv, receiptFrozen.KV...)
	logs = append(logs, receiptFrozen.Logs...)

	//should double check if current value in Bill changed
	dpdtToken.Used += c.Amount
	key := calcBillNewKey(c.Id)
	value := types.Encode(bill)
	kv = append(kv, &types.KeyValue{Key: key, Value: value})

	receipt := &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}
	return receipt, nil
}

//unApplyForFinancing action
func (action *Action) unApplyForFinancing(c *pty.UnApplyForFinancing) (*types.Receipt, error) {
	financelog.Debug("unApplyForFinancing")
	if c.UnApplyAddr != action.fromaddr {
		return nil, pty.ErrOperateAddrMismatch
	}
	if c == nil {
		return nil, types.ErrInvalidParam
	}
	//撤销申请融资的人必须是供应商
	if c.UnApplyIndentity != pty.Supplier {
		return nil, pty.ErrOperateAddrIsNotSupplier
	}

	bill, err := getBillFromDB(action.db, c.Id)
	if err != nil && err != types.ErrNotFound {
		financelog.Error("unApplyForFinancing", "getIouFromDB", err)
		return nil, err
	} else if err == types.ErrNotFound {
		return nil, pty.ErrBillNotFound
	}

	var logs []*types.ReceiptLog
	var kv []*types.KeyValue

	//从localDB中的table表中查找相应的申请融资信息
	local, err := action.getApplyFinancingItem(c.TxHash)
	if err != nil {
		financelog.Error("unApplyForFinancing", "getApplyFinancingItem", err)
		return nil, err
	}

	//操作人必须是白条持有人
	if local.ApplyAddr != c.UnApplyAddr {
		financelog.Error("unApplyForFinancing failed", "applyAddr", local.ApplyAddr, "action address", c.UnApplyAddr)
		return nil, pty.ErrOperateAddrMismatch
	}

	//Unfrozen token
	amount := local.Amount
	if amount > 0 {
		//Active relate token assets,解冻供应商相应的白条币资产
		receiptActive, err := action.activeFrozenAssets(pty.TokenExecer, local.ApplyAddr, local.Amount, bill.TokenSymbol)
		if err != nil {
			financelog.Error("Active frozen assets failed")
			return nil, err
		}
		kv = append(kv, receiptActive.KV...)
		logs = append(logs, receiptActive.Logs...)
	}

	//Update Bill
	var dpdtToken *pty.DpdtToken
	for _, v := range bill.DpdtTokens {
		if v.Symbol == local.DpdtSymbol {
			dpdtToken = v
		}
	}
	if dpdtToken != nil {
		dpdtToken.Used -= local.Amount
		key := calcBillNewKey(c.Id)
		value := types.Encode(bill)
		kv = append(kv, &types.KeyValue{Key: key, Value: value})
	}

	receipt := &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}
	return receipt, nil
}

//confirmFinancing action
func (action *Action) confirmFinancing(c *pty.ConfirmFinancing) (*types.Receipt, error) {
	financelog.Debug("confirmFinancing")
	receiptConfirmFinancing := &pty.ReceiptLogConfirmFinancing{}
	receiptConfirmFinancing.TransfItems = &pty.AssetTransferTags{Items: []*pty.AssetTransferTag{}}
	if c == nil {
		return nil, types.ErrInvalidParam
	}
	if c.ConfirmAddr != action.fromaddr {
		return nil, pty.ErrOperateAddrMismatch
	}
	//确认融资的人必须是资金方
	if c.ConfirmIndentity != pty.Funder {
		return nil, pty.ErrOperateAddrIsNotFunder
	}

	//从申请融资的条目中查看申请信息
	local, err := action.getApplyFinancingItem(c.TxHash)
	if err != nil {
		financelog.Error("confirmFinancing", "getApplyFinancingItem", err)
		return nil, err
	}
	//融资对象必须是授信人
	if local.CreditorAddr != c.ConfirmAddr {
		financelog.Error("confirmFinancing: CreditorAddr mismatch", "CreditorAddr", local.CreditorAddr, "operate addr", c.ConfirmAddr)
		return nil, pty.ErrOperateAddrMismatch
	}

	//获取授信Token信息
	credit, err := getCreditTokenFromDB(action.db, local.DpdtSymbol)
	if err != nil && err != types.ErrNotFound {
		financelog.Error("confirmFinancing", "getCreditTokenFromDB", err)
		return nil, err
	} else if err == types.ErrNotFound {
		return nil, pty.ErrCreditNotFound
	}

	//获取白条信息
	bill, err := getBillFromDB(action.db, local.Id)
	if err != nil && err != types.ErrNotFound {
		financelog.Error("confirmFinancing", "getBillFromDB", err)
		return nil, err
	} else if err == types.ErrNotFound {
		return nil, pty.ErrBillNotFound
	}

	//授信有效性检查
	if action.blocktime > (credit.AddTime + credit.Duration) {
		//授信已经过期, 拒绝融资
		return nil, pty.ErrCreditTokenExpired
	} else if action.blocktime > bill.RepayDate {
		//白条已逾期, 白条逾期之后则不允许融资
		return nil, pty.ErrBillHasExpired
	} else if bill.Status != pty.StatusReleased {
		//白条无效, 不能为其融资
		return nil, pty.ErrBillStatus
	}

	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	if c.Ack {
		//同意融资操作
		//开启币币交换
		//local.Amount申请融资的白条币金额,credit.Rate融資贴现率
		tokenValue := calcTokenValueForFinancing(local.Amount, credit.Rate, (bill.RepayDate-action.blocktime)/pty.OneDaySeconds)
		if tokenValue == 0 {
			return nil, pty.ErrWrongTokenValue
		}
		//查询资金方(确认申请方)的token.CCNY余额是否足够
		balance, err := action.getBalance(pty.CoinsExecer, bill.CoinsSymbol, c.ConfirmAddr)
		if err != nil {
			financelog.Error("confirmFinancing", "getBalance", err)
			return nil, err
		}
		//token.CCNY<tokenValue，没有足够的CCNY进行融資
		if balance < tokenValue {
			financelog.Error("confirmFinancing: no balance", "tokenValue", tokenValue, "balance", balance)
			return nil, pty.ErrAssetsNotEnough
		}
		//将资金方的CCNy资产转账到融资申请人(供应商)地址
		receiptTransferCoins, err := action.transferAssets(pty.CoinsExecer, c.ConfirmAddr, local.ApplyAddr, tokenValue, bill.CoinsSymbol)
		if err != nil {
			financelog.Error("confirmFinancing: transferAssets failed", "err", err)
			return nil, err
		}
		//Add transfer tag to log
		transferItem := &pty.AssetTransferTag{FromAddr: c.ConfirmAddr, ToAddr: local.ApplyAddr, Symbol: bill.CoinsSymbol, Amount: tokenValue, Remark: c.TxHash}
		receiptConfirmFinancing.TransfItems.Items = append(receiptConfirmFinancing.TransfItems.Items, transferItem)

		//融资申请人(供应商)地址将白条币转账到资金方的地址
		receiptTransferTokens, err := action.transferFrozenAssets(pty.TokenExecer, local.ApplyAddr, c.ConfirmAddr, local.Amount, bill.TokenSymbol)
		if err != nil {
			financelog.Error("confirmFinancing: transferFrozenAssets failed", "err", err)
			return nil, err
		}
		//Add transfer tag to log
		transferItem = &pty.AssetTransferTag{ToAddr: c.ConfirmAddr, Symbol: bill.TokenSymbol, FromAddr: local.ApplyAddr, Amount: local.Amount, Remark: c.TxHash}
		receiptConfirmFinancing.TransfItems.Items = append(receiptConfirmFinancing.TransfItems.Items, transferItem)

		//fill with kv and logs
		kv = append(kv, receiptTransferCoins.KV...)
		kv = append(kv, receiptTransferTokens.KV...)

		logs = append(logs, receiptTransferCoins.Logs...)
		logs = append(logs, receiptTransferTokens.Logs...)
	} else {
		//不同意融资操作
		//为申请人解冻白条资产
		receipt, err := action.activeFrozenAssets(pty.TokenExecer, local.ApplyAddr, local.Amount, bill.TokenSymbol)
		if err != nil {
			//解冻失败
			financelog.Error("confirmFinancing: Active frozen assets failed", "apply address", local.ApplyAddr, "apply symbol", bill.TokenSymbol, "apply amount", local.Amount)
			return nil, err
		}
		kv = append(kv, receipt.KV...)
		logs = append(logs, receipt.Logs...)

		//Update Bill
		var dpdtToken *pty.DpdtToken
		for _, v := range bill.DpdtTokens {
			if v.Symbol == local.DpdtSymbol {
				dpdtToken = v
			}
		}
		if dpdtToken != nil {
			dpdtToken.Used -= local.Amount
			//更新白条信息
			key := calcBillNewKey(local.Id)
			value := types.Encode(bill)
			kv = append(kv, &types.KeyValue{Key: key, Value: value})
		}
	}

	log := types.Encode(receiptConfirmFinancing)
	logs = append(logs, &types.ReceiptLog{Ty: pty.TyConfirmFinancingLog, Log: log})

	receipt := &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}
	return receipt, nil
}

// CashBill Action
func (action *Action) CashBill(c *pty.CashBill) (*types.Receipt, error) {
	financelog.Debug("CashBill")
	receiptCashBill := &pty.ReceiptLogCashBill{}
	receiptCashBill.TransfItems = &pty.AssetTransferTags{Items: []*pty.AssetTransferTag{}}
	if c == nil {
		return nil, types.ErrInvalidParam
	}
	if c.CashAddr != action.fromaddr {
		return nil, pty.ErrOperateAddrMismatch
	}

	//删除TokenSymbol的前缀"SLD",获取billID
	id := tokenSymbolToBillId(c.TokenSymbol)
	//获取白条信息
	bill, err := getBillFromDB(action.db, id)
	if err != nil && err != types.ErrNotFound {
		financelog.Error("CashBill", "getBillFromDB", err)
		return nil, err
	} else if err == types.ErrNotFound {
		return nil, pty.ErrBillNotFound
	}

	if bill.Status != pty.StatusReleased {
		financelog.Error("CashBill error happen, Bill status mismatch", "current status", bill.Status)
		return nil, pty.ErrNoPermissionAction
	}

	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	var tokenValue int64

	//获取兑现人的白条币余额
	balance, err := action.getBalance(pty.TokenExecer, bill.TokenSymbol, c.CashAddr)
	if err != nil {
		financelog.Error("Get token amount failed")
		return nil, err
	}
	if balance < c.Amount {
		financelog.Error("Token amount not enough", "balance", balance, "symbol", bill.TokenSymbol)
		return nil, pty.ErrAssetsNotEnough
	}
	//从当前账户中冻结token用于Token兑现，冻结兑现人相应数量的白条币
	receiptFrozen, err := action.frozenAssets(pty.TokenExecer, bill.TokenSymbol, c.CashAddr, c.Amount)
	if err != nil {
		return nil, err
	}

	kv = append(kv, receiptFrozen.KV...)
	logs = append(logs, receiptFrozen.Logs...)

	//填写兑现记录,写道Log中
	cashRecord := &pty.CashRecordLocalDB{
		CashAddr:   c.CashAddr,
		BillId:     bill.Id,
		TxHash:     "0x" + hex.EncodeToString(action.txhash),
		Value:      tokenValue,
		Amount:     c.Amount,
		Timestamp:  action.blocktime,
		BillStatus: bill.Status,
		BillType:   bill.BillType,
	}

	receiptCashBill.Record = cashRecord
	log := types.Encode(receiptCashBill)
	logs = append(logs, &types.ReceiptLog{Ty: pty.TyCashBillLog, Log: log})

	receipt := &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}
	return receipt, nil
}

// RepayBill Action
func (action *Action) RepayBill(c *pty.RepayBill) (*types.Receipt, error) {
	financelog.Debug("RepayBill")
	receiptRepayBill := &pty.ReceiptLogRepayBill{}
	receiptRepayBill.TransfItems = &pty.AssetTransferTags{Items: []*pty.AssetTransferTag{}}
	if c == nil {
		return nil, types.ErrInvalidParam
	}
	if c.RepayAddr != action.fromaddr {
		return nil, pty.ErrOperateAddrMismatch
	}

	bill, err := getBillFromDB(action.db, c.Id)
	if err != nil && err != types.ErrNotFound {
		financelog.Error("RepayBill", "getBillFromDB", err)
		return nil, err
	} else if err == types.ErrNotFound {
		return nil, pty.ErrBillNotFound
	}

	//Bill和Bond都会遵守的准则： 总的应还款期数和当前已还款期数相等时，认为还款结束
	//bill还款只有一期
	if bill.CurrentRepayedCount == bill.RepayCount {
		financelog.Warn("Borrower has repaid all the arrears, need not repay again", "repaied count", bill.CurrentRepayedCount)
		return nil, nil
	}

	//还款信息记录
	repayRecord := &pty.RepayRecordLocalDB{
		RepayAddr:   c.RepayAddr,
		BillId:      bill.Id,
		TxHash:      "0x" + hex.EncodeToString(action.txhash),
		Timestamp:   action.blocktime,
		BillCreator: bill.Borrower,
		BillType:    bill.BillType,
	}

	receiptRepayBill.Record = repayRecord

	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	var repayAmount int64
	//For support Bill
	if bill.BillType == pty.BILL {
		//根据白条id获取申请兑现列表中的信息
		cashList, err := action.getCashList(c.Id)
		if err != nil {
			return nil, err
		}
		var totalCashAmount int64
		for _, v := range cashList {
			totalCashAmount += v.Amount
		}
		//TotalRedeemToken 总的赎回的token数额,就是之前有人兑现白条币，这次还有人兑现白条币
		if bill.LoanAmount < bill.TotalRedeemToken+totalCashAmount {
			financelog.Error("The total token amount which need to be cashed mismatch.", "totalCashAmount", totalCashAmount)
			return nil, pty.ErrTotalCashAmountMismatch
		}

		//Cash each item 根据CashList兑现所有的申请兑现信息
		for _, item := range cashList {
			var repayInfo *pty.ReplyBillRepayInfo
			cashTime := action.blocktime

			//申请时间>白条发布时间+白条流通时间，证明白条已经过期
			if item.Timestamp > bill.IssueDate+bill.CirculationTime {
				//兑现时间=(操作时间-申请时间)+白条发布时间+白条流通时间
				cashTime = bill.IssueDate + bill.CirculationTime + action.blocktime - item.Timestamp
			}

			m := &BillRepayMethod{}
			//item.Amount兑现金额
			repayInfo = m.CalcIouRepayInfo(bill, cashTime, item.Amount)
			if repayInfo == nil {
				continue
			}
			repayValue := repayInfo.ActualRepayAmount
			if repayValue == 0 {
				//实际还款==0,这个兑现跳出
				continue
			}

			//CCNY transfer 将借款人的CCNY转账到白条币拥有者的地址
			receiptTransferCoins, err := action.transferAssets(pty.CoinsExecer, c.RepayAddr, item.TokenOwner, repayValue, bill.CoinsSymbol)
			if err != nil {
				financelog.Error("RepayBill", "receiptTransferCoins", err)
				return nil, err
			}
			//Add transfer tags to log 记录转账操作 申请兑现的item.TxHash
			transferItem := &pty.AssetTransferTag{FromAddr: c.RepayAddr, ToAddr: item.TokenOwner, Symbol: bill.CoinsSymbol, Amount: repayValue, Remark: item.TxHash}
			receiptRepayBill.TransfItems.Items = append(receiptRepayBill.TransfItems.Items, transferItem)
			kv = append(kv, receiptTransferCoins.KV...)
			logs = append(logs, receiptTransferCoins.Logs...)

			//Token transfer 将白条币转到中间地址MiddleAddr，等同于销毁白条币，在区块链后端的平行链上超级管理员进行白条的销毁操作
			receiptTransferTokens, err := action.transferFrozenAssets(pty.TokenExecer, item.TokenOwner, pty.MiddleAddr, item.Amount, bill.TokenSymbol)
			if err != nil {
				financelog.Error("RepayBill", "transferFrozenAssets", err)
				return nil, err
			}
			//Add transfer tag to log 申请兑现的item.TxHash
			transferItem = &pty.AssetTransferTag{ToAddr: pty.MiddleAddr, Symbol: bill.TokenSymbol, FromAddr: item.TokenOwner, Amount: item.Amount, Remark: item.TxHash}
			receiptRepayBill.TransfItems.Items = append(receiptRepayBill.TransfItems.Items, transferItem)
			kv = append(kv, receiptTransferTokens.KV...)
			logs = append(logs, receiptTransferTokens.Logs...)

			//Add cashHash to repay record 兑现的Hash列表
			repayRecord.CashHashes = append(repayRecord.CashHashes, item.TxHash)

			//Update Bill
			repayAmount += repayValue
			bill.TotalCashAmount += repayValue   //当前已兑现的总额度
			bill.TotalRepayAmount += repayValue  //总还款额度CCNY
			bill.TotalRedeemToken += item.Amount //总的赎回的token数额
			bill.TotalCashedToken += item.Amount //总的兑现的Token数

			//如果总的赎回Token==借款总欸，那摩还款结束
			if bill.TotalRedeemToken == bill.LoanAmount {
				//赎回所有token,完成还款
				bill.CurrentRepayedCount = 1
				bill.Status = pty.StatusBlocked
			}
		}
	}

	repayRecord.Value = repayAmount //还款额度(token价值)，还了多少CCNY
	//票据的还款期数为1, 所以当还清所有款项的时候设CurrentRepayedCount为1; 债券则按照实际还款期数来计算;
	if bill.CurrentRepayedCount == bill.RepayCount {
		//已还清款项
		bill.RepayDate = 0x7FFFFFFFFFFFFFFF
		repayRecord.IsArrearsCleared = true

		//解冻对应授信
		if bill.NeedDpdtToken {
			for _, v := range bill.DpdtTokens {
				if v.Amount <= 0 {
					continue
				}
				receipt, err := action.activeFrozenAssets(pty.TokenExecer, bill.Borrower, v.Amount, v.Symbol)
				if err != nil {
					//解冻失败
					financelog.Error("Active frozen assets for credit failed", "credit owner", bill.Borrower, "credit symbol", v.Symbol, "credit amount", v.Amount)
					return nil, err
				}
				kv = append(kv, receipt.KV...)
				logs = append(logs, receipt.Logs...)
			}
		}

	} else {
		repayRecord.IsArrearsCleared = false
	}

	key := calcBillNewKey(c.Id)
	value := types.Encode(bill)
	kv = append(kv, &types.KeyValue{Key: key, Value: value})

	log := types.Encode(receiptRepayBill)
	logs = append(logs, &types.ReceiptLog{Ty: pty.TyRepayBillLog, Log: log})

	receipt := &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}
	return receipt, nil
}

// ReportBroken Action
func (action *Action) ReportBroken(c *pty.ReportBroken) (*types.Receipt, error) {
	financelog.Debug("ReportBroken")
	if c == nil {
		return nil, types.ErrInvalidParam
	}
	if c.ReportAddr != action.fromaddr {
		return nil, pty.ErrOperateAddrMismatch
	}

	bill, err := getBillFromDB(action.db, c.Id)
	if err != nil && err !=
		types.ErrNotFound {
		financelog.Error("ReportBroken", "getBillFromDB", err)
		return nil, err
	} else if err == types.ErrNotFound {
		return nil, pty.ErrBillNotFound
	}

	//Bill
	if action.blocktime < bill.IssueDate+bill.CirculationTime {
		financelog.Error("User have not broken Bill", "repay due date", bill.IssueDate+bill.CirculationTime, "now time", action.blocktime)
		return nil, pty.ErrNoPermissionAction
	}

	var repayInfo *pty.ReplyBillRepayInfo
	var overdueAmount int64
	if bill.BillType == pty.BILL {
		//查询白条兑现申请list
		cashList, err := action.getCashList(c.Id)
		if err != nil {
			return nil, err
		}
		var totalCashAmount int64
		for _, v := range cashList {
			totalCashAmount += v.Amount
		}
		m := &BillRepayMethod{}
		//action.blocktime入参是否合理--Todo
		//查看总的兑现金额
		repayInfo = m.CalcIouRepayInfo(bill, action.blocktime, totalCashAmount)
		if repayInfo != nil {
			overdueAmount = repayInfo.ActualRepayAmount - repayInfo.NormalRepayAmount
		}
	}

	//举报记录，存在逾期金额
	brokenRecord := &pty.BrokenRecordStateDB{
		BorrowerAddr: bill.Borrower,
		Phone:        bill.Phone,
		Identifier:   bill.Identifier,
		Name:         bill.Name,
		BillId:       bill.Id,
		TxHash:       "0x" + hex.EncodeToString(action.txhash),
		OverdueDays:  (action.blocktime - bill.RepayDate) / (24 * 3600),
		OverdueValue: overdueAmount,
	}

	var logs []*types.ReceiptLog
	var kv []*types.KeyValue
	log := types.Encode(brokenRecord)
	logs = append(logs, &types.ReceiptLog{Ty: pty.TyReportBrokenLog, Log: log})

	receipt := &types.Receipt{Ty: types.ExecOk, KV: kv, Logs: logs}
	return receipt, nil
}
