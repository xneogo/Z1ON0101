package geoip

import (
	"fmt"
	"github.com/ip2location/ip2location-go/v9"
)

func main() {
	// 替换为你的 IP2Location 数据库文件路径
	dbPath := "./250317/IP2LOCATION-LITE-DB11.BIN"

	// 创建 IP2Location 客户端
	client, err := ip2location.OpenDB(dbPath)
	if err != nil {
		panic(err)
	}
	defer client.Close()

	// 要查询的 IP 地址
	ip := "49.95.191.187"

	// 获取 IP 地址的地理位置信息
	result, err := client.Get_all(ip)
	if err != nil {
		panic(err)
	}

	// 打印地理位置信息
	fmt.Printf("IP: %s\n", ip)
	fmt.Printf("Country: %s\n", result.Country_long)
	fmt.Printf("Region: %s\n", result.Region)
	fmt.Printf("City: %s\n", result.City)
	fmt.Printf("Latitude: %f\n", result.Latitude)
	fmt.Printf("Longitude: %f\n", result.Longitude)
	fmt.Printf("Zipcode: %s\n", result.Zipcode)
	fmt.Printf("Timezone: %s\n", result.Timezone)
	fmt.Printf("All: %+v\n", result)
}
