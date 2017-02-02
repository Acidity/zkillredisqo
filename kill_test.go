package zkillredisqo

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
	"time"
)

func TestUnmarshalKill(t *testing.T) {
	sample, err := os.OpenFile("./test/sample.json", os.O_RDONLY, 0644)
	if err != nil {
		t.Fatalf("Failed to open sample file: %v\n", err)
		return
	}
	defer sample.Close()

	var kill *Kill
	if err = json.NewDecoder(sample).Decode(&kill); err != nil {
		t.Fatalf("Failed to decode sample kill: %v\n", err)
		return
	}

	if kill.IsNullKill() {
		t.Error("Kill is null kill")
	}
	if kill.KillID != 59704390 || kill.KillMail.KillID != 59704390 {
		t.Errorf("Invalid killID: %+v, expected 59704390\n", kill.KillID)
	}
	refTime, err := time.Parse(time.RFC3339, "2017-02-02T19:51:15Z")
	if err != nil {
		t.Errorf("Failed to parse refTime: %v", err)
	}
	if !kill.KillMail.KillTime.Equal(refTime) {
		t.Errorf("Invalid killTime: %+v, expected %v\n", kill.KillMail.KillTime, refTime)
	}
	if kill.KillMail.SolarSystem.ID != 30000142 || !strings.EqualFold(kill.KillMail.SolarSystem.Name, "Jita") {
		t.Errorf("Invalid solarSystem: %+v, expected ID 30000142 and name \"Jita\"", kill.KillMail.SolarSystem)
	}
	if kill.KillMail.AttackerCount != 2 || len(kill.KillMail.Attackers) != 2 {
		t.Errorf("Invalid attackers: %+v, expected 2 entries", kill.KillMail.Attackers)
	}
	if kill.KillMail.Victim.Character.ID != 188689214 || !strings.EqualFold(kill.KillMail.Victim.Character.Name, "Zhivchik") {
		t.Errorf("Invalid victim: %+v, expected ID 188689214 and name \"Zhivchik\"", kill.KillMail.Victim.Character)
	}
	if kill.KillMail.Victim.Corporation.ID != 98470839 || !strings.EqualFold(kill.KillMail.Victim.Corporation.Name, "Hearts of the Void") {
		t.Errorf("Invalid victim corporation: %+v, expected ID 98470839 and name \"Hearts of the Void\"", kill.KillMail.Victim.Corporation)
	}
	if kill.KillMail.Victim.ShipType.ID != 670 || !strings.EqualFold(kill.KillMail.Victim.ShipType.Name, "Capsule") {
		t.Errorf("Invalid victim shipType: %+v, expected ID 670 and name \"Capsule\"", kill.KillMail.Victim.ShipType)
	}
	if kill.KillMail.Victim.DamageTaken != 438 {
		t.Errorf("Invalid victim damageTaken: %+v, expected 438", kill.KillMail.Victim.DamageTaken)
	}
	if len(kill.KillMail.Victim.Items) != 5 {
		t.Errorf("Invalid victim items: %+v, expected 5 entries", kill.KillMail.Victim.Items)
	}
	if kill.KillMail.War.ID != 0 {
		t.Errorf("Invalid war: %+v, expected 0", kill.KillMail.War)
	}
}
