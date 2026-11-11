package main

import (
	"database/sql"
	"log"
	"net/http"
	
	_ "github.com/lib/pq" // анонимный импорт: регистрирует драйвер "postgres"

	"service-xp/internal/handler"
	"service-xp/internal/repo"
	"service-xp/internal/service"

)


func main() {

	// 1. Подключение к БД.(поправить под мою бд)
	conn := "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable&connect_timeout=5"
	db, err := sql.Open("postgres", conn)
	log.Println("db start")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()


	//2. Проверка, что БД реально доступна.
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	// 3. Проводка слоёв: repo → service → handler.
	userRepo := repo.NewUserRepo(db)
	userSvc := service.NewUserService(userRepo)
	userHandler := handler.NewUserHandler(userSvc)


	// 4. Маршруты.
	mux := http.NewServeMux()
	mux.HandleFunc("/users", userHandler.HandleUsers)
	

	// 5. Запуск.
	log.Println("linstening on :8081")
	if err := http.ListenAndServe(":8081", mux); err != nil {
		log.Fatal(err)
	}


}


// Где делать Update — по файлам
// internal/repo/user.go — метод Update (SQL).
// internal/service/user.go — метод UpdateUser (валидация).
// internal/handler/user.go — метод UpdateUser (HTTP).
// internal/handler/user.go — метод HandleUserByID (switch по методу).
// main.go — маршрут "/users/" → HandleUserByID.
// Всё. Начни с internal/repo/user.go.

