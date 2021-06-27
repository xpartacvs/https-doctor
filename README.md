# HTTPS Doctor

A https health doctor

## Configuration

I was designed to get configured by environment variables.

| **Environment Variable** | **Type**  | **Req** | **Default**                           | **Description**                                                                                                                                                    |
| :---                     | :---      | :---:   | :---                                  | :---                                                                                                                                                               |
| `DISHOOK_URL`            | `string`  | √       |                                       | Discord webhook's URL.                                                                                                                                             |
| `DISHOOK_BOT_NAME`       | `string`  |         | `HTTPS Doctor`                        | Discord webhook bot's display name.                                                                                                                                |
| `DISHOOK_BOT_AVATAR`     | `string`  |         | default URL                           | Discord webhook bot's avatar URL.                                                                                                                                  |
| `DISHOOK_BOT_MESSAGE`    | `string`  |         | `Your HTTPS health monitoring result` | Discord webhook bot's alert message.                                                                                                                               |
| `GRACEPERIOD`            | `integer` |         | `14`                                  | Number of days before the host's SSL certificate get expired. When the current time is in the range, alert will get triggered.                                     |
| `HOSTS`                  | `string`  | √       | empty string                          | Coma separated list of hosts to check on.                                                                                                                          |
| `LOGLEVEL`               | `string`  |         | `disabled`                            | The logging mode: `debug`, `info`, `warn`, `error`, and `disabled`.                                                                                                |
| `SCHEDULE`               | `string`  |         | `0 0 * * *`                           | The HTTPS health checking schedule in CRON format.                                                                                                                 |
| `TZ`                     | `string`  |         | local system                          | The timezone. Must contain one of [IANA Time Zone database](https://en.wikipedia.org/wiki/List_of_tz_database_time_zones) associate to your preferred time format. |
