package hhclient

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const defaultUserAgent = "autoclick-hh-bot/1.0 (+contact:region-manager43@yandex.ru)"

func newHTTPClient() *http.Client {
	return &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) > 0 {
				prev := via[len(via)-1]
				req.Header.Set("Authorization", prev.Header.Get("Authorization"))
				req.Header.Set("User-Agent", prev.Header.Get("User-Agent"))
				req.Header.Set("Accept", prev.Header.Get("Accept"))
			}
			return nil
		},
	}
}

type Client struct {
	clientID     string
	clientSecret string
	username     string
	password     string
}

func New(clientID, clientSecret, username, password string) *Client {
	return &Client{
		clientID:     clientID,
		clientSecret: clientSecret,
		username:     username,
		password:     password,
	}
}

// -------------------- Структуры --------------------

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

type Vacancy struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Area struct {
		Name string `json:"name"`
	} `json:"area"`
	Employer struct {
		Name string `json:"name"`
	} `json:"employer"`
}

type vacanciesResponse struct {
	Items []Vacancy `json:"items"`
	Found int       `json:"found"`
	Page  int       `json:"page"`
	Pages int       `json:"pages"`
}

// -------------------- Методы --------------------

// GetToken — получает access_token с фоллбеком Basic Auth
func (c *Client) GetToken() (string, error) {
	// Попытка №1: параметры в body (классический вариант)
	token, triedBody, err := c.getTokenWithBody()
	if err == nil && token != "" {
		log.Printf("token prefix: %q", safePrefix(token))
		return token, nil
	}
	// Если сервер вернул unsupported_grant_type/400 — пробуем Basic
	if triedBody {
		fallback, err2 := c.getTokenWithBasic()
		if err2 == nil && fallback != "" {
			log.Printf("token prefix: %q", safePrefix(fallback))
			return fallback, nil
		}
		if err2 != nil {
			return "", err2
		}
	}
	return "", err
}

func safePrefix(s string) string {
	if len(s) >= 8 {
		return s[:8]
	}
	return s
}

func (c *Client) getTokenWithBody() (string, bool, error) {
	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("client_id", c.clientID)
	data.Set("client_secret", c.clientSecret)
	data.Set("username", c.username)
	data.Set("password", c.password)

	resp, err := http.Post(
		"https://hh.ru/oauth/token",
		"application/x-www-form-urlencoded",
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		return "", true, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 500))
		// Сообщаем вызывающему, что это была попытка с body (для решения о фоллбеке)
		return "", true, fmt.Errorf("token %s: %s", resp.Status, string(body))
	}

	var token tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return "", true, err
	}
	return token.AccessToken, true, nil
}

func (c *Client) getTokenWithBasic() (string, error) {
	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("username", c.username)
	data.Set("password", c.password)

	req, err := http.NewRequest(
		"POST",
		"https://hh.ru/oauth/token",
		strings.NewReader(data.Encode()),
	)
	if err != nil {
		return "", err
	}
	// Basic <base64(client_id:client_secret)>
	basic := base64.StdEncoding.EncodeToString([]byte(c.clientID + ":" + c.clientSecret))
	req.Header.Set("Authorization", "Basic "+basic)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", defaultUserAgent)

	client := newHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 500))
		return "", fmt.Errorf("token(basic) %s: %s", resp.Status, string(body))
	}

	var token tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return "", err
	}
	return token.AccessToken, nil
}

// RefreshToken — обновление access_token по refresh_token
func (c *Client) RefreshToken(refreshToken string) (string, string, error) {
	data := url.Values{}
	data.Set("grant_type", "refresh_token")
	data.Set("refresh_token", refreshToken)
	data.Set("client_id", c.clientID)
	data.Set("client_secret", c.clientSecret)

	req, err := http.NewRequest("POST", "https://hh.ru/oauth/token", strings.NewReader(data.Encode()))
	if err != nil {
		return "", "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", defaultUserAgent)

	client := newHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 500))
		return "", "", fmt.Errorf("refresh %s: %s", resp.Status, string(body))
	}

	var token tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return "", "", err
	}
	return token.AccessToken, token.RefreshToken, nil
}

// SearchVacancies — пагинированный поиск
func (c *Client) SearchVacancies(token string, keywords []string) ([]Vacancy, error) {
	perPage := 100
	page := 0
	var all []Vacancy
	client := newHTTPClient()

	for {
		req, err := http.NewRequest("GET", "https://api.hh.ru/vacancies", nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("User-Agent", defaultUserAgent)
		req.Header.Set("Accept", "application/json")

		q := req.URL.Query()
		q.Add("text", strings.Join(keywords, " OR "))
		q.Add("area", "113")
		q.Add("search_field", "name")
		// Фильтруем только удалённые вакансии
		q.Add("schedule", "remote")
		q.Add("per_page", strconv.Itoa(perPage))
		q.Add("page", strconv.Itoa(page))
		req.URL.RawQuery = q.Encode()

		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		if resp.StatusCode != http.StatusOK {
			body, _ := io.ReadAll(io.LimitReader(resp.Body, 500))
			resp.Body.Close()
			return nil, fmt.Errorf("vacancies %s: %s", resp.Status, string(body))
		}

		if resp.StatusCode == http.StatusForbidden {
			body, _ := io.ReadAll(io.LimitReader(resp.Body, 500))
			if strings.Contains(string(body), "token-expired") {
				return nil, fmt.Errorf("access_token истёк")
			}
		}

		var vr vacanciesResponse
		if err = json.NewDecoder(resp.Body).Decode(&vr); err != nil {
			resp.Body.Close()
			return nil, err
		}
		resp.Body.Close()

		if page == 0 {
			log.Printf("HH URL: %s", req.URL.String())
			log.Printf("Найдено всего (found): %d, pages: %d", vr.Found, vr.Pages)
		}

		all = append(all, vr.Items...)
		page++
		if page >= vr.Pages || len(vr.Items) == 0 {
			break
		}
	}

	log.Printf("Найдено вакансий (суммарно): %d", len(all))
	for i := 0; i < len(all) && i < 3; i++ {
		log.Printf("Пример #%d: %s", i+1, all[i].Name)
	}
	return all, nil
}

// ApplyVacancy — отправка отклика на вакансию
func (c *Client) ApplyVacancy(token, vacancyID, resumeID string) error {
	req, err := http.NewRequest("POST", "https://api.hh.ru/negotiations", nil)
	if err != nil {
		return err
	}

	q := req.URL.Query()
	q.Add("vacancy_id", vacancyID)
	q.Add("resume_id", resumeID)
	req.URL.RawQuery = q.Encode()

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", defaultUserAgent)
	req.Header.Set("Accept", "application/json")

	client := newHTTPClient()
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 500))
		return fmt.Errorf("ошибка отклика: %s: %s", resp.Status, string(body))
	}
	return nil
}
