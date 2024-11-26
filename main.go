package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"

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

func main() {
	e := echo.New()

	// ログインエンドポイントの設定
	e.PUT("/login", login)

	// フラグの受け取り口
	e.PUT("/flag", flag)

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
