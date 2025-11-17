package httpserver

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"authorization-vendor/internal/controller"
)

func NewRouter(ctrl *controller.AuthorizationVendorController) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/vendedores", listVendedoresHandler(ctrl))
	mux.HandleFunc("/vendedores/", vendorRouter(ctrl))
	mux.HandleFunc("/health", ctrl.Health)
	return withCORS(mux)
}

func listVendedoresHandler(ctrl *controller.AuthorizationVendorController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, errors.New("metodo no permitido"))
			return
		}
		vendors, err := ctrl.Service().ListarVendedores()
		if err != nil {
			writeError(w, http.StatusInternalServerError, err)
			return
		}
		writeJSON(w, http.StatusOK, vendors)
	}
}

func vendorRouter(ctrl *controller.AuthorizationVendorController) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodOptions {
			writeCORSHeaders(w)
			w.WriteHeader(http.StatusNoContent)
			return
		}
		id, action, err := parsePath(r.URL.Path)
		if err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		switch action {
		case "autorizar":
			usuario := r.Header.Get("X-User")
			estado, err := ctrl.Service().AutorizarVendedor(id, usuario)
			if err != nil {
				writeError(w, http.StatusBadRequest, err)
				return
			}
			writeJSON(w, http.StatusOK, map[string]any{"estado": estado})
		case "revocar":
			usuario := r.Header.Get("X-User")
			estado, err := ctrl.Service().RevocarVendedor(id, usuario)
			if err != nil {
				writeError(w, http.StatusBadRequest, err)
				return
			}
			writeJSON(w, http.StatusOK, map[string]any{"estado": estado})
		case "estado":
			estado, err := ctrl.Service().ConsultarEstado(id)
			if err != nil {
				writeError(w, http.StatusNotFound, err)
				return
			}
			writeJSON(w, http.StatusOK, map[string]any{"estado": estado})
		case "historial":
			list, err := ctrl.Service().Historial(id)
			if err != nil {
				writeError(w, http.StatusInternalServerError, err)
				return
			}
			writeJSON(w, http.StatusOK, list)
		default:
			writeError(w, http.StatusNotFound, errors.New("ruta no encontrada"))
		}
	}
}

func parsePath(path string) (int, string, error) {
	trimmed := strings.TrimPrefix(path, "/vendedores/")
	parts := strings.Split(trimmed, "/")
	if len(parts) < 2 {
		return 0, "", errors.New("ruta invalida")
	}
	id, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, "", err
	}
	return id, parts[1], nil
}

// withCORS wraps the handler to set CORS headers for browsers.
func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		writeCORSHeaders(w)
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func writeCORSHeaders(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type, X-User")
}
