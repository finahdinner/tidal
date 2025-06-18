package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func (g *GuiWrapper) getHelpSection() fyne.CanvasObject {

	c := container.New(
		layout.NewVBoxLayout(),
	)

	markdownLines := []string{
		"- Tidal allows you to dynamically and creatively update your Twitch stream title.",
		"- The application provides access to real-time data from your stream and channel, such as follower count and current viewer count - these values are known as **Stream Variables**.",
		"- In addition, Tidal allows you to leverage the power of Large Language Models to create **AI-Generated Variables**. These are values generated from custom prompts, which can include your own creative input along with substituted Stream Variables.",
		"- You can then construct a **Title Template** using specified Stream Variables and/or AI-Generated Variables, the values of which are substituted in for each new title. You can configure your stream title to update at fixed intervals, or update whenever specified Stream Variables change value.",
	}

	for _, line := range markdownLines {
		rt := widget.NewRichTextFromMarkdown(line)
		rt.Wrapping = fyne.TextWrapWord
		c.Objects = append(c.Objects, rt)
	}

	noteSegment := &widget.TextSegment{
		Text:  "Note: Each section within the application has a help button in the top right corner, which is able to provide more specific information.",
		Style: widget.RichTextStyleInline,
	}
	noteSegment.Style.TextStyle = fyne.TextStyle{
		Italic: true,
		Bold:   true,
	}
	noteRt := widget.NewRichText(noteSegment)
	noteRt.Wrapping = fyne.TextWrapWord

	c.Objects = append(c.Objects, noteRt)

	return sectionWrapper(
		"General Help",
		nil,
		nil,
		c,
		true,
		false,
		true,
	)
}
