package repository

func (s *Store) CreateTransaction(userID, typeID int, amount float64) error {
	_, err := s.DB.Exec(`INSERT INTO "transaction" (user_id, type_id, amount) VALUES ($1, $2, $3)`, userID, typeID, amount)
	return err
}
