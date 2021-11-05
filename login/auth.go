package login

import (
	"context"
	"fmt"
	"log"
	"net/http"

	spotifyauth "github.com/zmb3/spotify/v2/auth"

	"github.com/zmb3/spotify/v2"
)

// redirectURI is the OAuth redirect URI for the application.
// You must register an application at Spotify's developer portal
// and enter this value.
const redirectURI = "http://localhost:8080/callback"

var (
	ch    = make(chan *spotify.Client)
	state = "abc123"
)

type Acc struct {
	Client *spotify.Client
	Ctx    context.Context
	auth   *spotifyauth.Authenticator
}

func (acc *Acc) Auth() {
	acc.auth = spotifyauth.New(spotifyauth.WithRedirectURL(redirectURI), spotifyauth.WithScopes(spotifyauth.ScopeUserReadCurrentlyPlaying, spotifyauth.ScopeUserReadPlaybackState, spotifyauth.ScopeUserModifyPlaybackState, spotifyauth.ScopeUserReadRecentlyPlayed))
	// first start an HTTP server
	http.HandleFunc("/callback", acc.completeAuth)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Println("Got request for:", r.URL.String())
	})
	go func() {
		err := http.ListenAndServe(":8080", nil)
		if err != nil {
			log.Fatal(err)
		}
	}()

	url := acc.auth.AuthURL(state)
	fmt.Println("Please log in to Spotify by visiting the following page in your browser:", url)

	// wait for auth to complete
	client := <-ch

	acc.Client = client
	// use the client to make calls that require authorization
	user, err := client.CurrentUser(acc.Ctx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("You are logged in as:", user.ID)
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

	// use the token to get an authenticated client
	client := spotify.New(acc.auth.Client(r.Context(), tok))
	fmt.Fprintf(w, "Login Completed!")
	ch <- client
}

func MakeAcc() *Acc {
	return &Acc{Ctx: context.Background()}
}
