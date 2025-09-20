package main

import (
	"fmt"

	"github.com/findwannawhy/IAD-Semester5/internal/app/config"
	"github.com/findwannawhy/IAD-Semester5/internal/app/dsn"
	"github.com/findwannawhy/IAD-Semester5/internal/app/handler"
	"github.com/findwannawhy/IAD-Semester5/internal/app/repository"
	"github.com/findwannawhy/IAD-Semester5/internal/pkg"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func main() {
	router := gin.Default()
	conf, err := config.NewConfig()
	if err != nil {
		logrus.Fatalf("error loading config: %v", err)
	}

	postgresString := dsn.FromEnv()
	fmt.Println(postgresString)

	rep, errRep := repository.NewRepository(postgresString)
	if errRep != nil {
		logrus.Fatalf("error initializing repository: %v", errRep)
	}

	hand := handler.NewHandler(rep)

	application := pkg.NewApp(conf, router, hand)
	application.RunApp()
}