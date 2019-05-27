package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Player struct {
	Location *Room
	Inventory map[string][]string
	Actions map[string]func([]string)string
	Door bool
}

type Inside struct {
	Place string
	Things []string
}

type World []Room

type Room struct {
	Name string
	Items []Inside
	Comment string
	CanGo []string
	GoTo string
}

var Slowpoke Player
var Narnia = make([]Room, 4)
var Matching map[string][]string = map[string][]string{"ключи" : {"дверь", "дверь открыта"}}

func main() {
	initGame()
	in := bufio.NewScanner(os.Stdin)
	for in.Scan() {
		fmt.Println(handleCommand(in.Text()))
	}
}

func initGame() {
	Narnia = World{
		{
			Name: "кухня",
			Items: []Inside{
				{
					Place: "на столе",
					Things: []string{"чай"},
				},
			},
			Comment: "надо собрать рюкзак и идти в универ",
			CanGo: []string{"коридор"},
			GoTo: "кухня, ничего интересного. ",
		},
		{
			Name: "коридор",
			Items: []Inside{},
			Comment: "ничего интересного",
			CanGo: []string{"кухня", "комната", "улица"},
			GoTo: "ничего интересного. ",
		},
		{
			Name: "комната",
			Items: []Inside{
				{
					Place: "на столе",
					Things: []string{"ключи", "конспекты"},
				},
				{
					Place: "на стуле",
					Things: []string{"рюкзак"},
				},
			},
			Comment: "",
			CanGo: []string{"коридор"},
			GoTo: "ты в своей комнате. ",
		},
		{
			Name: "улица",
			Items: []Inside{},
			Comment: "на улице весна",
			CanGo: []string{"домой"},
			GoTo: "на улице весна. ",
		},
	}
	Slowpoke = Player{
		Location: &Narnia[0],
		Inventory: map[string][]string{},
		Actions: map[string]func([]string)string {"осмотреться" : lookAround, "идти" : goTo, "взять" : take, "надеть" : wear, "применить" : apply},
	}
}

func apply(param []string) string {
	var fl bool
	var res string
	for _, value := range Slowpoke.Inventory {
		for _, val := range value {
			if val == param[1] {
				fl = true
			}
		}
	}
	if !fl {
		res = "нет предмета в инвентаре - " + param[1]
		return res
	}
	tmp, exist := Matching[param[1]]
	if fl && exist && tmp[0] == param[2] {
		res = tmp[1]
		if param[2] == "дверь" {
			Slowpoke.Door = true
		}
		return res
	} else {
		res = "не к чему применить"
		return res
	}
}

func wear(param []string) string {
	fl := false
	var item string
	for _, val := range Slowpoke.Location.Items {
		for _, item = range val.Things {
			if item == param[1] {
				Slowpoke.Inventory = map[string][]string{param[1] : {}}
				fl = true
			}
		}
	}
	var res string
	if fl {
		deleteItem(param[1])
		res = "вы надели: " + param[1]
	} else {
		res = "нет такого"
	}
	return res
}

func deleteItem(item string) {
	for id, mas := range Slowpoke.Location.Items {
		for idx, val := range mas.Things {
			if strings.EqualFold(val, item) {
				if len(Slowpoke.Location.Items[id].Things) == 1 {
					Slowpoke.Location.Items = append(Slowpoke.Location.Items[:id], Slowpoke.Location.Items[id + 1:]... )
					return
				}
				Slowpoke.Location.Items[id].Things = append(Slowpoke.Location.Items[id].Things[:idx], Slowpoke.Location.Items[id].Things[idx + 1:]... )
			}
		}
	}
}

func take(param []string) string {
	fl := false
	var res string
	var item string
	if len(Slowpoke.Inventory) == 1 {
		for _, val := range Slowpoke.Location.Items {
			for _, item = range val.Things {
				if strings.EqualFold(item, param[1]) {
					Slowpoke.Inventory["рюкзак"] = append(Slowpoke.Inventory["рюкзак"], param[1])
					fl = true
				}
			}
		}
		if fl {
			deleteItem(param[1])
			res = "предмет добавлен в инвентарь: " + param[1]
		} else {
			res = "нет такого"
		}
	} else {
		res = "некуда класть"
	}
	return res
}

func goTo(param []string) string {
	var fl bool = false
	for _, val := range Slowpoke.Location.CanGo {
		if val == param[1] {
			fl = true
		}
	}
	if !fl {
		return "нет пути в " + param[1]
	}
	if param[1] == "улица" {
		if Slowpoke.Door == false {
			return "дверь закрыта"
		}
	}
	for id, val := range Narnia {
		if strings.EqualFold(val.Name, param[1]) {
			Slowpoke.Location = &Narnia[id]
		}
	}
	res := Slowpoke.Location.GoTo
	res += "можно пройти - "
	for _, val := range Slowpoke.Location.CanGo {
		res = res + val + ", "
	}
	res = res[:len(res) - 2]
	return res
}

func lookAround(param []string) string{
	loc := Slowpoke.Location.Name
	var res string
	if strings.EqualFold(loc, "кухня") {
		res = "ты находишься на кухне, "
		if _, exist := Slowpoke.Inventory["рюкзак"]; exist {
			Slowpoke.Location.Comment = "надо идти в универ"
		}
	}
	var fl bool
	if len(Slowpoke.Location.Items) == 0 {
		res += "пустая комната. "
	}
	for _, val := range Slowpoke.Location.Items {
		res = res + val.Place + ": "
		fl = true
		for _, v := range val.Things {
			res = res + v + ", "
		}
	}
	if fl && strings.EqualFold(Slowpoke.Location.Comment, "") {
		res = res[:len(res) - 2]
		res += ". "
	}
	if !strings.EqualFold(Slowpoke.Location.Comment, "") {
		res += Slowpoke.Location.Comment
		res += ". "
	}
	res += "можно пройти - "
	for _, val := range Slowpoke.Location.CanGo {
		res = res + val + ", "
	}
	res = res[:len(res) - 2]
	return res
}

func handleCommand(command string) string {
	com := strings.Split(command, " ")
	if val, exist := Slowpoke.Actions[com[0]]; exist {
		return val(com)
	} else {
		return "неизвестная команда"
	}
}
