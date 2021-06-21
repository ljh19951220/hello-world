package executor

import (
	"fmt"
	"github.com/33cn/chain33/common/db"
	"github.com/33cn/chain33/common/db/table"
	"github.com/33cn/chain33/types"
	financetypes "github.com/33cn/plugin/plugin/dapp/finance/types"
)

/***************************credit token info table*****************************/
var optCreditRecordLocalDB = &table.Option{
	Prefix:  KeyPrefixLocalDB,
	Name:    "CreditRecord",
	Primary: "symbol",
	Index:   []string{"creditAddr", "granteeAddr", "txHash"},
}

//CreditRecordRow table meta 结构
type CreditRecordRow struct {
	*financetypes.CreditRecordLocalDB
}

//NewCreditRecordRow 新建一个meta 结构
func NewCreditRecordRow() *CreditRecordRow {
	return &CreditRecordRow{CreditRecordLocalDB: &financetypes.CreditRecordLocalDB{}}
}

//CreateRow 新建数据行
func (tx *CreditRecordRow) CreateRow() *table.Row {
	return &table.Row{Data: &financetypes.CreditRecordLocalDB{}}
}

//SetPayload 设置数据
func (tx *CreditRecordRow) SetPayload(data types.Message) error {
	if txdata, ok := data.(*financetypes.CreditRecordLocalDB); ok {
		tx.CreditRecordLocalDB = txdata
		return nil
	}
	return types.ErrTypeAsset
}

//Get 按照indexName 查询 indexValue
func (tx *CreditRecordRow) Get(key string) ([]byte, error) {
	if key == "symbol" {
		return []byte(tx.Symbol), nil
	} else if key == "creditAddr" {
		return []byte(tx.CreditAddr), nil
	} else if key == "granteeAddr" {
		return []byte(tx.GranteeAddr), nil
	} else if key == "txHash" {
		return []byte(tx.TxHash), nil
	}

	return nil, types.ErrNotFound
}

//参考 NewCreditRecordTable 新建表
func NewCreditRecordTable(kvdb db.KV) *table.Table {
	rowmeta := NewCreditRecordRow()
	table, err := table.NewTable(rowmeta, kvdb, optCreditRecordLocalDB)
	if err != nil {
		panic(err)
	}
	return table
}

/***************************Deposit and withdrawal record table***************************/
//Table for store Deposit and withdrawal record
func getDWRecordOpt(addr string) *table.Option {
	name := fmt.Sprintf("DWRecord_%s", addr)
	return &table.Option{Prefix: KeyPrefixLocalDB, Name: name, Primary: "txHash", Index: []string{"symbol", "addr"}}
}

//DWRecordRow table meta 结构
type DWRecordRow struct {
	*financetypes.DWRecordLocalDB
}

//NewDWRecordRow 新建一个meta 结构
func NewDWRecordRow() *DWRecordRow {
	return &DWRecordRow{DWRecordLocalDB: &financetypes.DWRecordLocalDB{}}
}

//CreateRow 新建数据行
func (tx *DWRecordRow) CreateRow() *table.Row {
	return &table.Row{Data: &financetypes.DWRecordLocalDB{}}
}

//SetPayload 设置数据
func (tx *DWRecordRow) SetPayload(data types.Message) error {
	if txdata, ok := data.(*financetypes.DWRecordLocalDB); ok {
		tx.DWRecordLocalDB = txdata
		return nil
	}
	return types.ErrTypeAsset
}

func (tx *DWRecordRow) Get(key string) ([]byte, error) {
	if key == "txHash" {
		return []byte(tx.TxHash), nil
	} else if key == "symbol" {
		return []byte(tx.Symbol), nil
	} else if key == "addr" {
		return []byte(tx.Addr), nil
	}
	return nil, types.ErrNotFound
}

/***************************TransferCoins record table***************************/
//Table for storeTransferCoins record
func getTransferCoinsRecordOpt(addr string) *table.Option {
	name := fmt.Sprintf("TransferCoinsRecord_%s", addr)
	return &table.Option{Prefix: KeyPrefixLocalDB, Name: name, Primary: "txHash", Index: []string{"symbol", "fromAddr"}}
}

