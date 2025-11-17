package memory

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
	"time"

	"authorization-vendor/internal/model"
	repoiface "authorization-vendor/internal/repository"
)

// repository stores vendors in memory but persists changes to JSON files to simulate a DB.
type repository struct {
	mu             sync.RWMutex
	vendors        map[int]model.Vendedor
	historial      map[int][]model.HistorialAutorizacion
	vendorFilePath string
	histFilePath   string
}

var _ repoiface.AuthorizationVendorRepository = (*repository)(nil)

func NewJSONRepository(vendorFile, histFile string) repoiface.AuthorizationVendorRepository {
	repo := &repository{
		vendors:        make(map[int]model.Vendedor),
		historial:      make(map[int][]model.HistorialAutorizacion),
		vendorFilePath: vendorFile,
		histFilePath:   histFile,
	}
	_ = repo.loadFromDisk()
	return repo
}

func (r *repository) ObtenerEstado(idVendedor int) (bool, bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	v, ok := r.vendors[idVendedor]
	if !ok {
		return false, false
	}
	return v.EstadoAutorizado, true
}

func (r *repository) CambiarEstado(idVendedor int, estado bool) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	v, ok := r.vendors[idVendedor]
	if !ok {
		// create vendor if missing to keep simple
		v = model.Vendedor{IDVendedor: idVendedor}
	}
	v.EstadoAutorizado = estado
	now := time.Now()
	if estado {
		v.FechaAutorizado = now
		v.FechaRevocado = time.Time{}
	} else {
		v.FechaRevocado = now
	}
	r.vendors[idVendedor] = v

	return r.persistVendors()
}

func (r *repository) GuardarHistorial(h model.HistorialAutorizacion) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	h.FechaCambio = h.FechaCambio.Truncate(time.Second)
	r.historial[h.IDVendedor] = append(r.historial[h.IDVendedor], h)
	return r.persistHistorial()
}

func (r *repository) ListarHistorial(idVendedor int) ([]model.HistorialAutorizacion, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := r.historial[idVendedor]
	out := make([]model.HistorialAutorizacion, len(list))
	copy(out, list)
	return out, nil
}

func (r *repository) ListarVendedores() ([]model.Vendedor, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	list := make([]model.Vendedor, 0, len(r.vendors))
	for _, v := range r.vendors {
		list = append(list, v)
	}
	return list, nil
}

func (r *repository) loadFromDisk() error {
	if err := ensureDir(filepath.Dir(r.vendorFilePath)); err != nil {
		return err
	}
	if err := ensureDir(filepath.Dir(r.histFilePath)); err != nil {
		return err
	}
	if err := r.loadVendors(); err != nil {
		return err
	}
	if err := r.loadHistorial(); err != nil {
		return err
	}
	return nil
}

func (r *repository) loadVendors() error {
	bytes, err := os.ReadFile(r.vendorFilePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return r.persistVendors()
		}
		return err
	}
	var list []model.Vendedor
	if len(bytes) > 0 {
		if err := json.Unmarshal(bytes, &list); err != nil {
			return err
		}
	}
	for _, v := range list {
		r.vendors[v.IDVendedor] = v
	}
	return nil
}

func (r *repository) loadHistorial() error {
	bytes, err := os.ReadFile(r.histFilePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return r.persistHistorial()
		}
		return err
	}
	var list []model.HistorialAutorizacion
	if len(bytes) > 0 {
		if err := json.Unmarshal(bytes, &list); err != nil {
			return err
		}
	}
	for _, h := range list {
		r.historial[h.IDVendedor] = append(r.historial[h.IDVendedor], h)
	}
	return nil
}

func (r *repository) persistVendors() error {
	list := make([]model.Vendedor, 0, len(r.vendors))
	for _, v := range r.vendors {
		list = append(list, v)
	}
	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(r.vendorFilePath, data, 0o644)
}

func (r *repository) persistHistorial() error {
	list := make([]model.HistorialAutorizacion, 0)
	for _, h := range r.historial {
		list = append(list, h...)
	}
	data, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(r.histFilePath, data, 0o644)
}

func ensureDir(path string) error {
	if path == "." || path == "" {
		return nil
	}
	return os.MkdirAll(path, 0o755)
}
