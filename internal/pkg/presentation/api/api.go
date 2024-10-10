package api

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/diwise/iot-things/internal/app"
	"github.com/diwise/iot-things/internal/pkg/presentation/auth"
	"github.com/diwise/service-chassis/pkg/infrastructure/o11y"
	"github.com/diwise/service-chassis/pkg/infrastructure/o11y/logging"
	"github.com/diwise/service-chassis/pkg/infrastructure/o11y/tracing"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"go.opentelemetry.io/otel"
)

var tracer = otel.Tracer("iot-things/api/things")

func Register(ctx context.Context, app app.ThingsApp, policies io.Reader) (*chi.Mux, error) {
	log := logging.GetFromContext(ctx)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))

	authenticator, err := auth.NewAuthenticator(ctx, log, policies)
	if err != nil {
		return nil, fmt.Errorf("failed to create api authenticator: %w", err)
	}

	r.Route("/api/v0", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Use(authenticator)

			r.Route("/things", func(r chi.Router) {
				r.Get("/", queryHandler(log, app))
				r.Get("/{id}", getByIDHandler(log, app))
				r.Post("/", addHandler(log, app))
				r.Put("/{id}", updateHandler(log, app))
				r.Patch("/{id}", patchHandler(log, app))
				r.Get("/tags", getTagsHandler(log, app))
				r.Get("/types", getTypesHandler(log, app))
			})
		})
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	return r, nil
}

