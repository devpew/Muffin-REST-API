# Muffin API

## Authentication

### POST

All the Merhods in the API protected by JSON Web Tokens (JTW).

We must send field "Token" in Header. 

The value of the token shoud be generated based on the passphrase. The passphrase placed in the REST API server.

To get authenticate token we should use this request

```
/login
```

and send Username and Password, like that:

```
{
    "Username": "user",
    "Password": "pass"
}
```

If user and password are correct we will get 200 OK and reqponse with our key, and after that we should save it in localStorage:

```
HTTP/1.1 200 OK
Access-Control-Allow-Origin: *
Content-Type: application/json
Date: Fri, 10 Jul 2020 20:42:19 GMT
Content-Length: 163
Connection: close

"eykjlskdfijiI3djk3kd1idjkjfd.eysljFJejejf9DJlePppdjf8fkdjjdks.Dlfklksf183oowlsfid9edjoksrjf"
```

If user or pass is not correct we will get 200 OK with the error:

```
HTTP/1.1 200 OK
Access-Control-Allow-Origin: *
Content-Type: application/json
Date: Sun, 12 Jul 2020 15:20:25 GMT
Content-Length: 8
Connection: close

"error"
```


## Funds 

### GET

GET shares rub - `/funds/rub/shares`

GET bonds rub - `/funds/rub/bonds`

GET shares use -  `/funds/usd/shares`

GET bonds usd -  `/funds/usd/bonds` 

http://127.0.0.1:8000/funds/rub/shares

responce:

```
[
  {
    "id": 131,
    "name": "AGRO",
    "ticker": "AGRO",
    "amount": 32,
    "priceperitem": 103.48,
    "purchaseprice": 31360,
    "pricecurrent": 33113.6,
    "percentchanges": 5.5918367346938735,
    "yearlyinvestment": 10.327319061645275,
    "clearmoney": 1721.3631999999986,
    "datepurchase": "2019-12-09T13:59:53.66277+03:00",
    "datelastupdate": "2020-06-30T18:00:49+03:00",
    "type": "share"
  },
  {
    "id": 145,
    "name": "Детский мир",
    "ticker": "DSKY",
    "amount": 450,
    "priceperitem": 103.48,
    "purchaseprice": 41885,
    "pricecurrent": 46566,
    "percentchanges": 11.1758386057061,
    "yearlyinvestment": 40.81456061148672,
    "clearmoney": 4636.7745,
    "datepurchase": "2020-03-13T13:59:53.66277+03:00",
    "datelastupdate": "2020-06-19T23:54:26.655833+03:00",
    "type": "share"
  },
  {
    "id": 138,
    "name": "Газпром",
    "ticker": "GAZP",
    "amount": 300,
    "priceperitem": 195.81,
    "purchaseprice": 68100,
    "pricecurrent": 58743,
    "percentchanges": -13.740088105726873,
    "yearlyinvestment": -35.809394273127744,
    "clearmoney": -9420.4215,
    "datepurchase": "2020-01-31T13:59:53.66277+03:00",
    "datelastupdate": "2020-06-19T23:59:27.189249+03:00",
    "type": "share"
  }
]
```

### POST

`POST /funds/rub` – to add new fund

Need to send:

```
{
  "name": "РусАгро",
  "ticker": "AGRO",
  "amount": 32,
  "priceperitem": 660,
  "datepurchase": "2020-06-19T23:59:27.189249+03:00",
  "type": "bond"
}
```

### EDIT

`PUT /funds/rub/{id}` 

`PUT /funds/usd/{id}` 

Example URL – http://127.0.0.1:8000/funds/rub/131

```
{
  "Id": 131,
  "name": "РусАгро",
  "ticker": "AGRO",
  "amount": 32,
  "priceperitem": 660,
  "datepurchase": "2020-06-19T23:59:27.189249+03:00",
  "type": "share"
}
```

It's necessary to send field Id in the body. If the Id field in the body and {id} in the URL is not equal, it will not work.

It's necessary to pass TYPE field. In other way we will not understand how to save our data.

Fields `name`, `ticker`, `amount`, `priceperitem`, `datepurchase` can be send by one.

### DELETE

```
/funds/rub/{id}
```

```
/funds/usd/{id}
```

http://127.0.0.1:8000/funds/rub/239 – URL

In BODY we need to send:

```
{
  "id": 239
}
```

In will not work if `{id}` in URL and `"Id"` in body is not the same.
