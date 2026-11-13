package handler

import (
	"html/template"
	"net/http"
	"strconv"

	"online-courses-db/internal/repository"

	"golang.org/x/crypto/bcrypt"
)

type Handler struct {
	Store *repository.Store
}

func NewRouter(store *repository.Store) *http.ServeMux {
	h := &Handler{Store: store}
	mux := http.NewServeMux()

	mux.HandleFunc("/", h.dashboard)
	mux.HandleFunc("/register", h.register)
	mux.HandleFunc("/code/run", h.codeRun)
	mux.HandleFunc("/code/submit", h.codeSubmit)
	mux.HandleFunc("/test/start", h.startTest)
	mux.HandleFunc("/test/finish", h.finishTest)
	mux.HandleFunc("/problem/create", h.createProblem)
	mux.HandleFunc("/page/create", h.createPage)
	mux.HandleFunc("/question/create", h.createQuestion)
	mux.HandleFunc("/language/create", h.createLanguage)
	mux.HandleFunc("/transaction/create", h.createTransaction)

	return mux
}

func (h *Handler) dashboard(w http.ResponseWriter, r *http.Request) {
	data := h.Store.GetDashboardData(r.URL.Query().Get("msg"), r.URL.Query().Get("err"))
	tmpl := template.Must(template.ParseFiles("web/templates/index.html"))
	_ = tmpl.Execute(w, data)
}

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	hash, _ := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), bcrypt.DefaultCost)
	err := h.Store.CreateUserWithEntry(r.FormValue("login"), string(hash), r.FormValue("email"), r.FormValue("role"))
	if err != nil {
		http.Redirect(w, r, "/?err=Ошибка+регистрации+(проверьте+уникальность+email/login):+"+err.Error(), http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/?msg=Пользователь+зарегистрирован+(пароль+захэширован+bcrypt,+запись+добавлена+в+userentry)", http.StatusSeeOther)
}

func (h *Handler) codeRun(w http.ResponseWriter, r *http.Request) {
	uID, _ := strconv.Atoi(r.FormValue("user_id"))
	pID, _ := strconv.Atoi(r.FormValue("problem_id"))
	lID, _ := strconv.Atoi(r.FormValue("language_id"))
	if err := h.Store.RunCode(uID, pID, lID, r.FormValue("source_code")); err != nil {
		http.Redirect(w, r, "/?err=Нарушение+правила+7.1:+выбранный+язык+не+разрешен+для+данной+задачи+в+languagetoproblem!", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/?msg=Пробный+запуск+(coderun)+успешно+записан", http.StatusSeeOther)
}

func (h *Handler) codeSubmit(w http.ResponseWriter, r *http.Request) {
	uID, _ := strconv.Atoi(r.FormValue("user_id"))
	pID, _ := strconv.Atoi(r.FormValue("problem_id"))
	lID, _ := strconv.Atoi(r.FormValue("language_id"))
	if err := h.Store.SubmitCode(uID, pID, lID, r.FormValue("source_code")); err != nil {
		http.Redirect(w, r, "/?err=Отправка+отклонена+(правило+7.1):+язык+не+разрешен+для+задачи!", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/?msg=Решение+прошло+этапы+queued+->+accepted+в+таблице+codesubmit", http.StatusSeeOther)
}

func (h *Handler) startTest(w http.ResponseWriter, r *http.Request) {
	uID, _ := strconv.Atoi(r.FormValue("user_id"))
	tID, _ := strconv.Atoi(r.FormValue("test_id"))
	if err := h.Store.StartTest(uID, tID); err != nil {
		http.Redirect(w, r, "/?err="+err.Error(), http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/?msg=Старт+теста+зафиксирован+в+таблице+teststart", http.StatusSeeOther)
}

func (h *Handler) finishTest(w http.ResponseWriter, r *http.Request) {
	uID, _ := strconv.Atoi(r.FormValue("user_id"))
	tID, _ := strconv.Atoi(r.FormValue("test_id"))
	score, _ := strconv.Atoi(r.FormValue("score"))
	if err := h.Store.FinishTest(uID, tID, score); err != nil {
		http.Redirect(w, r, "/?err=Нарушение+правила+7.1:+нельзя+сохранить+testresult+без+записи+о+начале+теста+в+teststart!", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/?msg=Результат+теста+сохранен+в+testresult", http.StatusSeeOther)
}

func (h *Handler) createProblem(w http.ResponseWriter, r *http.Request) {
	tl, _ := strconv.Atoi(r.FormValue("time_limit"))
	ml, _ := strconv.Atoi(r.FormValue("memory_limit"))
	lID, _ := strconv.Atoi(r.FormValue("language_id"))
	cID, _ := strconv.Atoi(r.FormValue("company_id"))
	err := h.Store.CreateFullProblem(r.FormValue("title"), r.FormValue("description"), r.FormValue("difficulty"), tl, ml, lID, cID, r.FormValue("hint"))
	if err != nil {
		http.Redirect(w, r, "/?err=Ошибка+создания+задачи+(проверьте+лимит+времени+<=+5000+мс):+"+err.Error(), http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/?msg=Задача,+связь+с+языком,+компанией+и+подсказка+успешно+созданы", http.StatusSeeOther)
}

func (h *Handler) createPage(w http.ResponseWriter, r *http.Request) {
	if err := h.Store.CreatePage(r.FormValue("title"), r.FormValue("content")); err != nil {
		http.Redirect(w, r, "/?err="+err.Error(), http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/?msg=Страница+теории+(page)+добавлена", http.StatusSeeOther)
}

func (h *Handler) createQuestion(w http.ResponseWriter, r *http.Request) {
	tID, _ := strconv.Atoi(r.FormValue("test_id"))
	if err := h.Store.CreateQuestionWithAnswer(tID, r.FormValue("question"), r.FormValue("correct_answer"), r.FormValue("wrong_answer")); err != nil {
		http.Redirect(w, r, "/?err="+err.Error(), http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/?msg=Вопрос+(testquestion)+и+ответы+(testanswer)+добавлены", http.StatusSeeOther)
}

func (h *Handler) createLanguage(w http.ResponseWriter, r *http.Request) {
	if err := h.Store.CreateLanguage(r.FormValue("name"), r.FormValue("version")); err != nil {
		http.Redirect(w, r, "/?err="+err.Error(), http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/?msg=Новый+язык+добавлен+в+справочник+language", http.StatusSeeOther)
}

func (h *Handler) createTransaction(w http.ResponseWriter, r *http.Request) {
	uID, _ := strconv.Atoi(r.FormValue("user_id"))
	tID, _ := strconv.Atoi(r.FormValue("type_id"))
	amount, _ := strconv.ParseFloat(r.FormValue("amount"), 64)
	if err := h.Store.CreateTransaction(uID, tID, amount); err != nil {
		http.Redirect(w, r, "/?err=Нарушение+правила+7.1:+указанный+тип+транзакции+(transactiontype)+не+существует!", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/?msg=Финансовая+операция+записана+в+transaction", http.StatusSeeOther)
}
