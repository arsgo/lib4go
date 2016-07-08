package net

import "testing"

func TestSchemaExecute(t *testing.T) {

	http := NewHTTPClient()
	c := http.NewRequest("POST", "https://app.cloopen.com:8883/2013-12-26/Accounts/aaf98fda42c744c90142d505bfab0135/SMS/TemplateSMS?sig=7928964963931AD4538EB873DA454159")
	c.SetData("<?xml version='1.0' encoding='utf-8'?><TemplateSMS><to>13051880135</to><appId>8a48b5514e5298b9014e67a3f02f1411</appId><templateId>95042</templateId><datas><data>123454</data></datas></TemplateSMS>")
	c.SetHeader("Accept", "application/xml")
	c.SetHeader("Content-type", "application/xml")
	c.SetHeader("charset", "utf-8")
	c.SetHeader("Authorization", "YWFmOThmZGE0MmM3NDRjOTAxNDJkNTA1YmZhYjAxMzU6MjAxNjA3MDgwOTEzMjY=")
	r, s, e := c.Request()

	if e != nil {
		t.Error(e)
	}
	if s != 200 {
		t.Error("status error:", s)
	}
	t.Log(r)
}