//TransferCoinsRecordOpRow table meta 结构
type TransferCoinsRecordRow struct {
	*financetypes.TransferCoinsRecordLocalDB
}

//NewTransferCoinsRecordOp 新建一个meta 结构
func NewTransferCoinsRecordRow() *TransferCoinsRecordRow {
	return &TransferCoinsRecordRow{TransferCoinsRecordLocalDB: &financetypes.TransferCoinsRecordLocalDB{}}
}

//CreateRow 新建数据行
func (tx *TransferCoinsRecordRow) CreateRow() *table.Row {
	return &table.Row{Data: &financetypes.TransferCoinsRecordLocalDB{}}
}

//SetPayload 设置数据
func (tx *TransferCoinsRecordRow) SetPayload(data types.Message) error {
	if txdata, ok := data.(*financetypes.TransferCoinsRecordLocalDB); ok {
		tx.TransferCoinsRecordLocalDB = txdata
		return nil
	}
	return types.ErrTypeAsset
}

func (tx *TransferCoinsRecordRow) Get(key string) ([]byte, error) {
	if key == "txHash" {
		return []byte(tx.TxHash), nil
	} else if key == "symbol" {
		return []byte(tx.Symbol), nil
	} else if key == "fromAddr" {
		return []byte(tx.FromAddr), nil
	}
	return nil, types.ErrNotFound
}

/********************************************************************************/

/***************************Borrower's Bill info table***************************/
//Table for store Borrower's Iou info
func getBorrowerBillOpt(gAddr string) *table.Option {
	name := fmt.Sprintf("BorrowerBill_%s", gAddr)
	return &table.Option{Prefix: KeyPrefixLocalDB, Name: name, Primary: "ID", Index: []string{"status"}}
}

//BorrowerIouRow table meta 结构
type BorrowerBillRow struct {
	*financetypes.BorrowerBillLocalDB
}

//NewBorrowerIouRow 新建一个meta 结构
func NewBorrowerBillRow() *BorrowerBillRow {
	return &BorrowerBillRow{BorrowerBillLocalDB: &financetypes.BorrowerBillLocalDB{}}
}

//CreateRow 新建数据行
func (tx *BorrowerBillRow) CreateRow() *table.Row {
	return &table.Row{Data: &financetypes.BorrowerBillLocalDB{}}
}

//SetPayload 设置数据
func (tx *BorrowerBillRow) SetPayload(data types.Message) error {
	if txdata, ok := data.(*financetypes.BorrowerBillLocalDB); ok {
		tx.BorrowerBillLocalDB = txdata
		return nil
	}
	return types.ErrTypeAsset
}

func (tx *BorrowerBillRow) Get(key string) ([]byte, error) {
	if key == "ID" {
		return []byte(tx.Id), nil
	} else if key == "status" {
		return []byte(fmt.Sprintf("%d", tx.Status)), nil
	}
	return nil, types.ErrNotFound
}

/********************************************************************************/

/***************************DeliverIou info table*****************************/
var optDeliverBillList = &table.Option{
	Prefix:  KeyPrefixLocalDB,
	Name:    "DeliverBillList",
	Primary: "deliverID",
	Index:   []string{"deliverAddr", "toAddr"},
}

//DeliverBillListRow table meta 结构
type DeliverBillListRow struct {
	*financetypes.DeliverItemLocalDB
}

//NewDeliverIouListTable 新建表
func NewDeliverIouListTable(kvdb db.KV) *table.Table {
	rowmeta := NewDeliverBillListRow()
	table, err := table.NewTable(rowmeta, kvdb, optDeliverBillList)
	if err != nil {
		panic(err)
	}
	return table
}

//NewDeliverIouListRow 新建一个meta 结构
func NewDeliverBillListRow() *DeliverBillListRow {
	return &DeliverBillListRow{DeliverItemLocalDB: &financetypes.DeliverItemLocalDB{}}
}

//CreateRow 新建数据行
func (tx *DeliverBillListRow) CreateRow() *table.Row {
	return &table.Row{Data: &financetypes.DeliverItemLocalDB{}}
}

