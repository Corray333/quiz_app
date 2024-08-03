package main

import (
	"encoding/json"
	"os"

	"github.com/Corray333/quiz/internal/app"
	"github.com/Corray333/quiz/internal/config"
)

const (
	QuestionTypeText        = "text"
	QuestionTypeSelect      = "select"
	QuestionTypeMultiSelect = "multi_select"
	QuestionTypeFileUpload  = "file_upload"
)

func main() {
	// var questions []QuestionRaw
	// if err := json.Unmarshal([]byte(`
	// [
	// 		{
	// 			"type": "text",
	// 			"id":1,
	// 			"next":2,
	// 			"question":{
	// 				"question": "Who is the most beautiful girl in the observable universe?",
	// 				"answer": "Maria"
	// 			}
	// 		},
	// 		{
	// 			"type": "select",
	// 			"id":2,
	// 			"next":3,
	// 			"question":{
	// 				"question": "What is the name of the most perspective startup in Russia?",
	// 				"options": [
	// 					"OsaSoft",
	// 					"Otsosoft",
	// 					"Esoft"
	// 				],
	// 				"answer": "OsaSoft"
	// 			}
	// 		},
	// 		{
	// 			"type": "multi_select",
	// 			"id":3,
	// 			"next":4,
	// 			"question":{
	// 				"question": "Which companies are founded in Russia?",
	// 				"options": [
	// 					"OsaSoft",
	// 					"Google",
	// 					"Yandex",
	// 					"Alfa-Bank"
	// 				],
	// 				"answer": [
	// 					"OsaSoft",
	// 					"Yandex",
	// 					"Alfa-Bank"
	// 				]
	// 			}
	// 		}
	// 	]
	// `), &questions); err != nil {
	// 	fmt.Println(err)
	// }

	// for _, questionRaw := range questions {
	// 	switch questionRaw.Type {
	// 	case QuestionTypeText:
	// 		var question QuestionText
	// 		if err := json.Unmarshal([]byte(questionRaw.Question), &question); err != nil {
	// 			fmt.Println(err)
	// 		}
	// 		fmt.Println(question)
	// 	case QuestionTypeSelect:
	// 		var question QuestionSelect
	// 		if err := json.Unmarshal([]byte(questionRaw.Question), &question); err != nil {
	// 			fmt.Println(err)
	// 		}
	// 		fmt.Println(question)
	// 	case QuestionTypeMultiSelect:
	// 		var question QuestionMultiSelect
	// 		if err := json.Unmarshal([]byte(questionRaw.Question), &question); err != nil {
	// 			fmt.Println(err)
	// 		}
	// 		fmt.Println(question)
	// 	}
	// }

	config.MustInit(os.Args[1])
	app.New().Run()

}

type QuestionRaw struct {
	Type     string          `json:"type"`
	ID       int             `json:"id"`
	Next     int             `json:"next"`
	Question json.RawMessage `json:"question"`
}

type QuestionBase struct {
	Type string `json:"type"`
	ID   int    `json:"id"`
	Next int    `json:"next"`
}

type QuestionText struct {
	QuestionBase
	Type     string `json:"type"`
	Question string `json:"question"`
	Answer   string `json:"answer"`
}

type QuestionSelect struct {
	QuestionBase
	Type     string   `json:"type"`
	Question string   `json:"question"`
	Answer   string   `json:"answer"`
	Options  []string `json:"options"`
}

type QuestionMultiSelect struct {
	QuestionBase
	Type     string   `json:"type"`
	Question string   `json:"question"`
	Answer   []string `json:"answer"`
	Options  []string `json:"options"`
}
