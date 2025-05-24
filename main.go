// main.go
package main

import (
    "net/http"
    "strings"
)

// Context はリクエストごとのコンテキストを管理
type Context struct {
    Writer  http.ResponseWriter
    Request *http.Request
}

// JSON はレスポンスをJSON形式で返す
func (c *Context) JSON(status int, data interface{}) {
    c.Writer.Header().Set("Content-Type", "application/json")
    c.Writer.WriteHeader(status)
    // 簡易的にJSONを文字列で返す（実際はencoding/jsonを使う）
    c.Writer.Write([]byte(data.(string)))
}

// Framework はカスタムフレームワークの構造体
type Framework struct {
    routes map[string]func(*Context)
}

// NewFramework はフレームワークを初期化
func NewFramework() *Framework {
    return &Framework{
        routes: make(map[string]func(*Context)),
    }
}

// GET はGETリクエストのハンドラを登録
func (f *Framework) GET(path string, handler func(*Context)) {
    f.routes[path] = handler
}

// ServeHTTP はhttp.Handlerインターフェースを実装
func (f *Framework) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    if handler, ok := f.routes[r.URL.Path]; ok && strings.ToUpper(r.Method) == "GET" {
        ctx := &Context{Writer: w, Request: r}
        handler(ctx)
    } else {
        http.NotFound(w, r)
    }
}

// Run はサーバーを起動
func (f *Framework) Run(addr string) error {
    return http.ListenAndServe(addr, f)
}

func main() {
    // フレームワークを初期化
    f := NewFramework()

    // ルートを登録
    f.GET("/hello", func(c *Context) {
        c.JSON(200, `{"message":"Hello, World!"}`)
    })

    // サーバー起動
    if err := f.Run(":8081"); err != nil {
        panic(err)
    }
}