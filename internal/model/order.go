package model

import "time"

const (
	None   = iota
	OK     // 交易成功
	FAIL   // 交易失败
	CANCEL // 交易取消
)

// OID 订单ID alias
type OID string

// TransferOrder 转账交易记录
type TransferOrder struct {
	ID         OID       `json:"id" xorm:"not null default '' varchar(32) 'id' "`                  // 转账记录ID
	Ts         time.Time `json:"ts" xorm:"not null default 'CURRENT_TIMESTAMP' timestamp 'ts' "`   // 交易时间(order创建时间)
	Status     uint8     `json:"status" xorm:"not null default 0 int 'status' "`                   // 交易状态
	StatusInfo string    `json:"statusInfo" xorm:"not null default '' varchar(32) 'status_info' "` // 交易状态详情

	From     AID     `json:"from" xorm:"not null default '' varchar(16) 'from' "`          // 转出账户ID
	FromNum  float64 `json:"fromNum" xorm:"not null default 0 decimal(64,12) 'from_num' "` // 转出金额
	FromUnit int     `json:"fromUnit" xorm:"not null default 252 int(4) 'from_unit' "`     // 转出货币单位

	ExchangeRate float64 `json:"exchangeRate" xorm:"not null default 0 decimal(64,12) 'exchange_rate' "` // 转出货币汇率，1转出单位货币 = x转入单位货币

	To     AID     `json:"to" xorm:"not null default '' varchar(16) 'to' "`           // 转入账号ID
	ToUnit int     `json:"toUnit" xorm:"not null default 252 int(4) 'to_unit'" `      // 转入货币单位
	ToNum  float64 `json:"toNum"  xorm:"not null default 0 decimal(64,12) 'to_num' "` // 转入货币数量
}

func (TransferOrder) TableName() string {
	return "thunes_transfer_order"
}
