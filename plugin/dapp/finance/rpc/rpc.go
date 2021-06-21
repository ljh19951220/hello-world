package rpc

import (
	"context"
	"github.com/33cn/chain33/types"
	financetypes "github.com/33cn/plugin/plugin/dapp/finance/types"
)

/*
 * 实现json rpc和grpc service接口
 * json rpc用Jrpc结构作为接收实例
 * grpc使用channelClient结构作为接收实例
 */

/***********************************QueryAccountAsset查询合约子账户资产信息*******************************************/
//QueryAccountAsset query account assets info
func (c *channelClient) QueryAccountAsset(ctx context.Context, in *financetypes.ReqAccountAsset) (*financetypes.ReplyAccountAsset, error) {
	msg, err := c.Query(financetypes.FinanceX, "AccountAsset", in)
	if err != nil {
		return nil, err
	}
	if reply, ok := msg.(*financetypes.ReplyAccountAsset); ok {
		return reply, nil
	}
	return nil, types.ErrTypeAsset
}
func (j *Jrpc) QueryAccountAsset(in *financetypes.ReqAccountAsset, result *interface{}) error {
	//这里直接转发至grpc接口
	reply, err := j.cli.QueryAccountAsset(context.Background(), in)
	if err != nil {
		return err
	}
	*result = *reply
	return nil
}

//QueryUserAssets query user assets info
func (c *channelClient) QueryUserAssets(ctx context.Context, in *financetypes.ReqFinanceUserAssets) (*financetypes.ReplyFinanceUserAssets, error) {
	msg, err := c.Query(financetypes.FinanceX, "UserAssets", in)
	if err != nil {
		return nil, err
	}
	if reply, ok := msg.(*financetypes.ReplyFinanceUserAssets); ok {
		return reply, nil
	}
	return nil, types.ErrTypeAsset
}
func (j *Jrpc) QueryUserAssets(in *financetypes.ReqFinanceUserAssets, result *interface{}) error {
	//这里直接转发至grpc接口
	reply, err := j.cli.QueryUserAssets(context.Background(), in)
	if err != nil {
		return err
	}
	*result = *reply
	return nil
}

/************************查询出入金(DWRecord)相关操作信息存款操作*******************************************************/
//QueryDWRecord query Deposit and withdrawal record
func (c *channelClient) QueryDWRecord(ctx context.Context, in *financetypes.ReqDWRecord) (*financetypes.ReplyDWRecord, error) {
	msg, err := c.Query(financetypes.FinanceX, "DWRecord", in)
	if err != nil {
		return nil, err
	}
	if reply, ok := msg.(*financetypes.ReplyDWRecord); ok {
		return reply, nil
	}
	return nil, types.ErrTypeAsset
}
func (j *Jrpc) QueryDWRecord(in *financetypes.ReqDWRecord, result *interface{}) error {
	//这里直接转发至grpc接口
	reply, err := j.cli.QueryDWRecord(context.Background(), in)
	if err != nil {
		return err
	}
	*result = *reply
	return nil
}

//QueryDWRecordBySymbol query Deposit and withdrawal record
func (c *channelClient) QueryDWRecordBySymbol(ctx context.Context, in *financetypes.ReqDWRecord) (*financetypes.ReplyDWRecord, error) {
	msg, err := c.Query(financetypes.FinanceX, "DWRecordBySymbol", in)
	if err != nil {
		return nil, err
	}
	if reply, ok := msg.(*financetypes.ReplyDWRecord); ok {
		return reply, nil
	}
	return nil, types.ErrTypeAsset
}
func (j *Jrpc) QueryDWRecordBySymbol(in *financetypes.ReqDWRecord, result *interface{}) error {
	//这里直接转发至grpc接口
	reply, err := j.cli.QueryDWRecordBySymbol(context.Background(), in)
	if err != nil {
		return err
	}
	*result = *reply
	return nil
}

