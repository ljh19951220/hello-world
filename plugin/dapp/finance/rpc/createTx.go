package rpc

import (
	"context"
	"encoding/hex"
	"github.com/33cn/chain33/types"
	financetypes "github.com/33cn/plugin/plugin/dapp/finance/types"
)

func (c *channelClient) CreateDepositAsset(ctx context.Context, v *financetypes.DepositAsset) (*types.UnsignTx, error) {
	acct := &financetypes.FinanceAction{
		Ty:    financetypes.TyDepositAssetAction,
		Value: &financetypes.FinanceAction_DepositAsset{DepositAsset: v},
	}

	cfg := c.GetConfig()
	tx, err := types.CreateFormatTx(cfg, cfg.ExecName(financetypes.FinanceX), types.Encode(acct))
	if err != nil {
		return nil, err
	}

	data := types.Encode(tx)
	return &types.UnsignTx{Data: data}, nil
}
func (j *Jrpc) CreateDepositAsset(in *financetypes.DepositAsset, result *interface{}) error {
	//这里直接转发至grpc接口
	reply, err := j.cli.CreateDepositAsset(context.Background(), in)
	if err != nil {
		return err
	}
	*result = hex.EncodeToString(reply.Data)
	return nil
}

func (c *channelClient) CreateWithdrawlAsset(ctx context.Context, v *financetypes.WithdrawlAsset) (*types.UnsignTx, error) {
	acct := &financetypes.FinanceAction{
		Ty:    financetypes.TyWithdrawlAssetAction,
		Value: &financetypes.FinanceAction_WithdrawlAsset{WithdrawlAsset: v},
	}

	cfg := c.GetConfig()
	tx, err := types.CreateFormatTx(cfg, cfg.ExecName(financetypes.FinanceX), types.Encode(acct))
	if err != nil {
		return nil, err
	}

	data := types.Encode(tx)
	return &types.UnsignTx{Data: data}, nil
}
func (j *Jrpc) CreateWithdrawlAsset(in *financetypes.WithdrawlAsset, result *interface{}) error {
	//这里直接转发至grpc接口
	reply, err := j.cli.CreateWithdrawlAsset(context.Background(), in)
	if err != nil {
		return err
	}
	*result = hex.EncodeToString(reply.Data)
	return nil
}

func (c *channelClient) CreateTransferCoins(ctx context.Context, v *financetypes.TransferCoins) (*types.UnsignTx, error) {
	acct := &financetypes.FinanceAction{
		Ty:    financetypes.TyTransferCoinsAction,
		Value: &financetypes.FinanceAction_TransferCoins{TransferCoins: v},
	}

	cfg := c.GetConfig()
	tx, err := types.CreateFormatTx(cfg, cfg.ExecName(financetypes.FinanceX), types.Encode(acct))
	if err != nil {
		return nil, err
	}

	data := types.Encode(tx)
	return &types.UnsignTx{Data: data}, nil
}
func (j *Jrpc) CreateTransferCoins(in *financetypes.TransferCoins, result *interface{}) error {
	//这里直接转发至grpc接口
	reply, err := j.cli.CreateTransferCoins(context.Background(), in)
	if err != nil {
		return err
	}
	*result = hex.EncodeToString(reply.Data)
	return nil
}

func (c *channelClient) CreateAddCreditToken(ctx context.Context, v *financetypes.AddCreditToken) (*types.UnsignTx, error) {
	acct := &financetypes.FinanceAction{
		Ty:    financetypes.TyAddCreditTokenAction,
		Value: &financetypes.FinanceAction_AddCreditToken{AddCreditToken: v},
	}

	cfg := c.GetConfig()
	tx, err := types.CreateFormatTx(cfg, cfg.ExecName(financetypes.FinanceX), types.Encode(acct))
	if err != nil {
		return nil, err
	}

	data := types.Encode(tx)
	return &types.UnsignTx{Data: data}, nil
}
func (j *Jrpc) CreateAddCreditToken(in *financetypes.AddCreditToken, result *interface{}) error {
	//这里直接转发至grpc接口
	reply, err := j.cli.CreateAddCreditToken(context.Background(), in)
	if err != nil {
		return err
	}
	*result = hex.EncodeToString(reply.Data)
	return nil
}

func (c *channelClient) CreateFinanceBill(ctx context.Context, v *financetypes.CreateBill) (*types.UnsignTx, error) {
	acct := &financetypes.FinanceAction{
		Ty:    financetypes.TyCreateBillAction,
		Value: &financetypes.FinanceAction_CreateBill{CreateBill: v},
	}

	cfg := c.GetConfig()
	tx, err := types.CreateFormatTx(cfg, cfg.ExecName(financetypes.FinanceX), types.Encode(acct))
	if err != nil {
		return nil, err
	}

	data := types.Encode(tx)
	return &types.UnsignTx{Data: data}, nil
}
func (j *Jrpc) CreateFinanceBill(in *financetypes.CreateBill, result *interface{}) error {
	//这里直接转发至grpc接口
	reply, err := j.cli.CreateFinanceBill(context.Background(), in)
	if err != nil {
		return err
	}
	*result = hex.EncodeToString(reply.Data)
	return nil
}

