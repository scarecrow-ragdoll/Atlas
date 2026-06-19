// FILE: apps/api/cmd/server/main.go
// VERSION: 1.0.1
// START_MODULE_CONTRACT
//   PURPOSE: Start the API HTTP server and wire persistence, GraphQL, admin auth, Atlas module, health, and graceful shutdown.
//   SCOPE: Runtime dependency construction, admin bootstrap validation/seed, Atlas bootstrap, explicit route groups (public, admin, Atlas auth-public, Atlas guarded), middleware order, and process lifecycle.
//   DEPENDS: apps/api/internal/appconfig, apps/api/internal/graph, apps/api/internal/handler, apps/api/internal/middleware, apps/api/internal/atlas/..., apps/api/internal/repository/postgres, apps/api/internal/repository/redis, apps/api/internal/service, libs/go/config, libs/go/logger.
//   LINKS: M-API / M-GRAPHQL-SCHEMA / M-WEB-ADMIN / V-M-API.
//   ROLE: RUNTIME
//   MAP_MODE: SUMMARY
// END_MODULE_CONTRACT
// START_MODULE_MAP
//   main - Loads config, applies admin env/defaults, runs migrations, bootstraps initial admin and Atlas default user, wires route groups, and starts the HTTP server.
//   optionalEnvFile - Returns a local .env path only when present.
// END_MODULE_MAP
// START_CHANGE_SUMMARY
//   LAST_CHANGE: 1.0.2 - Wired WAVE-03 Atlas workout repository/service into API startup and Atlas resolver construction.
// END_CHANGE_SUMMARY

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	"monorepo-template/libs/go/config"
	"monorepo-template/libs/go/logger"

	"monorepo-template/apps/api/internal/appconfig"
	atlasGenerated "monorepo-template/apps/api/internal/atlas/graph/generated"
	atlasResolver "monorepo-template/apps/api/internal/atlas/graph/resolver"
	atlasMiddleware "monorepo-template/apps/api/internal/atlas/middleware"
	atlasPostgres "monorepo-template/apps/api/internal/atlas/repository/postgres"
	atlasRedis "monorepo-template/apps/api/internal/atlas/repository/redis"
	atlasService "monorepo-template/apps/api/internal/atlas/service"
	"monorepo-template/apps/api/internal/graph"
	healthHandler "monorepo-template/apps/api/internal/handler"
	"monorepo-template/apps/api/internal/middleware"
	"monorepo-template/apps/api/internal/repository/postgres"
	redisRepo "monorepo-template/apps/api/internal/repository/redis"
	"monorepo-template/apps/api/internal/service"
)

