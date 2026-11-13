package models

type DashboardData struct {
	Users        []UserRow
	Problems     []ProblemRow
	Languages    []LanguageRow
	Submissions  []SubmissionRow
	TestResults  []TestResultRow
	Transactions []TransactionRow
	Message      string
	Error        string
}
