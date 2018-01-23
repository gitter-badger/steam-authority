# Steam Authority

#### Running on local
export STEAM_GOOGLE_PROJECT="google-cloud-project-id"
export STEAM_GOOGLE_APPLICATION_CREDENTIALS="/path/to/key.json"
export STEAM_PICS_PATH="/path/to/steam-authority/steam-pics-api"
export STEAM_GITHUB_TOKEN="github-token" # https://github.com/settings/tokens

##### Setup
- git clone git@github.com:steam-authority/steam-authority.git
- dep ensure
- go build
