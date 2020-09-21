package controller

import (
	"errors"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-xorm/xorm"
	"github.com/thunes/internal/db"
	"github.com/thunes/internal/model"
)

// TransferReq 转账请求参数
type TransferReq struct {
	RequestID string    `json:"requestId"`                                                 // 调用方自定义请求ID，可以为空
	From      model.AID `json:"from" form:"from" binding:"required,min=16,max=16"`         // 转出方account_id
	FromUnit  int       `json:"fromUnit" form:"fromUnit" binding:"required,numeric,min=0"` // 转出方货币单位

	To     model.AID `json:"to" form:"to" binding:"required,min=16,max=16"`         // 转入方account_id
	ToUnit int       `json:"toUnit" form:"toUnit" binding:"required,numeric,min=0"` // 转入方货币单位
	ToNum  float64   `json:"toNum" form:"toNum" binding:"required,numeric,min=0"`   // 需要转入的货币数量，以转入方货币单位计算
}

// TransferResp 转账返回信息
type TransferResp struct {
	TransferReq
	TID  string  `json:"tid" binding:"required"`  // 转账ID
	Cost float64 `json:"cost" binding:"required"` // 转账消耗金额，货币单位以请求中FromUnit计算
	Rate float64 `json:"rate" binding:"required"` // 货币换算汇率 1转出方单位货币 = x转入方单位货币
}

// Transfer 转账
// @Summary 转账
// @Description 用户指定账户向指定账户转账
// @Produce  json
// @Param request body controller.TransferReq true "转账请求信息"
// @Param cid path int true "用户ID"
// @Param aid path string true "用户账户ID"
// @Success 200 {object} TransferResp
// @Failure 406 {object} restapi.Response "参数错误"
// @Failure 404 {object} restapi.Response "账户不存在"
// @Failure 451 {object} restapi.Response "账户余额不足"
// @Failure 500 {object} restapi.Response "服务器错误"
// @Router /v1/clients/{cid}/account/{aid}/transfer [post]
func (A *Account) Transfer(ctx *gin.Context) {
	var (
		req = TransferReq{}
	)

	// 参数检查
	if err := ctx.ShouldBindJSON(&req); err != nil {
		fmt.Println("Account: Transfer parse req error, ", err)
		ginInvalidParam(ctx)
		return
	}

	// 获取汇率
	rate, err := A.CurrentExchangeRate(req.FromUnit, req.ToUnit)
	if err != nil {
		fmt.Println("Account: Transfer error[get rate error] ", err)

		ginInternalErr(ctx)
		return
	}

	// 转账事务开始
	var (
		fromAccount = model.Account{}
		total       = req.ToNum / rate
		order       = model.TransferOrder{
			ID:           model.OID(_idGenerator.Generate().Base32()),
			Ts:           time.Now(),
			Status:       model.None,
			From:         req.From,
			FromUnit:     req.FromUnit,
			FromNum:      total,
			ExchangeRate: rate,

			To:     req.To,
			ToNum:  req.ToNum,
			ToUnit: req.ToUnit,
		}
	)

	// 检查from账户余额是否足够，目标货币种类额度是否足够
	if fromAccount, err = A.M.Get(req.From); err != nil {
		fmt.Println("Account: Transfer error [check from account error] ", err)

		if err == ErrAccountNotFound {
			ginInvalidParam(ctx, fmt.Sprintf("from account not exist"))
			return
		}

		ginInternalErr(ctx)
		return
	}

	// 余额检查
	if fromAccount.Balance < total {
		ginError(ctx, ErrBalanceNotOk)
		return
	}

	// 事务开始
	session := db.Cli().NewSession()
	if err = session.Begin(); err != nil {
		fmt.Println("Transfer: start tx error, ", err)
		ginInternalErr(ctx)
		return
	}

	defer func() {
		if err != nil {
			fmt.Println("will rollback", err)
			_ = session.Rollback()
		} else {
			fmt.Println("will commit", err)
			_ = session.Commit()
		}
	}()

	var (
		orderModel = &OrderModel{
			session: session.Table(&model.TransferOrder{}),
		}
		accountModel = &AccountModel{
			session: session.Table(&model.Account{}),
		}
	)

	// 记录交易
	if err = orderModel.NewOrder(order); err != nil {
		fmt.Println("Account: Transfer error [new order error], ", err)

		ginInternalErr(ctx)
		return
	}

	// 扣款
	if err = accountModel.DecrBalance(req.From, total); err != nil {
		fmt.Println("Account: Transfer error [decr from account balance error] ", err)

		ginInternalErr(ctx)
		return
	}

	// 转账
	if err = accountModel.AddBalance(req.To, req.ToNum, req.ToUnit); err != nil {
		fmt.Println("Account: Transfer error [add to account balance error] ", err)

		ginInternalErr(ctx)
		return
	}

	// 更新交易状态
	order.Status = model.OK
	if err = orderModel.SetOrderStatus(order); err != nil {
		fmt.Println("Account: Transfer error [set order status to ok error], ", err)

		ginInternalErr(ctx)
		return
	}

	ginSuccess(ctx, "tranfer success", TransferResp{
		TransferReq: req,
		TID:         string(order.ID),
		Cost:        total,
		Rate:        rate,
	})
}

//////////////////////////////////////////////////////////////////////////////////////

// CurrentExchangeRate 获取当前货币的转换汇率 1 fromUnit = rate toUnit
// fromUnit: 源货币
// toUnit: 目标货币

func (A Account) CurrentExchangeRate(fromUnit int, toUnit int) (rate float64, err error) {
	//TODO
	// 获取当前货币汇率
	return 1, nil
}

////////////////////////////////////////////////////////////////////////////////////////

type OrderModel struct {
	session *xorm.Session
}

func (M *OrderModel) NewOrder(order model.TransferOrder) (err error) {
	n, err := M.session.Table(&model.TransferOrder{}).InsertOne(&order)
	if err != nil {
		return err
	}

	if n == 0 {
		return errors.New("new order error, count is 0")
	}

	return
}

func (M *OrderModel) SetOrderStatus(order model.TransferOrder) (err error) {
	n, err := M.session.Table(&model.TransferOrder{}).
		Where("id = ?", order.ID).
		Cols("status").
		Update(&order)
	if err != nil {
		return err
	}

	// if n == 0 {
	// 	return errors.New("set order status error, count is 0")
	// }

	_ = n

	return
}
