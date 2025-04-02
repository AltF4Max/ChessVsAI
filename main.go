package main

import (
	"ChessVsAI/chessGame"
	"ChessVsAI/config"
	"ChessVsAI/models"
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/go-resty/resty/v2"
	"github.com/notnil/chess"
)

func main() {
	game := chess.NewGame()
	client := resty.New()
	requestBody := map[string]interface{}{
		"model": "deepseek/deepseek-r1:free",
		"messages": []map[string]string{
			//{"role": "user", "content": "Let's play chess, my move e4. Write only your move. Example: e4."},
			{"role": "user", "content": "Let's play chess, your move first. Write only your move. Example: e4."},
		},
	}
	for {
		newMessage := map[string]string{}
		resp, err := client.R().
			SetHeader("Authorization", "Bearer "+config.ApiKey).
			SetHeader("Content-Type", "application/json").
			SetHeader("HTTP-Referer", "YOUR_SITE_URL").
			SetHeader("X-Title", "YOUR_SITE_NAME").
			SetBody(requestBody).
			Post(config.Url)

		if err != nil {
			log.Fatalf("Request error: %v", err)
		}
		var result models.APIResponse
		if err := json.Unmarshal(resp.Body(), &result); err != nil {
			log.Fatalf("Parsing error JSON: %v", err)
		}
		if len(result.Choices) == 0 {
			fmt.Println("200 requests per day ended.")
			return
		}
		if result.Choices[0].FinishReason == "stop" && result.Choices[0].NativeFinishReason == "stop" {
			answerAI := strings.Fields(result.Choices[0].Message.Content)
			re := regexp.MustCompile(`[*.+]`)
			reAnswerAI := re.ReplaceAllString(answerAI[0], "")
			gameOver, err := chessGame.PlayChess(game, reAnswerAI)
			if err != nil {
				newMessage = map[string]string{"role": "user", "content": "move " + reAnswerAI + " is not possible."}
				requestBody["messages"] = append(requestBody["messages"].([]map[string]string), newMessage)
				continue
			}
			if gameOver == true {
				fmt.Println("Game over! Result: ", game.Outcome())
				return
			}
		} else { //Answer is empty
			continue
		}
		newMessage = map[string]string{"role": "assistant", "content": (result.Choices[0].Message.Content)}
		requestBody["messages"] = append(requestBody["messages"].([]map[string]string), newMessage)
		for {
			fmt.Println(game.Position().Board().Draw())
			scanner := bufio.NewScanner(os.Stdin)
			fmt.Print("Enter your move: ")

			if scanner.Scan() {
				gameOver, err := chessGame.PlayChess(game, scanner.Text())
				if err != nil {
					fmt.Printf("Wrong move: %s.", scanner.Text())
					continue
				}
				if gameOver == true {
					fmt.Println("Game over! Result: ", game.Outcome())
					return
				}
				newMessage = map[string]string{"role": "user", "content": (scanner.Text())}
				requestBody["messages"] = append(requestBody["messages"].([]map[string]string), newMessage)
				break
			}

			if err := scanner.Err(); err != nil {
				fmt.Println("Input error: ", err)
				continue
			}
		}
	}
}
