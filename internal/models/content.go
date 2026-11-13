package models

type PageRow struct {
	ID, Title, Content, CreatedAt string
}

type UserEntryRow struct {
	ID, UserLogin, Action, EntryTime string
}

type CodeRunRow struct {
	ID, UserLogin, ProblemTitle, LangName, Output, ExecTime, CreatedAt string
}

type CompanyStatRow struct {
	CompanyName, ProblemTitle, TotalSubmits, AcceptedSubmits string
}

type TestQuestionRow struct {
	ID, TestTitle, QuestionText, Answers string
}

type DashboardData struct {
	Users        []UserRow
	UserEntries  []UserEntryRow
	Pages        []PageRow
	Problems     []ProblemRow
	Languages    []LanguageRow
	CodeRuns     []CodeRunRow
	Submissions  []SubmissionRow
	Questions    []TestQuestionRow
	TestResults  []TestResultRow
	CompanyStats []CompanyStatRow
	Transactions []TransactionRow
	Message      string
	Error        string
}
