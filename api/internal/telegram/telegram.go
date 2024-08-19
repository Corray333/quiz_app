package telegram

import (
	"database/sql"
	"fmt"
	"log"
	"slices"
	"strconv"

	"github.com/Corray333/quiz/internal/types"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	StateWaitingFIO = iota + 1
	StateWaitingEmail
	StateWaitingDirection
	StateWaitingGroup
)

const (
	StateNothing = iota + 1
	StateWaitingUserTypePick
	StateWaitingMessageText
	StateWaitingMessageAttachment
	StateWaitingSending
)

const (
	MsgWelcome = iota
	MsgWelcomeToQuiz
	MsgError
	MsgFinishQuiz
)

type TelegramClient struct {
	service  Service
	bot      *tgbotapi.BotAPI
	messages map[int]string
}

type Service interface {
	CreateQuiz(quiz *types.Quiz) (int64, error)
	CreateQuestion(question *types.Question) (int64, error)
	GetQuestion(id int64) (*types.Question, error)
	GetUserByTG(id int64) (*types.User, error)
	UpdateUser(user *types.User) error
	CreateUser(user *types.User) (int64, error)
	GetQuiz(id int64) (*types.Quiz, error)

	SetAnswer(answer *types.Answer) (*types.Answer, error)
	GetFirstQuestion(quizID int64) (types.IQuestion, error)
	SetCurrentQuestion(uid, qid int64) error
	GetCurrentQuestion(uid int64) (types.IQuestion, error)
	GetNextQuestion(uid int64) (types.IQuestion, error)
	GetAnswers(userID int64, quizID int64) ([]types.Answer, error)
	GetUserAnswers(userID int64) ([]types.Answer, error)
	GetAnswer(uid, qid int64) (*types.Answer, error)
}

func NewClient(token string, service Service) *TelegramClient {
	bot, err := tgbotapi.NewBotAPI(token)
	if err != nil {
		log.Fatal("Failed to create bot: ", err)
	}

	bot.Debug = true

	return &TelegramClient{
		bot:     bot,
		service: service,
		messages: map[int]string{
			MsgWelcome:       "Добро пожаловать.",
			MsgWelcomeToQuiz: "%s\n\n%s",
			MsgError:         "Что-то пошло не так, попробуйте позже",
			MsgFinishQuiz:    "Поздравляем, квиз пройден!)",
		},
	}
}

func (tg *TelegramClient) Run() {
	defer func() {
		if r := recover(); r != nil {
			tg.HandleError("panic: "+r.(string), 0)
		}
	}()
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := tg.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil && update.CallbackQuery == nil {
			continue
		}

		if update.FromChat().Type != "private" {
			continue
		}

		user, err := tg.service.GetUserByTG(update.FromChat().ID)
		if err != nil {
			fmt.Println("WTG: ", err)
			if err == sql.ErrNoRows {
				if _, err := tg.service.CreateUser(&types.User{
					TgID:     update.FromChat().ID,
					Username: update.FromChat().UserName,
				}); err != nil {
					tg.HandleError("error while creating user: "+err.Error(), update.FromChat().ID, "update_id", update.UpdateID)
					continue
				}
				user, err = tg.service.GetUserByTG(update.FromChat().ID)
				if err != nil {
					tg.HandleError("error while getting user from db: "+err.Error(), update.FromChat().ID, "update_id", update.UpdateID)
					continue
				}
			} else {
				tg.HandleError("error while getting user from db: "+err.Error(), update.FromChat().ID, "update_id", update.UpdateID)
				continue
			}
		}

		switch {
		case user.IsAdmin():
			tg.handleAdminUpdate(user, update)
			continue
		default:
			tg.handleUserUpdate(user, &update)
			continue
		}

	}
}

func (tg *TelegramClient) handleUserUpdate(user *types.User, update *tgbotapi.Update) {

	fmt.Println(user)
	switch {
	case update.Message != nil && update.Message.Command() == "start":
		tg.welcomeToQuiz(user, update)
	case user.CurrentQuestion != 0:
		tg.handleNewAnswer(user, update)
	}

	if err := tg.service.UpdateUser(user); err != nil {
		tg.HandleError("error while updating user: "+err.Error(), update.FromChat().ID, "update_id", update.UpdateID)
	}
}

