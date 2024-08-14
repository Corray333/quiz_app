package main

import (
	"os"

	"github.com/Corray333/quiz/internal/app"
	"github.com/Corray333/quiz/internal/config"
)

func main() {

	config.MustInit(os.Args[1])
	// repo := repository.New()
	// s := service.NewService(repo)
	// t, err := s.GetAllAnswers(19, 0)
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
