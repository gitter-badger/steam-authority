package main

import (
	"github.com/steam-authority/steam-authority/datastore"
	"github.com/steam-authority/steam-authority/session"
	"github.com/steam-authority/steam-authority/steam"
	"github.com/yohcop/openid-go"
	"net/http"
	"os"
	"path"
	"strconv"
)

const (
	ID     = "id"
	NAME   = "name"
	AVATAR = "avatar"
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
		returnErrorTemplate(w, r, 500, err.Error())
		return
	}

	if loggedIn {
		http.Redirect(w, r, "/settings", 303)
		return
	}

	var url string
	url, err = openid.RedirectURL("http://steamcommunity.com/openid", os.Getenv("STEAM_DOMAIN")+"/login-callback", os.Getenv("STEAM_DOMAIN")+"/")
	if err != nil {
		returnErrorTemplate(w, r, 500, err.Error())
		return
	}

	http.Redirect(w, r, url, 303)
	return
}
func loginCallbackHandler(w http.ResponseWriter, r *http.Request) {

	openID, err := openid.Verify(os.Getenv("STEAM_DOMAIN")+r.URL.String(), discoveryCache, nonceStore)
	if err != nil {
		returnErrorTemplate(w, r, 500, err.Error())
		return
	}

	id, err := strconv.Atoi(path.Base(openID))
	if err != nil {
		returnErrorTemplate(w, r, 500, err.Error())
		return
	}

	resp, err := steam.GetPlayerSummaries([]int{id})
	if err != nil {
		returnErrorTemplate(w, r, 500, err.Error())
		return
	}

	session.WriteMany(w, r, map[string]string{
		ID:     strconv.Itoa(id),
		NAME:   resp.Response.Players[0].PersonaName,
		AVATAR: resp.Response.Players[0].AvatarMedium,
	})

	http.Redirect(w, r, "/settings", 302)
	return
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {

	session.Clear(w, r)
	http.Redirect(w, r, "/", 303)
	return
}

func settingsHandler(w http.ResponseWriter, r *http.Request) {

	loggedIn, err := session.IsLoggedIn(r)
	if err != nil {
		returnErrorTemplate(w, r, 500, err.Error())
		return
	}

	if !loggedIn {
		http.Redirect(w, r, "/login", 302)
		return
	}

	template := settingsTemplate{}
	template.SetSession(r)
	if err != nil {
		returnErrorTemplate(w, r, 500, err.Error())
		return
	}

	returnTemplate(w, r, "settings", template)

}

func saveSettingsHandler(w http.ResponseWriter, r *http.Request) {

}

type settingsTemplate struct {
	GlobalTemplate
	User datastore.Player
}