//QueryDWRecordByAddr query Deposit and withdrawal record
func (c *channelClient) QueryDWRecordByAddr(ctx context.Context, in *financetypes.ReqDWRecord) (*financetypes.ReplyDWRecord, error) {
	msg, err := c.Query(financetypes.FinanceX, "DWRecordByAddr", in)
	if err != nil {
		return nil, err
	}
	if reply, ok := msg.(*financetypes.ReplyDWRecord); ok {
		return reply, nil
	}
	return nil, types.ErrTypeAsset
}
func (j *Jrpc) QueryDWRecordByAddr(in *financetypes.ReqDWRecord, result *interface{}) error {
	//这里直接转发至grpc接口
	reply, err := j.cli.QueryDWRecordByAddr(context.Background(), in)
	if err != nil {
		return err
	}
	*result = *reply
	return nil
}

/***********************************TransferCoinsRecord查询合约内转账记录*******************************************/
//QueryTransferCoinsRecord
func (c *channelClient) QueryTransferCoinsRecord(ctx context.Context, in *financetypes.ReqTransferCoinsRecord) (*financetypes.ReplyTransferCoinsRecord, error) {
	msg, err := c.Query(financetypes.FinanceX, "TransferCoinsRecord", in)
	if err != nil {
		return nil, err
	}
	if reply, ok := msg.(*financetypes.ReplyTransferCoinsRecord); ok {
		return reply, nil
	}
	return nil, types.ErrTypeAsset
}
func (j *Jrpc) QueryTransferCoinsRecord(in *financetypes.ReqTransferCoinsRecord, result *interface{}) error {
	//这里直接转发至grpc接口
	reply, err := j.cli.QueryTransferCoinsRecord(context.Background(), in)
	if err != nil {
		return err
	}
	*result = *reply
	return nil
}

func (c *channelClient) QueryTransferCoinsRecordByAddr(ctx context.Context, in *financetypes.ReqTransferCoinsRecord) (*financetypes.ReplyTransferCoinsRecord, error) {
	msg, err := c.Query(financetypes.FinanceX, "TransferCoinsRecordByAddr", in)
	if err != nil {
		return nil, err
	}
	if reply, ok := msg.(*financetypes.ReplyTransferCoinsRecord); ok {
		return reply, nil
	}
	return nil, types.ErrTypeAsset
}
func (j *Jrpc) QueryTransferCoinsRecordByAddr(in *financetypes.ReqTransferCoinsRecord, result *interface{}) error {
	//这里直接转发至grpc接口
	reply, err := j.cli.QueryTransferCoinsRecordByAddr(context.Background(), in)
	if err != nil {
		return err
	}
	*result = *reply
	return nil
}

func (c *channelClient) QueryTransferCoinsRecordBySymbol(ctx context.Context, in *financetypes.ReqTransferCoinsRecord) (*financetypes.ReplyTransferCoinsRecord, error) {
	msg, err := c.Query(financetypes.FinanceX, "TransferCoinsRecordBySymbol", in)
	if err != nil {
		return nil, err
	}
	if reply, ok := msg.(*financetypes.ReplyTransferCoinsRecord); ok {
		return reply, nil
	}
	return nil, types.ErrTypeAsset
}
func (j *Jrpc) QueryTransferCoinsRecordBySymbol(in *financetypes.ReqTransferCoinsRecord, result *interface{}) error {
	//这里直接转发至grpc接口
	reply, err := j.cli.QueryTransferCoinsRecordBySymbol(context.Background(), in)
	if err != nil {
		return err
	}
	*result = *reply
	return nil
}

