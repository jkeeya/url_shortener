package transport

import (
	"context"
	"net/http"
	u "net/url"
	"regexp"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

type createAliasRequest struct {
	URL string `json:"url"`
}

type createAliasResponse struct {
	URL       string `json:"url"`
	ShortLink string `json:"short_link"`
	Existed   bool   `json:"existed"`
}

func (h *Handlers) CreateAlias(c echo.Context) error {
	var request createAliasRequest

	if err := c.Bind(&request); err != nil || strings.TrimSpace(request.URL) == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid or missing url"})
	}
	url := validateURL(request.URL)

	ctx, cancel := context.WithTimeout(c.Request().Context(), 3*time.Second)
	defer cancel()

	existed, alias, err := h.svc.CreateAlias(ctx, url)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	response := createAliasResponse{
		URL:       alias.URL,
		ShortLink: alias.ShortLink,
		Existed:   existed,
	}
	if existed {
		return c.JSON(http.StatusOK, response)
	}
	return c.JSON(http.StatusCreated, response)
}

type FindShortByURLRequest struct {
	URL string `json:"url"`
}
type FindShortByURLResponse struct {
	ShortLink string `json:"short_link"`
	Exist     bool   `json:"exist"`
}

func (h *Handlers) FindShortByURL(c echo.Context) error {
	var response FindShortByURLResponse
	url := c.QueryParam("url")
	if url == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid or missing url"})
	}

	ctx, cancel := context.WithTimeout(c.Request().Context(), 3*time.Second)
	defer cancel()

	shortLink, err := h.svc.GetShortByURL(ctx, url)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	exist := shortLink != ""

	response = FindShortByURLResponse{
		ShortLink: shortLink,
		Exist:     exist,
	}
	return c.JSON(http.StatusOK, response)
}

func validateURL(raw string) string {
	raw = strings.TrimSpace(raw)

	if !strings.Contains(raw, "://") {
		raw = "https://" + raw
	}

	parsed, err := u.ParseRequestURI(raw)
	if err != nil || parsed.Host == "" {
		return ""
	}

	// Проверяем что домен в нормальном формате
	host := parsed.Hostname()
	domainRE := regexp.MustCompile(`^(?:[A-Za-z0-9](?:[A-Za-z0-9-]{0,61}[A-Za-z0-9])?\.)+[A-Za-z]{2,63}$`)
	if !domainRE.MatchString(host) {
		return ""
	}

	// Проверяем доступность https, если недоступен - подставляем http
	client := &http.Client{Timeout: 1 * time.Second}
	resp, err := client.Head(raw)
	if err != nil {
		raw = strings.Replace(raw, "https://", "http://", 1)
	} else {
		resp.Body.Close()
	}

	return raw
}
