package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/gosuri/uilive"
	"github.com/guptarohit/asciigraph"
)

type BTC struct {
	Data struct {
		Buy_price string `json:"buy_price"`
	} `json:"BTC_USD"`
}

type ETH struct {
	Data struct {
		Buy_price string `json:"buy_price"`
	} `json:"ETH_USD"`
}
type LTC struct {
	Data struct {
		Buy_price string `json:"buy_price"`
	} `json:"LTC_USD"`
}

var gw = asciigraph.Width(100)
var gh = asciigraph.Height(10)

func getResp(respch chan *http.Response) {
	client := &http.Client{}

	req, err := http.NewRequest("POST", "https://api.exmo.com/v1.1/ticker", nil)
	if err != nil {
		fmt.Errorf("Невозможно создать запрос")
		return
	}
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	ctx := req.Context()
	ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
	defer cancel()
	resp, err := client.Do(req)
	if err != nil {
		fmt.Errorf("Невозможно отправить запрос")
		return
	}
	if ctx.Err() == context.DeadlineExceeded {
		log.Println("Request timeout", http.StatusRequestTimeout)
		time.Sleep(2 * time.Second)
	}
	respch <- resp
}
func getTime(ch chan string) {
	timeNow := time.Now()
	hour, min, sec := timeNow.Clock()
	year, mon, day := timeNow.Date()
	var hourStr, minStr, secStr, yearStr, monStr, dayStr string
	sliceInt := []int{hour, min, sec, year, int(mon), day}
	sliceStr := []string{hourStr, minStr, secStr, yearStr, monStr, dayStr}
	for i, num := range sliceInt {
		if num < 10 {
			sliceStr[i] = "0" + strconv.Itoa(num)
		} else {
			sliceStr[i] = strconv.Itoa(num)
		}
		ch <- sliceStr[i]
	}
}
func keyboarInGraphs(ch chan string) {
	for {
		_, key, err := keyboard.GetKey()
		if err != nil {
			fmt.Println("Ошибка чтения клавиши!")
		}
		if key == 0x7F {
			ch <- "goroutine is done"
			return
		}
	}
}
func mainMenu() {
	fmt.Print("\033[H\033[2J")
	fmt.Printf("1. BTC_USD\n2. LTC_USD\n3. ETH_USD\n\nPress 1-3 to change symbol, press q to exit\n")
	for {
		char, _, err := keyboard.GetKey()
		if err != nil {
			fmt.Println("Не удалось считать клавишу!")
		}
		if string(char) == "1" {
			BTC_USD()
		}
		if string(char) == "2" {
			LTC_USD()
		}
		if string(char) == "3" {
			ETH_USD()
		}
		if string(char) == "q" {
			keyboard.Close()
			log.Fatal("Программа завершенна!")
		}
	}

}
func BTC_USD() {
	defer mainMenu()
	fmt.Print("\033[H\033[2J")
	chResp := make(chan *http.Response)
	chTime := make(chan string, 6)
	chKeyboard := make(chan string)
	graphData := make([]float64, 0, 100)
	writer := uilive.New()
	writer.Start()
	go keyboarInGraphs(chKeyboard)
	defer writer.Stop()
	for {
		go getResp(chResp)
		go getTime(chTime)
		var btc BTC
		resp := <-chResp
		body, _ := io.ReadAll(resp.Body)
		json.Unmarshal(body, &btc)
		defer resp.Body.Close()
		value, _ := stringToFloat(btc.Data.Buy_price)
		if len(graphData) == 100 {
			graphData = append(graphData[1:], value)
		} else {
			graphData = append(graphData, value)
		}
		graph := asciigraph.Plot(graphData, gw, gh, asciigraph.SeriesColors(asciigraph.Red))
		fmt.Fprintf(writer, "\nBTC_USD: %.2f\n", value)
		fmt.Fprintf(writer, "%s", graph)
		hour, min, sec, year, mon, day := <-chTime, <-chTime, <-chTime, <-chTime, <-chTime, <-chTime
		fmt.Fprintf(writer, "\nТекущая дата: %s", strings.Join([]string{year, mon, day}, "-"))
		fmt.Fprintf(writer, "\nТекущее время: %s", strings.Join([]string{hour, min, sec}, ":"))
		select {
		case msg := <-chKeyboard:
			if msg == "goroutine is done" {
				return
			}
		case <-time.After(time.Second):
			continue
		}

	}
}
func LTC_USD() {
	defer mainMenu()
	fmt.Print("\033[H\033[2J")
	chResp := make(chan *http.Response)
	chTime := make(chan string, 6)
	chKeyboard := make(chan string)
	graphData := make([]float64, 0, 100)
	writer := uilive.New()
	writer.Start()
	go keyboarInGraphs(chKeyboard)
	defer writer.Stop()
	for {
		go getResp(chResp)
		go getTime(chTime)
		var ltc LTC
		resp := <-chResp
		body, _ := io.ReadAll(resp.Body)
		json.Unmarshal(body, &ltc)
		defer resp.Body.Close()
		value, _ := stringToFloat(ltc.Data.Buy_price)
		if len(graphData) == 100 {
			graphData = append(graphData[1:], value)
		} else {
			graphData = append(graphData, value)
		}
		graph := asciigraph.Plot(graphData, gw, gh, asciigraph.SeriesColors(asciigraph.Red))
		fmt.Fprintf(writer, "\nLTC_USD: %.2f\n", value)
		fmt.Fprintf(writer, "%s", graph)
		hour, min, sec, year, mon, day := <-chTime, <-chTime, <-chTime, <-chTime, <-chTime, <-chTime
		fmt.Fprintf(writer, "\nТекущая дата: %s", strings.Join([]string{year, mon, day}, "-"))
		fmt.Fprintf(writer, "\nТекущее время: %s", strings.Join([]string{hour, min, sec}, ":"))
		select {
		case msg := <-chKeyboard:
			if msg == "goroutine is done" {
				return
			}
		case <-time.After(time.Second):
			continue
		}

	}
}
func ETH_USD() {
	defer mainMenu()
	fmt.Print("\033[H\033[2J")
	chResp := make(chan *http.Response)
	chTime := make(chan string, 6)
	chKeyboard := make(chan string)
	graphData := make([]float64, 0, 100)
	writer := uilive.New()
	writer.Start()
	go keyboarInGraphs(chKeyboard)
	defer writer.Stop()
	for {
		go getResp(chResp)
		go getTime(chTime)
		var eth ETH
		resp := <-chResp
		body, _ := io.ReadAll(resp.Body)
		json.Unmarshal(body, &eth)
		defer resp.Body.Close()
		value, _ := stringToFloat(eth.Data.Buy_price)
		if len(graphData) == 100 {
			graphData = append(graphData[1:], value)
		} else {
			graphData = append(graphData, value)
		}
		graph := asciigraph.Plot(graphData, gw, gh, asciigraph.SeriesColors(asciigraph.Red))
		fmt.Fprintf(writer, "\nETH_USD: %.2f\n", value)
		fmt.Fprintf(writer, "%s", graph)
		hour, min, sec, year, mon, day := <-chTime, <-chTime, <-chTime, <-chTime, <-chTime, <-chTime
		fmt.Fprintf(writer, "\nТекущая дата: %s", strings.Join([]string{year, mon, day}, "-"))
		fmt.Fprintf(writer, "\nТекущее время: %s", strings.Join([]string{hour, min, sec}, ":"))
		select {
		case msg := <-chKeyboard:
			if msg == "goroutine is done" {
				return
			}
		case <-time.After(time.Second):
			continue
		}

	}
}
func stringToFloat(x string) (float64, error) {
	a, err := strconv.ParseFloat(x, 64)
	if err != nil {
		fmt.Println("Ошибка получения значения")
		return 0, err
	}
	return a, err
}
func main() {
	if err := keyboard.Open(); err != nil {
		panic("Не удалось взаимодействовать с клавиатурой")
	}
	defer func() {
		_ = keyboard.Close()
	}()
	mainMenu()

}
