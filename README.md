# Git Pull Agent

A lightweight Go agent designed to continuously synchronize a Git repository by performing periodic `git pull` operations. It's built to be easily configurable via environment variables, making it suitable for containerized environments like Kubernetes.

## Features

*   **Periodic Synchronization:** Automatically pulls the latest changes from a specified Git repository at a configurable interval.
*   **Flexible Authentication:** Supports various standard Git authentication methods:
    *   No authentication (for public repositories).
    *   SSH Key-based authentication (with optional passphrase).
    *   HTTP/HTTPS Basic Authentication (username and password).
*   **Environment Variable Configuration:** All operational parameters are configured exclusively through environment variables, simplifying deployment in container orchestration systems.
*   **Graceful Shutdown:** Handles `SIGINT` and `SIGTERM` signals for clean termination.
*   **Logging:** Provides informative logs for operations and errors.

## Configuration

The agent is configured entirely through the following environment variables:

*   `GIT_REPO_URL` (Required): The URL of the Git repository to synchronize (e.g., `https://github.com/your/repo.git` or `git@github.com:your/repo.git`).
*   `GIT_SYNC_PATH` (Optional): The local path where the repository will be cloned and kept synchronized. Defaults to `/git-repo`.
*   `GIT_PULL_INTERVAL_SECONDS` (Optional): The interval in seconds between `git pull` operations. Defaults to `60` seconds.
*   `GIT_AUTH_METHOD` (Optional): Specifies the authentication method to use.
    *   `none` (Default): No authentication is used. Suitable for public repositories.
    *   `ssh`: Uses an SSH private key for authentication. Requires `GIT_SSH_KEY_PATH`.
    *   `http-basic`: Uses HTTP Basic Authentication. Requires `GIT_HTTP_USERNAME` and `GIT_HTTP_PASSWORD`.

### Authentication-Specific Variables

#### For `GIT_AUTH_METHOD=ssh`
*   `GIT_SSH_KEY_PATH` (Required for SSH): The absolute path to the SSH private key file (e.g., `/etc/ssh/id_rsa`).
*   `GIT_SSH_KEY_PASSPHRASE` (Optional for SSH): The passphrase for the SSH private key, if it is encrypted.

#### For `GIT_AUTH_METHOD=http-basic`
*   `GIT_HTTP_USERNAME` (Required for HTTP Basic): The username for HTTP Basic Authentication.
*   `GIT_HTTP_PASSWORD` (Required for HTTP Basic): The password for HTTP Basic Authentication.

---

### Configuration Examples

#### 1. Public Repository (No Auth)

```bash
export GIT_REPO_URL="https://github.com/kubernetes/git-sync.git"
export GIT_SYNC_PATH="/tmp/my-synced-repo"
export GIT_PULL_INTERVAL_SECONDS="30"

go run main.go
```

#### 2. Private Repository (SSH Auth)

First, ensure your SSH private key is accessible, for instance, by mounting it into a container at `/etc/ssh/id_rsa`.

```bash
export GIT_REPO_URL="git@github.com:your/private-repo.git"
export GIT_SYNC_PATH="/app/repo"
export GIT_PULL_INTERVAL_SECONDS="120"
export GIT_AUTH_METHOD="ssh"
export GIT_SSH_KEY_PATH="/etc/ssh/id_rsa"
# export GIT_SSH_KEY_PASSPHRASE="your_ssh_passphrase" # Uncomment if your key has a passphrase

go run main.go
```

#### 3. Private Repository (HTTP Basic Auth)

```bash
export GIT_REPO_URL="https://github.com/your/private-repo.git"
export GIT_AUTH_METHOD="http-basic"
export GIT_HTTP_USERNAME="your_git_username"
export GIT_HTTP_PASSWORD="your_git_password_or_token"

go run main.go
```

## Building the Agent

To build the executable for your current operating system and architecture:

```bash
go build -o git-pull-agent main.go
```

To build a Linux executable (e.g., for Docker images):

```bash
GOOS=linux GOARCH=amd64 go build -o git-pull-agent main.go
```

## Running the Agent

After building, you can run the agent by executing the binary:

```bash
./git-pull-agent
```

Remember to set the necessary environment variables before running, as shown in the Configuration Examples.

## Building and Running with Docker

### Building the Docker Image

To build the Docker image using the provided `Dockerfile`:

```bash
docker build -t git-pull-agent:latest .
```

### Running the Docker Container

You can run the agent in a Docker container, providing the necessary environment variables.

#### Example: Public Repository (No Auth)

```bash
docker run -d \
  -e GIT_REPO_URL="https://github.com/kubernetes/git-sync.git" \
  -e GIT_SYNC_PATH="/app/repo" \
  -e GIT_PULL_INTERVAL_SECONDS="30" \
  -v "$(pwd)/local-repo:/app/repo" \
  git-pull-agent:latest
```
*   `-d`: Runs the container in detached mode.
*   `-e`: Sets environment variables inside the container.
*   `-v "$(pwd)/local-repo:/app/repo"`: Mounts a local directory `local-repo` to `/app/repo` inside the container. This allows the synchronized repository to persist on your host machine. Make sure the `local-repo` directory exists on your host before running.

#### Example: Private Repository (SSH Auth)

For SSH authentication, you'll need to mount your SSH private key into the container.

```bash
docker run -d \
  -e GIT_REPO_URL="git@github.com:your/private-repo.git" \
  -e GIT_SYNC_PATH="/app/repo" \
  -e GIT_PULL_INTERVAL_SECONDS="120" \
  -e GIT_AUTH_METHOD="ssh" \
  -e GIT_SSH_KEY_PATH="/etc/ssh/id_rsa" \
  -v "$(pwd)/local-repo:/app/repo" \
  -v "$HOME/.ssh/id_rsa:/etc/ssh/id_rsa:ro" \
  git-pull-agent:latest
```
*   `-v "$HOME/.ssh/id_rsa:/etc/ssh/id_rsa:ro"`: Mounts your local SSH private key into the container as a read-only file. Adjust `$HOME/.ssh/id_rsa` to the actual path of your private key.

#### Example: Private Repository (HTTP Basic Auth)

```bash
docker run -d \
  -e GIT_REPO_URL="https://github.com/your/private-repo.git" \
  -e GIT_AUTH_METHOD="http-basic" \
  -e GIT_HTTP_USERNAME="your_git_username" \
  -e GIT_HTTP_PASSWORD="your_git_password_or_token" \
  -v "$(pwd)/local-repo:/app/repo" \
  git-pull-agent:latest
```

This ensures the repository content is accessible outside the container and persists across container restarts.