//SetPayload 设置数据
func (tx *DeliverBillListRow) SetPayload(data types.Message) error {
	if txdata, ok := data.(*financetypes.DeliverItemLocalDB); ok {
		tx.DeliverItemLocalDB = txdata
		return nil
	}
	return types.ErrTypeAsset
}

//Get 按照indexName 查询 indexValue
func (tx *DeliverBillListRow) Get(key string) ([]byte, error) {
	if key == "deliverID" {
		return []byte(tx.DeliverID), nil
	} else if key == "deliverAddr" {
		return []byte(tx.DeliverAddr), nil
	} else if key == "toAddr" {
		return []byte(tx.ToAddr), nil
	}

	return nil, types.ErrNotFound
}

/********************************************************************************/

/***************************SplitBill info table*****************************/
var optSplitBillList = &table.Option{
	Prefix:  KeyPrefixLocalDB,
	Name:    "SplitBillList",
	Primary: "txHash",
	Index:   []string{"splitAddr", "toAddr"},
}

//DeliverBillListRow table meta 结构
type SplitBillListRow struct {
	*financetypes.SplitBillRecordLocalDB
}

//NewDeliverIouListTable 新建表
func NewSplitBillListTable(kvdb db.KV) *table.Table {
	rowmeta := NewSplitBillListRow()
	table, err := table.NewTable(rowmeta, kvdb, optSplitBillList)
	if err != nil {
		panic(err)
	}
	return table
}

//NewDeliverIouListRow 新建一个meta 结构
func NewSplitBillListRow() *SplitBillListRow {
	return &SplitBillListRow{SplitBillRecordLocalDB: &financetypes.SplitBillRecordLocalDB{}}
}

//CreateRow 新建数据行
func (tx *SplitBillListRow) CreateRow() *table.Row {
	return &table.Row{Data: &financetypes.SplitBillRecordLocalDB{}}
}

//SetPayload 设置数据
func (tx *SplitBillListRow) SetPayload(data types.Message) error {
	if txdata, ok := data.(*financetypes.SplitBillRecordLocalDB); ok {
		tx.SplitBillRecordLocalDB = txdata
		return nil
	}
	return types.ErrTypeAsset
}

//Get 按照indexName 查询 indexValue
func (tx *SplitBillListRow) Get(key string) ([]byte, error) {
	if key == "txHash" {
		return []byte(tx.TxHash), nil
	} else if key == "splitAddr" {
		return []byte(tx.SplitAddr), nil
	} else if key == "toAddr" {
		return []byte(tx.ToAddr), nil
	}

	return nil, types.ErrNotFound
}

/********************************************************************************/

/***************************ApplyForFinancing info table*****************************/
var optApplyForFinancingList = &table.Option{
	Prefix:  KeyPrefixLocalDB,
	Name:    "ApplyForFinancingList",
	Primary: "txHash",
	Index:   []string{"creditorAddr", "dpdtSymbol", "applyAddr", "id"},
}

//ApplyFnsListRow table meta 结构
type ApplyForFinancingListRow struct {
	*financetypes.ApplyFinancingItemLocalDB
}

//NewApplyForFinancingListTable 新建表
func NewApplyForFinancingListTable(kvdb db.KV) *table.Table {
	rowmeta := NewApplyForFinancingListRow()
	table, err := table.NewTable(rowmeta, kvdb, optApplyForFinancingList)
	if err != nil {
		panic(err)
	}
	return table
}

//NewApplyFnsListRow 新建一个meta 结构
func NewApplyForFinancingListRow() *ApplyForFinancingListRow {
	return &ApplyForFinancingListRow{ApplyFinancingItemLocalDB: &financetypes.ApplyFinancingItemLocalDB{}}
}

//CreateRow 新建数据行
func (tx *ApplyForFinancingListRow) CreateRow() *table.Row {
	return &table.Row{Data: &financetypes.ApplyFinancingItemLocalDB{}}
}

//SetPayload 设置数据
func (tx *ApplyForFinancingListRow) SetPayload(data types.Message) error {
	if txdata, ok := data.(*financetypes.ApplyFinancingItemLocalDB); ok {
		tx.ApplyFinancingItemLocalDB = txdata
		return nil
	}
	return types.ErrTypeAsset
}

