package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

var (
	secondBrainRoot   string
	todoListFile      string
	fullTodoListFile  string
	today             string
	dailyScheduleFile string
	prePrompt         string
)

func initConfig() {
	secondBrainRoot = os.Getenv("SECOND_BRAIN_ROOT")
	todoListFile = os.Getenv("TODO_LIST_FILE")
	fullTodoListFile = fmt.Sprintf("%s/%s", secondBrainRoot, todoListFile)
	today = time.Now().Format("2006-01-02")
	dailyScheduleFile = fmt.Sprintf("%s/%s.org", secondBrainRoot, today)
	prePrompt = fmt.Sprintf(`Given the following org-mode to-do list, create a structured daily org-mode formatted TODO-schedule for today (%s). Consider what date it is an any potential deadlines or just other contextual aspects of the date. Also be creative if you come up with something important to add to achieve the tasks in the list. I wont be able to make everything in one day. Prioritize what is important for the current day! Also, DO FOLLOW THESE RULES OR FEEL MY WRATH:
			* If you write a task for example: "Write a message" then provide an example for that message. Apply this thinking on all tasks you write.
			* Each title row, not each checkbox item, in the schedule should start with TODO to indicate actionable items.
			* The subtasks for a title row should be checkboxes.
			* Make each main task scheduled according to org mode standard, both date and timespan like this under the title row: 'SCHEDULED: <2015-02-20 Fri 15:15>'.
			* When writing a task for a meal provide an easy quick healthy vegetarian recipe and a shopping list.
			* If you see that some checkboxes are checked ('- [X]' looks like that instead of '- [ ]') in the provided to-do list, then please exclude them in the schedule.
			* The schedule should effectively balance work, personal tasks, and projects, including time for meals and breaks.
			* Return ONLY (!!!) the org-mode list. NOTHING else and NOTHING FROM the provided org-mode to-do list. Here is the org-mode to-do list to generate the org-mode formatted daily schedule from: `, today)
}

func readTodoList() string {
	fmt.Println("Reading org-mode to-do list from ", fullTodoListFile)
	content, err := os.ReadFile(fullTodoListFile)
	if err != nil {
		log.Fatalf("Could not read %s because: %v", fullTodoListFile, err)
	}

	return string(content[:])
}

func addPrePrompt(content string) string {
	fmt.Println("Adding pre-prompt to org-mode to-do list.")
	return prePrompt + content
}

func writeDailySchedule(content string) {
	fmt.Println("Writing daily schedule to", dailyScheduleFile)
	err := os.WriteFile(dailyScheduleFile, []byte(content), 0644)
	if err != nil {
		log.Fatalf("Could not write %s because: %v", dailyScheduleFile, err)
	}
}

type OpenAIRequestBody struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func callOpenAI(model, prompt string, temperature float64) (string, error) {
	url := "https://api.openai.com/v1/chat/completions"

	requestBody := OpenAIRequestBody{
		Model: model,
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Temperature: temperature,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("OPENAI_API_KEY"))

	fmt.Println("Calling OpenAI API...")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Read the response using bufio and os
	reader := bufio.NewReader(resp.Body)
	var response bytes.Buffer
	_, err = io.Copy(&response, reader)
	if err != nil {
		return "", err
	}

	return response.String(), nil
}

type OpenAIResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func extractContent(jsonData []byte) (string, error) {
	var response OpenAIResponse
	err := json.Unmarshal(jsonData, &response)
	if err != nil {
		return "", err
	}

	if len(response.Choices) > 0 && response.Choices[0].Message.Content != "" {
		return response.Choices[0].Message.Content, nil
	}

	return "", fmt.Errorf("no content found in the response")
}

func main() {
	err := godotenv.Load("/Users/olofjondelius/code/second-brain-updater/.env") // TODO: Fix this (Need absolute path for the cronjob to find it - use other solution)
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}

	initConfig()

	model := "gpt-4"
	prompt := addPrePrompt(readTodoList())
	temperature := 0.7

	response, err := callOpenAI(model, prompt, temperature)
	if err != nil {
		log.Fatalf("Error calling OpenAI API: %s", err)
	}

	jsonData := []byte(response)

	content, err := extractContent(jsonData)
	if err != nil {
		log.Fatalf("Error extracting content: %s", err)
	}

	writeDailySchedule(content)
}