func (tg *TelegramClient) handleAdminUpdate(user *types.User, update tgbotapi.Update) {

}

func (tg *TelegramClient) welcomeToQuiz(user *types.User, update *tgbotapi.Update) {
	chatID := update.FromChat().ID
	uid := user.ID
	quizID, err := strconv.ParseInt(update.Message.CommandArguments(), 10, 64)
	if err != nil {
		return
	}
	quiz, err := tg.service.GetQuiz(quizID)
	if err != nil {
		if err == sql.ErrNoRows {
			msg := tgbotapi.NewMessage(chatID, tg.messages[MsgError])
			if _, err := tg.bot.Send(msg); err != nil {
				tg.HandleError("error while sending message: "+err.Error(), chatID, "chat_id", chatID)
				return
			}
			return
		}
		tg.HandleError("error while getting quiz: "+err.Error(), chatID, "chat_id", chatID)
		return
	}

	question, err := tg.service.GetFirstQuestion(quizID)
	if err != nil {
		tg.HandleError("error while getting first question: "+err.Error(), chatID, "chat_id", chatID)
		return
	}
	fmt.Println(fmt.Sprintf("User: %+v\n", user))
	if err := tg.service.SetCurrentQuestion(uid, question.GetID()); err != nil {
		tg.HandleError("error while setting current question: "+err.Error(), chatID, "chat_id", chatID)
		return
	}

	var firstQuestion tgbotapi.Chattable

	switch q := question.(type) {
	case *types.QuestionText:
		if q.Image != "" {
			photo := tgbotapi.NewPhoto(chatID, tgbotapi.FileURL(q.Image))
			photo.ParseMode = "Markdown"
			photo.Caption = q.Question
			firstQuestion = photo
		} else {
			msg := tgbotapi.NewMessage(chatID, q.Question)
			msg.ParseMode = "Markdown"
			firstQuestion = msg
		}
	case *types.QuestionSelect:
		if q.Image != "" {
			photo := tgbotapi.NewPhoto(chatID, tgbotapi.FileURL(q.Image))
			photo.ParseMode = "Markdown"
			photo.Caption = q.Question
			markup := tgbotapi.InlineKeyboardMarkup{}
			for _, opt := range q.Options {
				fmt.Println(opt)
				btn := tgbotapi.NewInlineKeyboardButtonData(opt, opt)
				markup.InlineKeyboard = append(markup.InlineKeyboard, []tgbotapi.InlineKeyboardButton{btn})
			}
			photo.ReplyMarkup = &markup
			firstQuestion = photo
		} else {
			msg := tgbotapi.NewMessage(chatID, q.Question)
			msg.ParseMode = "Markdown"
			markup := tgbotapi.InlineKeyboardMarkup{}
			for _, opt := range q.Options {
				fmt.Println(opt)
				btn := tgbotapi.NewInlineKeyboardButtonData(opt, opt)
				markup.InlineKeyboard = append(markup.InlineKeyboard, []tgbotapi.InlineKeyboardButton{btn})
			}
			msg.ReplyMarkup = &markup
			firstQuestion = msg
		}
	case *types.QuestionMultiSelect:
		if q.Image != "" {
			photo := tgbotapi.NewPhoto(chatID, tgbotapi.FileURL(q.Image))
			photo.ParseMode = "Markdown"
			photo.Caption = q.Question
			markup := tgbotapi.InlineKeyboardMarkup{}

			answer, err := tg.service.GetAnswer(uid, question.GetID())
			if err != nil {
				if err == sql.ErrNoRows {
					answer = &types.Answer{}
				} else {
					tg.HandleError("error while getting answer: "+err.Error(), chatID, "chat_id", chatID)
					return
				}
			}

			for _, opt := range q.Options {
				text := opt
				if slices.Contains(answer.Answer, opt) {
					text = "✅" + text
				}
				btn := tgbotapi.NewInlineKeyboardButtonData(text, opt)
				markup.InlineKeyboard = append(markup.InlineKeyboard, []tgbotapi.InlineKeyboardButton{btn})
			}
			markup.InlineKeyboard = append(markup.InlineKeyboard, []tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData("Дать ответ", "Следующий вопрос➡")})
			photo.ReplyMarkup = &markup
			firstQuestion = photo
		} else {
			msg := tgbotapi.NewMessage(chatID, q.Question)
			msg.ParseMode = "Markdown"
			markup := tgbotapi.InlineKeyboardMarkup{}

			answer, err := tg.service.GetAnswer(uid, question.GetID())
			if err != nil {
				if err == sql.ErrNoRows {
					answer = &types.Answer{}
				} else {
					tg.HandleError("error while getting answer: "+err.Error(), chatID, "chat_id", chatID)
					return
				}
			}

			for _, opt := range q.Options {
				text := opt
				if slices.Contains(answer.Answer, opt) {
					text = "✅" + text
				}
				btn := tgbotapi.NewInlineKeyboardButtonData(text, opt)
				markup.InlineKeyboard = append(markup.InlineKeyboard, []tgbotapi.InlineKeyboardButton{btn})
			}

			markup.InlineKeyboard = append(markup.InlineKeyboard, []tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData("Дать ответ", "Следующий вопрос➡")})
			msg.ReplyMarkup = &markup
			firstQuestion = msg
		}
	}

	msgText := fmt.Sprintf(tg.messages[MsgWelcomeToQuiz], quiz.Title, quiz.Description)
	var msg tgbotapi.Chattable
	if quiz.Cover != "" {
		photo := tgbotapi.NewPhoto(chatID, tgbotapi.FileURL(quiz.Cover))
		photo.Caption = msgText
		msg = photo
	} else {
		msg = tgbotapi.NewMessage(chatID, msgText)
	}

	if _, err := tg.bot.Send(msg); err != nil {
		tg.HandleError("error while sending message: "+err.Error(), chatID, "chat_id", chatID)
		return
	}

	if _, err := tg.bot.Send(firstQuestion); err != nil {
		tg.HandleError("error while sending first question: "+err.Error(), chatID, "chat_id", chatID)
		return
	}
}

