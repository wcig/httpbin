package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/get", handler)
	_ = r.Run(":80")
}

func handler(c *gin.Context) {
	// args
	args := make(map[string]string)
	if err := c.BindQuery(&args); err != nil {
		fmt.Println(err)
	}

	// headers
	headers := make(map[string]string)
	rawHeaders := c.Request.Header
	for k, v := range rawHeaders {
		if len(v) > 0 {
			headers[k] = v[0]
		}
	}

	// origin
	origin := c.ClientIP()

	// url
	url := getRawRequestUrl(c)

	// response
	info := &httpBinInfo{
		Args:    args,
		Headers: headers,
		Origin:  origin,
		Url:     url,
	}
	var result string
	if data, err := jsonEncoding(info); err != nil {
		fmt.Println(err)
	} else {
		result = string(data)
	}
	c.String(200, "%s", result)
}

func getRawRequestUrl(c *gin.Context) string {
	scheme := "http"
	if c.Request.TLS != nil {
		scheme = "https"
	}
	return scheme + "://" + c.Request.Host + c.Request.RequestURI
}

func jsonEncoding(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}
	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)
	encoder.SetIndent("", "  ")
	err := encoder.Encode(t)
	return buffer.Bytes(), err
}

type httpBinInfo struct {
	Args    map[string]string `json:"args"`
	Headers map[string]string `json:"headers"`
	Origin  string            `json:"origin"`
	Url     string            `json:"url"`
}
