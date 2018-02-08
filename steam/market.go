package steam

import (
	"fmt"
	"net/url"
)

/**
https://partner.steamgames.com/doc/webapi/IEconMarketService#GetPopular
*/

func GetPopularItems() {

	options := url.Values{}
	options.Set("language", "")
	options.Set("rows", "20")
	options.Set("start", "")
	options.Set("filter_appid", "")
	options.Set("ecurrency", "")

	bytes, err := get("IEconMarketService/GetPopular/v1/", options)
	if err != nil {
		return
	}

	fmt.Println(string(bytes))

}
