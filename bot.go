package main

import (
	"context"
	"fmt"
	tgbotapi "gopkg.in/telegram-bot-api.v4"
	"net/http"
	"strconv"
	"strings"
)

var (
	BotToken   = "804024125:AAGmGaOBur6phuvFC-xoCMkFptqtPIrGAiQ"
	WebhookURL = "https://525f2cb5.ngrok.io"
)

type Task struct {
	Name     string
	Author   string
	Assignee string
}

var Users map[string]int64 = make(map[string]int64)

var Id int = 1

var TaskPool map[int]Task = make(map[int]Task)

func startTaskBot(ctx context.Context) error {
	bot, err := tgbotapi.NewBotAPI(BotToken)
	if err != nil {
		panic(err)
	}
	// bot.Debug = true
	fmt.Printf("Authorized on account %s\n", bot.Self.UserName)
	_, err = bot.SetWebhook(tgbotapi.NewWebhook(WebhookURL))
	if err != nil {
		panic(err)
	}
	updates := bot.ListenForWebhook("/")
	go http.ListenAndServe(":8081", nil)
	for update := range updates {
		UserName := update.Message.From.UserName
		ChatID := update.Message.Chat.ID
		Users[UserName] = ChatID
		Text := update.Message.Text
		txt := strings.Split(Text, " ")
		switch txt[0] {
		case "/tasks":
			reply := ""
			if len(TaskPool) == 0 {
				reply = "Нет задач"
			} else {
				for idx := 1; idx <= Id; idx++ {
					if value, isExists := TaskPool[idx]; isExists {
						reply += strconv.Itoa(idx) + `. ` + value.Name + ` by @` + value.Author
						if value.Assignee != "" {
							reply += "\nassignee: "
							if value.Assignee != UserName {
								reply += "@" + value.Assignee
							} else {
								reply += "я\n/unassign_" + strconv.Itoa(idx) + " /resolve_" + strconv.Itoa(idx)
							}
						} else {
							reply += "\n/assign_" + strconv.Itoa(idx)
						}
						reply += "\n\n"
					}
				}
				reply = reply[:len(reply)-2]
			}
			msg := tgbotapi.NewMessage(ChatID, reply)
			bot.Send(msg)
		case "/new":
			reply := `Задача "` + Text[5:] + `" создана, id=` + strconv.Itoa(Id)
			var t Task
			t.Author = UserName
			t.Name = Text[5:]
			TaskPool[Id] = t
			Id++
			msg := tgbotapi.NewMessage(ChatID, reply)
			bot.Send(msg)
		case "/my":
			fl := false
			reply := ""
			for idx := 1; idx <= Id; idx++ {
				if value, isExists := TaskPool[idx]; isExists {
					if value.Assignee == UserName {
						fl = true
						reply += strconv.Itoa(idx) + `. ` + value.Name + ` by @` + value.Author + "\n"
						reply += "/unassign_" + strconv.Itoa(idx) + " /resolve_" + strconv.Itoa(idx)
						reply += "\n\n"
					}
				}
			}
			if !fl {
				reply = "Нет задач"
			} else {
				reply = reply[:len(reply)-2]
			}
			msg := tgbotapi.NewMessage(ChatID, reply)
			bot.Send(msg)
		case "/owner":
			fl := false
			reply := ""
			for idx := 1; idx <= Id; idx++ {
				if value, isExists := TaskPool[idx]; isExists {
					if value.Author == UserName {
						fl = true
						reply += strconv.Itoa(idx) + `. ` + value.Name + ` by @` + value.Author + "\n"
						reply += "/assign_" + strconv.Itoa(idx)
						reply += "\n\n"
					}
				}
			}
			if !fl {
				reply = "Нет задач"
			} else {
				reply = reply[:len(reply)-2]
			}
			msg := tgbotapi.NewMessage(ChatID, reply)
			bot.Send(msg)
		default:
			com := strings.Split(txt[0], "_")
			switch com[0] {
			case "/assign":
				id, _ := strconv.Atoi(com[1])
				reply := `Задача "` + TaskPool[id].Name + `" назначена на вас`
				msg := tgbotapi.NewMessage(ChatID, reply)
				bot.Send(msg)
				if TaskPool[id].Author != UserName {
					reply = `Задача "` + TaskPool[id].Name + `" назначена на @` + UserName
					if TaskPool[id].Assignee == "" {
						msg = tgbotapi.NewMessage(Users[TaskPool[id].Author], reply)
					} else {
						msg = tgbotapi.NewMessage(Users[TaskPool[id].Assignee], reply)
					}
				}
				var t Task
				t.Assignee = UserName
				t.Author = TaskPool[id].Author
				t.Name = TaskPool[id].Name
				TaskPool[id] = t
				bot.Send(msg)
			case "/unassign":
				id, _ := strconv.Atoi(com[1])
				reply := ""
				fl := false
				if TaskPool[id].Assignee == UserName {
					fl = true
					reply = "Принято"
					var t Task
					t.Assignee = ""
					t.Author = TaskPool[id].Author
					t.Name = TaskPool[id].Name
					TaskPool[id] = t
				} else {
					reply = "Задача не на вас"
				}
				msg := tgbotapi.NewMessage(ChatID, reply)
				bot.Send(msg)
				if fl {
					reply = `Задача "` + TaskPool[id].Name + `" осталась без исполнителя`
					msg = tgbotapi.NewMessage(Users[TaskPool[id].Author], reply)
					bot.Send(msg)
				}
			case "/resolve":
				id, _ := strconv.Atoi(com[1])
				reply := `Задача "` + TaskPool[id].Name + `" выполнена`
				msg := tgbotapi.NewMessage(ChatID, reply)
				bot.Send(msg)
				reply = `Задача "` + TaskPool[id].Name + `" выполнена @` + TaskPool[id].Assignee
				msg = tgbotapi.NewMessage(Users[TaskPool[id].Author], reply)
				bot.Send(msg)
				delete(TaskPool, id)
			}
		}
	}
	return nil
}


func main() {
	err := startTaskBot(context.Background())
	if err != nil {
		panic(err)
	}
}
