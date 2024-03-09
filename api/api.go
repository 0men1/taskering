package api

import (
	"context"
	"encoding/json"
	"log"
	"fmt"
	"net/http"
	"os"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/option"
	"google.golang.org/api/tasks/v1"
	"regexp"
)



func GetSrvs() (*tasks.Service) {
	ctx := MakeContext()
	b := ReadCreds()
	config := GetConfig(b, tasks.TasksScope)
	client := GetClient(config)
	tasksrv, _ := tasks.NewService(ctx, option.WithHTTPClient(client))
	return tasksrv
}


func MakeContext() context.Context {
	return context.Background()	
}


func ReadCreds() ([]byte) {
	b, err := os.ReadFile("./api/credentials.json"); if err != nil {
		log.Fatalf("Could not parse credntials.json: %v", err)
	} 	
	return b
}


func GetTasksSrv(ctx context.Context, opts option.ClientOption) (*tasks.Service) {
	srv, err := tasks.NewService(ctx, opts); if err != nil {
		log.Fatalf("There was an error grabbing tasks service: %v", err)
	}
	return srv
}


func GetConfig(jsonKey []byte, scope ...string) (*oauth2.Config) {
	conf, err := google.ConfigFromJSON(jsonKey, scope[0]); if err != nil {
		log.Fatalf("There was an error fetching config: %v", err)
	}
	return conf
}

// Retrieve a token, saves the token, then returns the generated client.
func GetClient(config *oauth2.Config) *http.Client {
	// The file token.json stores the user's access and refresh tokens, and is
	// created automatically when the authorization flow completes for the first
	// time.
	tokFile := "./api/token.json"
	tok, err := TokenFromFile(tokFile)
	if err != nil {
		tok = GetTokenFromWeb(config)
		SaveToken(tokFile, tok)
	}
	return config.Client(context.Background(), tok)
}


// Request a token from the web, then returns the retrieved token.
func GetTokenFromWeb(config *oauth2.Config) *oauth2.Token {
	authURL := config.AuthCodeURL("state-token", oauth2.AccessTypeOffline)
	log.Printf("Go to the following link in your browser then copy paste the "+
		"url you were redirected to.: \n%v\n", authURL)
		
	var link string
	if _, err := fmt.Scan(&link); err != nil {
		log.Fatalf("Unable to read authorization code: %v", err)
	}
	

	//Regex to filter key
	r, _ := regexp.Compile("code=([^&]+)")

	authCode := r.FindStringSubmatch(link)

	fmt.Println(authCode[1])


	tok, err := config.Exchange(context.TODO(), authCode[1])

	if err != nil {
		log.Fatalf("Unable to retrieve token from web: %v", err)
	}
	return tok
}

// Retrieves a token from a local file.
func TokenFromFile(file string) (*oauth2.Token, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	tok := &oauth2.Token{}
	err = json.NewDecoder(f).Decode(tok)
	return tok, err
}

// Saves a token to a file path.
func SaveToken(path string, token *oauth2.Token) {
	fmt.Printf("Saving credential file to: %s\n", path)
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Fatalf("Unable to cache oauth token: %v", err)
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	enc.SetIndent("", "	")
	if err := enc.Encode(token); err != nil {
		panic(err)
	}
}

