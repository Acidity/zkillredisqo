package zkillredisqo

import (
	"strings"
	"time"
)

type Kill struct {
	KillPackage `json:"package"`
}

func (k *Kill) IsNullKill() bool {
	return k.KillID == 0
}

type KillPackage struct {
	KillID   int      `json:"killID"`
	KillMail KillMail `json:"killmail"`
}

type KillMail struct {
	KillID        int                      `json:"killID"`
	KillTime      KillMailTime             `json:"killTime"`
	SolarSystem   KillMailCommonAttributes `json:"solarSystem"`
	Attackers     []KillMailAttacker       `json:"attackers"`
	AttackerCount int                      `json:"attackerCount"`
	Victim        KillMailVictim           `json:"victim"`
	War           KillMailWar              `json:"war"`
}

type KillMailTime struct {
	time.Time
}

func (t *KillMailTime) UnmarshalJSON(buf []byte) error {
	parsedTime, err := time.Parse("2006.01.02 15:04:05", strings.Trim(string(buf), `"`))
	if err != nil {
		return err
	}

	t.Time = parsedTime
	return nil
}

type KillMailAttacker struct {
	Character      KillMailCommonAttributes `json:"character"`
	Corporation    KillMailCommonAttributes `json:"corporation"`
	Alliance       KillMailCommonAttributes `json:"alliance"`
	Faction        KillMailCommonAttributes `json:"faction"`
	ShipType       KillMailCommonAttributes `json:"shipType"`
	WeaponType     KillMailCommonAttributes `json:"weaponType"`
	DamageDone     int                      `json:"damageDone"`
	FinalBlow      bool                     `json:"finalBlow"`
	SecurityStatus float64                  `json:"securityStatus"`
}

type KillMailVictim struct {
	Character   KillMailCommonAttributes `json:"character"`
	Corporation KillMailCommonAttributes `json:"corporation"`
	Alliance    KillMailCommonAttributes `json:"alliance"`
	Faction     KillMailCommonAttributes `json:"faction"`
	ShipType    KillMailCommonAttributes `json:"shipType"`
	DamageTaken int                      `json:"damageTaken"`
	Items       []KillMailItem           `json:"items"`
	Position    struct {
		X float64 `json:"x"`
		Y float64 `json:"y"`
		Z float64 `json:"z"`
	}
}

type KillMailItem struct {
	ItemType          KillMailCommonAttributes `json:"itemType"`
	QuantityDropped   int                      `json:"quantityDropped"`
	QuantityDestroyed int                      `json:"quantityDestroyed"`
	Flag              int                      `json:"flag"`
	Singleton         int                      `json:"singleton"`
}

type KillMailWar struct {
	ID   int    `json:"id"`
	HRef string `json:"href"`
}

type KillMailCommonAttributes struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	HRef string `json:"href"`
	Icon struct {
		HRef string `json:"href"`
	} `json:"icon"`
}
