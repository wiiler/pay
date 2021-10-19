package wechat

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"pay/utils"
	"pay/wechat/model"
)

type Response struct {
	PrepayID string `json:"prepay_id"`
}

type ResponseUrl struct {
	CodeUrl string `json:"code_url"`
}

// jsapi、小程序下单
func (c *Client) JsApiPay(d *model.JsApiPay) (*model.JsApiPayResult, error) {
	ctx := context.Background()
	header := make(map[string]string)
	header[utils.ContentType] = utils.ApplicationJSON
	response, err := c.Post(ctx, JsapiPayUrl, header, d)
	if err != nil {
		return nil, err
	}
	if response.Body != nil {
		defer response.Body.Close()
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var res Response
	json.Unmarshal(body, &res)
	if len(res.PrepayID) == 0 {
		return nil, errors.New("失败")
	}
	nonce, _ := utils.GenerateNonce()
	jsRes := &model.JsApiPayResult{
		AppId:     c.AppID,
		TimeStamp: utils.TimeStamp(),
		NonceStr:  nonce,
		Package:   res.PrepayID,
		SignType:  "RSA",
		PaySign:   "",
	}
	msg := fmt.Sprintf("%s\n%s\n%s\n%s\n", jsRes.AppId, jsRes.TimeStamp, jsRes.NonceStr, jsRes.Package)
	jsRes.PaySign, _ = utils.SignSHA256WithRSA(msg, c.PrivateKey)
	return jsRes, nil
}

// App 下单
func (c *Client) AppPay(d *model.AppNativePay) (*model.AppPayResult, error) {
	ctx := context.Background()
	header := make(map[string]string)
	header[utils.ContentType] = utils.ApplicationJSON
	response, err := c.Post(ctx, AppPayUrl, header, d)
	if err != nil {
		return nil, err
	}
	if response.Body != nil {
		defer response.Body.Close()
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	var res Response
	json.Unmarshal(body, &res)
	if len(res.PrepayID) == 0 {
		return nil, errors.New("失败")
	}
	nonce, _ := utils.GenerateNonce()
	jsRes := &model.AppPayResult{
		AppId:     c.AppID,
		PartnerID: c.MchID,
		PrepayID:  res.PrepayID,
		Package:   "Sign=WXPay",
		NonceStr:  nonce,
		TimeStamp: utils.TimeStamp(),
		PaySign:   "",
	}
	msg := fmt.Sprintf("%s\n%s\n%s\n%s\n", jsRes.AppId, jsRes.TimeStamp, jsRes.NonceStr, jsRes.PrepayID)
	jsRes.PaySign, _ = utils.SignSHA256WithRSA(msg, c.PrivateKey)
	return jsRes, nil
}

// Native 下单
func (c *Client) NativePay(d *model.AppNativePay) (string, error) {
	ctx := context.Background()
	header := make(map[string]string)
	header[utils.ContentType] = utils.ApplicationJSON
	response, err := c.Post(ctx, NativePayUrl, header, d)
	if err != nil {
		return "", err
	}
	if response.Body != nil {
		defer response.Body.Close()
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", err
	}

	var res ResponseUrl
	json.Unmarshal(body, &res)
	if len(res.CodeUrl) == 0 {
		return "", errors.New("失败")
	}

	return res.CodeUrl, nil
}

// 微信支付订单号查询
func (c *Client) QueryTransactionID(transactionid string, submchid string) (*model.AppNativePay, error) {
	url := fmt.Sprintf(QueryTransactionIDUrl, transactionid, c.MchID, submchid)
	ctx := context.Background()
	res, err := c.Get(ctx, url)
	if err != nil {
		return nil, err
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	body, _ := ioutil.ReadAll(res.Body)
	var Res model.AppNativePay
	json.Unmarshal(body, &Res)
	return &Res, nil
}

// 商户订单号查询
func (c *Client) QueryOutTradeNo(outtradeno string, submchid string) (*model.AppNativePay, error) {
	url := fmt.Sprintf(QueryOutTradeNoUrl, outtradeno, c.MchID, submchid)
	ctx := context.Background()
	res, err := c.Get(ctx, url)
	if err != nil {
		return nil, err
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	body, _ := ioutil.ReadAll(res.Body)
	var Res model.AppNativePay
	json.Unmarshal(body, &Res)
	return &Res, nil
}

func (c *Client) CloseOutTradeNo(outtradeno, submchid string) (bool, error) {
	url := fmt.Sprintf(CloseOutTradeNoUrl, outtradeno)
	header := make(map[string]string)
	ctx := context.Background()
	header[utils.ContentType] = utils.ApplicationJSON
	d := map[string]string{
		"sp_mchid":  c.MchID,
		"sub_mchid": submchid,
	}
	response, err := c.Post(ctx, url, header, d)
	if err != nil {
		return false, err
	}
	if response.Body != nil {
		defer response.Body.Close()
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return false, err
	}
	if len(string(body)) == 0 && response.StatusCode == 204 {
		return true, nil
	}
	return false, errors.New("关闭失败")
}

//  退款申请
func (c *Client) Refund(d *model.Refund) (*model.RefundResult, error) {
	header := make(map[string]string)
	ctx := context.Background()
	header[utils.ContentType] = utils.ApplicationJSON

	response, err := c.Post(ctx, RefundUrl, header, d)
	if err != nil {
		return nil, err
	}
	if response.Body != nil {
		defer response.Body.Close()
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var res model.RefundResult
	json.Unmarshal(body, &res)
	return &res, nil
}

//  查询退款
func (c *Client) QueryRefund(out_refund_no, sub_mchid string) (*model.RefundResult, error) {
	ctx := context.Background()
	url := fmt.Sprintf(QueryRefundUrl, out_refund_no, sub_mchid)

	response, err := c.Get(ctx, url)
	if err != nil {
		return nil, err
	}
	if response.Body != nil {
		defer response.Body.Close()
	}
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var res model.RefundResult
	json.Unmarshal(body, &res)
	return &res, nil
}