/***************************query credit 查询授信相关操作********************************************************************/
//QueryCreditDetail query credit detail info
func (c *channelClient) QueryCreditDetail(ctx context.Context, in *financetypes.ReqCreditDetail) (*financetypes.ReplyCreditDetail, error) {
	msg, err := c.Query(financetypes.FinanceX, "CreditDetail", in)
	if err != nil {
		return nil, err
	}
	if reply, ok := msg.(*financetypes.ReplyCreditDetail); ok {
		return reply, nil
	}
	return nil, types.ErrTypeAsset
}
func (j *Jrpc) QueryCreditDetail(in *financetypes.ReqCreditDetail, result *interface{}) error {
	//这里直接转发至grpc接口
	reply, err := j.cli.QueryCreditDetail(context.Background(), in)
	if err != nil {
		return err
	}
	*result = *reply
	return nil
}

//QueryCreditForFunder query credit info for funder
func (c *channelClient) QueryCreditForFunder(ctx context.Context, in *financetypes.ReqCreditFunder) (*financetypes.ReplyCreditFunder, error) {
	msg, err := c.Query(financetypes.FinanceX, "CreditForFunder", in)
	if err != nil {
		return nil, err
	}
	if reply, ok := msg.(*financetypes.ReplyCreditFunder); ok {
		return reply, nil
	}
	return nil, types.ErrTypeAsset
}
func (j *Jrpc) QueryCreditForFunder(in *financetypes.ReqCreditFunder, result *interface{}) error {
	//这里直接转发至grpc接口
	reply, err := j.cli.QueryCreditForFunder(context.Background(), in)
	if err != nil {
		return err
	}
	*result = *reply
	return nil
}

//QueryCreditForCoreFirm query credit info for core firm
func (c *channelClient) QueryCreditForCoreFirm(ctx context.Context, in *financetypes.ReqCreditCoreFirm) (*financetypes.ReplyCreditCoreFirm, error) {
	msg, err := c.Query(financetypes.FinanceX, "CreditForCoreFirm", in)
	if err != nil {
		return nil, err
	}
	if reply, ok := msg.(*financetypes.ReplyCreditCoreFirm); ok {
		return reply, nil
	}
	return nil, types.ErrTypeAsset
}
func (j *Jrpc) QueryCreditForCoreFirm(in *financetypes.ReqCreditCoreFirm, result *interface{}) error {
	//这里直接转发至grpc接口
	reply, err := j.cli.QueryCreditForCoreFirm(context.Background(), in)
	if err != nil {
		return err
	}
	*result = *reply
	return nil
}

//QueryBillInfo query Bill info
func (c *channelClient) QueryBillInfo(ctx context.Context, in *financetypes.ReqBillInfo) (*financetypes.ReplyBillInfo, error) {
	msg, err := c.Query(financetypes.FinanceX, "BillInfo", in)
	if err != nil {
		return nil, err
	}
	if reply, ok := msg.(*financetypes.ReplyBillInfo); ok {
		return reply, nil
	}
	return nil, types.ErrTypeAsset
}
func (j *Jrpc) QueryBillInfo(in *financetypes.ReqBillInfo, result *interface{}) error {
	//这里直接转发至grpc接口
	reply, err := j.cli.QueryBillInfo(context.Background(), in)
	if err != nil {
		return err
	}
	*result = *reply
	return nil
}

//QueryBorrowerBill query borrower's relate Bills
func (c *channelClient) QueryBorrowerBill(ctx context.Context, in *financetypes.ReqBorrowerBill) (*financetypes.ReplyBorrowerBill, error) {
	msg, err := c.Query(financetypes.FinanceX, "BorrowerBill", in)
	if err != nil {
		return nil, err
	}
	if reply, ok := msg.(*financetypes.ReplyBorrowerBill); ok {
		return reply, nil
	}
	return nil, types.ErrTypeAsset
}
func (j *Jrpc) QueryBorrowerBill(in *financetypes.ReqBorrowerBill, result *interface{}) error {
	//这里直接转发至grpc接口
	reply, err := j.cli.QueryBorrowerBill(context.Background(), in)
	if err != nil {
		return err
	}
	*result = *reply
	return nil
}

