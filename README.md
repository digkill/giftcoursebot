# 🎁 GiftCourseBot — Telegram-бот с курсом на 8 дней

Добро пожаловать в проект подарочного Telegram-курса!  
Этот бот автоматически высылает пользователям по одному уроку в день в течение 8 дней.  
Простой, заботливый, с красивым приветствием и системой отзывов 💛

---

## 🚀 Возможности

✅ Автоматическая отправка уроков по дням  
✅ Приветственное сообщение с эмодзи  
✅ Инлайн-отзывы 👍 / 👎  
✅ Сохранение активности в MySQL  
✅ Админ-панель на React + Tailwind  
✅ Экспорт отзывов в CSV  
✅ Демон-режим через systemd  
✅ Поддержка Docker и Docker Compose

---

## 📦 Структура проекта

giftcoursebot/
├── cmd/ # main.go — точка входа
├── internal/
│ ├── components/
│ │ └── db/ # инициализация MySQL, функции работы с уроками
│ ├── handlers/ # обработчики сообщений и отзывов
│ ├── models/ # модели пользователей, уроков, отзывов
│ └── scheduler/ # фоновый шедулер отправки уроков
├── migrations/ # SQL-структура базы данных
├── docker-compose.yml # развёртывание всех компонентов
├── .env.example # пример переменных окружения
└── README.md # ты уже читаешь его


---

## ⚙️ Установка

### 1. Клонируй и настрой `.env`

```bash

cp .env.example .env
nano .env

📡 Системный демон (Linux)
Создай файл /etc/systemd/system/giftcoursebot.service:

ini
Copy
Edit

```bash
[Unit]
Description=Gift Course Telegram Bot
After=network.target

[Service]
ExecStart=/opt/giftcoursebot/giftcoursebot
WorkingDirectory=/opt/giftcoursebot
EnvironmentFile=/opt/giftcoursebot/.env
Restart=always
User=botuser

[Install]
WantedBy=multi-user.target

```

Затем:

```bash
sudo systemctl daemon-reload
sudo systemctl enable giftcoursebot
sudo systemctl start giftcoursebot
```

❤️ Благодарности
Создано с заботой и эмодзи для самых лучших пользователей 🫶
Отправляй знания как подарок — пусть учёба будет радостью ✨