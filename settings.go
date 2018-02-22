package main

import (
	"github.com/steam-authority/steam-authority/datastore"
	"github.com/steam-authority/steam-authority/session"
	"github.com/yohcop/openid-go"
	"fmt"
	"net/http"
)

const (
	OPENID   = "http://steamcommunity.com/openid"
	CALLBACK = "http://localhost:8085/login-callback"
	REALM    = "http://localhost:8085/"
)

// todo
// For the demo, we use in-memory infinite storage nonce and discovery
// cache. In your app, do not use this as it will eat up memory and never
// free it. Use your own implementation, on a better database system.
// If you have multiple servers for example, you may need to share at least
// the nonceStore between them.
var nonceStore = openid.NewSimpleNonceStore()
var discoveryCache = openid.NewSimpleDiscoveryCache()

func loginHandler(w http.ResponseWriter, r *http.Request) {

	loggedIn, err := session.IsLoggedIn(r)
	if err != nil {
		returnErrorTemplate(w, 500, err.Error())
		return
	}

	if loggedIn {
		http.Redirect(w, r, "/settings", 303)
		return
	}

	var url string
	url, err = openid.RedirectURL(OPENID, CALLBACK, REALM)
	if err != nil {
		returnErrorTemplate(w, 500, err.Error())
		return
	}

	http.Redirect(w, r, url, 303)
	return
}
func logoutHandler(w http.ResponseWriter, r *http.Request) {

	session.Clear(w, r)

	http.Redirect(w, r, "/", 303)
	return
}

func loginCallbackHandler(w http.ResponseWriter, r *http.Request) {

	fullUrl := "http://localhost:8085" + r.URL.String()
	id, err := openid.Verify(fullUrl, discoveryCache, nonceStore)
	if err != nil {
		returnErrorTemplate(w, 500, err.Error())
		return
	}

	// Save session
	err = session.Write(w, r, "player_id", id)
	if err != nil {
		returnErrorTemplate(w, 500, err.Error())
		return
	}

	http.Redirect(w, r, "/settings", 302)
	return
}

func settingsHandler(w http.ResponseWriter, r *http.Request) {

	loggedIn, err := session.IsLoggedIn(r)
	if err != nil {
		returnErrorTemplate(w, 500, err.Error())
		return
	}

	if !loggedIn {
		http.Redirect(w, r, "/login", 302)
		return
	}

	template := settingsTemplate{}
	template.Session, err = session.ReadAll(r)
	if err != nil {
		returnErrorTemplate(w, 500, err.Error())
		return
	}

	fmt.Println(template.Session)

	returnTemplate(w, "settings", template)

}

func saveSettingsHandler(w http.ResponseWriter, r *http.Request) {

}

type settingsTemplate struct {
	GlobalTemplate
	Session map[interface{}]interface{}
	User    datastore.Player
}
