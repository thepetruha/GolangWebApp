package apiserver

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"webapp/internal/app/model"
	"webapp/internal/app/store"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

const (
	sessionName        = "webappgoalng"
	ctxKeyUser  ctxKey = iota
)

type ctxKey int8

//структура сервера
type APIServer struct {
	config       *Config
	router       *mux.Router
	store        *store.Store
	sessionStore sessions.Store
}

//возвращает структуру сконфигурируимого сервера
func NewServer(c *Config, sessionStore sessions.Store) *APIServer {
	return &APIServer{
		config:       c,
		router:       mux.NewRouter(),
		sessionStore: sessionStore,
	}
}

// запуск сервера
func (s *APIServer) Start() error {
	s.ConfigureRouter()

	if err := s.ConfigureStore(); err != nil {
		fmt.Println(err)
		return err
	}

	defer s.store.Close()

	fmt.Println("Start API server...")
	return http.ListenAndServe(s.config.BindAddress, s.router)
}

//конфигурирование маршрутов
func (s *APIServer) ConfigureRouter() {
	s.router.HandleFunc("/session", s.HandleSessionUser()).Methods("POST")
	s.router.HandleFunc("/createuser", s.HandleCreateUser()).Methods("POST")

	private := s.router.PathPrefix("/private").Subrouter()
	private.Use(s.authenticateUser)
	private.HandleFunc("/whoami", s.HandleWhoami()).Methods("GET")
	private.HandleFunc("/createtodo", s.HandleCreateRowTodo()).Methods("POST")
	private.HandleFunc("/listtodo", s.HandleGetListTodo()).Methods("GET")
}

//конфигурирует окрытие соеденения с бд
func (s *APIServer) ConfigureStore() error {
	st := store.NewStore(s.config.Store)
	if err := st.Open(); err != nil {
		return err
	}

	s.store = st

	return nil
}

//middleware для входа на приватные страницы
func (s *APIServer) authenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.StatusError(w, r, http.StatusInternalServerError, "Could not find session")
			return
		}

		id, ok := session.Values["user_id"]
		if !ok {
			s.StatusError(w, r, http.StatusUnauthorized, "Error: not authenticated")
			return
		}

		u, err := s.store.User().FindUserId(id.(int))
		if err != nil {
			s.StatusError(w, r, http.StatusUnauthorized, "Error: not authenticated")
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, u)))
	})
}

func (s *APIServer) HandleGetListTodo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		customer_id := r.Context().Value(ctxKeyUser).(*model.User).ID
		rows, err := s.store.Todo().FindByUserId(customer_id)
		if err != nil {
			return
		}

		fmt.Println(rows)

		s.response(w, r, http.StatusOK, rows)
	}
}

func (s *APIServer) HandleCreateRowTodo() http.HandlerFunc {
	type requset struct {
		Text string `json:"text"`
		Date string `json:"date"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &requset{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.StatusError(w, r, http.StatusUnauthorized, "Incorrect format data")
			return
		}

		customer_id := r.Context().Value(ctxKeyUser).(*model.User).ID

		t := &model.ToDo{
			CustomerID: customer_id,
			Text:       req.Text,
			Date:       req.Date,
		}

		if err := s.store.Todo().CreateRow(t); err != nil {
			s.StatusError(w, r, http.StatusUnprocessableEntity, "Error created on user")
			return
		}

		s.response(w, r, http.StatusOK, t)
	}
}

//функция приватной страницы /private/whoami
func (s *APIServer) HandleWhoami() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.response(w, r, http.StatusOK, r.Context().Value(ctxKeyUser).(*model.User))
	}
}

//функция регистрации пользователя
func (s *APIServer) HandleCreateUser() func(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	//получение данных от клиента, дабавление в БД и ответ клиенту
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}

		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.StatusError(w, r, http.StatusUnauthorized, "Incorrect format data")
			return
		}

		u := &model.User{
			Email:    req.Email,
			Password: req.Password,
		}

		if err := s.store.User().Create(u); err != nil {
			s.StatusError(w, r, http.StatusUnprocessableEntity, "Error created on user")
			return
		}

		u.Snitized()
		s.response(w, r, http.StatusCreated, u)
	}
}

//авторизация пользователя
func (s *APIServer) HandleSessionUser() func(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	//получение данных от клиента, поиск клиента в БД, проверка пароля и ответ клиенту
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.StatusError(w, r, http.StatusUnauthorized, "Incorrect format data")
			return
		}

		u, err := s.store.User().FindByEmail(req.Email)
		if err != nil || !u.ComparePassword(req.Password) {
			s.StatusError(w, r, http.StatusUnauthorized, "Incorrect email of password")
			return
		}

		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.StatusError(w, r, http.StatusInternalServerError, "Could not find session")
			return
		}

		session.Values["user_id"] = u.ID
		if err := s.sessionStore.Save(r, w, session); err != nil {
			s.StatusError(w, r, http.StatusInternalServerError, "Could not find session")
			return
		}

		s.response(w, r, http.StatusOK, u)
	}
}

//функция ответа клиенту
func (s *APIServer) response(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	if data != nil {
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(data)
	}
}

//обработка ошибка на роутах
func (s *APIServer) StatusError(w http.ResponseWriter, r *http.Request, code int, err string) {
	s.response(w, r, code, map[string]string{"error": err})
}