//QueryBorrowerBillByStatus query borrower's relate Bills by status
func (c *channelClient) QueryBorrowerBillByStatus(ctx context.Context, in *financetypes.ReqBorrowerBill) (*financetypes.ReplyBorrowerBill, error) {
	msg, err := c.Query(financetypes.FinanceX, "BorrowerBillByStatus", in)
	if err != nil {
		return nil, err
	}
	if reply, ok := msg.(*financetypes.ReplyBorrowerBill); ok {
		return reply, nil
	}
	return nil, types.ErrTypeAsset
}
func (j *Jrpc) QueryBorrowerBillByStatus(in *financetypes.ReqBorrowerBill, result *interface{}) error {
	//这里直接转发至grpc接口
	reply, err := j.cli.QueryBorrowerBillByStatus(context.Background(), in)
	if err != nil {
		return err
	}
	*result = *reply
	return nil
}

//QueryBorrowerBillRBNP query borrower's relate Bills which released but not published.
func (c *channelClient) QueryBorrowerBillRBNP(ctx context.Context, in *financetypes.ReqBorrowerBill) (*financetypes.ReplyBorrowerBill, error) {
	msg, err := c.Query(financetypes.FinanceX, "BorrowerBillRBNP", in)
	if err != nil {
		return nil, err
	}
	if reply, ok := msg.(*financetypes.ReplyBorrowerBill); ok {
		return reply, nil
	}
	return nil, types.ErrTypeAsset
}
func (j *Jrpc) QueryBorrowerBillRBNP(in *financetypes.ReqBorrowerBill, result *interface{}) error {
	//这里直接转发至grpc接口
	reply, err := j.cli.QueryBorrowerBillRBNP(context.Background(), in)
	if err != nil {
		return err
	}
	*result = *reply
	return nil
}

//QuerySplitBillByTxHash
func (c *channelClient) QuerySplitBillByTxHash(ctx context.Context, in *financetypes.ReqSplitBillList) (*financetypes.ReplySplitBillList, error) {
	msg, err := c.Query(financetypes.FinanceX, "SplitBillByTxHash", in)
	if err != nil {
		return nil, err
	}
	if reply, ok := msg.(*financetypes.ReplySplitBillList); ok {
		return reply, nil
	}
	return nil, types.ErrTypeAsset
}
func (j *Jrpc) QuerySplitBillByTxHash(in *financetypes.ReqSplitBillList, result *interface{}) error {
	//这里直接转发至grpc接口
	reply, err := j.cli.QuerySplitBillByTxHash(context.Background(), in)
	if err != nil {
		return err
	}
	*result = *reply
	return nil
}

//QuerySplitBillByTxHash
func (c *channelClient) QuerySplitBillBySplitAddr(ctx context.Context, in *financetypes.ReqSplitBillList) (*financetypes.ReplySplitBillList, error) {
	msg, err := c.Query(financetypes.FinanceX, "SplitBillBySplitAddr", in)
	if err != nil {
		return nil, err
	}
	if reply, ok := msg.(*financetypes.ReplySplitBillList); ok {
		return reply, nil
	}
	return nil, types.ErrTypeAsset
}
func (j *Jrpc) QuerySplitBillBySplitAddr(in *financetypes.ReqSplitBillList, result *interface{}) error {
	//这里直接转发至grpc接口
	reply, err := j.cli.QuerySplitBillBySplitAddr(context.Background(), in)
	if err != nil {
		return err
	}
	*result = *reply
	return nil
}

//QuerySplitBillByTxHash
func (c *channelClient) QuerySplitBillByToAddr(ctx context.Context, in *financetypes.ReqSplitBillList) (*financetypes.ReplySplitBillList, error) {
	msg, err := c.Query(financetypes.FinanceX, "SplitBillByToAddr", in)
	if err != nil {
		return nil, err
	}
	if reply, ok := msg.(*financetypes.ReplySplitBillList); ok {
		return reply, nil
	}
	return nil, types.ErrTypeAsset
}
func (j *Jrpc) QuerySplitBillByToAddr(in *financetypes.ReqSplitBillList, result *interface{}) error {
	//这里直接转发至grpc接口
	reply, err := j.cli.QuerySplitBillByToAddr(context.Background(), in)
	if err != nil {
		return err
	}
	*result = *reply
	return nil
}

