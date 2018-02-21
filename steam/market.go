package steam

import (
	"net/url"

	"github.com/Jleagle/go-helpers/logger"
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

	logger.Info(string(bytes))
}