func queryHandler(log *slog.Logger, a app.ThingsApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		ctx, span := tracer.Start(r.Context(), "query-things")
		defer func() { tracing.RecordAnyErrorAndEndSpan(err, span) }()
		_, ctx, logger := o11y.AddTraceIDToLoggerAndStoreInContext(span, log, ctx)

		w.Header().Set("Content-Type", "application/vnd.api+json")

		result, err := a.QueryThings(ctx, r.URL.Query())
		if err != nil {
			logger.Error("could not query things", "err", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		if result.Count == 0 {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("[]"))
			return
		}

		data := make([]json.RawMessage, 0, len(result.Things))
		for _, thing := range result.Things {
			data = append(data, thing)
		}

		response := NewApiResponse(r, data, uint64(result.Count), uint64(result.TotalCount), uint64(result.Offset), uint64(result.Limit))

		b, err := json.Marshal(response)
		if err != nil {
			logger.Error("could not marshal query response", "err", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.Header().Set("Content-Type", "application/vnd.api+json")
		w.WriteHeader(http.StatusOK)
		w.Write(b)
	}
}

func getByIDHandler(log *slog.Logger, a app.ThingsApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		ctx, span := tracer.Start(r.Context(), "get-thing-byID")
		defer func() { tracing.RecordAnyErrorAndEndSpan(err, span) }()
		_, ctx, logger := o11y.AddTraceIDToLoggerAndStoreInContext(span, log, ctx)

		w.Header().Set("Content-Type", "application/vnd.api+json")

		thingId := chi.URLParam(r, "id")
		if thingId == "" {
			logger.Error("no id parameter found in request")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		result, err := a.QueryThings(ctx, map[string][]string{"id": {thingId}})
		if err != nil {
			logger.Debug("failed to query things", "err", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}
		if result.Count == 0 || result.Count > 1 {
			logger.Debug("thing not found", "id", thingId)
			w.WriteHeader(http.StatusNotFound)
			return
		}

		response := NewApiResponse(r, json.RawMessage(result.Things[0]), 1, 1, 0, 1)

		w.WriteHeader(http.StatusOK)
		w.Write(response.Byte())
	}
}

func addHandler(log *slog.Logger, a app.ThingsApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		defer r.Body.Close()

		w.Header().Set("Content-Type", "application/vnd.api+json")

		if isMultipartFormData(r) {
			ctx, span := tracer.Start(r.Context(), "seed")
			defer func() { tracing.RecordAnyErrorAndEndSpan(err, span) }()
			_, ctx, logger := o11y.AddTraceIDToLoggerAndStoreInContext(span, log, ctx)

			file, _, err := r.FormFile("fileupload")
			if err != nil {
				logger.Error("unable to get file from fileupload", "err", err.Error())
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			defer file.Close()

			err = a.Seed(ctx, file)
			if err != nil {
				logger.Error("could not seed", "err", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(err.Error()))
				return
			}

			w.WriteHeader(http.StatusCreated)
			return
		}

		ctx, span := tracer.Start(r.Context(), "create-thing")
		defer func() { tracing.RecordAnyErrorAndEndSpan(err, span) }()
		_, ctx, logger := o11y.AddTraceIDToLoggerAndStoreInContext(span, log, ctx)

		b, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Error("could not read body", "err", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		err = a.AddThing(ctx, b)
		if err != nil && errors.Is(err, app.ErrAlreadyExists) {
			w.WriteHeader(http.StatusConflict)
			return
		}
		if err != nil {
			logger.Error("could not create thing", "err", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusCreated)
	}
}

func updateHandler(log *slog.Logger, a app.ThingsApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		defer r.Body.Close()

		ctx, span := tracer.Start(r.Context(), "update-thing")
		defer func() { tracing.RecordAnyErrorAndEndSpan(err, span) }()
		_, ctx, logger := o11y.AddTraceIDToLoggerAndStoreInContext(span, log, ctx)

		w.Header().Set("Content-Type", "application/vnd.api+json")

		b, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Error("could not read body", "err", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		tenants := auth.GetAllowedTenantsFromContext(ctx)

		err = a.UpdateThing(ctx, b, tenants)
		if err != nil {
			logger.Error("could not update thing", "err", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func patchHandler(log *slog.Logger, a app.ThingsApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error
		defer r.Body.Close()

		ctx, span := tracer.Start(r.Context(), "patch-thing")
		defer func() { tracing.RecordAnyErrorAndEndSpan(err, span) }()
		_, ctx, logger := o11y.AddTraceIDToLoggerAndStoreInContext(span, log, ctx)

		w.Header().Set("Content-Type", "application/vnd.api+json")

		thingId := chi.URLParam(r, "id")
		if thingId == "" {
			logger.Error("no id parameter found in request")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		b, err := io.ReadAll(r.Body)
		if err != nil {
			logger.Error("could not read body", "err", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		tenants := auth.GetAllowedTenantsFromContext(ctx)

		err = a.MergeThing(ctx, thingId, b, tenants)
		if err != nil {
			logger.Error("could not patch thing", "err", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func getTagsHandler(log *slog.Logger, a app.ThingsApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		ctx, span := tracer.Start(r.Context(), "get-tags")
		defer func() { tracing.RecordAnyErrorAndEndSpan(err, span) }()
		_, ctx, logger := o11y.AddTraceIDToLoggerAndStoreInContext(span, log, ctx)

		w.Header().Set("Content-Type", "application/vnd.api+json")

		tenants := auth.GetAllowedTenantsFromContext(ctx)

		tags, err := a.GetTags(ctx, tenants)
		if err != nil {
			logger.Error("could not get tags", "err", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		response := NewApiResponse(r, tags, uint64(len(tags)), uint64(len(tags)), 0, uint64(len(tags)))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response.Byte())
	}
}

func getTypesHandler(log *slog.Logger, a app.ThingsApp) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var err error

		ctx, span := tracer.Start(r.Context(), "get-types")
		defer func() { tracing.RecordAnyErrorAndEndSpan(err, span) }()
		_, ctx, logger := o11y.AddTraceIDToLoggerAndStoreInContext(span, log, ctx)

		tenants := auth.GetAllowedTenantsFromContext(ctx)

		types, err := a.GetTypes(ctx, tenants)
		if err != nil {
			logger.Error("could not get tags", "err", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(err.Error()))
			return
		}

		response := NewApiResponse(r, types, uint64(len(types)), uint64(len(types)), 0, uint64(len(types)))

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(response.Byte())
	}
}

func isMultipartFormData(r *http.Request) bool {
	contentType := r.Header.Get("Content-Type")
	return strings.Contains(contentType, "multipart/form-data")
}
