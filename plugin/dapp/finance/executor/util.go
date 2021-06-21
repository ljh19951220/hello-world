package executor

import (
	"errors"
	pty "github.com/33cn/plugin/plugin/dapp/finance/types"
	"math"
	"strings"
)

//上链贷的首字母缩写作为BillID的分隔符
const iouIdDelimiter = "XYD"

func tokenSymbolToBillId(tokenSymbol string) string {
	subStr := strings.Split(tokenSymbol, iouIdDelimiter)
	return subStr[len(subStr)-1]
}

//calcOriginPriceForBill 获取票据的初始单价
func calcOriginPriceForBill(rate float64, circulationDays int64) int64 {
	if rate < 0 {
		return 0
	}
	return int64(1e8 / math.Pow(1+rate, float64(circulationDays)))
}

//BillRepayMethod Bill repay method
type BillRepayMethod struct {
}

//CalcIouRepayInfo 计算当前票据的还款情况, 最主要的是两个参数：是否逾期，应还款项
func (m *BillRepayMethod) CalcIouRepayInfo(bill *pty.Bill, timeNow int64, amount int64) *pty.ReplyBillRepayInfo {
	var reply pty.ReplyBillRepayInfo
	if bill.Status != pty.StatusReleased {
		financelog.Debug("Bill not be released", "Bill status", bill.Status)
		return nil
	}
	if bill.CurrentRepayedCount == bill.RepayCount {
		financelog.Debug("The borrower has repaid all the arrears", "Repayed Count", bill.CurrentRepayedCount)
		return nil
	}

	//对于票据来说只有一期
	reply.OverdueRate = bill.OverdueRate
	reply.CurrentCount = 1
	reply.TotalRepayCount = 1
	//RepayDueDate当前期的还款到期日
	reply.RepayDueDate = bill.IssueDate + bill.CirculationTime
	if timeNow > reply.RepayDueDate {
		//是否包含逾期还款
		reply.IncludeOverdueRepay = true
	} else {
		reply.IncludeOverdueRepay = false
	}

	//calculate ActualRepayAmount
	if reply.IncludeOverdueRepay == true {
		//逾期天数
		overDays := (timeNow - reply.RepayDueDate) / pty.OneDaySeconds
		//逾期利率
		overdueRate := float64(bill.OverdueRate) / 1e8
		base := amount //本金
		if overDays < bill.OverdueGracePeriod {
			//逾期宽限期之内还款都不计算逾期利息
			reply.ActualRepayAmount = amount
			//reply.ActualRepayAmount = int64((1 + overdueRate) * float64(base))
		} else {
			//超过逾期宽限日
			var caloverDays int64
			if overDays > bill.OverdueLimit {
				//逾期超过最大逾期记息时长, unit:days，就按照最大计息时长算
				//Make sure have a limitation to overdue date
				caloverDays = bill.OverdueLimit
			} else {
				caloverDays = overDays
			}
			reply.ActualRepayAmount = int64(math.Pow(1+overdueRate, float64(caloverDays)) * float64(base))
		}
		reply.NormalRepayAmount = amount
	} else {
		days := (timeNow - bill.IssueDate) / pty.OneDaySeconds
		//剩余天数
		remainDays := bill.CirculationTime/pty.OneDaySeconds - days
		//这个地方应该是：个人感觉
		//value := bill.LoanAmount
		//根据剩余天数计算白条币的价值,剩余天数为0的时候为loanAmount还款
		value, err := calcBillValue(amount, remainDays, bill.Rate)
		if err != nil {
			financelog.Error("calcBillValue failed", "error", err)
		}
		reply.NormalRepayAmount = value
		reply.ActualRepayAmount = value
	}

	return &reply
}

//GetTokenValue 获得token的价值数, 可兑现的CCNY数额; 对于票据而言，兑换价值是以最大面额作为可兑换的最大值
func (m *BillRepayMethod) GetTokenValue(tokenAmount int64, bill *pty.Bill, timeNow int64) (int64, error) {
	//逾期则按照最大面值算，未逾期则按照实际面值算
	var days int64
	if timeNow > bill.IssueDate+bill.CirculationTime {
		days = bill.CirculationTime / pty.OneDaySeconds
	} else {
		days = (timeNow - bill.IssueDate) / pty.OneDaySeconds
	}
	remainDays := bill.CirculationTime/pty.OneDaySeconds - days
	value, err := calcBillValue(tokenAmount, remainDays, bill.Rate)

	return value, err
}

