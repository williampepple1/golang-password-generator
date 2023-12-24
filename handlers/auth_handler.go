package handlers

import (
	"context"
	"net/http"
	"os"
	"password-generator/models"

	"encoding/json"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gorm.io/gorm"
)

var googleOauthConfig = &oauth2.Config{
	RedirectURL:  "http://localhost:8080/auth/google/callback",
	ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
	Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
	Endpoint:     google.Endpoint,
}

// GoogleUserInfo holds the structure of the user info we receive from Google
type GoogleUserInfo struct {
	Id            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"given_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
	Locale        string `json:"locale"`
}

func getUserInfo(accessToken string) (*GoogleUserInfo, error) {
	userInfoUrl := "https://www.googleapis.com/oauth2/v2/userinfo"

	// Prepare the request
	req, err := http.NewRequest("GET", userInfoUrl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	// Execute the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Decode the response
	var userInfo GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return nil, err
	}

	return &userInfo, nil
}
func GoogleLogin(c *gin.Context) {
	url := googleOauthConfig.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	c.Redirect(http.StatusTemporaryRedirect, url)
}

func GoogleCallback(c *gin.Context, db *gorm.DB) {
	code := c.Query("code")
	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	userInfo, err := getUserInfo(token.AccessToken)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// Create or update user in your database
	user := models.User{
		GoogleID: userInfo.Id,
		Email:    userInfo.Email,
		Name:     userInfo.Name,
	}
	db.FirstOrCreate(&user, models.User{GoogleID: userInfo.Id})

	// Implement session creation or JWT token generation
	// Set session or JWT token in response

	c.Redirect(http.StatusTemporaryRedirect, "/passwords")
}

// getUserInfo function is assumed to be implemented
// It should make a request to Google's userinfo endpoint and return the user information
