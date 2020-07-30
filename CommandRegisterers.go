package main

import "github.com/Lukaesebrot/dgc"

func bookCommandHandler(ctx *dgc.Ctx) {
	book := ctx.Arguments.Raw()
	ctx.Session.ChannelMessageDelete(ctx.Event.ChannelID, ctx.Event.ID)
	ctx.RespondEmbed(createBookEmbed(book))
}
