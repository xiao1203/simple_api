package main

import (
	"crypto/sha1"
	"encoding/csv"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

// リクエストボディの構造体
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type FlagRequest struct {
	Flag string `json:"flag"`
}

// レスポンスボディの構造体
type LoginResponse struct {
	Token string `json:"token"`
}

type CandleResponse struct {
	Open  uint64 `json:"ope"`
	Close uint64 `json:"close"`
	High  uint64 `json:"high"`
	Low   uint64 `json:"low"`
}

func main() {
	e := echo.New()

	// ログインエンドポイントの設定
	e.PUT("/login", login)

	// フラグの受け取り口
	e.PUT("/flag", flag)

	e.GET("/candle", candle)

	// サーバーを起動
	e.Logger.Fatal(e.Start(":8080"))
}

// ログイン処理
func login(c echo.Context) error {
	req := new(LoginRequest)

	// リクエストボディをバインド
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid request"})
	}

	// ユーザー名とパスワードを連結してトークンを生成
	token := req.Username + req.Password
	//sha1のchecksumを取得
	sha1 := sha1.New()
	io.WriteString(sha1, token)

	// レスポンスを返す
	return c.JSON(http.StatusOK, LoginResponse{Token: hex.EncodeToString(sha1.Sum(nil))})
}

func flag(c echo.Context) error {
	req := new(FlagRequest)

	// リクエストボディをバインド
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid request"})
	}

	// フラグを返す
	fmt.Println("-------------------------------")
	fmt.Println("Flag: " + req.Flag)
	fmt.Println("-------------------------------")

	return c.String(http.StatusOK, "flag{this_is_fake_flag}")
}

func candle(c echo.Context) error {
	code := c.QueryParam("code")
	yearStr := c.QueryParam("year")
	monthStr := c.QueryParam("month")
	dayStr := c.QueryParam("day")
	hourStr := c.QueryParam("hour")

	// string to uint16
	year, err := convertStrToUint(yearStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid request: year"})
	}
	month, err := convertStrToUint(monthStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid request: month"})
	}

	day, err := convertStrToUint(dayStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid request: day"})
	}

	hour, err := convertStrToUint(hourStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "invalid request: hour"})
	}

	// ローソク足の計算
	result, err := getCandleDate(code, year, month, day, hour)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"message": "calculation error"})
	}

	return c.JSON(http.StatusOK, result)
}

func convertStrToUint(str string) (uint16, error) {
	// string to uint16
	integer, err := strconv.ParseUint(str, 10, 16)
	if err != nil {
		return 0, err
	}

	return uint16(integer), nil
}

// データ構造体の定義
type Stock struct {
	Time  time.Time
	Code  string
	Price uint64
}

func getCandleDate(code string, year uint16, month uint16, day uint16, hour uint16) (CandleResponse, error) {
	// csvファイルからデータを取得
	csvFile, err := os.Open("order_books.csv")
	if err != nil {
		return CandleResponse{}, err
	}
	defer csvFile.Close()

	reader := csv.NewReader(csvFile)
	records, err := reader.ReadAll()
	if err != nil {
		return CandleResponse{}, err
	}
	if len(records) < 1 {
		return CandleResponse{}, fmt.Errorf("no data")
	}

	var stockList []Stock

	// データ行を構造体に変換
	for _, record := range records[1:] { // ヘッダー以降の行を処理
		timeValue, err := time.Parse("2006-01-02 15:04:05 -0700 MST", record[0])
		if err != nil {
			fmt.Println("Error parsing time:", err)
			continue
		}

		price, err := strconv.ParseUint(record[2], 10, 64)
		if err != nil {
			fmt.Println("Error parsing price:", err)
			continue
		}

		entry := Stock{
			Time:  timeValue,
			Code:  record[1],
			Price: price,
		}
		stockList = append(stockList, entry) // 構造体をスライスに追加
	}

	// ローソク足（1時間）の算出
	startTime := time.Date(int(year), time.Month(month), int(day), int(hour), 0, 0, 0, time.Local)
	endTime := startTime.Add(time.Hour)
	filtered := filterStockList(stockList, code, startTime, endTime)

	open := filtered[0].Price
	close := filtered[len(filtered)-1].Price

	prices := []uint64{}
	for _, stock := range filtered {
		prices = append(prices, stock.Price)
	}

	// ソート
	sort.Slice(prices, func(i, j int) bool {
		return prices[i] < prices[j]
	})

	high := prices[len(prices)-1]
	low := prices[0]

	// ローソク足の計算
	return CandleResponse{
		Open:  open,
		Close: close,
		High:  high,
		Low:   low,
	}, nil
}

// レコードをフィルタリングする関数
func filterStockList(stockList []Stock, code string, startTime, endTime time.Time) []Stock {
	var filtered []Stock
	for _, stock := range stockList {
		if stock.Code == code && stock.Time.After(startTime) && stock.Time.Before(endTime) {
			filtered = append(filtered, stock)
		}
	}
	return filtered
}