func (c *channelClient) CreateReleaseBill(ctx context.Context, v *financetypes.ReleaseBill) (*types.UnsignTx, error) {
	acct := &financetypes.FinanceAction{
		Ty:    financetypes.TyReleaseBillAction,
		Value: &financetypes.FinanceAction_ReleaseBill{ReleaseBill: v},
	}

	cfg := c.GetConfig()
	tx, err := types.CreateFormatTx(cfg, cfg.ExecName(financetypes.FinanceX), types.Encode(acct))
	if err != nil {
		return nil, err
	}

	data := types.Encode(tx)
	return &types.UnsignTx{Data: data}, nil
}
func (j *Jrpc) CreateReleaseBill(in *financetypes.ReleaseBill, result *interface{}) error {
	//这里直接转发至grpc接口
	reply, err := j.cli.CreateReleaseBill(context.Background(), in)
	if err != nil {
		return err
	}
	*result = hex.EncodeToString(reply.Data)
	return nil
}

func (c *channelClient) CreateUnReleaseBill(ctx context.Context, v *financetypes.UnReleaseBill) (*types.UnsignTx, error) {
	acct := &financetypes.FinanceAction{
		Ty:    financetypes.TyUnReleaseBillAction,
		Value: &financetypes.FinanceAction_UnReleaseBill{UnReleaseBill: v},
	}

	cfg := c.GetConfig()
	tx, err := types.CreateFormatTx(cfg, cfg.ExecName(financetypes.FinanceX), types.Encode(acct))
	if err != nil {
		return nil, err
	}

	data := types.Encode(tx)
	return &types.UnsignTx{Data: data}, nil
}
func (j *Jrpc) CreateUnReleaseBill(in *financetypes.UnReleaseBill, result *interface{}) error {
	//这里直接转发至grpc接口
	reply, err := j.cli.CreateUnReleaseBill(context.Background(), in)
	if err != nil {
		return err
	}
	*result = hex.EncodeToString(reply.Data)
	return nil
}

func (c *channelClient) CreateSplitBill(ctx context.Context, v *financetypes.SplitBill) (*types.UnsignTx, error) {
	acct := &financetypes.FinanceAction{
		Ty:    financetypes.TySplitBillAction,
		Value: &financetypes.FinanceAction_SplitBill{SplitBill: v},
	}

	cfg := c.GetConfig()
	tx, err := types.CreateFormatTx(cfg, cfg.ExecName(financetypes.FinanceX), types.Encode(acct))
	if err != nil {
		return nil, err
	}

	data := types.Encode(tx)
	return &types.UnsignTx{Data: data}, nil
}
func (j *Jrpc) CreateSplitBill(in *financetypes.SplitBill, result *interface{}) error {
	//这里直接转发至grpc接口
	reply, err := j.cli.CreateSplitBill(context.Background(), in)
	if err != nil {
		return err
	}
	*result = hex.EncodeToString(reply.Data)
	return nil
}

func (c *channelClient) CreateUnSplitBill(ctx context.Context, v *financetypes.UnSplitBill) (*types.UnsignTx, error) {
	acct := &financetypes.FinanceAction{
		Ty:    financetypes.TyUnSplitBillAction,
		Value: &financetypes.FinanceAction_UnSplitBill{UnSplitBill: v},
	}

	cfg := c.GetConfig()
	tx, err := types.CreateFormatTx(cfg, cfg.ExecName(financetypes.FinanceX), types.Encode(acct))
	if err != nil {
		return nil, err
	}

	data := types.Encode(tx)
	return &types.UnsignTx{Data: data}, nil
}
func (j *Jrpc) CreateUnSplitBill(in *financetypes.UnSplitBill, result *interface{}) error {
	//这里直接转发至grpc接口
	reply, err := j.cli.CreateUnSplitBill(context.Background(), in)
	if err != nil {
		return err
	}
	*result = hex.EncodeToString(reply.Data)
	return nil
}

