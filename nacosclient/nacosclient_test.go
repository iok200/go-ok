package nacosclient

import (
	"fmt"
	"github.com/iok200/go-ok/config"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	config.SetFilePath("../ok.yaml")
	conf, err := config.Load()
	if err != nil {
		return
	} else {
		fmt.Println(conf)
	}
	{
		client, err := LoadConfigClient("127.0.0.1:8848")
		if err != nil {
			fmt.Println(err)
			return
		}
		client.PublishConfig(vo.ConfigParam{
			DataId:  "666",
			Group:   "6",
			Content: "hello world!222222",
		})
		content, err := client.GetConfig(vo.ConfigParam{
			DataId: "666",
			Group:  "6",
		})
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(content)
		}
		content, err = client.GetConfig(vo.ConfigParam{
			DataId: "6661",
			Group:  "6",
		})
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(content)
		}
	}
	{
		_, err := LoadNamingClient("127.0.0.1:8848")
		if err != nil {
			fmt.Println(err)
			return
		}

	}
	time.Sleep(time.Hour)
}
