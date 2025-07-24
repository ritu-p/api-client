package apiclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type User struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type UserClient interface {
	CreateUser(user *User) (*User, error)
	UpdateUser(id string, user *User) (*User, error)
	GetUser(id string) (*User, error)
}

type userClient struct {
	baseURL string
	client  *http.Client
}

// NewClient creates a new user API client.
func NewClient(remoteService string) UserClient {
	return &userClient{
		baseURL: remoteService,
		client:  &http.Client{},
	}
}

func (u *userClient) CreateUser(user *User) (*User, error) {
	body, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}
	var resp *http.Response
	maxRetries := 3
	// Add retries to avoid failures when creating a user
	for attempt := 1; attempt <= maxRetries; attempt++ {
		resp, err = u.client.Post(fmt.Sprintf("%s/users", u.baseURL), "application/json", bytes.NewReader(body))
		if err == nil {
			break
		}
		if attempt < maxRetries {
			// Simple backoff
			time.Sleep(100 * time.Millisecond)
		}
	}
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}
	var created User
	if err := json.NewDecoder(resp.Body).Decode(&created); err != nil {
		return nil, err
	}
	return &created, nil
}

func (u *userClient) UpdateUser(id string, user *User) (*User, error) {
	body, err := json.Marshal(user)
	if err != nil {
		return nil, err
	}
	var req *http.Request
	maxRetries := 3
	for attempt := 1; attempt <= maxRetries; attempt++ {
		req, err = http.NewRequest(http.MethodPut, fmt.Sprintf("%s/users/%s", u.baseURL, id), bytes.NewReader(body))
		if err == nil {
			break
		}
		if attempt < maxRetries {
			time.Sleep(100 * time.Millisecond)
		}
	}
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := u.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}
	var updated User
	if err := json.NewDecoder(resp.Body).Decode(&updated); err != nil {
		return nil, err
	}
	return &updated, nil
}

func (u *userClient) GetUser(id string) (*User, error) {
	var resp *http.Response
	var err error
	maxRetries := 3
	for attempt := 1; attempt <= maxRetries; attempt++ {
		resp, err = u.client.Get(fmt.Sprintf("%s/users/%s", u.baseURL, id))
		if err == nil {
			break
		}
		if attempt < maxRetries {
			time.Sleep(100 * time.Millisecond)
		}
	}
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status: %s", resp.Status)
	}
	var user User
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, err
	}
	return &user, nil
}
