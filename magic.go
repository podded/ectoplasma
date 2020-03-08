package ectoplasma

const (
	ESIKillmailURLFmt = "https://esi.evetech.net/latest/killmails/%d/%s/"
	MagicHeader       = "X-Ecto-Update"

	RedisIngestQueue  = "ECTO_INGEST"
	RedisProcessQueue = "ECTO_PROCESS"
	RedisErrorQueue   = "ECTO_ERROR"
)
