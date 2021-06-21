package types

import "errors"

var (
	//ErrAssetsNotEnough error no enough assert
	ErrAssetsNotEnough = errors.New("ErrAssetsNotEnough")
	//ErrOperateAddrMismatch error operation address mismatch
	ErrOperateAddrMismatch = errors.New("ErrOperateAddrMismatch")
	//ErrOperateAddrIsNotCore error operation address is not Core
	ErrOperateAddrIsNotCore = errors.New("ErrOperateAddrIsNotCore")
	//ErrOperateAddrIsNotFunder error operation address is not Funder
	ErrOperateAddrIsNotFunder = errors.New("ErrOperateAddrIsNotFunder")
	//ErrOperateAddrIsNotSupplier error operation address is not Supplier
	ErrOperateAddrIsNotSupplier = errors.New("ErrOperateAddrIsNotSupplier")
	//ErrBillNotFound error Bill not found in DB
	ErrBillNotFound = errors.New("ErrBillNotFound")
	//ErrCreditNotFound error credit token not found
	ErrCreditNotFound = errors.New("ErrCreditTokenNotFound")
	//ErrBillStatus error Bill status
	ErrBillStatus = errors.New("ErrBillStatus")
	//ErrNoPermissionAction error have no permission to do something
	ErrNoPermissionAction = errors.New("ErrNoPermissionAction")
	//ErrUnReleaseIou error to UnRelease IOU
	ErrUnReleaseIou = errors.New("ErrUnReleaseIou")
	//ErrIOUMarketNotFound error IOU market not found in DB
	ErrIOUMarketNotFound = errors.New("ErrIOUMarketNotFound")
	//ErrApplyFinancingNotFound error applyFinancing not found in DB
	ErrApplyFinancingNotFound = errors.New("ErrApplyFinancingNotFound")
	//ErrInvalidRepayMethod error invalid repay method
	ErrInvalidRepayMethod = errors.New("ErrInvalidRepayMethod")
	//ErrInvalidRepayCount error invalid repay count
	ErrInvalidRepayCount = errors.New("ErrInvalidRepayCount")
	//ErrInvalidParameter error invalid parameter
	ErrInvalidParameter = errors.New("ErrInvalidParameter")
	//ErrDepositAsset error deposit ssset
	ErrDepositAsset = errors.New("ErrDepositAsset")
	//ErrNeedNotRepay error need't repay
	ErrNeedNotRepay = errors.New("ErrNeedNotRepay")
	//ErrTotalCashAmountMismatch error cash amount mismatch
	ErrTotalCashAmountMismatch = errors.New("ErrTotalCashAmountMismatch")
	//ErrCashHashNotFound error found cash item from cash list
	ErrCashHashNotFound = errors.New("ErrCashHashNotFound")
	//ErrNoGuarantees error no guarantees
	ErrNoGuarantees = errors.New("ErrNoGuarantees")
	//ErrNegativeNum error negative number
	ErrNegativeNum = errors.New("ErrNegativeNum")
	//ErrSystemErr system error
	ErrSystemErr = errors.New("ErrSystemErr")
	//ErrExchangeLimit exchange limit error
	ErrExchangeLimit = errors.New("ErrExchangeLimit")
	//ErrCreditTokenDismatch credit token dismatch error
	ErrCreditTokenDismatch = errors.New("ErrCreditTokenDismatch")
	//ErrUnusedCreditTNotEnough unused credit token not enough
	ErrUnusedCreditTNotEnough = errors.New("ErrUnusedCreditTNotEnough")
	//ErrCreditTokenExpired credit token expired
	ErrCreditTokenExpired = errors.New("ErrCreditTokenExpired")
	//ErrBillHasExpired iou has expired
	ErrBillHasExpired = errors.New("ErrBillHasExpired")
	//ErrWrongTokenValue wrong token value
	ErrWrongTokenValue = errors.New("ErrWrongTokenValue")
	//ErrDeliverTimeout deliver iou timeout
	ErrDeliverTimeout = errors.New("ErrDeliverTimeout")
	//ErrDeliverItemNotFound error deliver item not found in DB
	ErrDeliverItemNotFound = errors.New("ErrDeliverItemNotFound")
	//ErrBIllNotSplit error Bill can not Split
	ErrBillCanNotSplit = errors.New("ErrBillNotSplit")
	//ErrSplitTimeout split iou timeout
	ErrSplitTimeout = errors.New("ErrSplitTimeout")
	//ErrSplitItemNotFound error Split item not found in DB
	ErrSplitItemNotFound = errors.New("ErrSplitItemNotFound")
)
