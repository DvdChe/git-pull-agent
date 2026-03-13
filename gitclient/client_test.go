package gitclient

import (
	"net/url"
	"os"
	"path/filepath"
	"strings" 
	"testing"
	"time" 

	"github.com/DvdChe/git-pull-agent/config"

	"github.com/go-git/go-git/v5"
	gitconfig "github.com/go-git/go-git/v5/config" // Added import
	"github.com/go-git/go-git/v5/plumbing/object"
)

// setupBareRemote creates a bare Git repository in a temporary directory
func setupBareRemote(t *testing.T, remotePath string) string {
	_, err := git.PlainInit(remotePath, true)
	if err != nil {
		t.Fatalf("Failed to initialize bare remote repository: %v", err)
	}
	return remotePath
}

// commitFileToRemote creates a temporary worktree, commits a file, and pushes to the bare remote
func commitFileToRemote(t *testing.T, bareRemotePath, filename, content string) {
	tempLocalRepoPath := t.TempDir()
	
	// Initialize a non-bare repository
	r, err := git.PlainInit(tempLocalRepoPath, false)
	if err != nil {
		t.Fatalf("Failed to initialize local repository: %v", err)
	}

	w, err := r.Worktree()
	if err != nil {
		t.Fatalf("Failed to get worktree from local repository: %v", err)
	}

	// Create and commit a file
	filePath := filepath.Join(tempLocalRepoPath, filename)
	err = os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to write file to local repository: %v", err)
	}

	_, err = w.Add(filename)
	if err != nil {
		t.Fatalf("Failed to add file to local repository: %v", err)
	}

	_, err = w.Commit("Initial commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Test User",
			Email: "test@example.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		t.Fatalf("Failed to commit to local repository: %v", err)
	}

	// Add the bare remote as a remote to the local repository
	_, err = r.CreateRemote(&gitconfig.RemoteConfig{ // Corrected here
		Name: "origin",
		URLs: []string{(&url.URL{Scheme: "file", Path: bareRemotePath}).String()},
	})
	if err != nil {
		t.Fatalf("Failed to create remote: %v", err)
	}

	// Push to the bare remote
	err = r.Push(&git.PushOptions{
		RemoteName: "origin",
		Progress:  os.Stdout,
	})
	if err != nil {
		t.Fatalf("Failed to push to bare remote: %v", err)
	}
}

// updateBareRemote clones the bare remote, modifies a file, commits, and pushes back

func updateBareRemote(t *testing.T, bareRemotePath, filename, content string) {

	tempLocalRepoPath := t.TempDir()



	// Clone the bare remote into a temporary worktree

	r, err := git.PlainClone(tempLocalRepoPath, false, &git.CloneOptions{

		URL:      (&url.URL{Scheme: "file", Path: bareRemotePath}).String(),

		Progress: os.Stdout,

	})

	if err != nil {

		t.Fatalf("Failed to clone bare remote to temporary worktree for update: %v", err)

	}



	w, err := r.Worktree()

	if err != nil {

		t.Fatalf("Failed to get worktree from temporary clone for update: %v", err)

	}



	// Modify and commit a file

	filePath := filepath.Join(tempLocalRepoPath, filename)

	err = os.WriteFile(filePath, []byte(content), 0644)

	if err != nil {

		t.Fatalf("Failed to write file to temporary worktree for update: %v", err)

	}



	_, err = w.Add(filename)

	if err != nil {

		t.Fatalf("Failed to add file to temporary worktree for update: %v", err)

	}



	_, err = w.Commit("Update commit", &git.CommitOptions{

		Author: &object.Signature{

			Name:  "Test User",

			Email: "test@example.com",

			When:  time.Now(),

		},

	})

	if err != nil {

		t.Fatalf("Failed to commit to temporary worktree for update: %v", err)

	}



	// Push to the bare remote

	err = r.Push(&git.PushOptions{

		RemoteName: "origin", // Assuming "origin" is the remote name

		Progress:  os.Stdout,

	})

	if err != nil {

		t.Fatalf("Failed to push to bare remote for update: %v", err)

	}

}



