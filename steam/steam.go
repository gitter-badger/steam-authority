package steam

import (
	"errors"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/Jleagle/go-helpers/logger"
)

func get(path string, query url.Values, useKey ...bool) (bytes []byte, err error) {

	if path != "" {
		query.Add("format", "json")

		if !(len(useKey) > 0 && !useKey[0]) {
			query.Add("key", os.Getenv("STEAM_API_KEY"))
		}
		path = "http://api.steampowered.com/" + path
	} else {
		path = "http://store.steampowered.com/api/appdetails"
	}

	// Debug
	logger.Info("STEAM: " + path + "?" + strings.Replace(query.Encode(), os.Getenv("STEAM_API_KEY"), "_", 1))

	// Grab the JSON from node
	response, err := http.Get(path + "?" + query.Encode())
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	// Convert to bytes
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	// todo, test this works
	if string(contents) == "<html><head><title>Forbidden</title></head><body><h1>Forbidden</h1>Access is denied. Retrying will not help. Please verify your <pre>key=</pre> parameter.</body></html>" {
		errors.New("invalid api key")
	}

	return contents, nil
}

func GetID(in string) (out string, err error) {

	if regexp.MustCompile(`^STEAM_(0|1):(0|1):[0-9][0-9]{0,8}$`).MatchString(in) { // STEAM_0:0:4180232
		return convert0to64(in), nil
	} else if regexp.MustCompile(`(\[)?U:1:\d+(\])?`).MatchString(in) { // [U:1:8360464]
		return convert3to64(in), nil
	} else if regexp.MustCompile(`^\d{17}$`).MatchString(in) { // 76561197968626192
		return in, nil
	} else if regexp.MustCompile(`^\d{1,16}$`).MatchString(in) { // 8360464
		return convert32to64(in), nil
	} else {

		resp, err := ResolveVanityURL(in)
		if err != nil {
			return out, err
		}

		return resp.Response.SteamID, nil
	}
}

func convert3to64(in string) (out string) {
	parts := strings.Split(in, ":")
	part := parts[2]
	part = part[:len(part)-1] // Remove bracket
	return convert32to64(part)
}

func convert32to64(in string) (out string) {

	inBig, _ := new(big.Int).SetString(in, 10)
	mul, _ := new(big.Int).SetString("76561197960265728", 10)

	return inBig.Add(inBig, mul).String()
}

func convert0to64(in string) (out string) {

	parts := strings.Split(in, ":")
	add, _ := new(big.Int).SetString("76561197960265728", 10)
	level, _ := new(big.Int).SetString(parts[1], 10)

	ID64, _ := new(big.Int).SetString(parts[2], 10)
	ID64 = ID64.Mul(ID64, big.NewInt(2))
	ID64 = ID64.Add(ID64, add)
	ID64 = ID64.Add(ID64, level)

	return ID64.String()
}
