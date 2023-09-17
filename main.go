// Copyright 2023 BenderBlog Rodriguez.
// SPDX-License-Identifier: MIT
package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
)

type Payload struct {
	Account  string `json:"account"`
	Password string `json:"password"`
}

// GBK 转 UTF-8
func GbkToUtf8(s []byte) ([]byte, error) {
	reader := transform.NewReader(bytes.NewReader(s), simplifiedchinese.GB18030.NewDecoder())
	d, e := io.ReadAll(reader)
	if e != nil {
		return nil, e
	}
	fmt.Printf("%s", d)
	return d, nil
}

func get_data(data *Payload) ([]byte, error) {
	client := resty.New()
	client.SetRedirectPolicy(resty.NoRedirectPolicy())
	resp, err := client.R().
		SetBody(
			"__EVENTTARGET=&__EVENTARGUMENT=&"+
				"__VIEWSTATE=%2FwEPDwUKMTEzNzM0MjM0OWQYAQUeX19D"+
				"b250cm9sc1JlcXVpcmVQb3N0QmFja0tleV9fFgEFD2xvZ2luMSRidG5"+
				"Mb2dpbkOuzGVaztce4Ict7jsIJ0F5pUDb%2BsmSbCCrNVSBlPML&"+
				"__VIEWSTATEGENERATOR=EE008CD9&"+
				"__EVENTVALIDATION=%2FwEdAAcKecdPGDB%2BfW8Tyghx"+
				"7AeSpOzeiNZ7aaEg5p6LqSa9cODI2bZwNtRxUKPkisVLf8l"+
				"8Vv4WhRVIIhZlyYNJO%2BySrDKOhP%2B%2FYMNbVIh74hA2r"+
				"CYnBBSTsX9SjxiYNNk%2B5kglM%2B6pGIq22Oi5mNu6u6eC2W"+
				"EBfKAmATKwSpsOL%2FPNcRyi9l8Dnp6JamksyAzjhW4%3D&"+
				"login1%24StuLoginID="+data.Account+"&"+
				"login1%24StuPassword="+data.Password+"&"+
				"login1%24UserRole=Student&"+
				"login1%24btnLogin.x=28&"+
				"login1%24btnLogin.y=14").
		SetHeader("Content-Type", "application/x-www-form-urlencoded").Post("http://wlsy.xidian.edu.cn/PhyEws/default.aspx")

	fmt.Println("登录返回")
	if err != nil {
		fmt.Println(err.Error())
		// return err
	}

	var cookie = resp.Header()["Set-Cookie"]

	cookieStr := ""
	for i := range cookie {
		cookieStr += cookie[i]
	}

	resp, err = client.R().SetHeader("Cookie", cookieStr).Get("http://wlsy.xidian.edu.cn/PhyEws/student/select.aspx")

	if err != nil {
		return nil, err
	} else {
		return GbkToUtf8(resp.Body())
	}
}

func main() {
	app := fiber.New()

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("西电物理实验中转服务器")
	})

	app.Post("/", func(c *fiber.Ctx) error {
		payload := new(Payload)

		if err := c.BodyParser(payload); err != nil {
			return err
		} else {
			if len(strings.TrimSpace(payload.Account)) == 0 || len(strings.TrimSpace(payload.Password)) == 0 {
				return c.SendString("Error: Account or password have empty")
			}
			if tr, err := get_data(payload); err != nil {
				return c.SendString("Error: " + err.Error())
			} else {
				return c.SendString(string(tr))
			}
		}
	})

	app.Listen(":8080")
}
