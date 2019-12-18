package main

// Config represents the runtime configuration saved as a json.
type Config struct {
	Discord struct {
		Token        string `json:"token"`
		AdminRole    string `json:"adminRole"`
		DefaultGuild string `json:"defaultGuild"`
		Playing      string `json:"playing"`
	} `json:"discord"`
	Role struct {
		Casters     []string `json:"Casters"`
		RangedDPS   []string `json:"Ranged DPS"`
		MeleeDPS    []string `json:"Melee DPS"`
		Healers     []string `json:"Healers"`
		Tanks       []string `json:"Tanks"`
		Guest       []string `json:"Guest"`
		Chaos       []string `json:"Chaos"`
		Light       []string `json:"Light"`
		RaidLeaders []string `json:"Raid Leaders"`
	} `json:"role"`
	Throttle struct {
		StaticLfmLight struct {
			MaxTokens     int `json:"maxTokens"`
			TokenInterval int `json:"tokenInterval"`
		} `json:"static_lfm_light"`
		PlayerLfgLight struct {
			MaxTokens     int `json:"maxTokens"`
			TokenInterval int `json:"tokenInterval"`
		} `json:"player_lfg_light"`
		StaticLfmChaos struct {
			MaxTokens     int `json:"maxTokens"`
			TokenInterval int `json:"tokenInterval"`
		} `json:"static_lfm_chaos"`
		PlayerLfgChaos struct {
			MaxTokens     int `json:"maxTokens"`
			TokenInterval int `json:"tokenInterval"`
		} `json:"player_lfg_chaos"`
	} `json:"throttle"`
	Monitor struct {
		Output string   `json:"output"`
		Events []string `json:"events"`
	} `json:"monitor"`
	Welcome struct {
		Message string `json:"message"`
	} `json:"welcome"`
	Purge struct {
		PartyFinderChaos string `json:"party_finder_chaos"`
		PartyFinderLight string `json:"party_finder_light"`
	} `json:"purge"`
	Accuracy struct {
	} `json:"accuracy"`
	Help struct {
	} `json:"help"`
	Clear struct {
	} `json:"clear"`
	CommandPrefix string `json:"commandPrefix"`
}
