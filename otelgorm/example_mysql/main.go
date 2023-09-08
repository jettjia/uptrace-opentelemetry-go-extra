package main

import (
	"context"
	"fmt"
	"gorm.io/gorm/schema"

	"go.opentelemetry.io/otel"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	"github.com/uptrace/opentelemetry-go-extra/otelplay"
)

type User struct {
	Id     int    `gorm:"primary_key" json:"id"`
	Name   string `json:"name"`
	Age    int    `json:"age"`
	Gender int    `json:"gender"`
}

func main() {
	ctx := context.Background()

	shutdown := otelplay.ConfigureOpentelemetry(ctx)
	defer shutdown()

	dsn := "root:admin123@tcp(10.4.7.71:3306)/go_demo?charset=utf8mb4&parseTime=True&loc=Local"

	cfg := &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	}

	db, err := gorm.Open(mysql.Open(dsn), cfg)
	if err != nil {
		panic(err)
	}

	if err = db.Use(otelgorm.NewPlugin()); err != nil {
		panic(err)
	}

	tracer := otel.Tracer("app_or_package_name_mysql")

	ctx, span := tracer.Start(ctx, "root")
	defer span.End()

	var user User
	err = db.WithContext(ctx).Limit(1).Find(&user).Where("id = ?", 2).Error
	if err != nil {
		panic(err)
	}
	fmt.Println(user)

	otelplay.PrintTraceID(ctx)
}
