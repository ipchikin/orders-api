# Documentation

## Intro

This repo demos a simple order system with 3 APIs - place, take and list order.

## Setup

Please put the db password and Google Distance Matrix API key under `.env` file.

```
MYSQL_ROOT_PASSWORD=your_db_password
MAPS_API_KEY=your_api_key
```

Then run `./start.sh`

## Notes

### Order place

- `id` is in UUID V4 format
- lat, long stored to db will only up to 6 decimal places, even if the input have more decimal places than that

### Order take

- A 5s timeout is added for taking the order, which prevents `SELECT ... FOR UPDATE` to hold the request too long, if other request is taking the same order.

### Order list

- Since not sure about the `limit` range used in the test cases, this api implementation only check if it is >= 0. A 5s timeout is added for getting the order list though.
