package main

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gabrielgasp/go-taxa-tesouro/model"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httprate"
	"github.com/spf13/viper"
)

type Api interface {
	Run(context.Context)
}

type api struct {
	logger  *slog.Logger
	rwMutex *sync.RWMutex
	wg      *sync.WaitGroup
}

func NewApi(logger *slog.Logger, rxMutex *sync.RWMutex, wg *sync.WaitGroup) Api {
	return api{
		logger:  logger,
		rwMutex: rxMutex,
		wg:      wg,
	}
}

func (a api) Run(ctx context.Context) {
	defer a.wg.Done()

	r := chi.NewRouter()
	r.Use(a.realIpMiddleware)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(httprate.LimitByIP(viper.GetInt("RATE_LIMIT_PER_MINUTE"), time.Minute))

	r.Get("/health", a.health)
	r.Get("/bonds", a.listAllBonds)
	r.Get("/bonds/{bondName}", a.getBondByName)

	server := &http.Server{
		Addr:    ":" + viper.GetString("PORT"),
		Handler: r,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.logger.Error("Failed to start API", "error", err.Error())
			os.Exit(1)
		}
	}()

	a.logger.Info("Server is running on port " + viper.GetString("PORT"))

	<-ctx.Done()
	a.shutdown(server)
}

func (a api) shutdown(server *http.Server) {
	ctx, ctxCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer ctxCancel()

	if err := server.Shutdown(ctx); err != nil {
		a.logger.Error("Failed to shutdown API", "error", err.Error())
		return
	}

	a.logger.Info("API stopped")
}

func (a api) health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (a api) listAllBonds(w http.ResponseWriter, _ *http.Request) {
	a.rwMutex.RLock()
	defer a.rwMutex.RUnlock()

	var response model.ListAllBondsResponse
	response.Bonds = scrapperCache.BondsList
	response.UpdatedAt = scrapperCache.UpdatedAt

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		a.logger.Error("Failed to encode response", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a api) getBondByName(w http.ResponseWriter, r *http.Request) {
	a.rwMutex.RLock()
	defer a.rwMutex.RUnlock()

	bondName := chi.URLParam(r, "bondName")
	bondName = strings.ReplaceAll(bondName, "-", " ")

	cachedBond, found := scrapperCache.BondsMap[strings.ToLower(bondName)]
	if !found {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	var response model.GetBondByNameResponse
	response.Bond = cachedBond
	response.UpdatedAt = scrapperCache.UpdatedAt

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		a.logger.Error("Failed to encode response", "error", err.Error())
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (a api) realIpMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var ip string

		if cfip := r.Header.Get("CF-Connecting-IP"); cfip != "" {
			ip = cfip
		} else if tcip := r.Header.Get("True-Client-IP"); tcip != "" {
			ip = tcip
		} else if xrip := r.Header.Get("X-Real-IP"); xrip != "" {
			ip = xrip
		} else if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
			i := strings.Index(xff, ",")
			if i == -1 {
				i = len(xff)
			}
			ip = xff[:i]
		}

		if ip != "" && net.ParseIP(ip) != nil {
			r.RemoteAddr = ip
		}

		next.ServeHTTP(w, r)
	})
}
