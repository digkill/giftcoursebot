package models

import "database/sql"

type Lesson struct {
	ID        int64          `db:"id"`
	Title     string         `db:"title"`
	Content   string         `db:"content"`
	CreatedAt sql.NullString `db:"created_at"`
	UpdatedAt sql.NullString `db:"updated_at"`
}

var lessons = []string{
	"Урок 1: Введение",
	"Урок 2: Основы",
	"Урок 3: Продвинутые концепции",
	"Урок 4: Применение на практике",
	"Урок 5: Углубление",
	"Урок 6: Примеры из жизни",
	"Урок 7: Проверка знаний",
	"Урок 8: Заключение",
}
