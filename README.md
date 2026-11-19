# Образовательная платформа «Онлайн-Курсы» PostgreSQL + Go

Веб-приложение на языке Go для работы с базой данных образовательной платформы 18 таблиц, СУБД PostgreSQL.

---

## 1. Подключение к базе данных

Параметры подключения к PostgreSQL заданы по умолчанию в `internal/config/config.go` и `docker-compose.yml`:

* **Хост:** `localhost`
* **Порт:** `5433` *внешний порт Docker-контейнера; внутри контейнера — `5432`*
* **Имя базы данных:** `online_courses`
* **Пользователь:** `postgres`
* **Пароль:** `postgres`
* **Строка подключения (DSN):**
  ```text
  postgres://postgres:postgres@localhost:5433/online_courses?sslmode=disable
  ```

При необходимости строку подключения можно переопределить через переменную окружения `DATABASE_URL` (пример находится в `.env.example`).

Запуск контейнера с PostgreSQL:
```bash
docker compose up -d
```

---

## 2. Создание таблиц

Схема базы данных 18 таблиц, первичные и внешние ключи, ограничения `CHECK` и `UNIQUE` находится в файле `migrations/001_init_schema.up.sql`, а начальные тестовые данные — в файле `migrations/002_seed_data.sql`.

### Автоматически через Docker
При первом запуске `docker compose up -d` оба SQL-файла из папки `migrations/` выполняются автоматически внутри контейнера `online_courses_pg`.

### Вручную через `psql`
1. Создать таблицы `migrations/001_init_schema.up.sql`:
   ```bash
   psql -h localhost -p 5433 -U postgres -d online_courses -f migrations/001_init_schema.up.sql
   ```
2. Заполнить таблицы тестовыми данными `migrations/002_seed_data.sql`:
   ```bash
   psql -h localhost -p 5433 -U postgres -d online_courses -f migrations/002_seed_data.sql
   ```

Если локальный клиент `psql` не установлен, выполните скрипты внутри Docker-контейнера:
```bash
docker exec -i online_courses_pg psql -U postgres -d online_courses < migrations/001_init_schema.up.sql
docker exec -i online_courses_pg psql -U postgres -d online_courses < migrations/002_seed_data.sql
```

Для удаления всех таблиц предусмотрен скрипт `migrations/001_init_schema.down.sql`:
```bash
psql -h localhost -p 5433 -U postgres -d online_courses -f migrations/001_init_schema.down.sql
```

---

## 3. Запуск приложения

1. Убедитесь, что контейнер PostgreSQL запущен:
   ```bash
   docker compose up -d
   ```
2. Установите зависимости Go из корня репозитория:
   ```bash
   go mod tidy
   ```
3. Запустите веб-сервер:
   ```bash
   go run ./cmd/app/main.go
   ```
4. Откройте приложение в браузере:
   ```text
   http://localhost:8080
   ```