func (tg *TelegramClient) handleNewAnswer(user *types.User, update *tgbotapi.Update) {
	chatID := update.FromChat().ID

	question, err := tg.service.GetCurrentQuestion(user.ID)
	if err != nil {
		tg.HandleError("error while getting current question: "+err.Error(), chatID, "chat_id", chatID, "update_id", update.UpdateID)
		return
	}
	answer := &types.Answer{}

	if update.CallbackQuery != nil && (question.GetType() == types.QuestionTypeSelect || (question.GetType() == types.QuestionTypeMultiSelect && update.CallbackData() != "Следующий вопрос➡")) {

		cb := tgbotapi.NewCallback(update.CallbackQuery.ID, "")
		if _, err := tg.bot.Request(cb); err != nil {
			tg.HandleError("error while deleting callback: "+err.Error(), chatID, "chat_id", chatID, "update_id", update.UpdateID)
			return
		}

		if question.GetType() == types.QuestionTypeSelect {
			rm := tgbotapi.NewEditMessageReplyMarkup(chatID, update.CallbackQuery.Message.MessageID, tgbotapi.InlineKeyboardMarkup{
				InlineKeyboard: [][]tgbotapi.InlineKeyboardButton{},
			})
			if _, err := tg.bot.Request(rm); err != nil {
				tg.HandleError("error while deleting keyboard: "+err.Error(), chatID, "chat_id", chatID, "update_id", update.UpdateID)
				return
			}
		}

		answer = &types.Answer{
			QuestionID: user.CurrentQuestion,
			UserID:     user.ID,
			Answer:     []string{update.CallbackQuery.Data},
		}
	} else if question.GetType() == types.QuestionTypeText && update.Message != nil {
		fmt.Println("Ответ пользователя: ", update.Message.Text)

		answer = &types.Answer{
			QuestionID: user.CurrentQuestion,
			UserID:     user.ID,
			Answer:     []string{update.Message.Text},
		}

	} else if question.GetType() == types.QuestionTypeFileUpload && (update.Message.Document != nil || update.Message.Photo != nil) {
		// fmt.Println("Ответ пользователя: ", update.Message.Document.FileID)
		// fmt.Println("Правильный ответ: ", question.(*types.QuestionFileUpload).Answer)

		// answer = &types.Answer{
		// 	QuestionID: user.CurrentQuestion,
		// 	UserID:     user.ID,
		// 	Answer:     []string{update.Message.Document.FileID},
		// }

	} else if update.CallbackQuery != nil && update.CallbackData() == "Следующий вопрос➡" {

	} else {
		return
	}

	fmt.Println("Setting answer: ", answer)
	if answer.QuestionID != 0 {
		answer, err = tg.service.SetAnswer(answer)
		if err != nil {
			tg.HandleError("error while setting answer: "+err.Error(), chatID, "chat_id", chatID, "update_id", update.UpdateID)
			return
		}
	}

	if question.GetType() == types.QuestionTypeMultiSelect && update.CallbackQuery != nil && update.CallbackData() != "Следующий вопрос➡" {
		markup := [][]tgbotapi.InlineKeyboardButton{}
		q := question.(*types.QuestionMultiSelect)
		fmt.Println("Multi select answer: ", answer)
		for _, a := range q.Options {
			if slices.Contains(answer.Answer, a) {
				markup = append(markup, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("✅"+a, a)))
			} else {
				markup = append(markup, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData(a, a)))
			}
		}
		markup = append(markup, tgbotapi.NewInlineKeyboardRow(tgbotapi.NewInlineKeyboardButtonData("Дать ответ", "Следующий вопрос➡")))

		edit := tgbotapi.NewEditMessageReplyMarkup(update.FromChat().ID, update.CallbackQuery.Message.MessageID, tgbotapi.InlineKeyboardMarkup{InlineKeyboard: markup})
		if _, err := tg.bot.Request(edit); err != nil {
			tg.HandleError("error while editing message: "+err.Error(), chatID, "chat_id", chatID, "update_id", update.UpdateID)
			return
		}
		return
	}
	if question.GetType() == types.QuestionTypeMultiSelect && update.CallbackQuery != nil && update.CallbackData() == "Следующий вопрос➡" {
		markup := [][]tgbotapi.InlineKeyboardButton{}
		edit := tgbotapi.NewEditMessageReplyMarkup(update.FromChat().ID, update.CallbackQuery.Message.MessageID, tgbotapi.InlineKeyboardMarkup{InlineKeyboard: markup})
		if _, err := tg.bot.Request(edit); err != nil {
			tg.HandleError("error while editing message: "+err.Error(), chatID, "chat_id", chatID, "update_id", update.UpdateID)
			return
		}
	}

	quizID := question.GetQuizID()

	question, err = tg.service.GetNextQuestion(user.ID)
	if err != nil {
		tg.HandleError("error while getting next question: "+err.Error(), chatID, "chat_id", chatID, "update_id", update.UpdateID)
		return
	}
	if question.GetID() == 0 {
		answers, err := tg.service.GetUserAnswers(user.ID)
		if err != nil {
			tg.HandleError("error while getting answers: "+err.Error(), chatID, "chat_id", chatID, "update_id", update.UpdateID)
			return
		}
		quiz, err := tg.service.GetQuiz(quizID)
		if err != nil {
			tg.HandleError("error while getting quiz: "+err.Error(), chatID, "chat_id", chatID, "update_id", update.UpdateID)
			return
		}
		msgText := generateQuizCompletionMessage(answers, quiz.Type)
		msg := tgbotapi.NewMessage(chatID, msgText)
		msg.ParseMode = "Markdown"
		if _, err := tg.bot.Send(msg); err != nil {
			tg.HandleError("error while sending message: "+err.Error(), chatID, "chat_id", chatID, "update_id", update.UpdateID)
			return
		}
		return
	}
	tg.sendQuestion(question, user.ID, chatID)

}

