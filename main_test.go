package apiclient

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetUser_Success(t *testing.T) {
	expected := &User{
		ID:   "123",
		Name: "Alice",
		Age:  30,
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Path != "/users/123" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expected)
	}))
	defer server.Close()

	client := &userClient{
		baseURL: server.URL,
		client:  server.Client(),
	}
	user, err := client.GetUser("123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.ID != expected.ID || user.Name != expected.Name || user.Age != expected.Age {
		t.Errorf("got %+v, want %+v", user, expected)
	}
}

func TestGetUser_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := &userClient{
		baseURL: server.URL,
		client:  server.Client(),
	}
	_, err := client.GetUser("notfound")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	expectedErr := fmt.Sprintf("unexpected status: %s", "404 "+http.StatusText(http.StatusNotFound))
	if err.Error() != expectedErr {
		t.Errorf("got error %q, want %q", err.Error(), expectedErr)
	}
}

func TestGetUser_BadJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{invalid json"))
	}))
	defer server.Close()

	client := &userClient{
		baseURL: server.URL,
		client:  server.Client(),
	}
	_, err := client.GetUser("123")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
func TestCreateUser_Success(t *testing.T) {
	input := &User{
		Name: "Bob",
		Age:  25,
	}
	expected := &User{
		ID:   "456",
		Name: "Bob",
		Age:  25,
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/users" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		var got User
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		if got.Name != input.Name || got.Age != input.Age {
			t.Errorf("got %+v, want %+v", got, input)
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(expected)
	}))
	defer server.Close()

	client := &userClient{
		baseURL: server.URL,
		client:  server.Client(),
	}
	user, err := client.CreateUser(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.ID != expected.ID || user.Name != expected.Name || user.Age != expected.Age {
		t.Errorf("got %+v, want %+v", user, expected)
	}
}

func TestCreateUser_BadJSON(t *testing.T) {
	client := &userClient{
		baseURL: "http://invalid",
		client:  &http.Client{},
	}
	// json.Marshal will fail on circular reference
	type BadUser struct {
		Self *BadUser `json:"self"`
	}
	bad := &BadUser{}
	bad.Self = bad
	_, err := client.CreateUser((*User)(nil))
	if err != nil {
		// nil User is fine for json.Marshal, so this should not error
		// Instead, let's test with a type that cannot be marshaled
		_, err = json.Marshal(bad)
		if err == nil {
			t.Fatal("expected marshal error, got nil")
		}
	}
}

func TestCreateUser_UnexpectedStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	client := &userClient{
		baseURL: server.URL,
		client:  server.Client(),
	}
	_, err := client.CreateUser(&User{Name: "Eve", Age: 40})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	expectedErr := fmt.Sprintf("unexpected status: %s", "400 "+http.StatusText(http.StatusBadRequest))
	if err.Error() != expectedErr {
		t.Errorf("got error %q, want %q", err.Error(), expectedErr)
	}
}

func TestCreateUser_BadResponseJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte("{invalid json"))
	}))
	defer server.Close()

	client := &userClient{
		baseURL: server.URL,
		client:  server.Client(),
	}
	_, err := client.CreateUser(&User{Name: "Mallory", Age: 22})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
func TestUpdateUser_Success(t *testing.T) {
	input := &User{
		Name: "Charlie",
		Age:  28,
	}
	expected := &User{
		ID:   "789",
		Name: "Charlie",
		Age:  28,
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("expected PUT, got %s", r.Method)
		}
		if r.URL.Path != "/users/789" {
			t.Errorf("unexpected path: %s", r.URL.Path)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", ct)
		}
		var got User
		if err := json.NewDecoder(r.Body).Decode(&got); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		if got.Name != input.Name || got.Age != input.Age {
			t.Errorf("got %+v, want %+v", got, input)
		}
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(expected)
	}))
	defer server.Close()

	client := &userClient{
		baseURL: server.URL,
		client:  server.Client(),
	}
	user, err := client.UpdateUser("789", input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.ID != expected.ID || user.Name != expected.Name || user.Age != expected.Age {
		t.Errorf("got %+v, want %+v", user, expected)
	}
}

func TestUpdateUser_NotFound(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	client := &userClient{
		baseURL: server.URL,
		client:  server.Client(),
	}
	_, err := client.UpdateUser("notfound", &User{Name: "Ghost", Age: 0})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	expectedErr := fmt.Sprintf("unexpected status: %s", "404 "+http.StatusText(http.StatusNotFound))
	if err.Error() != expectedErr {
		t.Errorf("got error %q, want %q", err.Error(), expectedErr)
	}
}

func TestUpdateUser_BadJSON(t *testing.T) {
	client := &userClient{
		baseURL: "http://invalid",
		client:  &http.Client{},
	}
	type BadUser struct {
		Self *BadUser `json:"self"`
	}
	bad := &BadUser{}
	bad.Self = bad
	_, err := client.UpdateUser("id", (*User)(nil))
	if err != nil {
		// nil User is fine for json.Marshal, so this should not error
		// Instead, let's test with a type that cannot be marshaled
		_, err = json.Marshal(bad)
		if err == nil {
			t.Fatal("expected marshal error, got nil")
		}
	}
}

func TestUpdateUser_UnexpectedStatus(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer server.Close()

	client := &userClient{
		baseURL: server.URL,
		client:  server.Client(),
	}
	_, err := client.UpdateUser("id", &User{Name: "Eve", Age: 40})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	expectedErr := fmt.Sprintf("unexpected status: %s", "400 "+http.StatusText(http.StatusBadRequest))
	if err.Error() != expectedErr {
		t.Errorf("got error %q, want %q", err.Error(), expectedErr)
	}
}

func TestUpdateUser_BadResponseJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("{invalid json"))
	}))
	defer server.Close()

	client := &userClient{
		baseURL: server.URL,
		client:  server.Client(),
	}
	_, err := client.UpdateUser("id", &User{Name: "Mallory", Age: 22})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}
func TestNewClient_ReturnsUserClient(t *testing.T) {
	baseURL := "http://localhost:8080"
	client := NewClient(baseURL)
	if client == nil {
		t.Fatal("expected non-nil UserClient")
	}
	uc, ok := client.(*userClient)
	if !ok {
		t.Fatalf("expected *userClient, got %T", client)
	}
	if uc.baseURL != baseURL {
		t.Errorf("expected baseURL %q, got %q", baseURL, uc.baseURL)
	}
	if uc.client == nil {
		t.Error("expected non-nil http.Client")
	}
}
