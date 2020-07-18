package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/BurntSushi/toml"
	"github.com/Lukaesebrot/dgc"
	"github.com/bwmarrin/discordgo"
	"google.golang.org/api/books/v1"
	"google.golang.org/api/option"
)

var ctx = context.Background()

var config tomlConfig

var bookService *books.Service

func main() {

	loadConfig()
	authGoogleBooksAPI()

	discordSession := authDiscordAPI()

	fmt.Print("Connected to google and discord! ")

	handleCommands(discordSession)
	//Allow to terminate application with CTRL+C(SIGINT)
	defer func() {
		sc := make(chan os.Signal, 1)
		signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
		<-sc
		discordSession.Close()
	}()

}

func loadConfig() {
	fmt.Println("Loading config ")
	if _, err := toml.DecodeFile("config.toml", &config); err != nil {
		fmt.Println(err)
		return
	}
}

func authGoogleBooksAPI() {
	fmt.Println("Authenticating Google Books API. ")
	var err error
	bookService, err = books.NewService(ctx, option.WithAPIKey(config.Keys.GoogleToken))

	if err != nil {
		fmt.Println("Google doesn't want to authenticate us. Maybe something is wrong? ")
		return
	}
}

func authDiscordAPI() (discordSession *discordgo.Session) {
	fmt.Println("Authenticating discord. ")
	discordSession, err := discordgo.New("Bot " + config.Keys.DiscordToken)

	if err != nil {
		fmt.Println("Yeah, it didn't work. Authentication at discord failed :/ ")
		return
	}
	err = discordSession.Open()
	if err != nil {
		fmt.Println("There was an error while trying to open a Discord connection. ", err)
		return
	}
	return discordSession
}

func handleCommands(discordSession *discordgo.Session) {
	fmt.Println("Registering commands. ")
	router := dgc.Create(&dgc.Router{
		Prefixes: []string{config.Command.Prefix},
		Commands: []*dgc.Command{},
	})

	router.RegisterCmd(&dgc.Command{
		Name: config.Command.CommandName,

		Aliases:    config.Command.CommandAliases,
		IgnoreCase: true,
		Handler:    bookCommandHandler,
	})
	router.Initialize(discordSession)
}
