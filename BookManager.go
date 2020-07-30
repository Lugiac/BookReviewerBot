package main

import (
	"fmt"

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

	if bookPrice == "0" {
		return &discordgo.MessageEmbed{
			Author:      &discordgo.MessageEmbedAuthor{},
			Color:       0xffffff, // Green
			Description: bookDescription,
			Thumbnail: &discordgo.MessageEmbedThumbnail{
				URL: bookThumbnail,
			},
			Title: bookTitle,
		}

	} else {

		return &discordgo.MessageEmbed{
			Author:      &discordgo.MessageEmbedAuthor{},
			Color:       0xffffff, // Green
			Description: bookDescription,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   config.Command.EmbedPrice, // I wish I could insert IFs inside this xD
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
}
