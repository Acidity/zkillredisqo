package zkillredisqo

import (
	"strings"
	"time"
)

// Kill represents a kill as available on zKillboard
type Kill struct {
	KillPackage `json:"package"`
}

// IsNullKill checks whether a kill is an empty "null" kill
func (k *Kill) IsNullKill() bool {
	return k.ID == 0
}

// KillPackage represents the "package" layer kills are wrap in by RedisQ
type KillPackage struct {
	ID         int                  `json:"killID"`
	KillMail   KillMail             `json:"killmail"`
	ZKillboard ZKillboardAttributes `json:"zkb"`
}

// KillMail stores the actual information about an EVE killmail
type KillMail struct {
	ID            int                      `json:"killID"`
	Time          KillMailTime             `json:"killTime"`
	SolarSystem   KillMailCommonAttributes `json:"solarSystem"`
	Attackers     []KillMailAttacker       `json:"attackers"`
	AttackerCount int                      `json:"attackerCount"`
	Victim        KillMailVictim           `json:"victim"`
	War           struct {
		ID   int    `json:"id"`
		HRef string `json:"href"`
	} `json:"war"`
}

// KillMailTime is used to allow for proper time parsing from JSON
type KillMailTime struct {
	time.Time
}

// UnmarshalJSON tries to parse a timestamp provided by RedisQ to a Go time type
func (t *KillMailTime) UnmarshalJSON(buf []byte) error {
	parsedTime, err := time.Parse("2006.01.02 15:04:05", strings.Trim(string(buf), `"`))
	if err != nil {
		return err
	}

	t.Time = parsedTime
	return nil
}

// KillMailAttacker represents information about a single attacking party
type KillMailAttacker struct {
	Character      KillMailCommonAttributes `json:"character"`
	Corporation    KillMailCommonAttributes `json:"corporation"`
	Alliance       KillMailCommonAttributes `json:"alliance"`
	Faction        KillMailCommonAttributes `json:"faction"`
	Ship           KillMailCommonAttributes `json:"shipType"`
	Weapon         KillMailCommonAttributes `json:"weaponType"`
	DamageDone     int                      `json:"damageDone"`
	FinalBlow      bool                     `json:"finalBlow"`
	SecurityStatus float64                  `json:"securityStatus"`
}

// KillMailVictim represents information about the victim of a kill
type KillMailVictim struct {
	Character   KillMailCommonAttributes `json:"character"`
	Corporation KillMailCommonAttributes `json:"corporation"`
	Alliance    KillMailCommonAttributes `json:"alliance"`
	Faction     KillMailCommonAttributes `json:"faction"`
	Ship        KillMailCommonAttributes `json:"shipType"`
	DamageTaken int                      `json:"damageTaken"`
	Items       []KillMailItem           `json:"items"`
	Position    struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
		Z float64 `json:"z"`
	}
}

// KillMailItem represents information about a single item included in a kill, either dropped or destroyed
type KillMailItem struct {
	Item              KillMailCommonAttributes `json:"itemType"`
	QuantityDropped   int                      `json:"quantityDropped"`
	QuantityDestroyed int                      `json:"quantityDestroyed"`
	Flag              int                      `json:"flag"`
	Singleton         int                      `json:"singleton"`
}

// KillMailCommonAttributes represents information being reused for multiple parts of a killmail such as corporations, characters or ship types
type KillMailCommonAttributes struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	HRef string `json:"href"`
	Icon struct {
		HRef string `json:"href"`
	} `json:"icon"`
}

// ZKillboardAttributes stores additional information regarding the kill, provided by zKillboard
type ZKillboardAttributes struct {
	TotalValue float64 `json:"totalValue"`
	Points     int     `json:"points"`
	LocationID int     `json:"locationID"`
	Hash       string  `json:"hash"`
	HRef       string  `json:"href"`
}
