package main

import (
	"fmt"

	"github.com/Lukaesebrot/dgc"
	"github.com/bwmarrin/discordgo"
	"google.golang.org/api/books/v1"
)

func adjustDescriptionSize(bookDescription string) (adjustedBookDescription string) {
	if len(bookDescription) > 2000 {
		return bookDescription[:1996] + "..."
	}
	return bookDescription
}

func getBookPriceWithCurrencyCode(bookListPrice *books.VolumeSaleInfoListPrice) (bookPriceWithCurrencyCode string) {
	return fmt.Sprintf("%.2f",
		bookListPrice.Amount) + " " +
		bookListPrice.CurrencyCode
}

func getBookInfos(bookName string) (Title string, Description string, Price string, Thumbnail string) {
	bookSearchResults, err := bookService.Volumes.List().Q(bookName).Do()
	bookVolume := bookSearchResults.Items[0]

	if err != nil {
		fmt.Println(err)
	}

	bookVolumeInfo := bookVolume.VolumeInfo
	bookVolumeSaleInfo := bookVolume.SaleInfo

	bookPrice := "0"
	if bookVolumeSaleInfo.ListPrice != nil {
		bookPrice = getBookPriceWithCurrencyCode(bookVolumeSaleInfo.ListPrice)
	}

	var bookThumbnail string
	if bookVolumeInfo.ImageLinks != nil {
		bookThumbnail = bookVolumeInfo.ImageLinks.Thumbnail
	}

	bookDescription := adjustDescriptionSize(bookVolumeInfo.Description)

	return bookVolumeInfo.Title, bookDescription, bookPrice, bookThumbnail
}

func createBookEmbed(bookName string) (bookReviewEmbed *discordgo.MessageEmbed) {
	bookTitle, bookDescription, bookPrice, bookThumbnail := getBookInfos(bookName)
	return &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{},
		Color:       0x00ff00, // Green
		Description: bookDescription,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   config.Command.EmbedPrice,
				Value:  bookPrice,
				Inline: true,
			},
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: bookThumbnail,
		},
		Title: bookTitle,
	}
}

func bookCommandHandler(ctx *dgc.Ctx) {
	book := ctx.Arguments.Raw()
	ctx.Session.ChannelMessageDelete(ctx.Event.ChannelID, ctx.Event.ID)
	ctx.RespondEmbed(createBookEmbed(book))
}
