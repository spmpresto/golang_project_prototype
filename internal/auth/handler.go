package auth

import (
	"golang/advanced/configs"
	"golang/advanced/pkg/jwt"
	"golang/advanced/pkg/req"
	"golang/advanced/pkg/res"
	"net/http"
)

type AuthHandlerDepth struct {
	*configs.Config
	*AuthService
}

type AuthHandler struct {
	*configs.Config
	*AuthService
}

func NewAuthHandler(router *http.ServeMux, deps AuthHandlerDepth) {
	handler := &AuthHandler{
		Config:      deps.Config,
		AuthService: deps.AuthService,
	}
	router.HandleFunc("POST /auth/login", handler.Login())
	router.HandleFunc("POST /auth/register", handler.Register())
}

func (handler *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		// read body
		//var payload LoginRequest
		//err := json.NewDecoder(req.Body).Decode(&payload)
		//if err != nil {
		//	res.Json(w, err.Error(), http.StatusBadRequest)
		//	return
		//}

		//if payload.Email == "" {
		//	res.Json(w, "email is required", http.StatusBadRequest)
		//	return
		//}
		////match, _ := regexp.MatchString(`[A-Za-z0-9\._\-]+@[a-zA-Z0-9\-]+\.[a-zA-Z]{2,}`, payload.Email)
		//_, err = mail.ParseAddress(payload.Email)
		//if err != nil {
		//	res.Json(w, "email is invalid", http.StatusBadRequest)
		//	return
		//}
		//if payload.Password == "" {
		//	res.Json(w, "password is required", http.StatusBadRequest)
		//	return
		//}
		//fmt.Println(payload)
		////fmt.Println(handler.Config.Auth.Secret)

		//validate := validator.New()
		//err = validate.Struct(payload)
		//if err != nil {
		//	res.Json(w, err.Error(), http.StatusBadRequest)
		//	return
		//}

		body, err := req.HandleBody[LoginRequest](&w, r)
		if err != nil {
			return
		}
		email, err := handler.AuthService.Login(body.Email, body.Password)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		token, err := jwt.NewJWT(handler.Config.Auth.Secret).Create(jwt.JWTData{
			Email: email,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data := LoginResponse{
			Token: token,
		}
		res.Json(w, data, 200)
	}
}

func (handler *AuthHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		body, err := req.HandleBody[RegisterRequest](&w, r)
		if err != nil {
			return
		}
		email, err := handler.AuthService.Register(body.Email, body.Password, body.Name)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		token, err := jwt.NewJWT(handler.Config.Auth.Secret).Create(jwt.JWTData{
			Email: email,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		data := RegisterResponse{
			Token: token,
		}
		res.Json(w, data, http.StatusCreated)
	}
}
