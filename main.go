package main

import (
	"flag"
	"fmt"
	"github.com/robfig/cron/v3"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"net/http"
	"time"
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
	log.Println("dsn: ", *dsn)
	var err error
	db, err = gorm.Open(postgres.Open(*dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("数据库连接失败!")
		return
	}

	db.AutoMigrate(&Record{})
	addIPRecord()
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
	r.IP = ip
	r.DateTime = time.Now()
	db.Where(&Record{IP: ip}).FirstOrCreate(&r)
}

func getLocalIPByInternet() (string, error) {
	resp, err := http.Get("https://api.ipify.org")
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
