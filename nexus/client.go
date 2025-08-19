// nexus/client.go
package nexus

import (
    "bytes"
    "crypto/tls"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "nexuscli/config"
)

// Client represents the Nexus API client
type Client struct {
    BaseURL    string
    Token      string
    Username   string
    Password   string
    HTTPClient *http.Client
}

// NewClient creates a new Nexus API client
func NewClient(cfg config.Config) *Client {
    tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: cfg.Nexus.InsecureSkipVerify},
    }
    httpClient := &http.Client{
        Timeout:   cfg.GetTimeout(),
        Transport: tr,
    }
    return &Client{
        BaseURL:    cfg.Nexus.URL,
        Token:      cfg.Nexus.Token,
        Username:   cfg.Nexus.Username,
        Password:   cfg.Nexus.Password,
        HTTPClient: httpClient,
    }
}

// doRequest performs an HTTP request to the Nexus API
func (c *Client) doRequest(method, path string, body interface{}) ([]byte, error) {
    url := fmt.Sprintf("%s/service/rest/%s", c.BaseURL, path)

    var reqBody io.Reader
    if body != nil {
        jsonBody, err := json.Marshal(body)
        if err != nil {
            return nil, fmt.Errorf("failed to marshal request body: %w", err)
        }
        reqBody = bytes.NewBuffer(jsonBody)
    }
    
    req, err := http.NewRequest(method, url, reqBody)
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }

    req.Header.Set("Content-Type", "application/json")
    
    if c.Token != "" {
        req.Header.Set("Authorization", "NexusAPIKey "+c.Token)
    } else if c.Username != "" && c.Password != "" {
        req.SetBasicAuth(c.Username, c.Password)
    }

    resp, err := c.HTTPClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("failed to execute request to %s: %w", url, err)
    }
    defer resp.Body.Close()

    respBody, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read response body: %w", err)
    }

    if resp.StatusCode < 200 || resp.StatusCode >= 300 {
        return nil, fmt.Errorf("API request failed with status %d: %s", resp.StatusCode, string(respBody))
    }

    return respBody, nil
}

// --- Repository API Calls ---

// CreateRepository creates a new repository (simplified for Docker proxy/hosted)
func (c *Client) CreateRepository(repoType, repoName string) error {
    var payload map[string]interface{}

    switch repoType {
    case "docker-proxy":
        payload = map[string]interface{}{
            "name":    repoName,
            "format":  "docker",
            "type":    "proxy",
            "url":     "https://registry-1.docker.io", // Example proxy URL
            "online":  true,
            "storage": map[string]string{
                "blobStoreName": "default",
                "strictContentTypeValidation": "true",
            },
            "proxy": map[string]interface{}{
                "remoteUrl": "https://registry-1.docker.io",
            },
            "httpClient": map[string]interface{}{
                "blocked": false,
                "autoBlock": true,
            },
            "docker": map[string]interface{}{
                "v1Enabled":      false,
                "forceBasicAuth": false,
                "httpPort":       0, // Or your specific port
                "httpsPort":      0, // Or your specific port
            },
        }
    case "docker-hosted":
        payload = map[string]interface{}{
            "name":    repoName,
            "format":  "docker",
            "type":    "hosted",
            "online":  true,
            "storage": map[string]string{
                "blobStoreName": "default",
                "strictContentTypeValidation": "true",
            },
            "docker": map[string]interface{}{
                "v1Enabled":      false,
                "forceBasicAuth": true, // Usually true for hosted
                "httpPort":       0,
                "httpsPort":      0,
            },
            "cleanup": map[string]interface{}{
                "policyNames": []string{},
            },
            "component": map[string]interface{}{
                "proprietaryComponents": true,
            },
       }
    default:
        return fmt.Errorf("unsupported repository type: %s. Implement more types based on Nexus API.", repoType)
    }

    _, err := c.doRequest(http.MethodPost, "v1/repositories", payload)
    if err != nil {
        return fmt.Errorf("failed to create repository '%s': %w", repoName, err)
    }
    return nil
}

// DeleteRepository deletes a repository
func (c *Client) DeleteRepository(repoName string) error {
    _, err := c.doRequest(http.MethodDelete, fmt.Sprintf("v1/repositories/%s", repoName), nil)
    if err != nil {
        return fmt.Errorf("failed to delete repository '%s': %w", repoName, err)
    }
    return nil
}

// --- User API Calls ---

// CreateUser creates a new user
func (c *Client) CreateUser(username, password, firstName, lastName, email string, roles []string) error {
    payload := map[string]interface{}{
        "userId":    username,
        "firstName": firstName,
        "lastName":  lastName,
        "emailAddress": email,
        "password":  password,
        "status":    "active", // or "locked"
        "roles":     roles,
    }

    _, err := c.doRequest(http.MethodPost, "v1/security/users", payload)
    if err != nil {
        return fmt.Errorf("failed to create user '%s': %w", username, err)
    }
    return nil
}

// DeleteUser deletes a user
func (c *Client) DeleteUser(username string) error {
    _, err := c.doRequest(http.MethodDelete, fmt.Sprintf("v1/security/users/%s", username), nil)
    if err != nil {
        return fmt.Errorf("failed to delete user '%s': %w", username, err)
    }
    return nil
}

// ListRepositories lists all repositories
func (c *Client) ListRepositories() ([]map[string]interface{}, error) {
    body, err := c.doRequest(http.MethodGet, "v1/repositories", nil)
    if err != nil {
        return nil, fmt.Errorf("failed to list repositories: %w", err)
    }

    var repos []map[string]interface{}
    if err := json.Unmarshal(body, &repos); err != nil {
        return nil, fmt.Errorf("failed to unmarshal repositories response: %w", err)
    }
    return repos, nil
}

// ListUsers lists all users
func (c *Client) ListUsers() ([]map[string]interface{}, error) {
    body, err := c.doRequest(http.MethodGet, "v1/security/users", nil)
    if err != nil {
        return nil, fmt.Errorf("failed to list users: %w", err)
    }

    var users []map[string]interface{}
    if err := json.Unmarshal(body, &users); err != nil {
        return nil, fmt.Errorf("failed to unmarshal users response: %w", err)
    }
    return users, nil
}
