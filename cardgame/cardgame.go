package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

type drawResp struct {
	Success bool   `json:"success"`
	DeckID  string `json:"deck_id"`
	Cards   []struct {
		Value string `json:"value"`
	} `json:"cards"`
	Remaining int `json:"remaining"`
}

func main() {
	// Тут я юзаю json.NewDecoder(resp.Body).Decode(...). Почему? NewDecoder(resp.Body) сразу
	// создаёт декодер на основе буффера, и декод матчит респонз в структуру. Для анмаршалла нужен промежуточный буффер:
	// https://pkg.go.dev/encoding/json#Decoder.Decode
	// https://pkg.go.dev/encoding/json#Decoder

	in := bufio.NewReader(os.Stdin)
	var guess int
	fmt.Print("Введите позицию (1..52): ")
	_, err := fmt.Fscan(in, &guess)

	resp, err := http.Get("https://deckofcardsapi.com/api/deck/new/shuffle/?deck_count=1")
	if err != nil {
		fmt.Println("net err:", err)
		return
	}
	defer resp.Body.Close()

	var first drawResp
	if err := json.NewDecoder(resp.Body).Decode(&first); err != nil || !first.Success {
		fmt.Println("api err: shuffle")
		return
	}
	deckID := first.DeckID

	var queens []int
	for i := 1; i <= 52; i++ {
		url := "https://deckofcardsapi.com/api/deck/" + deckID + "/draw/?count=1"
		r, err := http.Get(url)
		if err != nil {
			fmt.Println("net err:", err)
			return
		}
		var d drawResp
		if err := json.NewDecoder(r.Body).Decode(&d); err != nil || !d.Success || len(d.Cards) != 1 {
			r.Body.Close()
			fmt.Println("api err: draw")
			return
		}
		fmt.Print(d)
		r.Body.Close()
		if d.Cards[0].Value == "QUEEN" {
			queens = append(queens, i)
		}
	}

	hit := false
	for _, q := range queens {
		if q == guess {
			hit = true
			break
		}
	}

	if hit {
		fmt.Printf("ДА! На позиции %d — Дама.\n", guess)
	} else {
		fmt.Printf("Нет. На позиции %d Дамы нет.\n", guess)
	}
	fmt.Printf("Позиции Дам: %v\n", queens)
}
