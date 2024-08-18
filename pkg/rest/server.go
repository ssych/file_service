package rest

import (
	"context"
	"log"
	"net/http"

	"github.com/ssych/file_service/pkg/config"
	"github.com/ssych/file_service/pkg/handlers"
	"github.com/ssych/file_service/pkg/rest/middleware"
	"github.com/ssych/file_service/pkg/storage"
	"github.com/ssych/file_service/pkg/store"
)

type Server struct {
	ctx    context.Context
	server *http.Server
}

func NewServer(
	ctx context.Context,
) (*Server, error) {
	connStr := "host=localhost user=dev dbname=file_service_development sslmode=disable password=dev"

	db, err := storage.NewDB(ctx, &config.DBOption{
		ConnectString: connStr,
	})
	if err != nil {
		return nil, err
	}

	st := store.NewStore(db)

	h := handlers.NewHandler(st)
	m := middleware.NewAuthMiddleware(st)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/auth", h.Login)

	mux.Handle("POST /api/upload-asset/{asset_name}", m.MiddlewareFunc(http.HandlerFunc(h.CreateAsset)))

	mux.Handle("GET /api/asset/{asset_name}", m.MiddlewareFunc(http.HandlerFunc(h.FindAsset)))

	srv := &http.Server{
		Addr:    ":8091",
		Handler: mux,
	}

	return &Server{
		ctx:    ctx,
		server: srv,
	}, nil
}

func (s *Server) Run() error {
	log.Println("run server")
	if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("failed to start HTTP/REST server: %v", err)
		return err
	}
	return nil
}

func (s *Server) Close() {
	if err := s.server.Shutdown(s.ctx); err != nil {
		log.Fatalf("failed to Shutdown HTTP/REST server: %v", err)
	}
}
