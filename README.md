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