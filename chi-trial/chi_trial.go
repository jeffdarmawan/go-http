package chi_trial

import (
	"context"
	"fmt"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

// Source https://github.com/go-chi/chi
// Additional
// - In chi, Router is an interface while Mux is the implementation of chi Router

// // HTTP-method routing along `pattern`
// Connect(pattern string, h http.HandlerFunc)
// Delete(pattern string, h http.HandlerFunc)
// Get(pattern string, h http.HandlerFunc)
// Head(pattern string, h http.HandlerFunc)
// Options(pattern string, h http.HandlerFunc)
// Patch(pattern string, h http.HandlerFunc)
// Post(pattern string, h http.HandlerFunc)
// Put(pattern string, h http.HandlerFunc)
// Trace(pattern string, h http.HandlerFunc)

// Custom method
// chi.RegisterMethod("JELLO")
// r.Method("JELLO", "/path", myJelloMethodHandler)

func InitChiServer() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)

	// 1. Simple hello world
	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	})

	// 2. In-URL parameter
	// can also use regex
	// r.Get("/articles/{rid:^[0-9]{5,6}}", getArticle)
	r.Get("/user/{user_id}", getUserProfile)

	// 3. Sub-router
	// 3.1. Mounting
	userAPI := chi.NewRouter()
	userAPI.Get("/{user_id}", getUserProfile)

	blogAPI := chi.NewRouter()
	blogAPI.Get("/article/{article_id}", getArticle)

	r.Mount("/user", userAPI)
	r.Mount("/api", blogAPI)

	// 3.2. Route()
	r.Route("/transaction", func(r chi.Router) {
		r.Get("/", transactionList)

		r.Route("/{tx_id}", func(r chi.Router) {
			r.Get("/", transactionDetail)
			r.Post("/", trasactionModify)
		})
	})

	// 4. Grouping
	// for partial middleware implementation
	r.Group(func(r chi.Router) {
		r.Get("/", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("Hello World!"))
		})
		r.Get("/all-products", func(w http.ResponseWriter, r *http.Request) { w.Write([]byte("your products")) })
	})

	r.Group(func(r chi.Router) {
		r.Use(authMiddleware)
		r.Get("/transaction", func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("your transactions"))
		})
	})

	http.ListenAndServe(":3000", r)
}

// example handler
func getUserProfile(w http.ResponseWriter, r *http.Request) {
	userIDStr := chi.URLParam(r, "user_id")
	if userIDStr == "" {
		// rewriting error message
		w.WriteHeader(422)
		w.Write([]byte("error: empty user id"))
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		w.WriteHeader(422)
		w.Write([]byte(fmt.Sprintf("error: user id is not numeric. user id: %s", userIDStr)))
		return
	}

	w.Write([]byte(fmt.Sprintf("This is profile of user id: %d", userID)))
}

// dummy handlers
func getArticle(w http.ResponseWriter, r *http.Request)        {}
func transactionList(w http.ResponseWriter, r *http.Request)   {}
func transactionDetail(w http.ResponseWriter, r *http.Request) {}
func trasactionModify(w http.ResponseWriter, r *http.Request)  {}

// example middleware
func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// insert auth mechanism here
		ctx := context.WithValue(r.Context(), "user", "123")

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// built-in middlewares
// r.Use(middleware.AllowContentEncoding("deflate", "gzip"))
// r.Use(middleware.AllowContentType("application/json","text/xml"))
// r.Use(middleware.CleanPath) // /users//1 or //users////1 will both be treated as: /users/1
// r.Use(middleware.Compress(5, "text/html", "text/css"))
// allowedCharsets := []string{"UTF-8", "Latin-1", ""}
// r.Use(middleware.ContentCharset(allowedCharsets...))
// Basic CORS
//   // for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
//   r.Use(cors.Handler(cors.Options{
//     // AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
//     AllowedOrigins:   []string{"https://*", "http://*"},
//     // AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
//     AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
//     AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
//     ExposedHeaders:   []string{"Link"},
//     AllowCredentials: false,
//     MaxAge:           300, // Maximum value not ignored by any of major browsers
//   }))
// r.Use(middleware.GetHead) // GetHead automatically route undefined HEAD requests to GET handlers.
// r.Use(middleware.Heartbeat("/"))
// r.Use(middleware.Logger)        // <--<< Logger should come before Recoverer
// r.Use(middleware.Recoverer)
// r.Use(middleware.NoCache)
