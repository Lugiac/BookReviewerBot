package main

type tomlConfig struct {
	Keys    keysConfig
	Command commandConfig
}

type keysConfig struct {
	DiscordToken string
	GoogleToken  string
}

type commandConfig struct {
	Prefix         string
	CommandName    string
	CommandAliases []string
	EmbedPrice     string
}
