package main

import (
	"fmt"
	"golang/advanced/configs"
	"golang/advanced/internal/auth"
	"golang/advanced/internal/link"
	"golang/advanced/internal/stat"
	"golang/advanced/internal/user"
	"golang/advanced/pkg/db"
	"golang/advanced/pkg/event"
	"golang/advanced/pkg/middleware"
	"net/http"
)

func App() http.Handler {
	conf := configs.LoadConfig()
	db := db.NewDb(conf)
	router := http.NewServeMux()
	eventBus := event.NewEventBus()

	// Repositories
	linkRepository := link.NewLinkRepository(db)
	userRepository := user.NewUserRepository(db)
	statRepository := stat.NewStatRepository(db)

	// Service
	authService := auth.NewAuthService(userRepository)
	statService := stat.NewStatService(&stat.StatServiceDeps{
		EventBus:       eventBus,
		StatRepository: statRepository,
	})

	// Handler
	auth.NewAuthHandler(router, auth.AuthHandlerDepth{
		Config:      conf,
		AuthService: authService,
	})

	link.NewLinkHandler(router, link.LinkHandlerDepth{
		LinkRepository: linkRepository,
		EventBus:       eventBus,
		Config:         conf,
	})

	stat.NewStatHandler(router, stat.StatHandlerDepth{
		StatRepository: statRepository,
		Config:         conf,
	})

	go statService.AddClick()

	// Middlewares
	stack := middleware.Chain(
		middleware.CORS,
		middleware.Logging,
	)
	return stack(router)
}

func main() {
	app := App()

	server := &http.Server{
		Addr:    ":8081",
		Handler: app,
	}

	fmt.Println("Server is listening on port 8081")
	server.ListenAndServe()
}

func ping(url string, respCh chan int, errCh chan error) {
	resp, err := http.Get(url)
	if err != nil {
		errCh <- err
		return
	}
	defer resp.Body.Close()
	respCh <- resp.StatusCode
}

func sumPart(arr []int, ch chan int) {
	sum := 0
	for _, num := range arr {
		sum += num
	}
	ch <- sum
}

func getHttpCode(codeCh chan int) {
	resp, err := http.Get("http://www.google.com")
	if err != nil {
		fmt.Println(err)
	}
	defer resp.Body.Close()
	codeCh <- resp.StatusCode
}
