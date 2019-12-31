package tools

import (
	"bytes"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

// WalletInfoResp 钱包信息
type WalletInfoResp struct {
	OrderId string `json:"order_id"`
	Address string `json:"address"`
	AppId   string `json:"app_id"`
	Hash    string `json:"hash"`
}

// WalletResp 钱包返回
type WalletResp struct {
	Code int            `json:"code"`
	Msg  string         `json:"msg"`
	Data WalletInfoResp `json:"data"`
}

// HttpRequestResp 请求返回
type HttpRequestResp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

/**
 * @param param array 字符串数组
 * @param string 签名密钥
 * @return string 签名
 */
func GenerateSign(param map[string]interface{}, secret string) string {
	var sign string
	var key []string
	// 按照参数数组的key升序排序
	for k := range param {
		key = append(key, k)
	}
	sort.Strings(key)
	// 生成签名
	for _, k := range key {
		// sign参数不参于签名
		if k == "hash" || k == "claims" {
			continue
		}
		sign += fmt.Sprintf("%s=%v&", k, param[k])
	}
	fmt.Println(sign + fmt.Sprintf("hash=%s", secret))
	md5HashInBytes := md5.Sum([]byte(sign + fmt.Sprintf("hash=%s", secret)))
	md5HashInString := hex.EncodeToString(md5HashInBytes[:])
	return strings.ToUpper(md5HashInString)
}

// RegisterWalletAddr 注册钱包地址
func RegisterWalletAddr(app_id, url, srekey string) (*WalletResp, error) {
	p := map[string]interface{}{
		"app_id":   app_id,
		"order_id": fmt.Sprintf("%s", uuid.NewV4()),
	}
	p["hash"] = GenerateSign(p, srekey)
	data, err := HttpPostBase(url, p)
	return data, err
}

// HttpPost 请求
func HttpPost(p map[string]interface{}, url, srekey string) (*WalletResp, error) {
	p["hash"] = GenerateSign(p, srekey)
	data, err := HttpPostBase(url, p)
	return data, err
}

// HttpGetBase get请求基础
func HttpGetBase(url string, param, headers map[string]string) (*HttpRequestResp, error) {
	//new request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.New("new request is fail ")
	}
	//add params
	q := req.URL.Query()
	if param != nil {
		for key, val := range param {
			q.Add(key, val)
		}
		req.URL.RawQuery = q.Encode()
	}
	//add headers
	if headers != nil {
		for key, val := range headers {
			req.Header.Add(key, val)
		}
	}
	//http client
	client := &http.Client{}
	log.Printf("Go %s URL : %s \n", http.MethodGet, req.URL.String())
	resp, err := client.Do(req)

	defer resp.Body.Close()
	if err != nil {
		logrus.Error("http get url : %s data : %v error:%v", url, param, err)
		return nil, err
	}
	result, _ := ioutil.ReadAll(resp.Body)
	data := &HttpRequestResp{}
	_ = json.Unmarshal([]byte(string(result)), data)
	return data, nil
}

// httpPostBase post请求基础
func HttpPostBase(url string, param map[string]interface{}) (*WalletResp, error) {
	jsonStr, _ := json.Marshal(param)
	req, err := http.NewRequest(`POST`, url, bytes.NewBuffer(jsonStr))
	req.Header.Add(`content-type`, "application/json")
	defer req.Body.Close()
	log.Printf("Go %s URL : %s \n", http.MethodPost, req.URL.String())
	if err != nil {
		logrus.Error("http post url : %s data : %v error:%v", url, param, err)
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)

	defer resp.Body.Close()
	if err != nil {
		logrus.Error("http post url : %s data : %v error:%v", url, param, err)
		return nil, err
	}
	result, _ := ioutil.ReadAll(resp.Body)
	data := &WalletResp{}
	_ = json.Unmarshal([]byte(string(result)), data)
	return data, nil
}

//密码相关以及格式验证
//生产用户密码盐
func NewPwdSalt(id string, retime int) string {
	return Hash256(id, strconv.Itoa(retime))
}

// Hash256 生成盐
func Hash256(pwd, salt string) string {
	s := pwd + salt
	h := sha256.New()
	h.Write([]byte(s))
	hs := h.Sum(nil)
	return fmt.Sprintf("%x", hs)
}