//QueryApplyFinancingListByTxHash query applyFinancing list by txHash
func (c *channelClient) QueryApplyFinancingListByTxHash(ctx context.Context, in *financetypes.ReqApplyList) (*financetypes.ReplyApplyList, error) {
	msg, err := c.Query(financetypes.FinanceX, "ApplyFinancingListByTxHash", in)
	if err != nil {
		return nil, err
	}
	if reply, ok := msg.(*financetypes.ReplyApplyList); ok {
		return reply, nil
	}
	return nil, types.ErrTypeAsset
}
func (j *Jrpc) QueryApplyFinancingListByTxHash(in *financetypes.ReqApplyList, result *interface{}) error {
	//这里直接转发至grpc接口
	reply, err := j.cli.QueryApplyFinancingListByTxHash(context.Background(), in)
	if err != nil {
		return err
	}
	*result = *reply
	return nil
}

//QueryApplyFinancingListByCreditor query applyFinancing list by Creditor
func (c *channelClient) QueryApplyFinancingListByCreditor(ctx context.Context, in *financetypes.ReqApplyList) (*financetypes.ReplyApplyList, error) {
	msg, err := c.Query(financetypes.FinanceX, "ApplyFinancingListByCreditor", in)
	if err != nil {
		return nil, err
	}
	if reply, ok := msg.(*financetypes.ReplyApplyList); ok {
		return reply, nil
	}
	return nil, types.ErrTypeAsset
}
func (j *Jrpc) QueryApplyFinancingListByCreditor(in *financetypes.ReqApplyList, result *interface{}) error {
	//这里直接转发至grpc接口
	reply, err := j.cli.QueryApplyFinancingListByCreditor(context.Background(), in)
	if err != nil {
		return err
	}
	*result = *reply
	return nil
}

//QueryApplyFinancingListByApplicant query applyFinancing list by Applicant
func (c *channelClient) QueryApplyFinancingListByApplicant(ctx context.Context, in *financetypes.ReqApplyList) (*financetypes.ReplyApplyList, error) {
	msg, err := c.Query(financetypes.FinanceX, "ApplyFinancingListByApplicant", in)
	if err != nil {
		return nil, err
	}
	if reply, ok := msg.(*financetypes.ReplyApplyList); ok {
		return reply, nil
	}
	return nil, types.ErrTypeAsset
}
func (j *Jrpc) QueryApplyFinancingListByApplicant(in *financetypes.ReqApplyList, result *interface{}) error {
	//这里直接转发至grpc接口
	reply, err := j.cli.QueryApplyFinancingListByApplicant(context.Background(), in)
	if err != nil {
		return err
	}
	*result = *reply
	return nil
}

//QueryApplyFinancingListByBillID query applyFinancing list by Bill-id
func (c *channelClient) QueryApplyFinancingListByBillID(ctx context.Context, in *financetypes.ReqApplyList) (*financetypes.ReplyApplyList, error) {
	msg, err := c.Query(financetypes.FinanceX, "ApplyFinancingListByBillID", in)
	if err != nil {
		return nil, err
	}
	if reply, ok := msg.(*financetypes.ReplyApplyList); ok {
		return reply, nil
	}
	return nil, types.ErrTypeAsset
}
func (j *Jrpc) QueryApplyFinancingListByBillID(in *financetypes.ReqApplyList, result *interface{}) error {
	//这里直接转发至grpc接口
	reply, err := j.cli.QueryApplyFinancingListByBillID(context.Background(), in)
	if err != nil {
		return err
	}
	*result = *reply
	return nil
}

