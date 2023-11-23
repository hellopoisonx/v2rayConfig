package dc

import (
	"encoding/base64"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

func GetSub() []byte {
	resp, err := http.Get("https://fatcat.trafficmanager.net/api/v1/client/subscribe?token=29ae95209f8cfc55e7774c66149cd0b2&flag=v2ray")
	count := 1
	for err != nil {
		if count == 6 {
			log.Fatalln("请稍后尝试", err)
		}
		log.Println("第", count, "次请求订阅内容失败", err)
		time.Sleep(time.Second * 3)
		count++
		resp, err = http.Get("https://fatcat.trafficmanager.net/api/v1/client/subscribe?token=29ae95209f8cfc55e7774c66149cd0b2&flag=v2ray")
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln("请求订阅内容失败", err)
	}
	return body
}
func Base64ToSs(code string) ([]string, error) {
	decodeBytes, err := base64.StdEncoding.DecodeString(code)
	if err != nil {
		return nil, err
	}
	address := strings.Split(string(decodeBytes), "\n")
	return address, nil
}
