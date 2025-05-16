package main

import (
	"fmt"
	"log"

	"github.com/MelleKoning/aifun/internal/terminal"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

func main() {
	app := tview.NewApplication()

	mdRenderer, err := terminal.New()
	if err != nil {
		fmt.Println(err)
		return
	}
	var flex *tview.Flex
	var textArea *tview.TextArea
	var dropDown *tview.DropDown
	var outputView *tview.TextView
	var submitButton *tview.Button
	var promptView = tview.NewTextArea()

	// Create a text view for displaying output
	outputView = tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetSize(0, 0).
		SetChangedFunc(func() {
			outputView.ScrollToEnd()
			app.Draw()
		}).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyTAB {
			app.SetFocus(dropDown)
		}
	})
	outputView.SetBorder(true)

	// Create an input field for user input
	textArea = tview.NewTextArea().
		SetLabel("Enter command: ")
	textArea.SetBorder(true)
	// Capture key events for the text area
	textArea.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTAB {
			app.SetFocus(submitButton) // Move focus to the submit button
			return nil                 // Consume the event
		}
		if event.Key() == tcell.KeyEnter {
			return event // nil // Consume other events - do nothing!
		}
		return event
	})
	submitButton = tview.NewButton("Submit").SetSelectedFunc(
		func() {
			command := textArea.GetText()
			txtRendered, err := mdRenderer.GetRendered(command)
			if err != nil {
				log.Print(err)
			}
			txtRendered = tview.TranslateANSI(txtRendered)
			outputView.SetText(outputView.GetText(false) + txtRendered)
			app.SetFocus(outputView)
		}).SetExitFunc(func(key tcell.Key) {
		if key == tcell.KeyTAB {
			app.SetFocus(outputView)
		}
	},
	)

	// Create a dropdown for selecting options
	dropDown = tview.NewDropDown().
		SetLabel("Select option: ").
		SetOptions([]string{
			"Continue",
			"OutputView",
			"PromptView",
			"Exit"}, func(option string, index int) {
			switch option {
			case "Exit":
				app.Stop()
			case "OutputView":

				flex.RemoveItem(promptView)
				flex.AddItem(outputView, 0, 10, true)
			case "PromptView":
				flex.RemoveItem(outputView)
				flex.AddItem(promptView, 0, 10, true)
			}
			if option == "Exit" {
				app.Stop()

			}
			outputView.SetText(outputView.GetText(true) + option)
		}).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyTAB {
			app.SetFocus(textArea)
		}
	})

	// Set initial focus to the dropdown
	app.SetFocus(textArea).EnablePaste(true)

	// Create a flex layout to arrange the views
	flex = tview.NewFlex().
		SetDirection(tview.FlexRow).
		AddItem(outputView, 0, 10, true).
		AddItem(dropDown, 1, 1, true).
		AddItem(textArea, 0, 3, true).
		AddItem(submitButton, 1, 1, true)

	flex.SetTitle("What a view")

	if err := app.SetRoot(flex, true).Run(); err != nil {
		panic(err)
	}
}
