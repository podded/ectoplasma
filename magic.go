package ectoplasma


const (
	ESI_KILLMAIL_URL_FMT = "https://esi.evetech.net/latest/killmails/%d/%s/"
	USER_AGENT           = "Ectoplasma by @Crypta Electrica"
	MAGIC_HEADER         = "X-Ecto-Update"

	REDIS_INGEST_QUEUE  = "ECTO_INGEST"
	REDIS_PROCESS_QUEUE = "ECTO_PROCESS"
	REDIS_ERROR_QUEUE   = "ECTO_ERROR"
)
