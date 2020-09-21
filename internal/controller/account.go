package controller

import (
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"github.com/thunes/internal/db"
	"github.com/thunes/internal/model"
)

func NewAccountHandler() *Account {
	return &Account{
		M: &AccountModel{
			session: db.Cli().Table(&model.Account{}),
		},
	}
}

// Account 用户账户金额相关操作
type Account struct {
	M *AccountModel
}

//////////////////////////////////////////////////////////////////////////////

type GetAccountReq struct {
	AID int64 `json:"aid" uri:"aid" binding:"required,numeric,min=1,max=9999999999999999"`
}

type GetAccountResp model.Account

// Get 获取指定账户信息
// @Summary 获取指定账户信息
// @Description 获取指定账户信息
// @Produce  json
// @Param aid path string true "账户"
// @Success 200 {object} controller.GetAccountResp
// @Failure 406 {object} restapi.Response "参数错误"
// @Failure 404 {object} restapi.Response "账户不存在"
// @Failure 401 {object} restapi.Response "授权失败/未授权"
// @Failure 500 {object} restapi.Response "服务器错误"
// @Router /v1/clients/{cid}/account/{aid} [get]
func (A *Account) Get(ctx *gin.Context) {
	var (
		req = GetAccountReq{}
	)

	if err := ctx.ShouldBindUri(&req); err != nil {
		fmt.Println("Account: Get error [invalid param], ", err)
		ginInvalidParam(ctx)
		return
	}

	account, err := A.M.Get(model.AID(req.AID))
	if err != nil {
		if err == ErrAccountNotFound {
			ginError(ctx, err)
			return
		}

		ginInternalErr(ctx)
		return
	}

	ginSuccess(ctx, "success", account)
}

//////////////////////////////////////////////////////////////////////

type AccountModel struct {
	session *xorm.Session
}

// Get 获取指定账户信息
func (M *AccountModel) Get(aid model.AID) (account model.Account, err error) {
	accounts := []model.Account{}

	if err = M.session.Where("account_id = ?", aid).Find(&accounts); err != nil {
		return
	}

	if len(accounts) == 0 {
		return account, ErrAccountNotFound
	}

	return accounts[0], nil
}

// AddBalance 账户增加余额
func (M *AccountModel) AddBalance(aid model.AID, delta float64, unit int) (err error) {
	if delta < 0 {
		return errors.New("AddBalance error: delta is negative")
	}

	res, err := M.session.Exec(`
		update thunes_account 
		set balance = balance + ?
		where account_id = ? and unit = ?
		`, delta, aid, unit)
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if n == 0 {
		return ErrAccountNotFound
	}

	return nil
}

// DecrBalance 账户余额减少
func (M *AccountModel) DecrBalance(aid model.AID, delta float64) (err error) {
	if delta < 0 {
		return errors.New("DecrBalance error: delta is negative")
	}

	res, err := M.session.Exec(`
		update thunes_account 
		set balance = balance - ? 
		where account_id = ?
		`, delta, aid)
	if err != nil {
		return err
	}

	n, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if n == 0 {
		return ErrAccountNotFound
	}

	return nil
}

// FindByUID 查询指定用户的账户
func (M *AccountModel) FindByUID(uid model.UID) (accounts []model.Account, err error) {
	accounts = []model.Account{}

	cols := []string{
		"id", "account_id", "balance", "unit", "created_at", "modified_at",
	}

	if err = M.session.Where("id = ?", uid).Cols(cols...).Find(&accounts); err != nil {
		return
	}

	if len(accounts) == 0 {
		return accounts, ErrAccountNotFound
	}

	return
}