func (c *channelClient) CreateConfirmSplitBill(ctx context.Context, v *financetypes.ConfirmSplitBill) (*types.UnsignTx, error) {
	acct := &financetypes.FinanceAction{
		Ty:    financetypes.TyConfirmSplitBillAction,
		Value: &financetypes.FinanceAction_ConfirmSplitBill{ConfirmSplitBill: v},
	}

	cfg := c.GetConfig()
	tx, err := types.CreateFormatTx(cfg, cfg.ExecName(financetypes.FinanceX), types.Encode(acct))
	if err != nil {
		return nil, err
	}

	data := types.Encode(tx)
	return &types.UnsignTx{Data: data}, nil
}
func (j *Jrpc) CreateConfirmSplitBill(in *financetypes.ConfirmSplitBill, result *interface{}) error {
	//这里直接转发至grpc接口
	reply, err := j.cli.CreateConfirmSplitBill(context.Background(), in)
	if err != nil {
		return err
	}
	*result = hex.EncodeToString(reply.Data)
	return nil
}

func (c *channelClient) CreateApplyForFinancing(ctx context.Context, v *financetypes.ApplyForFinancing) (*types.UnsignTx, error) {
	acct := &financetypes.FinanceAction{
		Ty:    financetypes.TyApplyForFinancingAction,
		Value: &financetypes.FinanceAction_ApplyForFinancing{ApplyForFinancing: v},
	}

	cfg := c.GetConfig()
	tx, err := types.CreateFormatTx(cfg, cfg.ExecName(financetypes.FinanceX), types.Encode(acct))
	if err != nil {
		return nil, err
	}

	data := types.Encode(tx)
	return &types.UnsignTx{Data: data}, nil
}
func (j *Jrpc) CreateApplyForFinancing(in *financetypes.ApplyForFinancing, result *interface{}) error {
	//这里直接转发至grpc接口
	reply, err := j.cli.CreateApplyForFinancing(context.Background(), in)
	if err != nil {
		return err
	}
	*result = hex.EncodeToString(reply.Data)
	return nil
}

func (c *channelClient) CreateUnApplyForFinancing(ctx context.Context, v *financetypes.UnApplyForFinancing) (*types.UnsignTx, error) {
	acct := &financetypes.FinanceAction{
		Ty:    financetypes.TyUnApplyForFinancingAction,
		Value: &financetypes.FinanceAction_UnApplyForFinancing{UnApplyForFinancing: v},
	}

	cfg := c.GetConfig()
	tx, err := types.CreateFormatTx(cfg, cfg.ExecName(financetypes.FinanceX), types.Encode(acct))
	if err != nil {
		return nil, err
	}

	data := types.Encode(tx)
	return &types.UnsignTx{Data: data}, nil
}
func (j *Jrpc) CreateUnApplyForFinancing(in *financetypes.UnApplyForFinancing, result *interface{}) error {
	//这里直接转发至grpc接口
	reply, err := j.cli.CreateUnApplyForFinancing(context.Background(), in)
	if err != nil {
		return err
	}
	*result = hex.EncodeToString(reply.Data)
	return nil
}

func (c *channelClient) CreateConfirmFinancing(ctx context.Context, v *financetypes.ConfirmFinancing) (*types.UnsignTx, error) {
	acct := &financetypes.FinanceAction{
		Ty:    financetypes.TyConfirmFinancingAction,
		Value: &financetypes.FinanceAction_ConfirmFinancing{ConfirmFinancing: v},
	}

	cfg := c.GetConfig()
	tx, err := types.CreateFormatTx(cfg, cfg.ExecName(financetypes.FinanceX), types.Encode(acct))
	if err != nil {
		return nil, err
	}

	data := types.Encode(tx)
	return &types.UnsignTx{Data: data}, nil
}
func (j *Jrpc) CreateConfirmFinancing(in *financetypes.ConfirmFinancing, result *interface{}) error {
	//这里直接转发至grpc接口
	reply, err := j.cli.CreateConfirmFinancing(context.Background(), in)
	if err != nil {
		return err
	}
	*result = hex.EncodeToString(reply.Data)
	return nil
}

func (c *channelClient) CreateDeliverBill(ctx context.Context, v *financetypes.DeliverBill) (*types.UnsignTx, error) {
	acct := &financetypes.FinanceAction{
		Ty:    financetypes.TyDeliverBillAction,
		Value: &financetypes.FinanceAction_DeliverBill{DeliverBill: v},
	}

	cfg := c.GetConfig()
	tx, err := types.CreateFormatTx(cfg, cfg.ExecName(financetypes.FinanceX), types.Encode(acct))
	if err != nil {
		return nil, err
	}

	data := types.Encode(tx)
	return &types.UnsignTx{Data: data}, nil
}
func (j *Jrpc) CreateDeliverBill(in *financetypes.DeliverBill, result *interface{}) error {
	//这里直接转发至grpc接口
	reply, err := j.cli.CreateDeliverBill(context.Background(), in)
	if err != nil {
		return err
	}
	*result = hex.EncodeToString(reply.Data)
	return nil
}