func TestCloneOrPull_CloneNewRepo(t *testing.T) {

	// Setup bare remote

	bareRemotePath := setupBareRemote(t, t.TempDir())

	commitFileToRemote(t, bareRemotePath, "testfile.txt", "initial content")



	// Setup config for cloning

	syncPath := t.TempDir()

	cfg := &config.Config{

		RepoURL:    (&url.URL{Scheme: "file", Path: bareRemotePath}).String(),

		SyncPath:   syncPath,

		AuthMethod: config.AuthMethodNone,

	}



	// Perform clone

	err := CloneOrPull(cfg)

	if err != nil {

		t.Fatalf("CloneOrPull failed on clone: %v", err)

	}



	// Verify cloned content

	content, err := os.ReadFile(filepath.Join(syncPath, "testfile.txt"))

	if err != nil {

		t.Fatalf("Failed to read cloned file: %v", err)

	}

	if string(content) != "initial content" {

		t.Errorf("Cloned file content mismatch: expected 'initial content', got '%s'", string(content))

	}

}



func TestCloneOrPull_PullExistingRepo(t *testing.T) {

	// Setup bare remote

	bareRemotePath := setupBareRemote(t, t.TempDir())

	commitFileToRemote(t, bareRemotePath, "testfile.txt", "initial content")



	// Setup config for initial clone

	syncPath := t.TempDir()

	cfg := &config.Config{

		RepoURL:    (&url.URL{Scheme: "file", Path: bareRemotePath}).String(),

		SyncPath:   syncPath,

		AuthMethod: config.AuthMethodNone,

	}



	// Initial clone

	err := CloneOrPull(cfg)

	if err != nil {

		t.Fatalf("Initial CloneOrPull failed: %v", err)

	}



	// Make a change in the remote using the new helper

	updateBareRemote(t, bareRemotePath, "testfile.txt", "updated content")



	// Perform pull

	err = CloneOrPull(cfg)

	if err != nil {

		t.Fatalf("CloneOrPull failed on pull: %v", err)

	}



	// Verify pulled content

	content, err := os.ReadFile(filepath.Join(syncPath, "testfile.txt"))

	if err != nil {

		t.Fatalf("Failed to read pulled file: %v", err)

	}

	if string(content) != "updated content" {

		t.Errorf("Pulled file content mismatch: expected 'updated content', got '%s'", string(content))

	}

}



func TestCloneOrPull_AlreadyUpToDate(t *testing.T) {

	// Setup bare remote

	bareRemotePath := setupBareRemote(t, t.TempDir())

	commitFileToRemote(t, bareRemotePath, "testfile.txt", "initial content")



	// Setup config for initial clone

	syncPath := t.TempDir()

	cfg := &config.Config{

		RepoURL:    (&url.URL{Scheme: "file", Path: bareRemotePath}).String(),

		SyncPath:   syncPath,

		AuthMethod: config.AuthMethodNone,

	}



	// Initial clone

	err := CloneOrPull(cfg)

	if err != nil {

		t.Fatalf("Initial CloneOrPull failed: %v", err)

	}



	// Second pull, should be up-to-date

	err = CloneOrPull(cfg)

	if err != nil { // git.NoErrAlreadyUpToDate is handled internally and not returned

		t.Fatalf("CloneOrPull failed on up-to-date check: %v", err)

	}

	// No explicit check for output, relying on no error being returned

}



func TestCloneOrPull_InvalidRepoURL(t *testing.T) {

	syncPath := t.TempDir()

	cfg := &config.Config{

		RepoURL:    "invalid-url",

		SyncPath:   syncPath,

		AuthMethod: config.AuthMethodNone,

	}



	err := CloneOrPull(cfg)

	if err == nil {

		t.Fatal("Expected error for invalid repo URL, got nil")

	}

	expectedErrorSubstring := "failed to clone repository"

	if !strings.Contains(err.Error(), expectedErrorSubstring) {

		t.Errorf("Expected error to contain '%s', got '%s'", expectedErrorSubstring, err.Error())

	}

}