//Get 按照indexName 查询 indexValue
func (tx *ApplyForFinancingListRow) Get(key string) ([]byte, error) {
	if key == "txHash" {
		return []byte(tx.TxHash), nil
	} else if key == "creditorAddr" {
		return []byte(tx.CreditorAddr), nil
	} else if key == "dpdtSymbol" {
		return []byte(tx.DpdtSymbol), nil
	} else if key == "applyAddr" {
		return []byte(tx.ApplyAddr), nil
	} else if key == "id" {
		return []byte(tx.Id), nil
	}

	return nil, types.ErrNotFound
}

/********************************************************************************/

/***************************CashRecordLocalDB info table*****************************/
var optCashRecordLocalDB = &table.Option{
	Prefix:  KeyPrefixLocalDB,
	Name:    "CashRecordLocalDB",
	Primary: "txHash",
	Index:   []string{"cashAddr", "billId", "timestamp"},
}

//CashRecordRow table meta 结构
type CashRecordLocalDBRow struct {
	*financetypes.CashRecordLocalDB
}

//NewCashRecordTable 新建表
func NewCashRecordTable(kvdb db.KV) *table.Table {
	rowmeta := NewCashRecordLocalDBRow()
	table, err := table.NewTable(rowmeta, kvdb, optCashRecordLocalDB)
	if err != nil {
		panic(err)
	}
	return table
}

//NewCashRecordRow 新建一个meta 结构
func NewCashRecordLocalDBRow() *CashRecordLocalDBRow {
	return &CashRecordLocalDBRow{CashRecordLocalDB: &financetypes.CashRecordLocalDB{}}
}

//CreateRow 新建数据行
func (tx *CashRecordLocalDBRow) CreateRow() *table.Row {
	return &table.Row{Data: &financetypes.CashRecordLocalDB{}}
}

//SetPayload 设置数据
func (tx *CashRecordLocalDBRow) SetPayload(data types.Message) error {
	if txdata, ok := data.(*financetypes.CashRecordLocalDB); ok {
		tx.CashRecordLocalDB = txdata
		return nil
	}
	return types.ErrTypeAsset
}

//Get 按照indexName 查询 indexValue
func (tx *CashRecordLocalDBRow) Get(key string) ([]byte, error) {
	if key == "txHash" {
		return []byte(tx.TxHash), nil
	} else if key == "cashAddr" {
		return []byte(tx.CashAddr), nil
	} else if key == "timestamp" {
		return []byte(fmt.Sprintf("%d", tx.Timestamp)), nil
	} else if key == "billId" {
		return []byte(tx.BillId), nil
	}

	return nil, types.ErrNotFound
}

/********************************************************************************/

/***************************Cash List table***************************/
//Table for store Cash list info
func getCashListOpt(billID string) *table.Option {
	name := fmt.Sprintf("CashList_%s", billID)
	return &table.Option{Prefix: KeyPrefixLocalDB, Name: name, Primary: "txHash", Index: []string{"tokenOwner"}}
}

//CashListRow table meta 结构
type CashListRow struct {
	*financetypes.CashList
}

//NewCashListTable 新建表
func NewCashListTable(kvdb db.KV, iouID string) *table.Table {
	rowmeta := NewCashListRow()
	table, err := table.NewTable(rowmeta, kvdb, getCashListOpt(iouID))
	if err != nil {
		panic(err)
	}
	return table
}

//NewCashListRow 新建一个meta 结构
func NewCashListRow() *CashListRow {
	return &CashListRow{CashList: &financetypes.CashList{}}
}

//CreateRow 新建数据行
func (tx *CashListRow) CreateRow() *table.Row {
	return &table.Row{Data: &financetypes.CashList{}}
}

//SetPayload 设置数据
func (tx *CashListRow) SetPayload(data types.Message) error {
	if txdata, ok := data.(*financetypes.CashList); ok {
		tx.CashList = txdata
		return nil
	}
	return types.ErrTypeAsset
}

func (tx *CashListRow) Get(key string) ([]byte, error) {
	if key == "txHash" {
		return []byte(tx.TxHash), nil
	} else if key == "tokenOwner" {
		return []byte(tx.TokenOwner), nil
	}
	return nil, types.ErrNotFound
}

/********************************************************************************/

