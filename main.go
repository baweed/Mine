package main

import (
	"html/template"
	"net/http"
	"strconv"
)

var (
	game *Game
	tmpl *template.Template
)

func main() {
	// Инициализация игры
	game = NewGame(10, 10, 10)

	// Загрузка шаблонов с пользовательской функцией CellWithGame
	var err error
	tmpl, err = template.New("").Funcs(template.FuncMap{
		"CellWithGame": func(cell *Cell, game *Game) CellWithGame {
			return CellWithGame{Cell: cell, Game: game}
		},
	}).ParseGlob("templates/*.html")
	if err != nil {
		panic("Failed to load templates: " + err.Error())
	}

	// Роуты
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/click", handleClick)
	http.HandleFunc("/new-game", handleNewGame)
	http.HandleFunc("/flag", handleFlag)

	// Обработка статических файлов
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Запуск сервера
	http.ListenAndServe(":8080", nil)
}

// handleIndex отображает главную страницу
func handleIndex(w http.ResponseWriter, r *http.Request) {
	err := tmpl.ExecuteTemplate(w, "index.html", game)
	if err != nil {
		http.Error(w, "Failed to render template: "+err.Error(), http.StatusInternalServerError)
		return
	}
}

// handleClick обрабатывает клик по клетке
func handleClick(w http.ResponseWriter, r *http.Request) {
	x, _ := strconv.Atoi(r.URL.Query().Get("x"))
	y, _ := strconv.Atoi(r.URL.Query().Get("y"))

	// Запускаем обработку клика в горутине
	go func() {
		game.OpenCell(x, y)
	}()

	// Возвращаем обновлённое поле
	tmpl.ExecuteTemplate(w, "game.html", game)
}

// handleNewGame начинает новую игру
func handleNewGame(w http.ResponseWriter, r *http.Request) {
	width, _ := strconv.Atoi(r.URL.Query().Get("width"))
	height, _ := strconv.Atoi(r.URL.Query().Get("height"))
	mines, _ := strconv.Atoi(r.URL.Query().Get("mines"))

	game = NewGame(width, height, mines)
	tmpl.ExecuteTemplate(w, "game.html", game)
}

// handleFlag обрабатывает установку флажка
func handleFlag(w http.ResponseWriter, r *http.Request) {
	x, _ := strconv.Atoi(r.URL.Query().Get("x"))
	y, _ := strconv.Atoi(r.URL.Query().Get("y"))

	cell := &game.Grid[y][x]

	// Если клетка уже открыта, ничего не делаем
	if cell.IsOpen {
		tmpl.ExecuteTemplate(w, "game.html", game)
		return
	}

	// Устанавливаем или снимаем флажок
	if cell.IsFlagged {
		cell.IsFlagged = false
		game.FlagsRemaining++ // Увеличиваем количество оставшихся флагов
	} else if game.FlagsRemaining > 0 { // Проверяем, есть ли ещё флаги
		cell.IsFlagged = true
		game.FlagsRemaining-- // Уменьшаем количество оставшихся флагов
	}

	// Возвращаем обновлённое поле
	tmpl.ExecuteTemplate(w, "game.html", game)
}
