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
		return &Store{}, err
	}
	return &Store{DB: db}, db.Ping()
}

// Преподаватель: создание задачи + привязка языка (languagetoproblem) + компании (problem_to_company) + подсказки (tips)
func (s *Store) CreateFullProblem(title, desc, diff string, timeLimit, memLimit, langID, compID int, hint string) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var problemID int
	err = tx.QueryRow(`
		INSERT INTO problem (title, description, difficulty, time_limit_ms, memory_limit_mb)
		VALUES ($1, $2, $3, $4, $5) RETURNING id`,
		title, desc, diff, timeLimit, memLimit).Scan(&problemID)
	if err != nil {
		return err
	}

	if _, err := tx.Exec(`INSERT INTO languagetoproblem (problem_id, language_id) VALUES ($1, $2)`, problemID, langID); err != nil {
		return err
	}
	if compID > 0 {
		if _, err := tx.Exec(`INSERT INTO problem_to_company (problem_id, company_id) VALUES ($1, $2)`, problemID, compID); err != nil {
			return err
		}
	}
	if hint != "" {
		if _, err := tx.Exec(`INSERT INTO tips (problem_id, hint_text) VALUES ($1, $2)`, problemID, hint); err != nil {
			return err
		}
	}
	return tx.Commit()
}

// Преподаватель: добавление страницы теории (page)
func (s *Store) CreatePage(title, content string) error {
	_, err := s.DB.Exec(`INSERT INTO page (title, content) VALUES ($1, $2)`, title, content)
	return err
}

// Преподаватель: добавление вопроса (testquestion) и вариантов ответа (testanswer)
func (s *Store) CreateQuestionWithAnswer(testID int, questionText, correctAns, wrongAns string) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var qID int
	if err := tx.QueryRow(`INSERT INTO testquestion (test_id, question_text) VALUES ($1, $2) RETURNING id`, testID, questionText).Scan(&qID); err != nil {
		return err
	}
	if _, err := tx.Exec(`INSERT INTO testanswer (question_id, answer_text, is_correct) VALUES ($1, $2, true), ($1, $3, false)`, qID, correctAns, wrongAns); err != nil {
		return err
	}
	return tx.Commit()
}

// Администратор: добавление нового языка программирования (language)
func (s *Store) CreateLanguage(name, version string) error {
	_, err := s.DB.Exec(`INSERT INTO language (name, compiler_version) VALUES ($1, $2)`, name, version)
	return err
}

func (s *Store) GetDashboardData(msg, errStr string) models.DashboardData {
	data := models.DashboardData{Message: msg, Error: errStr}
	if s == nil || s.DB == nil {
		data.Error = "Нет подключения к базе данных PostgreSQL!"
		return data
	}

	// 1. Пользователи (users)
	if rows, err := s.DB.Query(`SELECT id, login, email, role, TO_CHAR(registered_at, 'YYYY-MM-DD HH24:MI') FROM users ORDER BY id`); err == nil {
		defer rows.Close()
		for rows.Next() {
			var u models.UserRow
			rows.Scan(&u.ID, &u.Login, &u.Email, &u.Role, &u.RegisteredAt)
			data.Users = append(data.Users, u)
		}
	}

	// 2. Лог входов (userentry)
	if rows, err := s.DB.Query(`
		SELECT ue.id, u.login, ue.action, TO_CHAR(ue.entry_time, 'YYYY-MM-DD HH24:MI:SS')
		FROM userentry ue JOIN users u ON ue.user_id = u.id ORDER BY ue.id DESC LIMIT 8`); err == nil {
		defer rows.Close()
		for rows.Next() {
			var ue models.UserEntryRow
			rows.Scan(&ue.ID, &ue.UserLogin, &ue.Action, &ue.EntryTime)
			data.UserEntries = append(data.UserEntries, ue)
		}
	}

	// 3. Страницы теории (page)
	if rows, err := s.DB.Query(`SELECT id, title, content, TO_CHAR(created_at, 'YYYY-MM-DD') FROM page ORDER BY id`); err == nil {
		defer rows.Close()
		for rows.Next() {
			var pg models.PageRow
			rows.Scan(&pg.ID, &pg.Title, &pg.Content, &pg.CreatedAt)
			data.Pages = append(data.Pages, pg)
		}
	}

	// 4. Языки (language)
	if rows, err := s.DB.Query(`SELECT id, name, compiler_version FROM language ORDER BY id`); err == nil {
		defer rows.Close()
		for rows.Next() {
			var l models.LanguageRow
			rows.Scan(&l.ID, &l.Name, &l.Version)
			data.Languages = append(data.Languages, l)
		}
	}

	// 5. Задачи (problem + languagetoproblem + problem_to_company + tips)
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

	// 6. Пробные запуски (coderun)
	if rows, err := s.DB.Query(`
		SELECT cr.id, u.login, p.title, l.name, COALESCE(cr.output, ''), COALESCE(cr.exec_time_ms, 0), TO_CHAR(cr.created_at, 'HH24:MI:SS')
		FROM coderun cr
		JOIN users u ON cr.user_id = u.id
		JOIN problem p ON cr.problem_id = p.id
		JOIN language l ON cr.language_id = l.id
		ORDER BY cr.id DESC LIMIT 8`); err == nil {
		defer rows.Close()
		for rows.Next() {
			var cr models.CodeRunRow
			rows.Scan(&cr.ID, &cr.UserLogin, &cr.ProblemTitle, &cr.LangName, &cr.Output, &cr.ExecTime, &cr.CreatedAt)
			data.CodeRuns = append(data.CodeRuns, cr)
		}
	}

	// 7. Отправки решений (codesubmit)
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

	// 8. Вопросы и ответы тестов (test + testquestion + testanswer)
	if rows, err := s.DB.Query(`
		SELECT tq.id, t.title, tq.question_text,
			COALESCE(STRING_AGG(ta.answer_text || CASE WHEN ta.is_correct THEN ' (✔)' ELSE '' END, ' / '), 'Нет ответов')
		FROM testquestion tq
		JOIN test t ON tq.test_id = t.id
		LEFT JOIN testanswer ta ON tq.id = ta.question_id
		GROUP BY tq.id, t.title ORDER BY tq.id`); err == nil {
		defer rows.Close()
		for rows.Next() {
			var q models.TestQuestionRow
			rows.Scan(&q.ID, &q.TestTitle, &q.QuestionText, &q.Answers)
			data.Questions = append(data.Questions, q)
		}
	}

	// 9. Результаты тестов (testresult)
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

	// 10. Статистика для Представителя компании (company + problem_to_company + codesubmit)
	if rows, err := s.DB.Query(`
		SELECT c.name, p.title,
			COUNT(cs.id) AS total_submits,
			COUNT(CASE WHEN cs.status = 'accepted' THEN 1 END) AS accepted_submits
		FROM company c
		JOIN problem_to_company ptc ON c.id = ptc.company_id
		JOIN problem p ON ptc.problem_id = p.id
		LEFT JOIN codesubmit cs ON p.id = cs.problem_id
		GROUP BY c.name, p.title ORDER BY c.name`); err == nil {
		defer rows.Close()
		for rows.Next() {
			var cs models.CompanyStatRow
			rows.Scan(&cs.CompanyName, &cs.ProblemTitle, &cs.TotalSubmits, &cs.AcceptedSubmits)
			data.CompanyStats = append(data.CompanyStats, cs)
		}
	}

	// 11. Транзакции (transaction + transactiontype)
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
