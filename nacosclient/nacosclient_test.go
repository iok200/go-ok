package nacosclient

import (
	"fmt"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"testing"
	"time"
)

func TestClient(t *testing.T) {
	{
		client, err := LoadConfigClient()
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
		_, err := LoadNamingClient()
		if err != nil {
			fmt.Println(err)
			return
		}

	}
	time.Sleep(time.Hour)
}
