package ectoplasma

<<<<<<< HEAD
const (
	ESI_KILLMAIL_URL_FMT = "https://esi.evetech.net/latest/killmails/%d/%s/"
=======
const(
	REDIS_INGEST_QUEUE   = "ECTO_INGEST"
	REDIS_UPDATE_QUEUE   = "ECTO_UPDATE"
	ESI_KILLMAIL_URL_FMT = "https://esi.evetech.net/v1/killmails/%d/%s/"
>>>>>>> 9380b0ff3060b97a19d7215ab6e1220c4ec86acb
	USER_AGENT           = "Ectoplasma by @Crypta Electrica"
	MAGIC_HEADER         = "X-Ecto-Update"

	REDIS_INGEST_QUEUE  = "ECTO_INGEST"
	REDIS_PROCESS_QUEUE = "ECTO_PROCESS"
	REDIS_ERROR_QUEUE   = "ECTO_ERROR"
)
