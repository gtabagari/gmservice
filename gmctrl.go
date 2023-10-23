package main

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/tuya/tuya-connector-go/connector"
	"github.com/tuya/tuya-connector-go/connector/env"
	"github.com/tuya/tuya-connector-go/connector/httplib"
	"github.com/tuya/tuya-connector-go/connector/logger"
	"github.com/tuya/tuya-connector-go/example/model"
	"net/http"
	"time"
)

type CustomContext struct {
	echo.Context
}

func (c *CustomContext) Foo() {
	fmt.Println("foo")
}

func (c *CustomContext) Bar() {
	fmt.Println("bar")
}

func (c *CustomContext) Test1() {
	fmt.Println("Test1")
}

func PostDeviceCmd(s bool) {
	device_id := "bf4fa8c44a11e0a7f6t6ia"
	body := []byte(fmt.Sprintf("{\"commands\":[{\"code\":\"switch_1\",\"value\":%s}]}", fmt.Sprintf("%t", s)))
	resp := &model.PostDeviceCmdResponse{}
	err := connector.MakePostRequest(
		context.Background(),
		connector.WithAPIUri(fmt.Sprintf("/v1.0/devices/%s/commands", device_id)),
		connector.WithPayload(body),
		connector.WithResp(resp))

	if err != nil {
		logger.Log.Errorf("err:%s", err.Error())
		return
	}
}

func main() {

	connector.InitWithOptions(env.WithApiHost(httplib.URL_EU),
		env.WithMsgHost(httplib.MSG_EU),
		env.WithAccessID("t3dsut8mrgpw7tudmxny"),
		env.WithAccessKey("f75cf4f71f38433588dc33a6a6279a58"),
		env.WithAppName("tuyaSDK"),
		env.WithDebugMode(false))
	PostDeviceCmd(false)
	//echo initialization
	e := echo.New()
	e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			cc := &CustomContext{c}
			return next(cc)
		}
	})
	e.File("/", "public/index.html")
	//e.GET("/", func(c echo.Context) error {
	//	cc := c.(*CustomContext)
	//	cc.Foo()
	//	cc.Bar()
	//	return cc.String(200, "OK")
	//})
	e.GET("/test1", func(c echo.Context) error {
		cc := c.(*CustomContext)
		name := cc.QueryParam("name")
		PostDeviceCmd(true)
		time.Sleep(50 * time.Millisecond)
		PostDeviceCmd(false)
		fmt.Println(name)
		return cc.String(http.StatusOK, name)
	})

	e.Logger.Fatal(e.Start(":1323"))
}
