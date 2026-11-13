package repository

import "time"

func (s *Store) StartTest(userID, testID int) error {
	_, err := s.DB.Exec(`
		INSERT INTO teststart (user_id, test_id, started_at) VALUES ($1, $2, $3)
		ON CONFLICT (user_id, test_id) DO UPDATE SET started_at = EXCLUDED.started_at`, userID, testID, time.Now())
	return err
}

func (s *Store) FinishTest(userID, testID, score int) error {
	_, err := s.DB.Exec(`INSERT INTO testresult (user_id, test_id, score, duration_sec) VALUES ($1, $2, $3, 180)`, userID, testID, score)
	return err
}
