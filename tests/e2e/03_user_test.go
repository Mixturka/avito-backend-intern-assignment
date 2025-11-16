package e2e_test

import (
	"avito-backend-intern-assignment/internal/app/api"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUser_SetIsActive(t *testing.T) {
	body := map[string]any{
		"user_id":   "u2",
		"is_active": false,
	}
	data, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/users/setIsActive", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	testRouter.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d, body=%s", rec.Code, rec.Body.String())
	}

	var resp struct {
		User api.User `json:"user"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid json: %v", err)
	}

	if resp.User.UserId != "u2" {
		t.Fatalf("expected user_id u2, got %s", resp.User.UserId)
	}
	if resp.User.IsActive != false {
		t.Fatalf("expected is_active false, got %v", resp.User.IsActive)
	}
}

func TestUser_SetIsActive_NotFound(t *testing.T) {
	body := map[string]any{
		"user_id":   "no_such_user",
		"is_active": true,
	}
	data, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/users/setIsActive", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	testRouter.ServeHTTP(rec, req)

	if rec.Code != 404 {
		t.Fatalf("expected 404, got %d, body=%s", rec.Code, rec.Body.String())
	}
}

func TestUser_GetReview(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/users/getReview?user_id=u2", nil)
	rec := httptest.NewRecorder()
	testRouter.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d, body=%s", rec.Code, rec.Body.String())
	}

	var resp struct {
		UserID       string                 `json:"user_id"`
		PullRequests []api.PullRequestShort `json:"pull_requests"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid json: %v", err)
	}

	if resp.UserID != "u2" {
		t.Fatalf("expected user_id u2, got %s", resp.UserID)
	}

	if len(resp.PullRequests) == 0 {
		t.Fatalf("expected at least 1 pull request for u2")
	}
}

func TestUser_GetReview_NoPRs(t *testing.T) {
	teamBody := api.Team{
		TeamName: "no_pr_team",
		Members: []api.TeamMember{
			{UserId: "no_pr_user", Username: "No PR User", IsActive: true},
		},
	}

	teamData, _ := json.Marshal(teamBody)
	teamReq := httptest.NewRequest(http.MethodPost, "/team/add", bytes.NewBuffer(teamData))
	teamReq.Header.Set("Content-Type", "application/json")
	teamRec := httptest.NewRecorder()
	testRouter.ServeHTTP(teamRec, teamReq)
	if teamRec.Code != 201 {
		t.Fatalf("failed to create team with user no_pr_user, got %d", teamRec.Code)
	}

	req := httptest.NewRequest(http.MethodGet, "/users/getReview?user_id=no_pr_user", nil)
	rec := httptest.NewRecorder()
	testRouter.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d, body=%s", rec.Code, rec.Body.String())
	}

	var resp struct {
		UserID       string                 `json:"user_id"`
		PullRequests []api.PullRequestShort `json:"pull_requests"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid json: %v", err)
	}

	if resp.UserID != "no_pr_user" {
		t.Fatalf("expected user_id no_pr_user, got %s", resp.UserID)
	}

	if len(resp.PullRequests) != 0 {
		t.Fatalf("expected 0 pull requests for no_pr_user, got %d", len(resp.PullRequests))
	}
}
