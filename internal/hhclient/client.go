package hhclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

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
}

// -------------------- Методы --------------------

// GetToken — получает access_token
func (c *Client) GetToken() (string, error) {
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
		return "", err
	}
	defer resp.Body.Close()

	var token tokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&token); err != nil {
		return "", err
	}

	return token.AccessToken, nil
}

// SearchVacancies — ищет вакансии и возвращает список
func (c *Client) SearchVacancies(token, text string) ([]Vacancy, error) {
	req, err := http.NewRequest(
		"GET",
		"https://api.hh.ru/vacancies?text="+url.QueryEscape(text),
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var vacanciesResp vacanciesResponse
	if err := json.NewDecoder(resp.Body).Decode(&vacanciesResp); err != nil {
		return nil, err
	}

	return vacanciesResp.Items, nil
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

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("ошибка отклика: %s", resp.Status)
	}

	return nil
}
