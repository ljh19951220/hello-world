package executor

import (
	"fmt"
)
/*
 * 用户合约存取kv数据时，key值前缀需要满足一定规范
 * 即key = keyPrefix + userKey
 * 需要字段前缀查询时，使用’-‘作为分割符号
 */

var (
	//KeyPrefixStateDB state db key必须前缀
	KeyPrefixStateDB = "mavl-finance-"
	//keyPrefixBrokenList prefix of credit
	keyPrefixCredit = "mavl-finance-credit-"
	//keyPrefixBrokenList prefix of Bill
	keyPrefixBill = "mavl-finance-bill-"

	//KeyPrefixLocalDB local db的key必须前缀
	KeyPrefixLocalDB = "LODB-finance"
	//TokenAssetsPrefix prefix of account token assets
	TokenAssetsPrefix = "LODB-finance-assets:"
)

//授信前缀
func calcCreditNewKey(symbol string) (key []byte) {
	return []byte(fmt.Sprintf(keyPrefixCredit+"%s", symbol))
}

//用户Token资产，记录到LocalDB中
func calcTokenAssetsKey(addr string) []byte {
	return []byte(fmt.Sprintf(TokenAssetsPrefix+"%s", addr))
}

func calcBillNewKey(id string) (key []byte) {
	return []byte(fmt.Sprintf(keyPrefixBill+"%s", id))
}
