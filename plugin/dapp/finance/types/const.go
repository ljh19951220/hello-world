package types

//execer name
const (
	//TokenExecer 代币执行器名
	TokenExecer = "token"
	//CoinsExecer 货币执行器名
	CoinsExecer = "paracross"
)

//权值设置
const (
	//E8Weight 1E8
	E8Weight      = 100000000
	OneDaySeconds = 86400
	NegtiveNum    = -1
)

//Support out coin/token symbol
const (
	CCNY = "token.CCNY"
)

//出入金动作类型
const (
	//入金
	DEPOSIT = 1
	//出金
	WITHDRAWAL = 2
)

//单据类型
const (
	//票据
	BILL = iota + 1
	//白条
	BOND
)

const (
	//DefaultOverdueGracePeriod 默认的逾期宽限天数
	DefaultOverdueGracePeriod = 1
	//DefaultOverdueLimit 默认逾期记息最大时长
	DefaultOverdueLimit = 10
)

//IOU status
const (
	//StatusCreated ...
	StatusCreated = iota + 1
	//StatusReleased 只有发行的白条才有价值，token才能兑现
	StatusReleased
	//StatusPublished 处于该种状态，说明相关token已经全部被赎回
	StatusBlocked
	//白条被删除
	StatusDeleted
)

//在借款人眼里，白条的状态, for RepayInfo
const (
	//Holded 白条持有中(已创建待发行)
	IouHolded = iota + 1
	//Fulfilling 正在履行的白条
	IouFulfilling
	//Finished 已完成履行的白条(还款结束)
	IouFinished
	//IouDeleted 被删除的白条
	IouDeleted
)

//地址身份
const (
	//供应商
	Supplier = iota + 1
	//资金方
	Funder
	//核心企业
	Core
)

//兑付地址账户
const (
	//MiddleAddr Address of middle account, current it's static.
	MiddleAddr = "1HrWVgGEuGw6E1BrntqMaCmJUA58ArX7FN"
	//MiddlePrivateKey ...
	MiddlePrivateKey = "0x8dec0ccc4bd7aed795b531457a07a4cad49f665ab93fd7848c6ac3527fa5f9b0"
)