func calcBillValue(faceValue int64, remainDays int64, billDayRate int64) (int64, error) {
	rate := float64(billDayRate) / 1e8
	if rate < 0 {
		err := errors.New("illegal bill rate")
		return 0, err
	}

	return int64(float64(faceValue) / math.Pow(1+rate, float64(remainDays))), nil
}

//ACAPRepayMethod ACAP repay method
type ACAPRepayMethod struct {
}

//CalcIouRepayInfo 计算当前白条的还款情况, 最主要的是两个参数：是否逾期，应还款项
func (m *ACAPRepayMethod) CalcIouRepayInfo(bill *pty.Bill, timeNow int64, amount int64) *pty.ReplyBillRepayInfo {
	var reply pty.ReplyBillRepayInfo
	if bill.Status != pty.StatusReleased {
		financelog.Debug("Iou not be released", "Iou status", bill.Status)
		return nil
	}
	if bill.CurrentRepayedCount == bill.RepayCount {
		financelog.Debug("The borrower has repaid all the arrears", "Repayed Count", bill.CurrentRepayedCount)
		return nil
	}
	period := bill.CirculationTime / int64(bill.RepayCount)
	reply.CurrentCount = int32((timeNow-bill.IssueDate)/period) + 1
	reply.RepayDueDate = bill.IssueDate + period*int64(reply.CurrentCount)
	if reply.CurrentCount > (bill.CurrentRepayedCount + 1) {
		reply.IncludeOverdueRepay = true
	}
	repayItem := getRepayPlanItem(bill.RepayItems, reply.CurrentCount)

	if bill.CurrentRepayedCount == reply.CurrentCount {
		//当前期已经完成还款，所以应还额度为0
		financelog.Debug("The borrower has repaid the arrears for current phase", "Repayed Count", bill.CurrentRepayedCount)
		reply.ActualRepayAmount = 0
	} else {
		var sn float64
		A := float64(repayItem.RepayAmount)
		r := float64(bill.OverdueRate) / 1e8
		t := float64(period / pty.OneDaySeconds)
		t1 := float64((period - (reply.RepayDueDate - timeNow)) / (3600 * 24))
		k := float64(reply.CurrentCount - bill.CurrentRepayedCount)
		//利用等比数列求和公式计算逾期的本息和; 当K<2时，sn等于0;
		if k >= 2 {
			if r != 0 {
				sn = math.Pow(A, 2) * math.Pow(1+r, t1) * (math.Pow(1+r, (k-1)*t) - math.Pow(1+r, t)) / (math.Pow(1+r, t) - 1)
			} else {
				sn = (k - 1) * A
			}
		}
		//实际应还款项额度
		reply.ActualRepayAmount = int64(sn) + repayItem.RepayAmount
	}

	//白条总共需还款的期数
	reply.TotalRepayCount = bill.RepayCount
	//正常情况下，当期应还款额度
	reply.NormalRepayAmount = repayItem.RepayAmount
	reply.OverdueRate = bill.OverdueRate
	return &reply
}
func getRepayPlanItem(r []*pty.RepayPlan, repayPhase int32) *pty.RepayPlan {
	if len(r) == 0 {
		financelog.Error("Repay plan item not exist", "repay phase", repayPhase)
		return nil
	}

	for _, v := range r {
		if v.RepayPhase == repayPhase {
			return v
		}
	}

	return nil
}

//calcTokenValueForFns 计算贴现之后的代币价值 rate:贴现率
func calcTokenValueForFinancing(amount int64, rate int64, days int64) int64 {
	if rate < 0 || amount <= 0 {
		financelog.Error("calcTokenValueForFinancing: wrong parameter", "rate", rate, "amount", amount)
		return 0
	}

	value := int64(float64(amount) * (1 - float64(rate*days)/1e8))
	if value < 0 {
		financelog.Error("calcTokenValueForFinancing: nagetive token value", "rate", rate, "amount", amount, "days", days)
		return 0
	}

	return value
}
