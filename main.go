package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"os/exec"
	"strconv"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go"
	"github.com/joho/godotenv"
)

type SpeedTestResult struct {
	Download struct {
		Bandwidth float64 `json:"bandwidth"`
		Bytes     float64 `json:"bytes"`
		Elapsed   float64 `json:"elapsed"`
		Latency   struct {
			Jitter  float64 `json:"jitter"`
			Latency float64 `json:"latency"`
			Low     float64 `json:"low"`
			Iqm     float64 `json:"iqm"`
		} `json:"latency"`
	} `json:"download"`
	Upload struct {
		Bandwidth float64 `json:"bandwidth"`
		Bytes     float64 `json:"bytes"`
		Elapsed   float64 `json:"elapsed"`
		Latency   struct {
			Jitter  float64 `json:"jitter"`
			Latency float64 `json:"latency"`
			Low     float64 `json:"low"`
			Iqm     float64 `json:"iqm"`
		} `json:"latency"`
	} `json:"upload"`
	PacketLoss float64 `json:"packetLoss"`
	Isp        string  `json:"isp"`
	Server     struct {
		ID       int    `json:"id"`
		Name     string `json:"name"`
		Location string `json:"location"`
		Country  string `json:"country"`
		Host     string `json:"host"`
		Port     int    `json:"port"`
		IP       string `json:"ip"`
	} `json:"server"`
	Result struct {
		ID        string `json:"id"`
		URL       string `json:"url"`
		Persisted bool   `json:"persisted"`
	} `json:"result"`
	Ping struct {
		Jitter  float64 `json:"jitter"`
		Latency float64 `json:"latency"`
		Low     float64 `json:"low"`
		Iqm     float64 `json:"iqm"`
	} `json:"ping"`
	Interface struct {
		InternalIP string `json:"internalIp"`
		Name       string `json:"name"`
		MacAddr    string `json:"macAddr"`
		IsVpn      bool   `json:"isVpn"`
		ExternalIP string `json:"externalIp"`
	} `json:"interface"`
}

func prettyByteSize(b int) string {
	bf := float64(b)
	for _, unit := range []string{"", "K", "M", "G", "T", "P", "E", "Z"} {
		if math.Abs(bf) < 1024.0 {
			return fmt.Sprintf("%3.1f%sbps", bf, unit)
		}
		bf /= 1024.0
	}
	return fmt.Sprintf("%.1fYbps", bf)
}

// pass data to influxdb
func sendToInfluxDB(data SpeedTestResult) {
	host := os.Getenv("INFLUX_HOST")
	token := os.Getenv("INFLUX_TOKEN")
	org := os.Getenv("INFLUX_ORG")
	bucket := os.Getenv("INFLUX_BUCKET")
	client := influxdb2.NewClient(host, token)
	writeAPI := client.WriteAPIBlocking(org, bucket)
	// create point using fluent style

	download := (data.Download.Bandwidth / .125)
	upload := (data.Upload.Bandwidth / .125)
	ping := data.Ping.Latency
	log.Printf("Download: %s Upload: %s Ping: %.0f ms\n", prettyByteSize(int(download)), prettyByteSize(int(upload)), ping)

	p := influxdb2.NewPointWithMeasurement("speedtest").
		AddTag("server", strconv.Itoa(data.Server.ID)).
		AddTag("server_name", data.Server.Name).
		AddTag("server_country", data.Server.Country).
		AddField("download_speed", download).
		AddField("upload_speed", upload).
		AddField("ping", ping).
		AddField("link", data.Result.URL).
		SetTime(time.Now())
	writeAPI.WritePoint(context.Background(), p)
	client.Close()
}

// func loadFile() string {
// 	data, err := os.ReadFile("out.json")
// 	if err != nil {
// 		log.Fatal(err)
// 	}

// 	return string(data)
// }

func worker() {
	out, err := exec.Command("speedtest", "--accept-license", "--accept-gdpr", "--format=json").Output()

	if err != nil {
		log.Fatal(err)
	}
	//out := loadFile()
	var parsed SpeedTestResult
	parseErr := json.Unmarshal([]byte(out), &parsed)
	if parseErr != nil {
		log.Fatal(parseErr)
	}
	log.Println("Running speedtest")
	sendToInfluxDB(parsed)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	wait, err := strconv.ParseInt(os.Getenv("APP_INTERVAL"), 10, 0)
	if err != nil {
		log.Fatal(err)
	}
	tick := time.Duration(wait) * time.Minute
	ticker := time.NewTicker(tick)
	log.Println("Starting speedtest")
	worker()
	log.Printf("Next run in %d minutes", wait)
	for range ticker.C {
		worker()
		log.Printf("Next run in %d minutes", wait)
	}
}
