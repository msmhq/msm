package player

import (
	"encoding/json"
	"fmt"
	"os"
)

type AllowlistEntry struct {
	UUID string `json:"uuid,omitempty"`
	Name string `json:"name"`
}

type OpsEntry struct {
	UUID                string `json:"uuid,omitempty"`
	Name                string `json:"name"`
	Level               int    `json:"level"`
	BypassesPlayerLimit bool   `json:"bypassesPlayerLimit"`
}

type BannedPlayerEntry struct {
	UUID    string `json:"uuid,omitempty"`
	Name    string `json:"name"`
	Created string `json:"created,omitempty"`
	Source  string `json:"source,omitempty"`
	Expires string `json:"expires,omitempty"`
	Reason  string `json:"reason,omitempty"`
}

type BannedIPEntry struct {
	IP      string `json:"ip"`
	Created string `json:"created,omitempty"`
	Source  string `json:"source,omitempty"`
	Expires string `json:"expires,omitempty"`
	Reason  string `json:"reason,omitempty"`
}

func LoadAllowlist(path string) ([]AllowlistEntry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []AllowlistEntry{}, nil
		}
		return nil, err
	}

	var entries []AllowlistEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}
	return entries, nil
}

func SaveAllowlist(path string, entries []AllowlistEntry) error {
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func AddToAllowlist(path, name string) error {
	entries, err := LoadAllowlist(path)
	if err != nil {
		return err
	}

	for _, e := range entries {
		if e.Name == name {
			return fmt.Errorf("player %q is already in the allowlist", name)
		}
	}

	entries = append(entries, AllowlistEntry{Name: name})
	return SaveAllowlist(path, entries)
}

func RemoveFromAllowlist(path, name string) error {
	entries, err := LoadAllowlist(path)
	if err != nil {
		return err
	}

	var newEntries []AllowlistEntry
	found := false
	for _, e := range entries {
		if e.Name != name {
			newEntries = append(newEntries, e)
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("player %q is not in the allowlist", name)
	}

	return SaveAllowlist(path, newEntries)
}

func LoadOps(path string) ([]OpsEntry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []OpsEntry{}, nil
		}
		return nil, err
	}

	var entries []OpsEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}
	return entries, nil
}

func SaveOps(path string, entries []OpsEntry) error {
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func AddOp(path, name string, level int) error {
	entries, err := LoadOps(path)
	if err != nil {
		return err
	}

	for i, e := range entries {
		if e.Name == name {
			entries[i].Level = level
			return SaveOps(path, entries)
		}
	}

	entries = append(entries, OpsEntry{
		Name:                name,
		Level:               level,
		BypassesPlayerLimit: false,
	})
	return SaveOps(path, entries)
}

func RemoveOp(path, name string) error {
	entries, err := LoadOps(path)
	if err != nil {
		return err
	}

	var newEntries []OpsEntry
	found := false
	for _, e := range entries {
		if e.Name != name {
			newEntries = append(newEntries, e)
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("player %q is not an operator", name)
	}

	return SaveOps(path, newEntries)
}

func LoadBannedPlayers(path string) ([]BannedPlayerEntry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []BannedPlayerEntry{}, nil
		}
		return nil, err
	}

	var entries []BannedPlayerEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}
	return entries, nil
}

func SaveBannedPlayers(path string, entries []BannedPlayerEntry) error {
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func BanPlayer(path, name, reason string) error {
	entries, err := LoadBannedPlayers(path)
	if err != nil {
		return err
	}

	for _, e := range entries {
		if e.Name == name {
			return fmt.Errorf("player %q is already banned", name)
		}
	}

	entries = append(entries, BannedPlayerEntry{
		Name:   name,
		Reason: reason,
	})
	return SaveBannedPlayers(path, entries)
}

func UnbanPlayer(path, name string) error {
	entries, err := LoadBannedPlayers(path)
	if err != nil {
		return err
	}

	var newEntries []BannedPlayerEntry
	found := false
	for _, e := range entries {
		if e.Name != name {
			newEntries = append(newEntries, e)
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("player %q is not banned", name)
	}

	return SaveBannedPlayers(path, newEntries)
}

func LoadBannedIPs(path string) ([]BannedIPEntry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []BannedIPEntry{}, nil
		}
		return nil, err
	}

	var entries []BannedIPEntry
	if err := json.Unmarshal(data, &entries); err != nil {
		return nil, err
	}
	return entries, nil
}

func SaveBannedIPs(path string, entries []BannedIPEntry) error {
	data, err := json.MarshalIndent(entries, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

func BanIP(path, ip, reason string) error {
	entries, err := LoadBannedIPs(path)
	if err != nil {
		return err
	}

	for _, e := range entries {
		if e.IP == ip {
			return fmt.Errorf("IP %q is already banned", ip)
		}
	}

	entries = append(entries, BannedIPEntry{
		IP:     ip,
		Reason: reason,
	})
	return SaveBannedIPs(path, entries)
}

func UnbanIP(path, ip string) error {
	entries, err := LoadBannedIPs(path)
	if err != nil {
		return err
	}

	var newEntries []BannedIPEntry
	found := false
	for _, e := range entries {
		if e.IP != ip {
			newEntries = append(newEntries, e)
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("IP %q is not banned", ip)
	}

	return SaveBannedIPs(path, newEntries)
}
