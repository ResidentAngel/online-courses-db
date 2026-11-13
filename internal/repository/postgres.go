package repository

import (
	"database/sql"
	"online-courses-db/internal/models"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Store struct {
	DB *sql.DB
}

func NewPostgresDB(dsn string) (*Store, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &Store{DB: db}, nil
}

func (s *Store) GetDashboardData(msg, errStr string) models.DashboardData {
	data := models.DashboardData{Message: msg, Error: errStr}
	if s == nil || s.DB == nil {
		data.Error = "Нет подключения к базе данных PostgreSQL!"
		return data
	}
	// Пользователи
	if rows, err := s.DB.Query(`SELECT id, login, email, role, TO_CHAR(registered_at, 'YYYY-MM-DD HH24:MI') FROM users ORDER BY id`); err == nil {
		defer rows.Close()
		for rows.Next() {
			var u models.UserRow
			rows.Scan(&u.ID, &u.Login, &u.Email, &u.Role, &u.RegisteredAt)
			data.Users = append(data.Users, u)
		}
	}

	// Языки
	if rows, err := s.DB.Query(`SELECT id, name, compiler_version FROM language ORDER BY id`); err == nil {
		defer rows.Close()
		for rows.Next() {
			var l models.LanguageRow
			rows.Scan(&l.ID, &l.Name, &l.Version)
			data.Languages = append(data.Languages, l)
		}
	}

	// Задачи + связи M:N и подсказки
	if rows, err := s.DB.Query(`
		SELECT p.id, p.title, p.difficulty, p.time_limit_ms, p.memory_limit_mb,
			COALESCE(STRING_AGG(DISTINCT l.name || ' (ID:' || l.id || ')', ', '), 'Нет'),
			COALESCE(STRING_AGG(DISTINCT c.name, ', '), 'Общая'),
			COALESCE(STRING_AGG(DISTINCT t.hint_text, ' | '), 'Нет подсказок')
		FROM problem p
		LEFT JOIN languagetoproblem ltp ON p.id = ltp.problem_id
		LEFT JOIN language l ON ltp.language_id = l.id
		LEFT JOIN problem_to_company ptc ON p.id = ptc.problem_id
		LEFT JOIN company c ON ptc.company_id = c.id
		LEFT JOIN tips t ON p.id = t.problem_id
		GROUP BY p.id ORDER BY p.id`); err == nil {
		defer rows.Close()
		for rows.Next() {
			var p models.ProblemRow
			rows.Scan(&p.ID, &p.Title, &p.Difficulty, &p.TimeLimit, &p.MemoryLimit, &p.AllowedLangs, &p.Companies, &p.Hints)
			data.Problems = append(data.Problems, p)
		}
	}

	// Отправки решений (codesubmit)
	if rows, err := s.DB.Query(`
		SELECT cs.id, u.login, p.title, l.name, cs.status, COALESCE(cs.exec_time_ms, 0), TO_CHAR(cs.submitted_at, 'HH24:MI:SS')
		FROM codesubmit cs
		JOIN users u ON cs.user_id = u.id
		JOIN problem p ON cs.problem_id = p.id
		JOIN language l ON cs.language_id = l.id
		ORDER BY cs.id DESC LIMIT 10`); err == nil {
		defer rows.Close()
		for rows.Next() {
			var sub models.SubmissionRow
			rows.Scan(&sub.ID, &sub.UserLogin, &sub.ProblemTitle, &sub.LangName, &sub.Status, &sub.ExecTime, &sub.SubmittedAt)
			data.Submissions = append(data.Submissions, sub)
		}
	}

	// Результаты тестов (testresult)
	if rows, err := s.DB.Query(`
		SELECT tr.id, u.login, t.title, tr.score, tr.duration_sec, TO_CHAR(tr.completed_at, 'YYYY-MM-DD HH24:MI')
		FROM testresult tr
		JOIN users u ON tr.user_id = u.id
		JOIN test t ON tr.test_id = t.id
		ORDER BY tr.id DESC`); err == nil {
		defer rows.Close()
		for rows.Next() {
			var tr models.TestResultRow
			rows.Scan(&tr.ID, &tr.UserLogin, &tr.TestTitle, &tr.Score, &tr.Duration, &tr.CompletedAt)
			data.TestResults = append(data.TestResults, tr)
		}
	}

	// Транзакции (transaction)
	if rows, err := s.DB.Query(`
		SELECT tx.id, u.login, tt.name, tx.amount, TO_CHAR(tx.created_at, 'YYYY-MM-DD HH24:MI')
		FROM "transaction" tx
		JOIN users u ON tx.user_id = u.id
		JOIN transactiontype tt ON tx.type_id = tt.id
		ORDER BY tx.id DESC`); err == nil {
		defer rows.Close()
		for rows.Next() {
			var tx models.TransactionRow
			rows.Scan(&tx.ID, &tx.UserLogin, &tx.TypeName, &tx.Amount, &tx.CreatedAt)
			data.Transactions = append(data.Transactions, tx)
		}
	}

	return data
}
