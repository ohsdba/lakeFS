package nessie

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/go-openapi/runtime"
	"github.com/go-openapi/swag"
	"github.com/stretchr/testify/require"
	"github.com/treeverse/lakefs/api/gen/client/branches"
	"github.com/treeverse/lakefs/api/gen/client/commits"
	"github.com/treeverse/lakefs/api/gen/client/objects"
	"github.com/treeverse/lakefs/api/gen/client/refs"
	"github.com/treeverse/lakefs/api/gen/models"
)

func TestHooks(t *testing.T) {
	server := startWebhookServer()
	defer func() {
		_ = server.s.Close()
	}()

	ctx, logger, repo := setupTest(t)
	const branch = "feature-1"

	logger.WithField("branch", branch).Info("Create branch")
	ref, err := client.Branches.CreateBranch(
		branches.NewCreateBranchParamsWithContext(ctx).
			WithRepository(repo).
			WithBranch(&models.BranchCreation{
				Name:   swag.String(branch),
				Source: swag.String(masterBranch),
			}), nil)
	require.NoError(t, err, "failed to create branch")
	logger.WithField("branchRef", ref).Info("Branch created")

	logger.WithField("branch", branch).Info("Upload initial content")
	preMergeAction := fmt.Sprintf(`
name: Test Merge
description: set of checks to verify that branch is good
on:
  pre-merge:
	branches:
	  - master
hooks:
  - id: test_webhook
    type: webhook
    description: Check webhooks for pre-merge works
    properties:
      url: "https://%s/pre-merge"
`, webhookServerAddr)
	_, err = client.Objects.UploadObject(
		objects.NewUploadObjectParamsWithContext(ctx).
			WithRepository(repo).
			WithBranch(branch).
			WithPath("_lakefs_actions/testing_pre_merge").
			WithContent(runtime.NamedReader("content", strings.NewReader(preMergeAction))), nil)
	require.NoError(t, err)

	preCommitAction := fmt.Sprintf(`
name: Test Commit
description: set of checks to verify that branch is good
on:
  pre-commit:
    branches:
      - feature-*
hooks:
  - id: test_webhook
    type: webhook
    description: Check webhooks for pre-commit works
    properties:
      url: "https://%s/pre-commit"
`, webhookServerAddr)
	_, err = client.Objects.UploadObject(
		objects.NewUploadObjectParamsWithContext(ctx).
			WithRepository(repo).
			WithBranch(branch).
			WithPath("_lakefs_actions/testing_pre_commit").
			WithContent(runtime.NamedReader("content", strings.NewReader(preCommitAction))), nil)
	require.NoError(t, err)
	logger.WithField("branch", branch).Info("Commit initial content")
	stats, err := client.Commits.Commit(
		commits.NewCommitParamsWithContext(ctx).
			WithRepository(repo).
			WithBranch(branch).
			WithCommit(&models.CommitCreation{Message: swag.String("Initial content")}),
		nil)
	webhookData := <-server.commit
	require.NoError(t, err, "failed to commit initial content")
	select {
	case err = <-server.errCh:
	default:
		err = nil
	}
	require.NoError(t, err, "error on pre commit serving")
	decoder := json.NewDecoder(bytes.NewReader(webhookData))
	var commitEvent, mergeEvent webhookEventInfo
	require.NoError(t, decoder.Decode(&commitEvent), "reading pre-commit data")

	require.Equal(t, "pre-commit", commitEvent.EventType)
	require.Equal(t, "Test Commit", commitEvent.ActionName)
	require.Equal(t, "test_webhook", commitEvent.HookID)
	require.Equal(t, repo, commitEvent.RepositoryID)
	require.Equal(t, branch, commitEvent.BranchID)
	require.Equal(t, stats.Payload.Committer, commitEvent.Committer)
	require.Equal(t, stats.Payload.Message, commitEvent.CommitMessage)
	require.Equal(t, branch, commitEvent.SourceRef)
	require.Equal(t, stats.Payload.Metadata, commitEvent.Metadata)

	mergeRes, err := client.Refs.MergeIntoBranch(
		refs.NewMergeIntoBranchParamsWithContext(ctx).WithRepository(repo).WithDestinationBranch(masterBranch).WithSourceRef(branch), nil)

	webhookData = <-server.merge
	logger.WithField("mergeResult", mergeRes).Info("Merged successfully")
	require.NoError(t, err, "failed to merge branches")
	select {
	case err = <-server.errCh:
	default:
		err = nil
	}

	require.NoError(t, err, "error on pre commit serving")
	decoder = json.NewDecoder(bytes.NewReader(webhookData))
	require.NoError(t, decoder.Decode(&mergeEvent), "reading pre-merge data")

	require.Equal(t, "pre-merge", mergeEvent.EventType)
	require.Equal(t, "Test Merge", mergeEvent.ActionName)
	require.Equal(t, "test_webhook", mergeEvent.HookID)
	require.Equal(t, repo, mergeEvent.RepositoryID)
	require.Equal(t, branch, mergeEvent.BranchID)
	require.Equal(t, branch, mergeEvent.SourceRef)
}

const webhookServerAddr = "127.0.0.1:8081"

type server struct {
	s      *http.Server
	commit chan []byte
	merge  chan []byte
	errCh  chan error
}

func startWebhookServer() *server {
	commit := make(chan []byte)
	merge := make(chan []byte)
	errCh := make(chan error, 2)
	http.HandleFunc("/pre-commit", hookHandlerFunc(commit, errCh))
	http.HandleFunc("/pre-merge", hookHandlerFunc(merge, errCh))
	s := &http.Server{
		Addr: webhookServerAddr,
	}
	go func() {
		if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			fmt.Printf("server failed to listen on %s: %v\n", webhookServerAddr, err)
			os.Exit(1)
		}
	}()

	return &server{
		s:      s,
		commit: commit,
		merge:  merge,
		errCh:  errCh,
	}
}

func hookHandlerFunc(resCh chan<- []byte, errCh chan<- error) func(http.ResponseWriter, *http.Request) {
	return func(writer http.ResponseWriter, request *http.Request) {
		bytes, err := ioutil.ReadAll(request.Body)
		if err != nil {
			errCh <- err
			_, _ = io.WriteString(writer, "Failed")
			writer.WriteHeader(http.StatusInternalServerError)
			return
		}
		resCh <- bytes
		_, _ = io.WriteString(writer, "OK")
		writer.WriteHeader(http.StatusOK)
		return
	}

}

type webhookEventInfo struct {
	EventType     string            `json:"event_type"`
	EventTime     string            `json:"event_time"`
	ActionName    string            `json:"action_name"`
	HookID        string            `json:"hook_id"`
	RepositoryID  string            `json:"repository_id"`
	BranchID      string            `json:"branch_id"`
	SourceRef     string            `json:"source_ref"`
	CommitMessage string            `json:"commit_message"`
	Committer     string            `json:"committer"`
	Metadata      map[string]string `json:"metadata"`
}