/***************************RepayRecordLocalDB table***************************/
var optRepayRecordLocalDB = &table.Option{
	Prefix:  KeyPrefixLocalDB,
	Name:    "RepayRecordLocalDB",
	Primary: "txHash",
	Index:   []string{"repayAddr", "billId", "timestamp"},
}

//RepayRecordRow table meta 结构
type RepayRecordLocalDBRow struct {
	*financetypes.RepayRecordLocalDB
}

//NewRepayRecordTable 新建表
func NewRepayRecordLocalDBTable(kvdb db.KV) *table.Table {
	rowmeta := NewRepayRecordLocalDBRow()
	table, err := table.NewTable(rowmeta, kvdb, optRepayRecordLocalDB)
	if err != nil {
		panic(err)
	}
	return table
}

//NewRepayRecordRow 新建一个meta 结构
func NewRepayRecordLocalDBRow() *RepayRecordLocalDBRow {
	return &RepayRecordLocalDBRow{RepayRecordLocalDB: &financetypes.RepayRecordLocalDB{}}
}

//CreateRow 新建数据行
func (tx *RepayRecordLocalDBRow) CreateRow() *table.Row {
	return &table.Row{Data: &financetypes.RepayRecordLocalDB{}}
}

//SetPayload 设置数据
func (tx *RepayRecordLocalDBRow) SetPayload(data types.Message) error {
	if txdata, ok := data.(*financetypes.RepayRecordLocalDB); ok {
		tx.RepayRecordLocalDB = txdata
		return nil
	}
	return types.ErrTypeAsset
}

//Get 按照indexName 查询 indexValue
func (tx *RepayRecordLocalDBRow) Get(key string) ([]byte, error) {
	if key == "txHash" {
		return []byte(tx.TxHash), nil
	} else if key == "repayAddr" {
		return []byte(tx.RepayAddr), nil
	} else if key == "timestamp" {
		return []byte(fmt.Sprintf("%d", tx.Timestamp)), nil
	} else if key == "billId" {
		return []byte(tx.BillId), nil
	}

	return nil, types.ErrNotFound
}

/********************************************************************************/

/***************************RepayRecordLocalDB table***************************/
var optBrokenRecord = &table.Option{
	Prefix:  KeyPrefixLocalDB,
	Name:    "BrokenRecord",
	Primary: "billId",
	Index:   []string{"phone", "borrowerAddr", "identifier", "name"},
}

//BrokenRecordRow table meta 结构
type BrokenRecordStateDBRow struct {
	*financetypes.BrokenRecordStateDB
}

//NewBrokenRecordTable 新建表
func NewBrokenRecordStateDBTable(kvdb db.KV) *table.Table {
	rowmeta := NewBrokenRecordStateDBRow()
	table, err := table.NewTable(rowmeta, kvdb, optBrokenRecord)
	if err != nil {
		panic(err)
	}
	return table
}

//NewBrokenRecordRow 新建一个meta 结构
func NewBrokenRecordStateDBRow() *BrokenRecordStateDBRow {
	return &BrokenRecordStateDBRow{BrokenRecordStateDB: &financetypes.BrokenRecordStateDB{}}
}

//CreateRow 新建数据行
func (tx *BrokenRecordStateDBRow) CreateRow() *table.Row {
	return &table.Row{Data: &financetypes.BrokenRecordStateDB{}}
}

//SetPayload 设置数据
func (tx *BrokenRecordStateDBRow) SetPayload(data types.Message) error {
	if txdata, ok := data.(*financetypes.BrokenRecordStateDB); ok {
		tx.BrokenRecordStateDB = txdata
		return nil
	}
	return types.ErrTypeAsset
}

//Get 按照indexName 查询 indexValue
func (tx *BrokenRecordStateDBRow) Get(key string) ([]byte, error) {
	if key == "billId" {
		return []byte(tx.BillId), nil
	} else if key == "borrowerAddr" {
		return []byte(tx.BorrowerAddr), nil
	} else if key == "phone" {
		return []byte(tx.Phone), nil
	} else if key == "identifier" {
		return []byte(tx.Identifier), nil
	} else if key == "name" {
		return []byte(tx.Name), nil
	}

	return nil, types.ErrNotFound
}
