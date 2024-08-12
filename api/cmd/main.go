package main

import (
	"os"

	"github.com/Corray333/quiz/internal/app"
	"github.com/Corray333/quiz/internal/config"
)

func main() {

	config.MustInit(os.Args[1])
	// repo := repository.New()
	// t, _ := repo.GetQuizAnswers(1, 16)
	// fmt.Printf("%+v\n", t)
	// fmt.Println(repo.SetAnswer(
	// 	&types.Answer{
	// 		UserID:     1,
	// 		QuestionID: 13,
	// 		Answer:     []string{"Test", "Jopa"},
	// 	},
	// ))

	// q, _ := repo.GetFirstQuestion(14)
	// switch q.(type) {
	// case *types.QuestionText:
	// 	question := q.(*types.QuestionText)
	// 	question.ID = 100
	// }
	// fmt.Printf("%+v", q)
	app.New().Run()

}
