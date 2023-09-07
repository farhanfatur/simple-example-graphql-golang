package main

import (
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/farhanfatur/simple-example-graphql-golang/graph"
	"github.com/farhanfatur/simple-example-graphql-golang/internal/auth"
	database "github.com/farhanfatur/simple-example-graphql-golang/internal/pkg/db/migrations/postgres"
	"github.com/go-chi/chi"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	app := chi.NewRouter()

	app.Use(auth.Middleware())

	database.InitDB()
	defer database.CloseDB()
	database.Migrate()

	srv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{}}))
	app.Handle("/", playground.Handler("GraphQL playground", "/query"))
	app.Handle("/query", srv)

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, app))
}
