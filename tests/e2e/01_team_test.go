package e2e_test

import (
	"avito-backend-intern-assignment/internal/app/api"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTeam_AddAndGet(t *testing.T) {
	body := api.Team{
		TeamName: "payments",
		Members: []api.TeamMember{
			{UserId: "u1", Username: "Alice", IsActive: true},
			{UserId: "u2", Username: "Bob", IsActive: true},
		},
	}

	data, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/team/add", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()

	testRouter.ServeHTTP(rec, req)
	if rec.Code != 201 {
		t.Fatalf("expected 201, got %d, body=%s", rec.Code, rec.Body.String())
	}

	var created api.PostTeamAdd201JSONResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &created); err != nil {
		t.Fatalf("invalid json: %v", err)
	}

	if created.Team.TeamName != body.TeamName {
		t.Fatalf("team name mismatch: expected %s, got %s", body.TeamName, created.Team.TeamName)
	}

	getReq := httptest.NewRequest(http.MethodGet, "/team/get?team_name=payments", nil)
	getRec := httptest.NewRecorder()
	testRouter.ServeHTTP(getRec, getReq)

	if getRec.Code != 200 {
		t.Fatalf("expected 200, got %d", getRec.Code)
	}

	var got api.GetTeamGet200JSONResponse
	if err := json.Unmarshal(getRec.Body.Bytes(), &got); err != nil {
		t.Fatalf("invalid json: %v", err)
	}

	if len(got.Members) != 2 {
		t.Fatalf("expected 2 members, got %d", len(got.Members))
	}
}

func TestTeam_Add_UpdateMembers(t *testing.T) {
	initial := api.Team{
		TeamName: "backend",
		Members: []api.TeamMember{
			{UserId: "u3", Username: "Charlie", IsActive: true},
		},
	}

	b1, _ := json.Marshal(initial)
	req1 := httptest.NewRequest(http.MethodPost, "/team/add", bytes.NewBuffer(b1))
	req1.Header.Set("Content-Type", "application/json")

	rec1 := httptest.NewRecorder()
	testRouter.ServeHTTP(rec1, req1)

	if rec1.Code != 201 {
		t.Fatalf("expected 201, got %d", rec1.Code)
	}

	updated := api.Team{
		TeamName: "backend",
		Members: []api.TeamMember{
			{UserId: "u3", Username: "Charlie", IsActive: true},
			{UserId: "u4", Username: "Dave", IsActive: true},
		},
	}

	b2, _ := json.Marshal(updated)
	req2 := httptest.NewRequest(http.MethodPost, "/team/add", bytes.NewBuffer(b2))
	req2.Header.Set("Content-Type", "application/json")

	rec2 := httptest.NewRecorder()
	testRouter.ServeHTTP(rec2, req2)

	if rec2.Code != 201 {
		t.Fatalf("expected 201 on update, got %d", rec2.Code)
	}

	getReq := httptest.NewRequest(http.MethodGet, "/team/get?team_name=backend", nil)
	getRec := httptest.NewRecorder()
	testRouter.ServeHTTP(getRec, getReq)

	if getRec.Code != 200 {
		t.Fatalf("expected 200, got %d", getRec.Code)
	}

	var got api.GetTeamGet200JSONResponse
	_ = json.Unmarshal(getRec.Body.Bytes(), &got)

	if len(got.Members) != 2 {
		t.Fatalf("expected 2 members after update, got %d", len(got.Members))
	}
}

func TestTeam_Add_InvalidBody(t *testing.T) {
	invalidJSON := []byte(`{"team_name": 123}`)

	req := httptest.NewRequest(http.MethodPost, "/team/add", bytes.NewBuffer(invalidJSON))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	testRouter.ServeHTTP(rec, req)

	if rec.Code != 400 {
		t.Fatalf("expected 400, got %d (%s)", rec.Code, rec.Body.String())
	}
}

func TestTeam_Get_NotFound(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/team/get?team_name=no_such_team", nil)
	rec := httptest.NewRecorder()
	testRouter.ServeHTTP(rec, req)

	if rec.Code != 404 {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestTeam_Add_EmptyMembers(t *testing.T) {
	body := api.Team{
		TeamName: "emptyteam",
		Members:  []api.TeamMember{},
	}

	data, _ := json.Marshal(body)

	req := httptest.NewRequest(http.MethodPost, "/team/add", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	testRouter.ServeHTTP(rec, req)

	if rec.Code != 201 {
		t.Fatalf("expected 201, got %d", rec.Code)
	}

	getReq := httptest.NewRequest(http.MethodGet, "/team/get?team_name=emptyteam", nil)
	getRec := httptest.NewRecorder()
	testRouter.ServeHTTP(getRec, getReq)

	if getRec.Code != 200 {
		t.Fatalf("expected 200, got %d", getRec.Code)
	}

	var got api.GetTeamGet200JSONResponse
	_ = json.Unmarshal(getRec.Body.Bytes(), &got)

	if len(got.Members) != 0 {
		t.Fatalf("expected 0 members, got %d", len(got.Members))
	}
}
