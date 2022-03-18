package apiserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"webapp/internal/app/model"
	"webapp/internal/app/store"

	"github.com/gorilla/mux"
)

//структура сервера
type APIServer struct {
	config *Config
	router *mux.Router
	store  *store.Store
}

//возвращает структуру сконфигурируимого сервера
func NewServer(c *Config) *APIServer {
	return &APIServer{
		config: c,
		router: mux.NewRouter(),
	}
}

// запуск сервера
func (s *APIServer) Start() error {
	s.ConfigureRouter()

	if err := s.ConfigureStore(); err != nil {
		return err
	}

	defer s.store.Close()

	fmt.Println("Start API server...")
	return http.ListenAndServe(s.config.BindAddress, s.router)
}

//конфигурирование маршрутов
func (s *APIServer) ConfigureRouter() {
	s.router.HandleFunc("/", s.HandleStart())
	s.router.HandleFunc("/todo", s.HandleTodo())
	s.router.HandleFunc("/about", s.HandleAbout())
	s.router.HandleFunc("/session", s.HandleSessionUser()).Methods("POST")
	s.router.HandleFunc("/createuser", s.HandleCreateUser()).Methods("POST")
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

//функция обработки стартового маршрута
func (s *APIServer) HandleStart() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Welcome")
	}
}

//функция обработки маршрута вывода списка
func (s *APIServer) HandleTodo() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "this list")
	}
}

//функция обработки маршрута контакты
func (s *APIServer) HandleAbout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "+79516864862")
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

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(u)
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
