package main

import (
	"context"
	"io/ioutil"
	"log"

	"github.com/chromedp/chromedp"
)

const URL = "https://booking.stockholm.se/"

func main() {
	// create context
	ctx, cancel := chromedp.NewContext(
		context.Background(),
		// chromedp.WithDebugf(log.Printf),
	)
	defer cancel()

	if err := chromedp.Run(ctx,
		chromedp.Tasks{
			chromedp.Navigate(URL),
		},
	); err != nil {
		log.Fatal(err)
	}

	if err := screenshot(ctx, "data/first.jpeg"); err != nil {
		log.Fatal(err)
	}

	if err := chromedp.Run(ctx,
		chromedp.Tasks{
			chromedp.Click("#hplSearch", chromedp.NodeVisible),
			chromedp.WaitVisible("#drplFacility_box2"),
		},
	); err != nil {
		log.Fatal(err)
	}

	if err := screenshot(ctx, "data/click.jpeg"); err != nil {
		log.Fatal(err)
	}

	// Set attribute for campo
	// Lista os nossos
	// drplFacility_box2

	// Set "Fotboll"
	// #drplActivity_box2

	// Clicka #btnSearch

	// Get data from form

}

func screenshot(ctx context.Context, file string) error {
	var buf []byte

	if err := chromedp.Run(ctx,
		chromedp.Tasks{
			chromedp.FullScreenshot(&buf, 90),
		},
	); err != nil {
		return err
	}

	if err := ioutil.WriteFile(file, buf, 0o644); err != nil {
		return err
	}

	return nil
}
