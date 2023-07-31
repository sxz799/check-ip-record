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
)

type Record struct {
	ID       uint      `gorm:"primaryKey"`
	IP       string    `gorm:"not null"`
	DateTime time.Time `gorm:"not null"`
}

var db *gorm.DB

func main() {
	// Connect to the SQLite database using GORM postgres://render_w6mb_user:5aj1sSgABmaoAPRTdoB1G605n9ePQzhI@dpg-cj3hstdiuie55pl5816g-a/render_w6mb
	dsn := "postgres://render_w6mb_user:5aj1sSgABmaoAPRTdoB1G605n9ePQzhI@dpg-cj3hstdiuie55pl5816g-a.singapore-postgres.render.com/render_w6mb"
	db, _ = gorm.Open(postgres.Open(dsn), &gorm.Config{})

	// Auto-migrate the User model to create the table in the database
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
