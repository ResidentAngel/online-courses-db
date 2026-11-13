package repository

func (s *Store) CreateUserWithEntry(login, hash, email, role string) error {
	tx, err := s.DB.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var userID int
	err = tx.QueryRow(`INSERT INTO users (login, password_hash, email, role) VALUES ($1, $2, $3, $4) RETURNING id`,
		login, hash, email, role).Scan(&userID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(`INSERT INTO userentry (user_id, action) VALUES ($1, 'Регистрация и первый вход')`, userID)
	if err != nil {
		return err
	}
	return tx.Commit()
}
