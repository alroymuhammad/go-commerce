package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

type ProductHandler struct {
	db *sql.DB
}

type ProductRequest struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Stock       int     `json:"stock"`
}

type ProductResponse struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	CreatedAt   time.Time `json:"created_at"`
}

func NewProductHandler(db *sql.DB) *ProductHandler {
	return &ProductHandler{db: db}
}

func getProductIDFromPath(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) > 2 {
		return parts[2]
	}
	return ""
}

func (h *ProductHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	switch {
	case r.Method == http.MethodGet && r.URL.Path == "/products":
		h.ListProducts(w, r)
	case r.Method == http.MethodGet && strings.HasPrefix(r.URL.Path, "/products/"):
		h.GetProduct(w, r)
	case r.Method == http.MethodPost && r.URL.Path == "/products":
		h.CreateProduct(w, r)
	case r.Method == http.MethodPut && strings.HasPrefix(r.URL.Path, "/products/"):
		h.UpdateProduct(w, r)
	case r.Method == http.MethodDelete && strings.HasPrefix(r.URL.Path, "/products/"):
		h.DeleteProduct(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (h *ProductHandler) ListProducts(w http.ResponseWriter, r *http.Request) {
	rows, err := h.db.Query("SELECT id, name, description, price, stock, created_at FROM products")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var products []ProductResponse
	for rows.Next() {
		var product ProductResponse
		if err := rows.Scan(&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock, &product.CreatedAt); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		products = append(products, product)
	}

	json.NewEncoder(w).Encode(products)
}

func (h *ProductHandler) GetProduct(w http.ResponseWriter, r *http.Request) {
	id := getProductIDFromPath(r.URL.Path)

	var product ProductResponse
	err := h.db.QueryRow("SELECT id, name, description, price, stock, created_at FROM products WHERE id = $1", id).Scan(
		&product.ID, &product.Name, &product.Description, &product.Price, &product.Stock, &product.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(product)
}

func (h *ProductHandler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var req ProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var id int64
	err := h.db.QueryRow(
		"INSERT INTO products (name, description, price, stock) VALUES ($1, $2, $3, $4) RETURNING id",
		req.Name, req.Description, req.Price, req.Stock,
	).Scan(&id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int64{"id": id})
}

func (h *ProductHandler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	id := getProductIDFromPath(r.URL.Path)

	var req ProductRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	_, err := h.db.Exec("UPDATE products SET name = $1, description = $2, price = $3, stock = $4 WHERE id = $5",
		req.Name, req.Description, req.Price, req.Stock, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ProductHandler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	id := getProductIDFromPath(r.URL.Path)

	_, err := h.db.Exec("DELETE FROM products WHERE id = $1", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
