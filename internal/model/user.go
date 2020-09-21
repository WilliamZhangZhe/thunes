package model

import (
	"time"
)

// UID 用户ID类型别名
type UID int64

// User 用户信息，用于用户身份验证
type User struct {
	ID    UID    `json:"id" xorm:"not null int(16) auto_increment 'id' "`
	Name  string `json:"name" xorm:"not null default '' varchar(32) 'name' "`
	Email string `json:"email" xorm:"not null default '' varchar(64) 'email' "`
	PWD   string `json:"pwd" xorm:"not null default '' varchar(64) 'pwd' "`
}

func (User) TableName() string {
	return "thunes_user"
}

////////////////////////////////////////////////////////////////////////////////

type AID string

// Account 用户账户信息，用于账户余额管理记录
type Account struct {
	ID         UID       `json:"uid" xorm:"not null int(16) auto_increment 'id' "`                              // 账户属主ID，同User.ID相同
	AccountID  AID       `json:"accountId" xorm:"not null default '' varchar(16) 'account_id' "`                // 账户ID
	Balance    float64   `json:"balance" xorm:"not null default 0 decimal(64,12) 'balance' "`                   // 账户余额
	Unit       int       `json:"unit" xorm:"not null default 252 int(4) 'unit' "`                               // 账户余额货币单位
	CreatedAt  time.Time `json:"createdAt" xorm:"not null default 'CURRENT_TIMESTAMP' created 'created_at' "`   // 创建时间
	ModifiedAt time.Time `json:"modifiedAt" xorm:"not null default 'CURRENT_TIMESTAMP' updated 'modified_at' "` // 创建时间
}

func (Account) TableName() string {
	return "thunes_account"
}