//QueryApplyFinancingListByDpdtSmybol query applyFinancing list by dpdtSymbol
func (c *channelClient) QueryApplyFinancingListByDpdtSmybol(ctx context.Context, in *financetypes.ReqApplyList) (*financetypes.ReplyApplyList, error) {
	msg, err := c.Query(financetypes.FinanceX, "ApplyFinancingListByDpdtSmybol", in)
	if err != nil {
		return nil, err
	}
	if reply, ok := msg.(*financetypes.ReplyApplyList); ok {
		return reply, nil
	}
	return nil, types.ErrTypeAsset
}
func (j *Jrpc) QueryApplyFinancingListByDpdtSmybol(in *financetypes.ReqApplyList, result *interface{}) error {
	//这里直接转发至grpc接口
	reply, err := j.cli.QueryApplyFinancingListByDpdtSmybol(context.Background(), in)
	if err != nil {
		return err
	}
	*result = *reply
	return nil
}

//QueryTokenValue query Token value
func (c *channelClient) QueryTokenValue(ctx context.Context, in *financetypes.ReqTokenValue) (*financetypes.ReplyTokenValue, error) {
	msg, err := c.Query(financetypes.FinanceX, "TokenValue", in)
	if err != nil {
		return nil, err
	}
	if reply, ok := msg.(*financetypes.ReplyTokenValue); ok {
		return reply, nil
	}
	return nil, types.ErrTypeAsset
}
func (j *Jrpc) QueryTokenValue(in *financetypes.ReqTokenValue, result *interface{}) error {
	//这里直接转发至grpc接口
	reply, err := j.cli.QueryTokenValue(context.Background(), in)
	if err != nil {
		return err
	}
	*result = *reply
	return nil
}

//QueryCashList query cash list
func (c *channelClient) QueryCashList(ctx context.Context, in *financetypes.ReqCashList) (*financetypes.ReplyCashList, error) {
	msg, err := c.Query(financetypes.FinanceX, "CashList", in)
	if err != nil {
		return nil, err
	}
	if reply, ok := msg.(*financetypes.ReplyCashList); ok {
		return reply, nil
	}
	return nil, types.ErrTypeAsset
}
func (j *Jrpc) QueryCashList(in *financetypes.ReqCashList, result *interface{}) error {
	//这里直接转发至grpc接口
	reply, err := j.cli.QueryCashList(context.Background(), in)
	if err != nil {
		return err
	}
	*result = *reply
	return nil
}

//QueryCashListByTokenOwner query cash list filter by tokenOwner
func (c *channelClient) QueryCashListByTokenOwner(ctx context.Context, in *financetypes.ReqCashList) (*financetypes.ReplyCashList, error) {
	msg, err := c.Query(financetypes.FinanceX, "CashListByTokenOwner", in)
	if err != nil {
		return nil, err
	}
	if reply, ok := msg.(*financetypes.ReplyCashList); ok {
		return reply, nil
	}
	return nil, types.ErrTypeAsset
}
func (j *Jrpc) QueryCashListByTokenOwner(in *financetypes.ReqCashList, result *interface{}) error {
	//这里直接转发至grpc接口
	reply, err := j.cli.QueryCashListByTokenOwner(context.Background(), in)
	if err != nil {
		return err
	}
	*result = *reply
	return nil
}

//QueryBillCashRecord query Bill cash record
func (c *channelClient) QueryBillCashRecord(ctx context.Context, in *financetypes.ReqCashRecord) (*financetypes.ReplyCashRecord, error) {
	msg, err := c.Query(financetypes.FinanceX, "BillCashRecord", in)
	if err != nil {
		return nil, err
	}
	if reply, ok := msg.(*financetypes.ReplyCashRecord); ok {
		return reply, nil
	}
	return nil, types.ErrTypeAsset
}
func (j *Jrpc) QueryBillCashRecord(in *financetypes.ReqCashRecord, result *interface{}) error {
	//这里直接转发至grpc接口
	reply, err := j.cli.QueryBillCashRecord(context.Background(), in)
	if err != nil {
		return err
	}
	*result = *reply
	return nil
}

