package gitclient

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/DvdChe/git-pull-agent/config"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/plumbing/transport/ssh"
	go_ssh "golang.org/x/crypto/ssh"
)

// CloneOrPull performs a git clone if the repository does not exist locally,
// otherwise it performs a git pull.
func CloneOrPull(cfg *config.Config) error {
	repoPath := cfg.SyncPath
	_, err := os.Stat(filepath.Join(repoPath, ".git"))

	var auth transport.AuthMethod
	if cfg.AuthMethod != config.AuthMethodNone {
		auth, err = setupAuth(cfg)
		if err != nil {
			return fmt.Errorf("failed to set up authentication: %w", err)
		}
	}

	if os.IsNotExist(err) {
		// Repository does not exist, so clone it
		fmt.Printf("Cloning repository %s into %s...\n", cfg.RepoURL, repoPath)
		_, err := git.PlainClone(repoPath, false, &git.CloneOptions{
			URL:      cfg.RepoURL,
			Auth:     auth,
			Progress: os.Stdout,
		})
		if err != nil {
			return fmt.Errorf("failed to clone repository: %w", err)
		}
		fmt.Println("Repository cloned successfully.")
	} else if err != nil {
		return fmt.Errorf("failed to check repository path: %w", err)
	} else {
		// Repository exists, so pull
		fmt.Printf("Pulling updates for repository in %s...\n", repoPath)
		r, err := git.PlainOpen(repoPath)
		if err != nil {
			return fmt.Errorf("failed to open repository: %w", err)
		}

		w, err := r.Worktree()
		if err != nil {
			return fmt.Errorf("failed to get worktree: %w", err)
		}

		err = w.Pull(&git.PullOptions{
			Auth:     auth,
			Progress: os.Stdout,
		})
		if err != nil && err != git.NoErrAlreadyUpToDate {
			return fmt.Errorf("failed to pull repository: %w", err)
		}
		if err == git.NoErrAlreadyUpToDate {
			fmt.Println("Repository already up-to-date.")
		} else {
			fmt.Println("Repository pulled successfully.")
		}
	}
	return nil
}

func setupAuth(cfg *config.Config) (transport.AuthMethod, error) {
	switch cfg.AuthMethod {
	case config.AuthMethodSSH:
		sshKey, err := os.ReadFile(cfg.SSHKeyPath)
		if err != nil {
			return nil, fmt.Errorf("failed to read SSH key from %s: %w", cfg.SSHKeyPath, err)
		}

		signer, err := go_ssh.ParsePrivateKeyWithPassphrase(sshKey, []byte(cfg.SSHKeyPassphrase))
		if err != nil {
			return nil, fmt.Errorf("failed to parse SSH private key: %w", err)
		}
		
		auth := &ssh.PublicKeys{Signer: signer}
		// Consider adding support for known_hosts or InsecureSkipHostKeyCheck in the future
		return auth, nil

	case config.AuthMethodHTTPBasic:
		return &http.BasicAuth{
			Username: cfg.HTTPUsername,
			Password: cfg.HTTPPassword,
		}, nil
	default:
		return nil, fmt.Errorf("unsupported authentication method: %s", cfg.AuthMethod)
	}
}
