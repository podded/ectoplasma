package killmail

import "time"

type (
	IDHashPair struct {
		ID   int32  `json:"id" bson:"_id"`
		Hash string `json:"hash" bson:"hash"`
	}

	PoddedKillmail struct {
		Killmail  Mail   `json:"esi_v1" bson:"esi_v1"`
		Hash      string `json:"hash" bson:"hash"`
		ID        int32  `json:"id" bson:"_id"`
		Published bool   `json:"published" bson:"published"` // TODO Implement the published marker
	}

	Mail struct {
		Victim        Victim     `json:"victim,omitempty"`
		Attackers     []Attacker `json:"attackers,omitempty"`       /* attackers array */
		KillmailTime  time.Time  `json:"killmail_time,omitempty"`   /* Time that the victim was killed and the killmail generated  */
		KillmailID    int32      `json:"killmail_id,omitempty"`     /* ID of the killmail */
		MoonID        int32      `json:"moon_id,omitempty"`         /* Moon if the kill took place at one */
		SolarSystemID int32      `json:"solar_system_id,omitempty"` /* Solar system that the kill took place in  */
		WarID         int32      `json:"war_id,omitempty"`          /* War if the killmail is generated in relation to an official war  */
	}

	Attacker struct {
		AllianceID     int32   `json:"alliance_id,omitempty"`     /* alliance_id integer */
		CharacterID    int32   `json:"character_id,omitempty"`    /* character_id integer */
		CorporationID  int32   `json:"corporation_id,omitempty"`  /* corporation_id integer */
		DamageDone     int32   `json:"damage_done,omitempty"`     /* damage_done integer */
		FactionID      int32   `json:"faction_id,omitempty"`      /* faction_id integer */
		FinalBlow      bool    `json:"final_blow,omitempty"`      /* Was the attacker the one to achieve the final blow  */
		SecurityStatus float32 `json:"security_status,omitempty"` /* Security status for the attacker  */
		ShipTypeID     int32   `json:"ship_type_id,omitempty"`    /* What ship was the attacker flying  */
		WeaponTypeID   int32   `json:"weapon_type_id,omitempty"`  /* What weapon was used by the attacker for the kill  */
	}

	Victim struct {
		Items         []VictimItems `json:"items,omitempty"` /* items array */
		Position      Position      `json:"position,omitempty"`
		AllianceID    int32         `json:"alliance_id,omitempty"`    /* alliance_id integer */
		CharacterID   int32         `json:"character_id,omitempty"`   /* character_id integer */
		CorporationID int32         `json:"corporation_id,omitempty"` /* corporation_id integer */
		DamageTaken   int32         `json:"damage_taken,omitempty"`   /* How much total damage was taken by the victim  */
		FactionID     int32         `json:"faction_id,omitempty"`     /* faction_id integer */
		ShipTypeID    int32         `json:"ship_type_id,omitempty"`   /* The ship that the victim was piloting and was destroyed  */
	}

	VictimItems struct {
		Flag              int32          `json:"flag,omitempty"`               /* Flag for the location of the item  */
		ItemTypeID        int32          `json:"item_type_id,omitempty"`       /* item_type_id integer */
		Items             []ItemSubItems `json:"items,omitempty"`              /* items array */
		QuantityDestroyed int64          `json:"quantity_destroyed,omitempty"` /* How many of the item were destroyed if any  */
		QuantityDropped   int64          `json:"quantity_dropped,omitempty"`   /* How many of the item were dropped if any  */
		Singleton         int32          `json:"singleton,omitempty"`          /* singleton integer */
	}

	ItemSubItems struct {
		Flag              int32 `json:"flag,omitempty"`               /* flag integer */
		ItemTypeID        int32 `json:"item_type_id,omitempty"`       /* item_type_id integer */
		QuantityDestroyed int64 `json:"quantity_destroyed,omitempty"` /* quantity_destroyed integer */
		QuantityDropped   int64 `json:"quantity_dropped,omitempty"`   /* quantity_dropped integer */
		Singleton         int32 `json:"singleton,omitempty"`          /* singleton integer */
	}

	Position struct {
		X float64 `json:"x,omitempty"` /* x number */
		Y float64 `json:"y,omitempty"` /* y number */
		Z float64 `json:"z,omitempty"` /* z number */
	}
)
