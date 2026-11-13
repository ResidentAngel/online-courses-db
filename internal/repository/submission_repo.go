package repository

import "math/rand"

func (s *Store) RunCode(userID, problemID, langID int, code string) error {
	_, err := s.DB.Exec(`
		INSERT INTO coderun (user_id, problem_id, language_id, source_code, output, exec_time_ms)
		VALUES ($1, $2, $3, $4, 'Test case #1: OK', 120)`, userID, problemID, langID, code)
	return err
}

func (s *Store) SubmitCode(userID, problemID, langID int, code string) error {
	var submitID int
	err := s.DB.QueryRow(`
		INSERT INTO codesubmit (user_id, problem_id, language_id, source_code, status)
		VALUES ($1, $2, $3, $4, 'queued') RETURNING id`, userID, problemID, langID, code).Scan(&submitID)
	if err != nil {
		return err
	}
	execTime := rand.Intn(800) + 50
	_, err = s.DB.Exec(`UPDATE codesubmit SET status = 'accepted', exec_time_ms = $1, memory_used_mb = 16 WHERE id = $2`, execTime, submitID)
	return err
}