func main() {
	cfg, err := config.Load[appconfig.Config](config.Options{
		ConfigPath: "config/config.yml",
		EnvFile:    optionalEnvFile(".env"),
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}
	if err := appconfig.ApplyAdminEnvOverlay(&cfg, os.LookupEnv); err != nil {
		fmt.Fprintf(os.Stderr, "failed to apply admin env: %v\n", err)
		os.Exit(1)
	}
	if err := appconfig.ApplyAdminDefaults(&cfg); err != nil {
		fmt.Fprintf(os.Stderr, "failed to apply admin defaults: %v\n", err)
		os.Exit(1)
	}

	l, err := logger.New(logger.Config{
		Level:  cfg.Log.Level,
		Format: cfg.Log.Format,
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to init logger: %v\n", err)
		os.Exit(1)
	}
	defer func() { _ = l.Sync() }()

	db, err := postgres.New(cfg.Postgres, l)
	if err != nil {
		l.Fatal("failed to connect to postgres", zap.Error(err))
	}
	defer db.Close()

	if err := postgres.RunMigrations(cfg.Postgres.DSN(), l); err != nil {
		l.Fatal("failed to run migrations", zap.Error(err))
	}

	rdb, err := redisRepo.New(cfg.Redis, l)
	if err != nil {
		l.Fatal("failed to connect to redis", zap.Error(err))
	}
	defer func() { _ = rdb.Close() }()

	userRepo := postgres.NewUserRepo(db.Pool)
	userService := service.NewUserService(userRepo)
	adminRepo := postgres.NewAdminRepo(db.Pool)
	adminSessionStore := redisRepo.NewAdminSessionStore(rdb.RDB, []byte(cfg.AdminSession.KeySecret), cfg.AdminSession.TTL)
	adminAuthService := service.NewAdminAuthService(adminRepo, adminSessionStore)

	l.Info("[AdminAuth][seed][BLOCK_BOOTSTRAP_ADMIN] seed bootstrap starting")
	adminCount, err := adminRepo.Count(context.Background())
	if err != nil {
		l.Fatal("failed to count admin users", zap.Error(err))
	}
	if err := appconfig.ValidateAdminBootstrapEnv(os.LookupEnv, adminCount == 0); err != nil {
		l.Fatal("[AdminAuth][seed][BLOCK_BOOTSTRAP_ADMIN] missing bootstrap env", zap.Error(err))
	}
	if err := appconfig.ValidateAdminBootstrap(cfg, adminCount == 0); err != nil {
		l.Fatal("[AdminAuth][seed][BLOCK_BOOTSTRAP_ADMIN] invalid bootstrap config", zap.Error(err))
	}
	seeded, err := adminAuthService.SeedInitialAdmin(context.Background(), service.InitialAdminInput{
		Email:    cfg.Admin.InitialEmail,
		Name:     cfg.Admin.InitialName,
		Password: cfg.Admin.InitialPassword,
	})
	if err != nil {
		l.Fatal("[AdminAuth][seed][BLOCK_BOOTSTRAP_ADMIN] seed failed", zap.Error(err))
	}
	if seeded {
		l.Info("[AdminAuth][seed][BLOCK_BOOTSTRAP_ADMIN] seeded initial admin")
	} else {
		l.Info("[AdminAuth][seed][BLOCK_BOOTSTRAP_ADMIN] seed skipped")
	}

	atlasSettingsRepo := atlasPostgres.NewSettingsRepository(db.Pool)
	atlasPinSessionStore := atlasRedis.NewPinSessionStore(rdb.RDB, []byte(cfg.AdminSession.KeySecret))
	atlasPinAttemptStore := atlasRedis.NewPinAttemptStore(
		rdb.RDB,
		cfg.AtlasPinAttempt.MaxFailures,
		cfg.AtlasPinAttempt.LockoutDuration,
		cfg.AtlasPinAttempt.EscalatedDuration,
	)
	atlasBootstrapService := atlasService.NewBootstrapService(db.Pool)
	atlasSettingsService := atlasService.NewSettingsService(atlasSettingsRepo)
	atlasPinService := atlasService.NewPinService(
		atlasSettingsRepo,
		atlasPinSessionStore,
		atlasService.Argon2Params{
			Memory:      cfg.AtlasPin.Argon2Memory,
			Iterations:  cfg.AtlasPin.Argon2Iterations,
			Parallelism: cfg.AtlasPin.Argon2Parallelism,
			KeyLength:   cfg.AtlasPin.Argon2KeyLength,
		},
		cfg.AtlasPin.MinLength,
		cfg.AtlasPin.MaxLength,
	)
	atlasExerciseRepo := atlasPostgres.NewExerciseRepository(db.Pool)
	atlasExerciseService := atlasService.NewExerciseService(atlasExerciseRepo)
	atlasWorkoutRepo := atlasPostgres.NewWorkoutRepository(db.Pool)
	atlasWorkoutService := atlasService.NewWorkoutService(atlasWorkoutRepo, atlasExerciseService)
	atlasMediaHandler := healthHandler.NewAtlasMediaHandler(atlasExerciseService, cfg.Media.BasePath)

	l.Info("[Atlas][bootstrap] ensuring default user and settings")
	atlasUserID, err := atlasBootstrapService.EnsureDefaultUser(context.Background())
	if err != nil {
		l.Fatal("[Atlas][bootstrap] failed to ensure default user", zap.Error(err))
	}
	if err := atlasBootstrapService.EnsureDefaultSettings(context.Background(), atlasUserID); err != nil {
		l.Fatal("[Atlas][bootstrap] failed to ensure default settings", zap.Error(err))
	}
	l.Info("[Atlas][bootstrap] default user ready", zap.String("user_id", atlasUserID))

	atlasRes := &atlasResolver.Resolver{
		SettingsService: atlasSettingsService,
		PinService:      atlasPinService,
		ExerciseService: atlasExerciseService,
		WorkoutService:  atlasWorkoutService,
	}
	atlasSrv := handler.NewDefaultServer(atlasGenerated.NewExecutableSchema(atlasGenerated.Config{Resolvers: atlasRes}))

	adminResolver := &graph.Resolver{
		UserService:      userService,
		AdminAuthService: adminAuthService,
	}
	adminSrv := handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: adminResolver}))

	r := chi.NewRouter()
	r.Use(logger.RequestID(l))
	r.Use(logger.Logging())

	publicCORS := middleware.CORS(middleware.CORSConfig{
		AllowedOrigins: cfg.Server.CORSOrigins,
		AllowedMethods: []string{"GET", "POST", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"Content-Type"},
	})
	adminCORS := middleware.CORS(middleware.CORSConfig{
		AllowedOrigins:   cfg.Admin.Origins,
		AllowedMethods:   []string{"GET", "POST", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})
	adminCookieConfig := middleware.AdminCookieConfigFromConfig(cfg.AdminSession, cfg.Server.Env)

	r.Group(func(public chi.Router) {
		public.Use(publicCORS)
		public.Get("/healthz", healthHandler.Healthz())
		public.Get("/readyz", healthHandler.Readyz(db, rdb))
		public.Get("/api/v1/healthz", healthHandler.AtlasHealthz())
		public.Get("/api/v1/readyz", healthHandler.AtlasReadyz(db, rdb))

		usersHandler := healthHandler.NewUsersHandler(userService, l)
		public.Mount("/api/users", usersHandler.Routes())
	})

	r.Group(func(admin chi.Router) {
		admin.Use(adminCORS)
		admin.Use(middleware.AdminOriginGuard(cfg.Admin.Origins))
		admin.Use(middleware.AdminSessionMiddleware(adminAuthService, cfg.AdminSession.CookieName))
		admin.Handle("/graphql", middleware.WithAdminCookieBridge(adminSrv, adminCookieConfig))
		if cfg.Server.Env != "production" {
			admin.Handle("/playground", playground.Handler("GraphQL Playground", "/graphql"))
		}
	})

	env := strings.ToLower(strings.TrimSpace(cfg.Server.Env))
	atlasCookieSecure := cfg.AtlasPinSession.CookieSecure == "true" || (cfg.AtlasPinSession.CookieSecure == "auto" && env == "production")
	atlasSameSite := atlasSameSiteMode(cfg.AtlasPinSession.SameSite)

	pinAuthHandler := healthHandler.NewPinAuthHandler(
		atlasPinService,
		atlasPinSessionStore,
		atlasPinAttemptStore,
		cfg.AtlasPinSession.CookieName,
		cfg.AtlasPinSession.IdleTTL,
		cfg.AtlasPinSession.AbsoluteTTL,
		atlasCookieSecure,
		atlasSameSite,
	)

	_ = r.Group(func(atlasAuth chi.Router) {
		atlasAuth.Use(atlasMiddleware.AtlasUserContext(atlasBootstrapService))
		atlasAuth.Post("/api/v1/auth/pin/unlock", pinAuthHandler.Unlock)
		atlasAuth.Post("/api/v1/auth/pin/lock", pinAuthHandler.Lock)
		atlasAuth.Get("/api/v1/auth/session", pinAuthHandler.SessionCheck)
	})

	_ = r.Group(func(atlas chi.Router) {
		atlas.Use(atlasMiddleware.AtlasUserContext(atlasBootstrapService))
		atlas.Use(atlasMiddleware.AtlasPinGuard(atlasPinService, atlasPinSessionStore, cfg.AtlasPinSession.CookieName))
		atlas.Handle("/graphql/atlas", atlasSrv)
		atlas.Post("/api/v1/media/upload", atlasMediaHandler.Upload)
		atlas.Get("/api/v1/media/{id}", atlasMediaHandler.Download)
		atlas.Delete("/api/v1/media/{id}", atlasMediaHandler.Delete)
	})

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:      r,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	go func() {
		l.Info("starting server", zap.Int("port", cfg.Server.Port))
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			l.Fatal("server failed", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	l.Info("shutting down server")
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	if err := httpServer.Shutdown(ctx); err != nil {
		l.Fatal("server forced to shutdown", zap.Error(err))
	}

	l.Info("server stopped")
}

func optionalEnvFile(path string) string {
	if _, err := os.Stat(path); err != nil {
		return ""
	}
	return path
}

func atlasSameSiteMode(value string) http.SameSite {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "strict":
		return http.SameSiteStrictMode
	case "none":
		return http.SameSiteNoneMode
	default:
		return http.SameSiteLaxMode
	}
}
