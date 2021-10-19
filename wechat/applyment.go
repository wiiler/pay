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

type ApplymentResponse struct {
	ApplymentID int `json:"applyment_id"`
}

// 提交申请
func (c *Client) Applyment(d *model.ApplymentData) (int, error) {
	ctx := context.Background()
	header := make(map[string]string)
	header[utils.ContentType] = utils.ApplicationJSON
	header["Wechatpay-Serial"] = c.WeChatNo
	response, err := c.Post(ctx, ApplymentUrl, header, d)
	if err != nil {
		return 0, err
	}
	body, _ := ioutil.ReadAll(response.Body)

	if response.Body != nil {
		defer response.Body.Close()
	}
	var Res ApplymentResponse
	json.Unmarshal(body, &Res)
	return Res.ApplymentID, nil
}

// 查询申请单状态 id 为微信返回id
func (c *Client) QueryApplyment(id string) (*model.ApplymentResult, error) {
	url := fmt.Sprintf(QueryApplymentUrl, id)
	ctx := context.Background()
	response, err := c.Get(ctx, url)
	if err != nil {
		return nil, err
	}
	body, _ := ioutil.ReadAll(response.Body)

	if response.Body != nil {
		defer response.Body.Close()
	}
	var res model.ApplymentResult
	json.Unmarshal(body, &res)
	return &res, nil
}

// 修改结算账号
func (c *Client) UpdateMchBank(mchId string, d *model.UpdateMchBankData) (bool, error) {
	url := fmt.Sprintf(UpdateMchBankUrl, mchId)
	ctx := context.Background()
	header := make(map[string]string)
	header[utils.ContentType] = utils.ApplicationJSON
	response, err := c.Post(ctx, url, header, d)
	if err != nil {
		return false, err
	}
	if response.Body != nil {
		defer response.Body.Close()
	}

	if response.StatusCode == 204 {
		return true, nil
	} else {
		body, _ := ioutil.ReadAll(response.Body)
		return false, errors.New(string(body))
	}

}

// 查询结算账户 mchId 为商户号
func (c *Client) QueryMchBank(mchId string) (*model.ApplymentBankResult, error) {
	url := fmt.Sprintf(QueryMchBankUrl, mchId)
	ctx := context.Background()
	response, err := c.Get(ctx, url)
	if err != nil {
		return nil, err
	}
	body, _ := ioutil.ReadAll(response.Body)

	if response.Body != nil {
		defer response.Body.Close()
	}
	var res model.ApplymentBankResult
	json.Unmarshal(body, &res)
	return &res, nil
}
