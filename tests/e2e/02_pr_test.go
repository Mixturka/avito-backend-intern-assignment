package e2e_test

import (
	"avito-backend-intern-assignment/internal/app/api"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestPullRequest_Create_Success(t *testing.T) {
	body := api.PostPullRequestCreateJSONBody{
		PullRequestId:   "pr-test",
		PullRequestName: "Add search",
		AuthorId:        "u1",
	}

	data, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/pullRequest/create", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	testRouter.ServeHTTP(rec, req)

	if rec.Code != 201 {
		t.Fatalf("expected 201, got %d, body=%s", rec.Code, rec.Body.String())
	}

	var resp api.PostPullRequestCreate201JSONResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid json: %v", err)
	}

	if resp.Pr.PullRequestId != body.PullRequestId {
		t.Fatalf("expected pr_id %s, got %s", body.PullRequestId, resp.Pr.PullRequestId)
	}
	if len(resp.Pr.AssignedReviewers) > 2 {
		t.Fatalf("assigned reviewers exceed 2: %v", resp.Pr.AssignedReviewers)
	}
}

func TestPullRequest_Create_AlreadyExists(t *testing.T) {
	body := api.PostPullRequestCreateJSONBody{
		PullRequestId:   "pr-test",
		PullRequestName: "Add search",
		AuthorId:        "u1",
	}
	data, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/pullRequest/create", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	testRouter.ServeHTTP(rec, req)

	if rec.Code != 409 {
		t.Fatalf("expected 409, got %d", rec.Code)
	}
}

func TestPullRequest_Merge_Success(t *testing.T) {
	body := api.PostPullRequestMergeJSONBody{
		PullRequestId: "pr-test",
	}

	data, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/pullRequest/merge", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	testRouter.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	var resp api.PostPullRequestMerge200JSONResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid json: %v", err)
	}

	if resp.Pr.Status != api.PullRequestStatusMERGED {
		t.Fatalf("expected MERGED, got %s", resp.Pr.Status)
	}
}

func TestPullRequest_Merge_NotFound(t *testing.T) {
	body := api.PostPullRequestMergeJSONBody{
		PullRequestId: "pr-404",
	}
	data, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/pullRequest/merge", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	testRouter.ServeHTTP(rec, req)

	if rec.Code != 404 {
		t.Fatalf("expected 404, got %d", rec.Code)
	}
}

func TestPullRequest_Reassign_Success(t *testing.T) {
	updateTeamBody := api.Team{
		TeamName: "payments",
		Members: []api.TeamMember{
			{UserId: "u1", Username: "Alice", IsActive: true},
			{UserId: "u2", Username: "Bob", IsActive: true},
			{UserId: "u3", Username: "Charlie", IsActive: true},
			{UserId: "u4", Username: "Dave", IsActive: true},
		},
	}

	updateTeamData, _ := json.Marshal(updateTeamBody)
	updateTeamReq := httptest.NewRequest(http.MethodPost, "/team/add", bytes.NewBuffer(updateTeamData))
	updateTeamReq.Header.Set("Content-Type", "application/json")
	updateTeamRec := httptest.NewRecorder()
	testRouter.ServeHTTP(updateTeamRec, updateTeamReq)
	if updateTeamRec.Code != 201 {
		t.Fatalf("failed to update team, got %d", updateTeamRec.Code)
	}

	newPRBody := api.PostPullRequestCreateJSONBody{
		PullRequestId:   "pr-reassign-test",
		PullRequestName: "Reassign test PR",
		AuthorId:        "u1",
	}

	newPRData, _ := json.Marshal(newPRBody)
	newPRReq := httptest.NewRequest(http.MethodPost, "/pullRequest/create", bytes.NewBuffer(newPRData))
	newPRReq.Header.Set("Content-Type", "application/json")
	newPRRec := httptest.NewRecorder()
	testRouter.ServeHTTP(newPRRec, newPRReq)
	if newPRRec.Code != 201 {
		t.Fatalf("failed to create new PR for reassignment test, got %d", newPRRec.Code)
	}

	body := api.PostPullRequestReassignJSONBody{
		PullRequestId: "pr-reassign-test",
		OldUserId:     "u2",
	}

	data, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/pullRequest/reassign", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	testRouter.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Fatalf("expected 200, got %d, body=%s", rec.Code, rec.Body.String())
	}

	var resp api.PostPullRequestReassign200JSONResponse
	if err := json.Unmarshal(rec.Body.Bytes(), &resp); err != nil {
		t.Fatalf("invalid json: %v", err)
	}

	if resp.ReplacedBy == "" {
		t.Fatalf("expected replaced_by user_id, got empty")
	}

	t.Logf("Reassigned reviewer: %s", resp.ReplacedBy)
}

func TestPullRequest_Reassign_NotAssigned(t *testing.T) {
	body := api.PostPullRequestReassignJSONBody{
		PullRequestId: "pr-test",
		OldUserId:     "u404",
	}

	data, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/pullRequest/reassign", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	testRouter.ServeHTTP(rec, req)

	if rec.Code != 409 {
		t.Fatalf("expected 409, got %d", rec.Code)
	}
}

func TestPullRequest_Reassign_MergedPR(t *testing.T) {
	body := api.PostPullRequestReassignJSONBody{
		PullRequestId: "pr-test",
		OldUserId:     "u3",
	}

	data, _ := json.Marshal(body)
	req := httptest.NewRequest(http.MethodPost, "/pullRequest/reassign", bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	rec := httptest.NewRecorder()
	testRouter.ServeHTTP(rec, req)

	if rec.Code != 409 {
		t.Fatalf("expected 409 for merged PR, got %d", rec.Code)
	}
}
