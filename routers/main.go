package routers

import (
	"encoding/json"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/mswatermelon/GB_backend_1_project/helpers"
	"github.com/mswatermelon/GB_backend_1_project/models"
	"github.com/mswatermelon/GB_backend_1_project/url_shorterer"
	"gorm.io/gorm"
	"net/http"
	"time"
)

type Handler struct {
	router *chi.Mux
	db     *gorm.DB
}

type GenerateHashPostParams struct {
	Url string `json:"url"`
}

type GetStatsResponseParams struct {
	url    string
	access []map[string]time.Time
}

func (h *Handler) Welcome(w http.ResponseWriter, r *http.Request) {
	helpers.RespondWithJSON(w, http.StatusOK, map[string]string{
		"message": "You can put your URL after slash and get it's short version",
	})
}

func (h *Handler) GenerateHash(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var post GenerateHashPostParams
	err := json.NewDecoder(r.Body).Decode(&post)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	shortUrl := url_shorterer.ShortifyUrl(post.Url)
	hash := models.Hash{
		Hash:      shortUrl,
		Url:       post.Url,
		CreatedAt: time.Now(),
	}
	h.db.Create(&hash)

	requestIp, err := helpers.GetIP(r)
	if err != nil {
		helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
	}

	hit := models.Hit{
		HashID:   hash.ID,
		Ip:       requestIp,
		AccessAt: time.Now(),
	}
	h.db.Create(&hit)

	helpers.RespondWithJSON(w, http.StatusCreated, map[string]string{"short_url": shortUrl})
}

func (h *Handler) DeleteHash(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")

	var hashRow = models.Hash{}
	foundedHash := h.db.Table("hashes").Where("hash = ?", hash)
	foundedHash.First(&hashRow)

	var hitRows []models.Hit
	h.db.Table("hits").Where("hash_id", hashRow.ID).Delete(&hitRows)
	foundedHash.Delete(&hashRow)

	helpers.RespondWithJSON(w, http.StatusNoContent, map[string]string{"message": "successfully deleted"})
}

func (h *Handler) Redirect(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")

	var hashRow = models.Hash{}
	h.db.Table("hashes").Where("hash = ?", hash).First(&hashRow)

	http.Redirect(w, r, hashRow.Url, http.StatusMovedPermanently)
}

func (h *Handler) GetStats(w http.ResponseWriter, r *http.Request) {
	hash := chi.URLParam(r, "hash")
	since := chi.URLParam(r, "since")
	var sinceFilter time.Time

	if len(since) == 0 {
		sinceFilter = time.Now().Add(-24 * time.Hour)
	} else {
		var err error
		sinceFilter, err = time.Parse("2006-01-02", since)
		if err != nil {
			helpers.RespondWithError(w, http.StatusBadRequest, err.Error())
		}
	}

	var hashRow = models.Hash{}
	h.db.Table("hashes").Where("hash = ?", hash).Where("created_at > ?", sinceFilter).First(&hashRow)

	var hitRows []models.Hit
	h.db.Table("hits").Where("hash_id = ?", hashRow.ID).Find(&hitRows)

	var responseParams = GetStatsResponseParams{url: hashRow.Url, access: make([]map[string]time.Time, len(hitRows))}
	for _, hit := range hitRows {
		if len(responseParams.access) != 0 {
			value := make(map[string]time.Time)
			value[hit.Ip] = hit.AccessAt
			responseParams.access = append(responseParams.access, value)
		}
	}
	jsonResp, err := json.Marshal(responseParams)
	helpers.Catch(err)

	helpers.RespondWithJSON(w, http.StatusOK, jsonResp)
}

func (h *Handler) SetupRouter(db *gorm.DB) *chi.Mux {
	h.db = db
	h.router = chi.NewRouter()
	h.router.Use(middleware.Recoverer)

	h.router.Get("/", h.Welcome)
	h.router.Get("/{hash}", h.Redirect)
	h.router.Get("/v1/{hash}/stats", h.GetStats)
	h.router.Post("/v1/hash", h.GenerateHash)
	h.router.Delete("/v1/{hash}", h.DeleteHash)

	return h.router
}
