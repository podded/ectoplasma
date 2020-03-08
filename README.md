# Ectoplasma

This module is the ingest engine for killmails into the database.

There is a single endpoint exposed, which is the `/submit` endpoint. This accepts `POST` requests with the following body:

```json
{
  "hash": "2103ec9811c2de80a2d4d53b2e9df2144b054627",
  "id": 34257049
}
```

When a request with a valid body is made, the first thing done is to fetch the kill from ESI.
If ESI returns a valid killmail then the original post request should be met one of the following responses
 - 201 - A new killmail has been created in the system
 - 402 - The killmail already exists but will be double checked to ensure accuracy
 
Once the killmail is stored in the DB the axiom attributes and KM stats are generated for that killmail.