//QueryBillRepayInfo query user assets
func (c *channelClient) QueryBillRepayInfo(ctx context.Context, in *financetypes.ReqBillRepayInfo) (*financetypes.ReplyBillRepayInfo, error) {
	msg, err := c.Query(financetypes.FinanceX, "BillRepayInfo", in)
	if err != nil {
		return nil, err
	}
	if reply, ok := msg.(*financetypes.ReplyBillRepayInfo); ok {
		return reply, nil
	}
	return nil, types.ErrTypeAsset
}
func (j *Jrpc) QueryBillRepayInfo(in *financetypes.ReqBillRepayInfo, result *interface{}) error {
	//这里直接转发至grpc接口
	reply, err := j.cli.QueryBillRepayInfo(context.Background(), in)
	if err != nil {
		return err
	}
	*result = *reply
	return nil
}

//QueryRepayInfo query user assets
func (c *channelClient) QueryRepayInfo(ctx context.Context, in *financetypes.ReqBillRepayInfo) (*financetypes.ReplyBillRepayInfo, error) {
	msg, err := c.Query(financetypes.FinanceX, "RepayInfo", in)
	if err != nil {
		return nil, err
	}
	if reply, ok := msg.(*financetypes.ReplyBillRepayInfo); ok {
		return reply, nil
	}
	return nil, types.ErrTypeAsset
}
func (j *Jrpc) QueryRepayInfo(in *financetypes.ReqBillRepayInfo, result *interface{}) error {
	//这里直接转发至grpc接口
	reply, err := j.cli.QueryRepayInfo(context.Background(), in)
	if err != nil {
		return err
	}
	*result = *reply
	return nil
}

//QueryBillRepayRecord query Bill repay record
func (c *channelClient) QueryBillRepayRecord(ctx context.Context, in *financetypes.ReqRepayRecord) (*financetypes.ReplyRepayRecord, error) {
	msg, err := c.Query(financetypes.FinanceX, "BillRepayRecord", in)
	if err != nil {
		return nil, err
	}
	if reply, ok := msg.(*financetypes.ReplyRepayRecord); ok {
		return reply, nil
	}
	return nil, types.ErrTypeAsset
}
func (j *Jrpc) QueryBillRepayRecord(in *financetypes.ReqRepayRecord, result *interface{}) error {
	//这里直接转发至grpc接口
	reply, err := j.cli.QueryBillRepayRecord(context.Background(), in)
	if err != nil {
		return err
	}
	*result = *reply
	return nil
}

//QueryBillBrokenList query broken list
func (c *channelClient) QueryBillBrokenList(ctx context.Context, in *financetypes.ReqBrokenList) (*financetypes.ReplyBrokenList, error) {
	msg, err := c.Query(financetypes.FinanceX, "BillBrokenList", in)
	if err != nil {
		return nil, err
	}
	if reply, ok := msg.(*financetypes.ReplyBrokenList); ok {
		return reply, nil
	}
	return nil, types.ErrTypeAsset
}
func (j *Jrpc) QueryBillBrokenList(in *financetypes.ReqBrokenList, result *interface{}) error {
	//这里直接转发至grpc接口
	reply, err := j.cli.QueryBillBrokenList(context.Background(), in)
	if err != nil {
		return err
	}
	*result = *reply
	return nil
}

//QueryBillBrokenListByPhone query broken list by phone
func (c *channelClient) QueryBillBrokenListByPhone(ctx context.Context, in *financetypes.ReqBrokenList) (*financetypes.ReplyBrokenList, error) {
	msg, err := c.Query(financetypes.FinanceX, "BillBrokenListByPhone", in)
	if err != nil {
		return nil, err
	}
	if reply, ok := msg.(*financetypes.ReplyBrokenList); ok {
		return reply, nil
	}
	return nil, types.ErrTypeAsset
}
func (j *Jrpc) QueryBillBrokenListByPhone(in *financetypes.ReqBrokenList, result *interface{}) error {
	//这里直接转发至grpc接口
	reply, err := j.cli.QueryBillBrokenListByPhone(context.Background(), in)
	if err != nil {
		return err
	}
	*result = *reply
	return nil
}

