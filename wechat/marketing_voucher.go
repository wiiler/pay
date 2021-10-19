package wechat

import (
	"context"
	"fmt"
	"io/ioutil"
	"pay/utils"
	"pay/wechat/model"
)

//  代金劵
//  CreateVoucher 创建优惠卷
func (c *Client) CreateVoucher(data *model.Stocks) (string, error) {
	ctx := context.Background()
	header := make(map[string]string)
	// ContentType
	header[utils.ContentType] = utils.ApplicationJSON
	response, err := c.Post(ctx, VoucherStocksUrl, header, data)
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
	return string(body), nil
}

// StartVoucher 激活优惠卷
// 创建批次的商户号 mchid
// 批次号 stockid
func (c *Client) StartVoucher(stockid string) (string, error) {
	param := map[string]string{
		"stock_creator_mchid": c.MchID,
	}
	ctx := context.Background()
	header := make(map[string]string)
	header[utils.ContentType] = utils.ApplicationJSON
	url := fmt.Sprintf(VoucherStartStocksUrl, stockid)
	response, err := c.Post(ctx, url, header, param)

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
	return string(body), nil
}

// PauseVoucher 暂停优惠卷
func (c *Client) PauseVoucher(stockid string) (string, error) {
	param := map[string]string{
		"stock_creator_mchid": c.MchID,
	}
	ctx := context.Background()
	url := fmt.Sprintf(VoucherPauseStocksUrl, stockid)
	header := make(map[string]string)
	header[utils.ContentType] = utils.ApplicationJSON
	response, err := c.Post(ctx, url, header, param)

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
	return string(body), nil
}

// PauseVoucher 暂停优惠卷
func (c *Client) RestartVoucher(stockid string) (string, error) {
	param := map[string]string{
		"stock_creator_mchid": c.MchID,
	}
	ctx := context.Background()
	url := fmt.Sprintf(VoucherRestartStocksUrl, stockid)
	header := make(map[string]string)
	header[utils.ContentType] = utils.ApplicationJSON
	response, err := c.Post(ctx, url, header, param)

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
	return string(body), nil
}

// GetVoucher发放优惠卷
func (c *Client) GetVoucher(openid string, merchantClient *Client, stockid string) (string, error) {
	ctx := context.Background()
	url := fmt.Sprintf(VoucherGetStocksUrl, openid)
	header := make(map[string]string)
	header[utils.ContentType] = utils.ApplicationJSON

	param := map[string]string{
		"stock_creator_mchid": merchantClient.MchID,
		"stock_id":            stockid,
		"out_request_no":      utils.OrderNo(),
		"appid":               merchantClient.AppID,
	}
	response, err := merchantClient.Post(ctx, url, header, param)
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
	return string(body), nil
}