func (tg *TelegramClient) sendQuestion(question types.IQuestion, uid, chatID int64) {
	fmt.Printf("New Question: %+v\n", question)

	if err := tg.service.SetCurrentQuestion(uid, question.GetID()); err != nil {
		tg.HandleError("error while setting current question: "+err.Error(), chatID, "chat_id", chatID)
		return
	}

	var newQuestion tgbotapi.Chattable

	switch q := question.(type) {
	case *types.QuestionText:
		if q.Image != "" {
			photo := tgbotapi.NewPhoto(chatID, tgbotapi.FileURL(q.Image))
			photo.ParseMode = "Markdown"
			photo.Caption = q.Question
			newQuestion = photo
		} else {
			msg := tgbotapi.NewMessage(chatID, q.Question)
			msg.ParseMode = "Markdown"
			newQuestion = msg
		}
	case *types.QuestionSelect:
		if q.Image != "" {
			photo := tgbotapi.NewPhoto(chatID, tgbotapi.FileURL(q.Image))
			photo.ParseMode = "Markdown"
			photo.Caption = q.Question
			markup := tgbotapi.InlineKeyboardMarkup{}
			for _, opt := range q.Options {
				fmt.Println(opt)
				btn := tgbotapi.NewInlineKeyboardButtonData(opt, opt)
				markup.InlineKeyboard = append(markup.InlineKeyboard, []tgbotapi.InlineKeyboardButton{btn})
			}
			photo.ReplyMarkup = &markup
			newQuestion = photo
		} else {
			msg := tgbotapi.NewMessage(chatID, q.Question)
			msg.ParseMode = "Markdown"
			markup := tgbotapi.InlineKeyboardMarkup{}
			for _, opt := range q.Options {
				fmt.Println(opt)
				btn := tgbotapi.NewInlineKeyboardButtonData(opt, opt)
				markup.InlineKeyboard = append(markup.InlineKeyboard, []tgbotapi.InlineKeyboardButton{btn})
			}
			msg.ReplyMarkup = &markup
			newQuestion = msg
		}
	case *types.QuestionMultiSelect:
		if q.Image != "" {
			photo := tgbotapi.NewPhoto(chatID, tgbotapi.FileURL(q.Image))
			photo.ParseMode = "Markdown"
			photo.Caption = q.Question
			markup := tgbotapi.InlineKeyboardMarkup{}

			answer, err := tg.service.GetAnswer(uid, question.GetID())
			if err != nil {
				if err == sql.ErrNoRows {
					answer = &types.Answer{}
				} else {
					tg.HandleError("error while getting answer: "+err.Error(), chatID, "chat_id", chatID)
					return
				}
			}

			for _, opt := range q.Options {
				text := opt
				if slices.Contains(answer.Answer, opt) {
					text = "✅" + text
				}
				btn := tgbotapi.NewInlineKeyboardButtonData(text, opt)
				markup.InlineKeyboard = append(markup.InlineKeyboard, []tgbotapi.InlineKeyboardButton{btn})
			}
			markup.InlineKeyboard = append(markup.InlineKeyboard, []tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData("Дать ответ", "Следующий вопрос➡")})
			photo.ReplyMarkup = &markup
			newQuestion = photo
		} else {
			msg := tgbotapi.NewMessage(chatID, q.Question)
			msg.ParseMode = "Markdown"
			markup := tgbotapi.InlineKeyboardMarkup{}

			answer, err := tg.service.GetAnswer(uid, question.GetID())
			if err != nil {
				if err == sql.ErrNoRows {
					answer = &types.Answer{}
				} else {
					tg.HandleError("error while getting answer: "+err.Error(), chatID, "chat_id", chatID)
					return
				}
			}

			for _, opt := range q.Options {
				text := opt
				if slices.Contains(answer.Answer, opt) {
					text = "✅" + text
				}
				btn := tgbotapi.NewInlineKeyboardButtonData(text, opt)
				markup.InlineKeyboard = append(markup.InlineKeyboard, []tgbotapi.InlineKeyboardButton{btn})
			}

			markup.InlineKeyboard = append(markup.InlineKeyboard, []tgbotapi.InlineKeyboardButton{tgbotapi.NewInlineKeyboardButtonData("Дать ответ", "Следующий вопрос➡")})
			msg.ReplyMarkup = &markup
			newQuestion = msg
		}
	}

	if _, err := tg.bot.Send(newQuestion); err != nil {
		tg.HandleError("error while sending first question: "+err.Error(), chatID, "chat_id", chatID)
		return
	}
}

func generateQuizCompletionMessage(answers []types.Answer, quizType string) string {
	fmt.Println(answers)
	if quizType == "quiz" {
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
	} else {
		return "Спасибо за уделенное время, ваши денные были сохранены)"
	}

}
