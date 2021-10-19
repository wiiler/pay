package wechat

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"pay/utils"
	"pay/wechat/model"
)

func (c *Client) GetCert() ([]*model.ItemCerts, error) {
	ctx := context.Background()
	res, err := c.Get(ctx, GetCertsUrl)
	if err != nil {
		return nil, err
	}
	if res.Body != nil {
		defer res.Body.Close()
	}
	body, _ := ioutil.ReadAll(res.Body)
	var certs model.CertsResp
	json.Unmarshal(body, &certs)
	for i, v := range certs.Data {
		c, e := utils.DecryptAES256GCM(c.MchAPIv3Key, v.EncryptCertificate.AssociatedData, v.EncryptCertificate.Nonce, v.EncryptCertificate.Ciphertext)
		if e != nil {
			return nil, e
		}
		certs.Data[i].PublicKey = c
	}
	return certs.Data, nil
}
