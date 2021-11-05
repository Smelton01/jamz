package login

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"

	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"

	"github.com/zmb3/spotify/v2"
)

// redirectURI is the OAuth redirect URI for the application.
// You must register an application at Spotify's developer portal
// and enter this value.
const redirectURI = "http://localhost:8080/callback"

var (
	ch        = make(chan *oauth2.Token)
	state     = "abc123"
	tokenPath string
)

type Acc struct {
	Client *spotify.Client
	Ctx    context.Context
	auth   *spotifyauth.Authenticator
}

func (acc *Acc) Auth() (*spotify.Client, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	tokenPath = filepath.Join(home, ".jamz")

	acc.auth = spotifyauth.New(spotifyauth.WithRedirectURL(redirectURI), spotifyauth.WithScopes(spotifyauth.ScopeUserReadCurrentlyPlaying, spotifyauth.ScopeUserReadPlaybackState, spotifyauth.ScopeUserModifyPlaybackState, spotifyauth.ScopeUserReadRecentlyPlayed))
	token, err := readToken()
	if err != nil {
		if os.IsNotExist(err) || err == ErrInvalidToken {
			if err := acc.login(); err != nil {
				return nil, err
			}
			token, err = readToken()
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}
	client := spotify.New(acc.auth.Client(context.Background(), token))
	return client, err
}

func (acc *Acc) login() error {
	// first start an HTTP server
	http.HandleFunc("/callback", acc.completeAuth)
	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	url := acc.auth.AuthURL(state)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)

	// wait for auth to complete
	token := <-ch
	if err := saveToken(token); err != nil {
		return err
	}
	return nil

}

func (acc *Acc) completeAuth(w http.ResponseWriter, r *http.Request) {
	tok, err := acc.auth.Token(r.Context(), state, r)
	if err != nil {
		http.Error(w, "Couldn't get token", http.StatusForbidden)
		log.Fatal(err)
	}
	if st := r.FormValue("state"); st != state {
		http.NotFound(w, r)
		log.Fatalf("State mismatch: %s != %s\n", st, state)
	}

	ch <- tok
}

func MakeAcc() *Acc {
	return &Acc{Ctx: context.Background()}
}

func saveToken(tok *oauth2.Token) error {
	f, err := os.OpenFile(tokenPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		log.Println(tokenPath, err, tok.TokenType)
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	return enc.Encode(tok)
}

func readToken() (*oauth2.Token, error) {
	content, err := ioutil.ReadFile(tokenPath)
	if err != nil {
		return nil, err
	}

	var tok oauth2.Token
	if err := json.Unmarshal(content, &tok); err != nil {
		return nil, err
	}
	if !tok.Valid() {
		return nil, ErrInvalidToken
	}

	return &tok, nil
}
