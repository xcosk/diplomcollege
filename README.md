
## Запуск

### Backend

```bash
cd backend
go mod download
GOLEARN_JWT_SECRET=change-me go run ./cmd/server
```

Backend стартует на `http://localhost:8080` и создаст `app.db` в папке `backend`.

### Frontend

```bash
cd frontend
npm install
npm run dev
```

Frontend доступен на `http://localhost:5173`.

## Возможности

- Регистрация/вход, доступ к личному кабинету только с авторизацией.
- JWT access + refresh токены.
- Тест подбора курса из 10 вопросов на главной странице.
- Курсы: базовый, средний, профессиональный.
- Контент по урокам и тест после каждого урока.
- Без успешного прохождения теста следующий урок закрыт.
- Админ‑панель для управления курсами, уроками и тестами.

## Как дать права администратора

Откройте SQLite и выставьте флаг `is_admin=1` нужному пользователю:

```bash
cd backend
sqlite3 app.db "UPDATE users SET is_admin = 1 WHERE email = 'you@example.com';"
```
