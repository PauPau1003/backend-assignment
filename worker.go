package main

import (
    "fmt"
    "sync"
    "github.com/go-resty/resty/v2"
    "github.com/spf13/cobra"
)

type User struct {
    ID        string `json:"id"`
    Title     string `json:"title"`
    FirstName string `json:"firstName"`
    LastName  string `json:"lastName"`
    Email     string `json:"email"`
    Gender    string `json:"gender"`
}

type UsersResponse struct {
    Data []User `json:"data"`
}

type Post struct {
    ID         string   `json:"id"`
    Text       string   `json:"text"`
    Likes      int      `json:"likes"`
    Tags       []string `json:"tags"`
    PublishDate string  `json:"publishDate"`
    Owner      User     `json:"owner"`
}

type PostsResponse struct {
    Data []Post `json:"data"`
}

var workerCmd = &cobra.Command{
    Use:   "worker",
    Short: "Scrape data from Dummy API",
    Run: func(cmd *cobra.Command, args []string) {
        scrapeData()
    },
}

const (
    apiKey  = "667137bfd4f1cc794a6e4c41" // replace with your actual app id
    baseURL = "https://dummyapi.io/data/v1/"
)

func fetchUsers(client *resty.Client, page int, wg *sync.WaitGroup) {
    defer wg.Done()
    url := fmt.Sprintf("%s/user?page=%d&limit=10", baseURL, page)
    var result UsersResponse
    resp, err := client.R().
        SetHeader("app-id", apiKey).
        SetResult(&result).
        Get(url)
    if err != nil {
        fmt.Println("Error fetching users:", err)
        return
    }

    if resp.IsError() {
        fmt.Println("Error response fetching users:", resp.Status())
        return
    }

    var userWg sync.WaitGroup
    for _, user := range result.Data {
        userWg.Add(1)
        go fetchUserDetails(client, user.ID, &userWg)
    }
    userWg.Wait()
}
func fetchUserDetails(client *resty.Client, userID string, wg *sync.WaitGroup) {
    defer wg.Done()
    url := fmt.Sprintf("%s/user/%s", baseURL, userID)
    var user User
    resp, err := client.R().
        SetHeader("app-id", apiKey).
        SetResult(&user).
        Get(url)
    if err != nil {
        fmt.Println("Error fetching user details:", err)
        return
    }

    if resp.IsError() {
        fmt.Println("Error response fetching user details:", resp.Status())
        return
    }

    fmt.Printf("User name %s %s %s %s %s\n", user.Title, user.FirstName, user.LastName, user.Email, user.Gender)
}
func fetchPosts(client *resty.Client, page int, wg *sync.WaitGroup) {
    defer wg.Done()
    url := fmt.Sprintf("%s/post?page=%d&limit=10", baseURL, page)
    var result PostsResponse
    resp, err := client.R().
        SetHeader("app-id", apiKey).
        SetResult(&result).
        Get(url)
    if err != nil {
        fmt.Println("Error fetching posts:", err)
        return
    }

    if resp.IsError() {
        fmt.Println("Error response fetching posts:", resp.Status())
        return
    }

    for _, post := range result.Data {
        fmt.Printf("Posted by %s %s:\n%s\n\nLikes %d Tags %v\nDate posted %s\n", post.Owner.FirstName, post.Owner.LastName, post.Text, post.Likes, post.Tags, post.PublishDate)
    }
}

func scrapeData() {
    client := resty.New()
    var wg sync.WaitGroup

    for i := 1; i <= 10; i++ {
        wg.Add(1)
        go fetchUsers(client, i, &wg)

        wg.Add(1)
        go fetchPosts(client, i, &wg)
    }

    wg.Wait()
}
