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
	mux.HandleFunc("/transaction/create", h.createTransaction)

	return mux
}

func (h *Handler) dashboard(w http.ResponseWriter, r *http.Request) {
	data := h.Store.GetDashboardData(r.URL.Query().Get("msg"), r.URL.Query().Get("err"))
	tmpl := template.Must(template.ParseFiles("web/templates/index.html"))
	tmpl.Execute(w, data)
}

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	hash, _ := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), bcrypt.DefaultCost)
	err := h.Store.CreateUserWithEntry(r.FormValue("login"), string(hash), r.FormValue("email"), r.FormValue("role"))
	if err != nil {
		http.Redirect(w, r, "/?err=Ошибка+регистрации:+"+err.Error(), http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/?msg=Пользователь+успешно+зарегистрирован", http.StatusSeeOther)
}

func (h *Handler) codeRun(w http.ResponseWriter, r *http.Request) {
	uID, _ := strconv.Atoi(r.FormValue("user_id"))
	pID, _ := strconv.Atoi(r.FormValue("problem_id"))
	lID, _ := strconv.Atoi(r.FormValue("language_id"))
	if err := h.Store.RunCode(uID, pID, lID, r.FormValue("source_code")); err != nil {
		http.Redirect(w, r, "/?err=Нарушение+правила+7.1:+язык+не+разрешен+для+задачи!", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/?msg=Пробный+запуск+(coderun)+выполнен", http.StatusSeeOther)
}

func (h *Handler) codeSubmit(w http.ResponseWriter, r *http.Request) {
	uID, _ := strconv.Atoi(r.FormValue("user_id"))
	pID, _ := strconv.Atoi(r.FormValue("problem_id"))
	lID, _ := strconv.Atoi(r.FormValue("language_id"))
	if err := h.Store.SubmitCode(uID, pID, lID, r.FormValue("source_code")); err != nil {
		http.Redirect(w, r, "/?err=Отправка+отклонена+(правило+7.1):+язык+не+поддерживается!", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/?msg=Решение+принято+(status:+accepted)", http.StatusSeeOther)
}

func (h *Handler) startTest(w http.ResponseWriter, r *http.Request) {
	uID, _ := strconv.Atoi(r.FormValue("user_id"))
	tID, _ := strconv.Atoi(r.FormValue("test_id"))
	if err := h.Store.StartTest(uID, tID); err != nil {
		http.Redirect(w, r, "/?err="+err.Error(), http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/?msg=Попытка+зафиксирована+в+teststart", http.StatusSeeOther)
}

func (h *Handler) finishTest(w http.ResponseWriter, r *http.Request) {
	uID, _ := strconv.Atoi(r.FormValue("user_id"))
	tID, _ := strconv.Atoi(r.FormValue("test_id"))
	score, _ := strconv.Atoi(r.FormValue("score"))
	if err := h.Store.FinishTest(uID, tID, score); err != nil {
		http.Redirect(w, r, "/?err=Нарушение+правила+7.1:+сначала+нажмите+Начать+тест+(teststart)!", http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/?msg=Результат+теста+сохранен+в+testresult", http.StatusSeeOther)
}

func (h *Handler) createTransaction(w http.ResponseWriter, r *http.Request) {
	uID, _ := strconv.Atoi(r.FormValue("user_id"))
	tID, _ := strconv.Atoi(r.FormValue("type_id"))
	amount, _ := strconv.ParseFloat(r.FormValue("amount"), 64)
	if err := h.Store.CreateTransaction(uID, tID, amount); err != nil {
		http.Redirect(w, r, "/?err=Ошибка+транзакции:+"+err.Error(), http.StatusSeeOther)
		return
	}
	http.Redirect(w, r, "/?msg=Транзакция+проведена", http.StatusSeeOther)
}
