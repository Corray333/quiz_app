package main

import (
	"fmt"
	"os"

	"github.com/Corray333/quiz/internal/app"
	"github.com/Corray333/quiz/internal/config"
	"github.com/Corray333/quiz/internal/types"
)

func main() {

	config.MustInit(os.Args[1])
	// repo := repository.New()
	// s := service.NewService(repo)
	// t, err := s.GetQuizAnswers(19, 0)
	// fmt.Println(err)
	// for _, ans := range t {
	// 	fmt.Printf("%+v\n", ans)
	// 	fmt.Println()
	// }
	// fmt.Println(generateQuizCompletionMessage(t))

	// q, _ := repo.GetFirstQuestion(14)
	// switch q.(type) {
	// case *types.QuestionText:
	// 	question := q.(*types.QuestionText)
	// 	question.ID = 100
	// }
	// fmt.Printf("%+v", q)
	app.New().Run()

}

func generateQuizCompletionMessage(answers []types.Answer) string {
	fmt.Println(answers)
	result := "Отлично, квиз пройден. Вот твои результаты:\n\n"
	for _, answer := range answers {

		your := ""
		for _, v := range answer.Answer {
			your += fmt.Sprintf("%s,", v)
		}
		your = your[:len(your)-1]

		result += fmt.Sprintf("Ваш ответ: %s\n", your)
		correct := ""
		for _, v := range answer.Correct {
			correct += fmt.Sprintf("%s,", v)
		}
		correct = correct[:len(correct)-1]
		result += fmt.Sprintf("Правильный ответ: %s\n", correct)
		if answer.IsCorrect {
			result += fmt.Sprintf("Ответ верный✅\n\n")
		} else {
			result += fmt.Sprintf("Ответ неверный❌\n\n")
		}
	}
	return result
}
