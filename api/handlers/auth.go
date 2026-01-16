package handlers

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

var db *sql.DB
func SetDB(database *sql.DB) {
	db = database
}

var (
	githubOauthConfig *oauth2.Config
	googleOauthConfig *oauth2.Config
)

func getGithubOauthConfig() *oauth2.Config {
	if githubOauthConfig != nil {
		return githubOauthConfig
	}
	apiURL := os.Getenv("API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:8080"
	}
	githubOauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GITHUB_CLIENT_ID"),
		ClientSecret: os.Getenv("GITHUB_CLIENT_SECRET"),
		RedirectURL:  apiURL + "/auth/github/callback",
		Scopes:       []string{"user:email"},
		Endpoint:     github.Endpoint,
	}
	return githubOauthConfig
}

func getGoogleOauthConfig() *oauth2.Config {
	if googleOauthConfig != nil {
		return googleOauthConfig
	}
	apiURL := os.Getenv("API_URL")
	if apiURL == "" {
		apiURL = "http://localhost:8080"
	}
	googleOauthConfig = &oauth2.Config{
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  apiURL + "/auth/google/callback",
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}
	return googleOauthConfig
}

func generateStateOauthCookie(c *gin.Context) string {
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	c.SetCookie("oauthstate", state, 600, "/", "", false, true)
	return state
}
func AuthGitHubHandler(c *gin.Context) {
	state := generateStateOauthCookie(c)
	url := getGithubOauthConfig().AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.SetAuthURLParam("prompt", "consent"))
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func AuthGoogleHandler(c *gin.Context) {
	state := generateStateOauthCookie(c)
	url := getGoogleOauthConfig().AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.SetAuthURLParam("prompt", "consent"))
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func verifyState(c *gin.Context) error {
	stateQuery := c.Query("state")
	stateCookie, err := c.Cookie("oauthstate")
	if err != nil {
		return fmt.Errorf("oauth state cookie missing")
	}
	if stateQuery != stateCookie {
		return fmt.Errorf("invalid oauth state")
	}
	return nil
}

func AuthGitHubCallbackHandler(c *gin.Context) {
	if err := verifyState(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid state parameter. Please try again."})
		return
	}

	// Delete state cookie
	c.SetCookie("oauthstate", "", -1, "/", "", false, true)

	code := c.Query("code")
	cliPort := c.Query("cli_port")
	if cliPort == "" {
		cliPort = "10999"
	}

	config := getGithubOauthConfig()
	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to exchange token"})
		return
	}
	client := config.Client(context.Background(), token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user info"})
		return
	}
	defer resp.Body.Close()

	var githubUser struct {
		ID    int    `json:"id"`
		Login string `json:"login"`
		Email string `json:"email"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&githubUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse user info"})
		return
	}

	email := githubUser.Email
	if email == "" {
		email = fmt.Sprintf("%s@github.com", githubUser.Login)
	}

	completeFlow(c, email, "github", fmt.Sprintf("%d", githubUser.ID), cliPort)
}

func AuthGoogleCallbackHandler(c *gin.Context) {
	if err := verifyState(c); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid state parameter"})
		return
	}
	c.SetCookie("oauthstate", "", -1, "/", "", false, true)

	code := c.Query("code")
	cliPort := "10999" // Simplified

	config := getGoogleOauthConfig()
	token, err := config.Exchange(context.Background(), code)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to exchange token"})
		return
	}
	client := config.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user info"})
		return
	}
	defer resp.Body.Close()

	var googleUser struct {
		ID    string `json:"id"`
		Email string `json:"email"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to parse user info"})
		return
	}

	completeFlow(c, googleUser.Email, "google", googleUser.ID, cliPort)
}

func completeFlow(c *gin.Context, email, provider, providerID, cliPort string) {
	var userID int
	err := db.QueryRow(`
		INSERT INTO users (email, provider, provider_id, password_hash) 
		VALUES ($1, $2, $3, NULL)
		ON CONFLICT (email) DO UPDATE 
		SET provider = $2, provider_id = $3
		RETURNING id`,
		email, provider, providerID).Scan(&userID)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
		return
	}

	jwtToken, err := generateJWT(userID, email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	_, err = db.Exec("INSERT INTO tokens (user_id, token) VALUES ($1, $2)", userID, jwtToken)
	if err != nil {
		fmt.Printf("Failed to save token: %v\n", err)
	}
	redirectURL := fmt.Sprintf("http://localhost:%s/callback?token=%s&email=%s", cliPort, jwtToken, email)
	c.Redirect(http.StatusTemporaryRedirect, redirectURL)
}

func generateJWT(userID int, email string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "default-secret-change-me"
	}

	claims := jwt.MapClaims{
		"user_id": userID,
		"email":   email,
		"exp":     time.Now().Add(time.Hour * 24 * 365).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}
