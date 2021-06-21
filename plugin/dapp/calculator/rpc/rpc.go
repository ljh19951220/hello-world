package rpc

import (
	"context"

	"github.com/33cn/chain33/types"
	calculatortypes "github.com/33cn/plugin/plugin/dapp/calculator/types"
)

/*
 * 实现json rpc和grpc service接口
 * json rpc用Jrpc结构作为接收实例
 * grpc使用channelClient结构作为接收实例
 */
func (c *channelClient) QueryCalcCount(ctx context.Context, in *calculatortypes.ReqQueryCalcCount) (*calculatortypes.ReplyQueryCalcCount, error) {

	msg, err := c.Query(calculatortypes.CalculatorX, "CalcCount", in)
	if err != nil {
		return nil, err
	}
	if reply, ok := msg.(*calculatortypes.ReplyQueryCalcCount); ok {
		return reply, nil
	}
	return nil, types.ErrTypeAsset
}
func (j *Jrpc) QueryCalcCount(in *calculatortypes.ReqQueryCalcCount, result *interface{}) error {

	//这里直接调用已实现的grpc接口
	reply, err := j.cli.QueryCalcCount(context.Background(), in)
	if err != nil {
		return err
	}
	*result = *reply
	return nil
}
//func (c *channelClient) Add(ctx context.Context, in *calculatortypes.Add) (*types.UnsignTx, error) {
//	acct := &calculatortypes.CalculatorAction{
//		Ty:    calculatortypes.TyAddAction,
//		Value: calculatortypes.CalculatorAction_Add{Add:in},
//	}
//
//	cfg := c.GetConfig()
//	tx, err := types.CreateFormatTx(cfg, cfg.ExecName(loantypes.LoanX), types.Encode(acct))
//	if err != nil {
//		return nil, err
//	}
//
//	data := types.Encode(tx)
//	return &types.UnsignTx{Data: data}, nil
//}
//func (j *Jrpc) Add(in *calculatortypes.ReqQueryCalcCount, result *interface{}) error {
//
//	//这里直接调用已实现的grpc接口
//	reply, err := j.cli.QueryCalcCount(context.Background(), in)
//	if err != nil {
//		return err
//	}
//	*result = *reply
//	return nil
//}
