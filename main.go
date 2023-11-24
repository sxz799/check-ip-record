package main

import (
	"fmt"
	"github.com/robfig/cron/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"net/http"
	"time"
	"flag"
)

type Record struct {
	ID       uint      `gorm:"primaryKey"`
	IP       string    `gorm:"not null"`
	DateTime time.Time `gorm:"not null"`
}

var db *gorm.DB

func main() {

	dsn := flag.String("dsn", "", "指定dsn")
	flag.Parse()
	log.Println("dsn: ",*dsn)
	db, err := gorm.Open(postgres.Open(*dsn), &gorm.Config{})
	if err!=nil{
		fmt.Println("数据库连接失败!")
		return
	}

	db.AutoMigrate(&Record{})

	c := cron.New()
	c.AddFunc("@every 10m", addIPRecord)
	c.Start()
	select {}

}

func addIPRecord() {
	ip, err := getLocalIPByInternet()
	if err != nil {
		log.Println(err)
		return
	}
	var r Record
	db.Where("ip=?", ip).First(&r)
	if r.ID < 1 {
		db.Create(&Record{
			IP:       ip,
			DateTime: time.Now(),
		})
	}
}

func getLocalIPByInternet() (string, error) {
	resp, err := http.Get("https://ip.sb")
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	ip, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error:", err)
		return "", err
	}
	return string(ip), nil
}
