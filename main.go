package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"time"
)

const (
	expectedToken = "4321"
	updateURL     = "http://localhost:8000/request/result/"
)

type HRResult struct {
	OrderID    string `json:"order_id"`
	Is_success int    `json:"is_success"`
	Token      string `json:"token"`
}

func main() {
	http.HandleFunc("/calc", handleProcess)
	fmt.Println("Server running at port :8080")
	http.ListenAndServe(":8080", nil)
}

func handleProcess(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не разрешен", http.StatusMethodNotAllowed)
		return
	}

	orderid := r.FormValue("order_id")
	token := r.FormValue("token")
	fmt.Println(orderid, token)

	if token == "" || token != expectedToken {
		http.Error(w, "Токены не совпадают", http.StatusForbidden)
		fmt.Println("Токены не совпадают")
		fmt.Println(token, expectedToken)
		return
	}

	w.WriteHeader(http.StatusOK)

	go func() {
		delay := 10
		time.Sleep(time.Duration(delay) * time.Second)

		result := rand.Intn(101)

		// Отправка результата на другой сервер
		expResult := HRResult{
			OrderID:    orderid,
			Is_success: result,
			// Token:  token,
		}
		fmt.Println("json", expResult)
		jsonValue, err := json.Marshal(expResult)
		if err != nil {
			fmt.Println("Ошибка при маршализации JSON:", err)
			return
		}

		req, err := http.NewRequest(http.MethodPut, updateURL, bytes.NewBuffer(jsonValue))
		if err != nil {
			fmt.Println("Ошибка при создании запроса на обновление:", err)
			return
		}
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		order, err := client.Do(req)
		if err != nil {
			fmt.Println("Ошибка при отправке запроса на обновление:", err)
			return
		}
		defer order.Body.Close()

		fmt.Println("Ответ от сервера обновления:", order.Status)
	}()
}
