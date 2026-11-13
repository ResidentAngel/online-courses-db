package models

type ProblemRow struct {
	ID, Title, Difficulty, TimeLimit, MemoryLimit, AllowedLangs, Companies, Hints string
}

type LanguageRow struct {
	ID, Name, Version string
}
