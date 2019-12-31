package tools

import (
	"encoding/json"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"time"
)

//生成4位随机验证码
func RandStr() string {
	var code string
	rand.Seed(time.Now().Unix())
	for i := 0; i < 4; i++ {
		code = code + strconv.Itoa(rand.Intn(10))
	}
	return code
}

func SendHttpPost(urls string, api string, data map[string]string, token string) (error, *ResponseData) {
	form := url.Values{}
	for k, v := range data {
		form.Set(k, v)
	}
	u, err := url.ParseRequestURI(urls)
	if err != nil {
		return err, nil
	}
	u.Path = api
	urlStr := u.String()

	client := &http.Client{}
	r, _ := http.NewRequest("POST", urlStr, strings.NewReader(form.Encode())) // URL-encoded payload
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if token != "" {
		r.Header.Set("Authorization", token)
	}
	resp, err := client.Do(r)
	if err != nil {
		return err, nil
	}
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	var rep ResponseData
	json.Unmarshal(body, &rep)
	return nil, &rep
}

type ResponseData struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

//格式验证
func VerifyEmailFormat(email string) bool {
	//pattern := `\w+([-+.]\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*` //匹配电子邮箱
	pattern := `^[0-9a-z][_.0-9a-z-]{0,31}@([0-9a-z][0-9a-z-]{0,30}[0-9a-z]\.){1,4}[a-z]{2,4}$`

	reg := regexp.MustCompile(pattern)
	return reg.MatchString(email)
}