func (c *channelClient) CreateUnDeliverBill(ctx context.Context, v *financetypes.UnDeliverBill) (*types.UnsignTx, error) {
	acct := &financetypes.FinanceAction{
		Ty:    financetypes.TyUnDeliverBillAction,
		Value: &financetypes.FinanceAction_UnDeliverBill{UnDeliverBill: v},
	}

	cfg := c.GetConfig()
	tx, err := types.CreateFormatTx(cfg, cfg.ExecName(financetypes.FinanceX), types.Encode(acct))
	if err != nil {
		return nil, err
	}

	data := types.Encode(tx)
	return &types.UnsignTx{Data: data}, nil
}
func (j *Jrpc) CreateUnDeliverBill(in *financetypes.UnDeliverBill, result *interface{}) error {
	//这里直接转发至grpc接口
	reply, err := j.cli.CreateUnDeliverBill(context.Background(), in)
	if err != nil {
		return err
	}
	*result = hex.EncodeToString(reply.Data)
	return nil
}

func (c *channelClient) CreateConfirmDeliverBill(ctx context.Context, v *financetypes.ConfirmDeliverBill) (*types.UnsignTx, error) {
	acct := &financetypes.FinanceAction{
		Ty:    financetypes.TyConfirmDeliverBillAction,
		Value: &financetypes.FinanceAction_ConfirmDeliverBill{ConfirmDeliverBill: v},
	}

	cfg := c.GetConfig()
	tx, err := types.CreateFormatTx(cfg, cfg.ExecName(financetypes.FinanceX), types.Encode(acct))
	if err != nil {
		return nil, err
	}

	data := types.Encode(tx)
	return &types.UnsignTx{Data: data}, nil
}
func (j *Jrpc) CreateConfirmDeliverBill(in *financetypes.ConfirmDeliverBill, result *interface{}) error {
	//这里直接转发至grpc接口
	reply, err := j.cli.CreateConfirmDeliverBill(context.Background(), in)
	if err != nil {
		return err
	}
	*result = hex.EncodeToString(reply.Data)
	return nil
}

func (c *channelClient) CreateCashBill(ctx context.Context, v *financetypes.CashBill) (*types.UnsignTx, error) {
	acct := &financetypes.FinanceAction{
		Ty:    financetypes.TyCashBillAction,
		Value: &financetypes.FinanceAction_CashBill{CashBill: v},
	}

	cfg := c.GetConfig()
	tx, err := types.CreateFormatTx(cfg, cfg.ExecName(financetypes.FinanceX), types.Encode(acct))
	if err != nil {
		return nil, err
	}

	data := types.Encode(tx)
	return &types.UnsignTx{Data: data}, nil
}
func (j *Jrpc) CreateCashBill(in *financetypes.CashBill, result *interface{}) error {
	//这里直接转发至grpc接口
	reply, err := j.cli.CreateCashBill(context.Background(), in)
	if err != nil {
		return err
	}
	*result = hex.EncodeToString(reply.Data)
	return nil
}

func (c *channelClient) CreateRepayBill(ctx context.Context, v *financetypes.RepayBill) (*types.UnsignTx, error) {
	acct := &financetypes.FinanceAction{
		Ty:    financetypes.TyRepayBillAction,
		Value: &financetypes.FinanceAction_RepayBill{RepayBill: v},
	}

	cfg := c.GetConfig()
	tx, err := types.CreateFormatTx(cfg, cfg.ExecName(financetypes.FinanceX), types.Encode(acct))
	if err != nil {
		return nil, err
	}

	data := types.Encode(tx)
	return &types.UnsignTx{Data: data}, nil
}
func (j *Jrpc) CreateRepayBill(in *financetypes.RepayBill, result *interface{}) error {
	//这里直接转发至grpc接口
	reply, err := j.cli.CreateRepayBill(context.Background(), in)
	if err != nil {
		return err
	}
	*result = hex.EncodeToString(reply.Data)
	return nil
}

func (c *channelClient) CreateReportBroken(ctx context.Context, v *financetypes.ReportBroken) (*types.UnsignTx, error) {
	acct := &financetypes.FinanceAction{
		Ty:    financetypes.TyReportBrokenAction,
		Value: &financetypes.FinanceAction_ReportBroken{ReportBroken: v},
	}

	cfg := c.GetConfig()
	tx, err := types.CreateFormatTx(cfg, cfg.ExecName(financetypes.FinanceX), types.Encode(acct))
	if err != nil {
		return nil, err
	}

	data := types.Encode(tx)
	return &types.UnsignTx{Data: data}, nil
}
func (j *Jrpc) CreateReportBroken(in *financetypes.ReportBroken, result *interface{}) error {
	//这里直接转发至grpc接口
	reply, err := j.cli.CreateReportBroken(context.Background(), in)
	if err != nil {
		return err
	}
	*result = hex.EncodeToString(reply.Data)
	return nil
}
