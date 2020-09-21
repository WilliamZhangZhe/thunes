package main

import (
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/thunes/internal/controller"
	"github.com/thunes/internal/model"
	"github.com/thunes/pkg/restapi"
	"gopkg.in/resty.v1"
)

func TestTransfer(t *testing.T) {
	req := controller.TransferReq{
		RequestID: time.Now().String(),
		From:      "6666999988881111",
		FromUnit:  1,
		To:        "6666999988882222",
		ToNum:     10,
		ToUnit:    2,
	}

	ret, err := Transfer(req)
	if err != nil {
		t.Fail()
	}

	t.Log(err, ret)
}

// TestTransferRace 并发操作统一账户
func TestTransferRace(t *testing.T) {
	n := 10
	wg := sync.WaitGroup{}
	wg.Add(n)

	for i := 0; i < n; i++ {
		go func() {
			req := controller.TransferReq{
				RequestID: time.Now().String(),
				From:      "6666999988881111",
				FromUnit:  1,
				To:        "6666999988882222",
				ToNum:     10,
				ToUnit:    2,
			}

			ret, err := Transfer(req)
			if err != nil {
				t.Fail()
			}

			t.Log(err, ret)

			wg.Done()
		}()
	}

	wg.Wait()
}

// TestTransferFromNotEnough 转出账号余额不够
func TestTransferFromNotEnough(t *testing.T) {
	req := controller.TransferReq{
		RequestID: time.Now().String(),
		From:      "6666999988883333",
		FromUnit:  1,
		To:        "6666999988884444",
		ToNum:     1000000000000000,
		ToUnit:    2,
	}

	// data, _ := json.Marshal(req)
	// fmt.Println(string(data))

	ret, err := Transfer(req)
	if err == nil {
		t.Fail()
	}

	t.Log(err, "\n", ret)
}

func Transfer(req controller.TransferReq) (ret controller.TransferResp, err error) {
	var (
		api = "http://localhost:8099/v1/clients/1/account/-1/transfer"
		res = struct {
			restapi.Response
			Data controller.TransferResp `json:"data"`
		}{}
	)

	resp, err := resty.R().SetBody(&req).SetResult(&res).Post(api)
	if err != nil {
		return res.Data, err
	}

	if resp.StatusCode() != http.StatusOK {
		err = fmt.Errorf("request error: http status = %d", resp.StatusCode())
		return
	}

	if res.Code != http.StatusOK {
		err = fmt.Errorf("server return error, %d => %s", res.Code, res.Msg)
		return
	}

	return res.Data, err
}

func GetClient(clientID model.UID) (client controller.GetClientResp, err error) {
	var (
		api = fmt.Sprintf("http://localhost:8099/v1/clients/%d", clientID)
		res = struct {
			restapi.Response
			Data controller.GetClientResp `json:"data"`
		}{}
	)

	resp, err := resty.R().SetResult(&res).Get(api)
	if err != nil {
		return res.Data, err
	}

	if resp.StatusCode() != http.StatusOK {
		err = fmt.Errorf("request error: http status = %d", resp.StatusCode())
	}

	return res.Data, err
}
