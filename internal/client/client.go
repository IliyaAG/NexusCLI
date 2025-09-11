package client

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "time"
)

type NexusClient struct {
    baseURL  string
    username string
    password string
    token    string
    timeout  time.Duration
    client   *http.Client
    verbose  int
}

func NewNexusClient(url, username, password, token string, timeoutSec, verbose int) *NexusClient {
    timeout := time.Duration(timeoutSec) * time.Second
    return &NexusClient{
        baseURL: url,
        username: username,
        password: password,
        token:   token,
        timeout: timeout,
        verbose: verbose,
        client: &http.Client{
            Timeout: timeout,
        },
    }
}

// ---------------- USER ---------------- //

func (c *NexusClient) CreateUser(username, password, firstName, lastName, email string, roles []string) error {
    body := map[string]interface{}{
        "userId":    username,
        "firstName": firstName,
        "lastName":  lastName,
        "emailAddress": email,
        "password":  password,
        "status":    "active",
        "roles":     roles,
    }

    return c.post("/service/rest/v1/security/users", body)
}

func (c *NexusClient) DeleteUser(username string) error {
    return c.delete("/service/rest/v1/security/users/" + username)
}

func (c *NexusClient) ListUsers() ([]map[string]interface{}, error) {
    data, err := c.get("/service/rest/v1/security/users")
    if err != nil {
        return nil, err
    }
    var users []map[string]interface{}
    if err := json.Unmarshal(data, &users); err != nil {
        return nil, err
    }
    return users, nil
}

// ---------------- REPO ---------------- //

func (c *NexusClient) CreateRepository(repoType, name string) error {
    body := map[string]interface{}{
        "name":   name,
        "online": true,
        "recipe": fmt.Sprintf("%s-hosted", repoType),
    }
    return c.post("/service/rest/v1/repositories/"+repoType+"/hosted", body)
}

func (c *NexusClient) DeleteRepository(name string) error {
    return c.delete("/service/rest/v1/repositories/" + name)
}

func (c *NexusClient) ListRepositories() ([]map[string]interface{}, error) {
    data, err := c.get("/service/rest/v1/repositories")
    if err != nil {
        return nil, err
    }
    var repos []map[string]interface{}
    if err := json.Unmarshal(data, &repos); err != nil {
        return nil, err
    }
    return repos, nil
}

// ---------------- LOW LEVEL ---------------- //

func (c *NexusClient) addAuth(req *http.Request) {
    if c.token != "" {
        req.Header.Set("Authorization", "Bearer "+c.token)
    }else if c.username != "" && c.password != "" {
        req.SetBasicAuth(c.username, c.password)
    }
}

func (c *NexusClient) logRequest(req *http.Request, body []byte) {
    if c.verbose >= 1 {
        fmt.Printf("[HTTP] %s %s\n", req.Method, req.URL.String())
    }
    if c.verbose >= 2 {
        fmt.Println("Headers:")
        for k, v := range req.Header {
            fmt.Printf("  %s: %v\n", k, v)
        }
        if len(body) > 0 {
            fmt.Println("Body:")
            fmt.Println(string(body))
        }
    }
}

func (c *NexusClient) logResponse(resp *http.Response, data []byte) {
    if c.verbose >= 1 {
        fmt.Printf("[HTTP] Response %s\n", resp.Status)
    }
    if c.verbose >= 2 {
        fmt.Println("Response Headers:")
        for k, v := range resp.Header {
            fmt.Printf("  %s: %v\n", k, v)
        }
        if len(data) > 0 {
            fmt.Println("Response Body:")
            fmt.Println(string(data))
        }
    }
}

func (c *NexusClient) get(path string) ([]byte, error) {
    req, err := http.NewRequest("GET", c.baseURL+path, nil)
    if err != nil {
        return nil, err
    }
    c.addAuth(req)
    c.logRequest(req, nil)

    resp, err := c.client.Do(req)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    data, _ := ioutil.ReadAll(resp.Body)

    c.logResponse(resp, data)
    return data, nil
}

func (c *NexusClient) post(path string, body map[string]interface{}) error {
    data, _ := json.Marshal(body)
    req, err := http.NewRequest("POST", c.baseURL+path, bytes.NewBuffer(data))
    if err != nil {
        return err
    }
    c.addAuth(req)
    req.Header.Set("Content-Type", "application/json")
    c.logRequest(req, data)

    resp, err := c.client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    respData, _ := ioutil.ReadAll(resp.Body)

    c.logResponse(resp, respData)
    if resp.StatusCode >= 300 {
        return fmt.Errorf("Nexus returned status %s", resp.Status)
    }
    return nil
}

func (c *NexusClient) delete(path string) error {
    req, err := http.NewRequest("DELETE", c.baseURL+path, nil)
    if err != nil {
        return err
    }
    c.addAuth(req)
    c.logRequest(req, nil)

    resp, err := c.client.Do(req)
    if err != nil {
        return err
    }
    defer resp.Body.Close()
    data, _ := ioutil.ReadAll(resp.Body)

    c.logResponse(resp, data)
    if resp.StatusCode >= 300 {
        return fmt.Errorf("Nexus returned status %s", resp.Status)
    }
    return nil
}
