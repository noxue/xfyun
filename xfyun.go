// xfyun project xfyun.go
package xfyun

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

type Xfyun struct {
	url string
	id  string
	key string
}

func New(id, key string) *Xfyun {
	return &Xfyun{
		url: "http://api.xfyun.cn/v1/service/v1/iat",
		id:  id,
		key: key,
	}
}

func (this *Xfyun) RunAsFile(file string) (string, error) {
	b, err := ioutil.ReadFile(file)
	if err != nil {
		return "", err
	}
	return this.RunAsrStream(b)
}

func (this *Xfyun) RunAsrStream(wav []byte) (string, error) {
	base64_audio := base64.StdEncoding.EncodeToString(wav)
	u := url.Values{}
	u.Set("audio", base64_audio)
	body := u.Encode()

	x_param := base64.StdEncoding.EncodeToString([]byte(`{"aue":"raw","engine_type":"sms8k"}`))
	x_time := time.Now().Unix()

	h := md5.New()
	h.Write([]byte(fmt.Sprintf("%s%d%s", this.key, x_time, x_param)))
	x_checksum := hex.EncodeToString(h.Sum(nil))

	client := &http.Client{}

	req, _ := http.NewRequest("POST", this.url, bytes.NewBuffer([]byte(body)))

	req.Header.Set("Accept-Encoding", "identity")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("X-Appid", this.id)
	req.Header.Add("X-CurTime", fmt.Sprintf("%d", x_time))
	req.Header.Add("X-Param", x_param)
	req.Header.Add("X-CheckSum", x_checksum)

	resp, err := client.Do(req)
	if err != nil {
		return "", nil
	}

	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", nil
	}

	return string(html), nil
}