//QueryBillBrokenListByAddr query broken list by address
func (c *channelClient) QueryBillBrokenListByAddr(ctx context.Context, in *financetypes.ReqBrokenList) (*financetypes.ReplyBrokenList, error) {
	msg, err := c.Query(financetypes.FinanceX, "BillBrokenListByAddr", in)
	if err != nil {
		return nil, err
	}
	if reply, ok := msg.(*financetypes.ReplyBrokenList); ok {
		return reply, nil
	}
	return nil, types.ErrTypeAsset
}
func (j *Jrpc) QueryBillBrokenListByAddr(in *financetypes.ReqBrokenList, result *interface{}) error {
	//这里直接转发至grpc接口
	reply, err := j.cli.QueryBillBrokenListByAddr(context.Background(), in)
	if err != nil {
		return err
	}
	*result = *reply
	return nil
}

//QueryDeliverListByDeliverID query deliver list by txHash
func (c *channelClient) QueryDeliverListByDeliverID(ctx context.Context, in *financetypes.ReqDeliverList) (*financetypes.ReplyDeliverList, error) {
	msg, err := c.Query(financetypes.FinanceX, "DeliverListByDeliverID", in)
	if err != nil {
		return nil, err
	}
	if reply, ok := msg.(*financetypes.ReplyDeliverList); ok {
		return reply, nil
	}
	return nil, types.ErrTypeAsset
}
func (j *Jrpc) QueryDeliverListByDeliverID(in *financetypes.ReqDeliverList, result *interface{}) error {
	//这里直接转发至grpc接口
	reply, err := j.cli.QueryDeliverListByDeliverID(context.Background(), in)
	if err != nil {
		return err
	}
	*result = *reply
	return nil
}

//QueryDeliverListByDeliverAddr query deliver list by deliver address
func (c *channelClient) QueryDeliverListByDeliverAddr(ctx context.Context, in *financetypes.ReqDeliverList) (*financetypes.ReplyDeliverList, error) {
	msg, err := c.Query(financetypes.FinanceX, "DeliverListByDeliverAddr", in)
	if err != nil {
		return nil, err
	}
	if reply, ok := msg.(*financetypes.ReplyDeliverList); ok {
		return reply, nil
	}
	return nil, types.ErrTypeAsset
}
func (j *Jrpc) QueryDeliverListByDeliverAddr(in *financetypes.ReqDeliverList, result *interface{}) error {
	//这里直接转发至grpc接口
	reply, err := j.cli.QueryDeliverListByDeliverAddr(context.Background(), in)
	if err != nil {
		return err
	}
	*result = *reply
	return nil
}

//QueryDeliverListByPeerAddr query deliver list by peer address
func (c *channelClient) QueryDeliverListByPeerAddr(ctx context.Context, in *financetypes.ReqDeliverList) (*financetypes.ReplyDeliverList, error) {
	msg, err := c.Query(financetypes.FinanceX, "DeliverListByPeerAddr", in)
	if err != nil {
		return nil, err
	}
	if reply, ok := msg.(*financetypes.ReplyDeliverList); ok {
		return reply, nil
	}
	return nil, types.ErrTypeAsset
}
func (j *Jrpc) QueryDeliverListByPeerAddr(in *financetypes.ReqDeliverList, result *interface{}) error {
	//这里直接转发至grpc接口
	reply, err := j.cli.QueryDeliverListByPeerAddr(context.Background(), in)
	if err != nil {
		return err
	}
	*result = *reply
	return nil
}
