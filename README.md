```bash

docker-compose up -d

```

Influx password reset

```bash

# login into influx container
docker exec -it <container id> bash
cat ./var/lib/influxdb2/influxd.bolt

# find there admin token and use it to reset password
docker exec -it <container id> influx user password -n admin -t <token>

```

To add annotations the following query can be used
```
from(bucket: "default")
  |> range(start: v.timeRangeStart, stop:v.timeRangeStop)
  |> filter(fn: (r) =>
    r._measurement == "events" and
    (r._field == "text" or
    r._field == "title")
  )
```
# Thermostat Monitoring
- [go-thermostat-monitor](https://github.com/blakehartshorn/go-thermostat-monitor)
- [pronestheus](https://github.com/grdl/pronestheus)

## Nest Dev account
- Follow instructions [here](https://developers.google.com/nest/device-access/authorize#google_hasnt_verified_this_app)
- [Google Cloud Console](https://console.cloud.google.com/apis/credentials)
- [Device Access Console](https://console.nest.google.com/device-access/project-list)
    - [Device Access Registration](https://developers.google.com/nest/device-access/registration)


First, we need to set the project ID and the client ID and secret. To do that, run the following command:
```bash
NEST_CLIENT_ID= && \
NEST_CLIENT_SECRET= && \
NEST_PROJECT_ID=
```
or use following when .env file contains the variables:
```bash
set -o allexport; source .env; set +o allexport
```

To get authorization code run the following command:
```bash
echo "https://nestservices.google.com/partnerconnections/${NEST_PROJECT_ID}/auth?redirect_uri=https://www.google.com&access_type=offline&prompt=consent&client_id=${NEST_CLIENT_ID}&response_type=code&scope=https://www.googleapis.com/auth/sdm.service" | pbcopy
```

Open a browser and paste the copied URL. You will be redirected to a page with an authorization code. Copy the code and run the following command to get the tokens:

```bash
AUTHZ_CODE=... && \
TOKENS_JSON=$(curl -L -X POST "https://www.googleapis.com/oauth2/v4/token?client_id=${CLIENT_ID}&client_secret=${CLIENT_SECRET}&code=${AUTHZ_CODE}&grant_type=authorization_code&redirect_uri=https://www.google.com") && \
ACCESS_TOKEN=$(echo "$TOKENS_JSON" | jq -r '.access_token') && \
REFRESH_TOKEN=$(echo "$TOKENS_JSON" | jq -r '.refresh_token')
```

Now we can get the list of devices:
```bash
curl -X GET "https://smartdevicemanagement.googleapis.com/v1/enterprises/${PROJECT_ID}/devices" \
    -H 'Content-Type: application/json' \
    -H "Authorization: Bearer ${ACCESS_TOKEN}"
```

To refresh token:
```bash
ACCESS_TOKEN=$(curl -X POST "https://www.googleapis.com/oauth2/v4/token?client_id=${CLIENT_ID}&client_secret=${CLIENT_SECRET}&refresh_token=${REFRESH_TOKEN}&grant_type=refresh_token" | jq -r '.access_token' )
```


# Telegram Bots

## Webhook To Use in IFTTT

Send a message to the bot to get the chat id

```bash
# set bot token
WH_TGBOT_TOKEN=...

# get chat id
CHAT_ID=$(curl -X GET "https://api.telegram.org/bot${WH_TGBOT_TOKEN}/getUpdates" | jq -r '.result[0].channel_post.sender_chat.id')


curl -X POST "https://api.telegram.org/bot${WH_TGBOT_TOKEN}/sendMessage" \
    -H 'Content-Type: application/json' \
    -d '{"chat_id":"<chat id>", "text":"Hi"}'
```

# Open Weather

```bash
# set api key
OW_API_KEY=... && \
LON=... && \
LAT=...

# get weather
curl -X GET "https://api.openweathermap.org/data/2.5/weather?lat=${LAT}&lon=${LON}&appid=${OW_API_KEY}"
```
