package tviewview

import (
	"log"

	"github.com/MelleKoning/aifun/internal/genaimodel"
	"github.com/MelleKoning/aifun/internal/terminal"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type tviewApp struct {
	app          *tview.Application
	mdRenderer   terminal.GlamourRenderer // can render markdown colours
	flex         *tview.Flex
	textArea     *tview.TextArea
	dropDown     *tview.DropDown
	outputView   *tview.TextView
	submitButton *tview.Button
	promptView   *tview.TextArea

	aimodel genaimodel.Action
}

type TviewApp interface {
	Run() error
	SetDefaultView()
}

// New will create a new VIEW on the terminal
// Always call a view, for example "SetDefaultView"
// to initialize the view container with a default view
// TODO Expose a good interface for this
func New(mdrenderer terminal.GlamourRenderer,
	aimodel genaimodel.Action) TviewApp {
	tv := &tviewApp{
		app:        tview.NewApplication(),
		mdRenderer: mdrenderer,
		aimodel:    aimodel,
	}
	tv.createOutputView()
	tv.createTextArea()
	tv.createSubmitButton()
	tv.createDropDown()

	// Create a flex layout to arrange the views
	tv.flex = tview.NewFlex()
	// Now arrange the views on the created Flex:
	tv.SetDefaultView()
	tv.app.SetRoot(tv.flex, true)

	return tv
}

func (tv *tviewApp) Run() error {
	err := tv.app.Run()
	if err != nil {
		return err
	}

	return nil
}
func (tv *tviewApp) createOutputView() {
	// Create a text view for displaying output
	// contains the logic for rendering
	tv.outputView = tview.NewTextView().
		SetDynamicColors(true).
		SetScrollable(true).
		SetSize(0, 0).
		SetChangedFunc(func() {
			tv.outputView.ScrollToEnd()
			tv.app.Draw()
		}).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyTAB {
			tv.app.SetFocus(tv.dropDown)
		}
	})

	tv.outputView.SetBorder(true)
}

func (tv *tviewApp) createTextArea() {
	// Create an input field for user input
	tv.textArea = tview.NewTextArea().
		SetLabel("Enter command: ")
	tv.textArea.SetBorder(true)
	// Capture key events for the text area
	tv.textArea.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		if event.Key() == tcell.KeyTAB {
			tv.app.SetFocus(tv.submitButton) // Move focus to the submit button
			return nil                       // Consume the event
		}
		if event.Key() == tcell.KeyEnter {
			return event // nil // Consume other events - do nothing!
		}
		return event
	})

	// we have to enable pasting for the user
	// for the whole app so that user can
	// paste into the Text Area
	tv.app.SetFocus(tv.textArea).EnablePaste(true)
}

func (tv *tviewApp) createSubmitButton() {
	tv.submitButton = tview.NewButton("Submit").SetSelectedFunc(
		func() {
			command := tv.textArea.GetText()
			txtRendered, err := tv.mdRenderer.GetUserRendered(command,
				tv.aimodel.GetHistoryLength())
			if err != nil {
				log.Print(err)
			}
			// First add the user command to the output
			txtRendered = tview.TranslateANSI(txtRendered)
			tv.outputView.SetText(tv.outputView.GetText(false) + txtRendered)
			// TODO: execute model
			result, err := tv.aimodel.ChatMessage(command)
			if err != nil {
				tv.outputView.SetText(tv.outputView.GetText(false) + err.Error())
			} else {
				renderedResult, _ := tv.mdRenderer.GetRendered(result)
				txtRendered = tview.TranslateANSI(renderedResult)
				tv.outputView.SetText(tv.outputView.GetText(false) + txtRendered)
			}
			tv.app.SetFocus(tv.outputView)
		}).
		SetExitFunc(func(key tcell.Key) {
			if key == tcell.KeyTAB {
				tv.app.SetFocus(tv.outputView)
			}
		},
		)
}

// we create a dropdown, but it should be fed with
// some model data instead of hardcoded static data
func (tv *tviewApp) createDropDown() {
	// Create a dropdown for selecting options
	tv.dropDown = tview.NewDropDown().
		SetLabel("Select option: ").
		SetOptions([]string{
			"Continue",
			"OutputView",
			"PromptView",
			"SystemPrompt",
			"Exit"}, func(option string, index int) {
			switch option {
			case "Exit":
				tv.app.Stop()
			case "OutputView":
				tv.SetDefaultView()
			case "PromptView":
				// TODO: move switching of view
				// to a function
				tv.flex.RemoveItem(tv.outputView)
				tv.flex.AddItem(tv.promptView, 0, 10, true)
			case "SystemPrompt":
				response := tv.aimodel.SendSystemPrompt()
				renderedResult, _ := tv.mdRenderer.GetRendered(response)
				txtRendered := tview.TranslateANSI(renderedResult)
				tv.outputView.SetText(tv.outputView.GetText(false) + txtRendered)

			}
			if option == "Exit" {
				tv.app.Stop()
			}
			tv.outputView.SetText(tv.outputView.GetText(true) + option)
		}).SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyTAB {
			tv.app.SetFocus(tv.textArea)
		}
	})
}

// SetDefaultView will set the default view
// of the tviewApp
func (tv *tviewApp) SetDefaultView() {
	tv.flex.Clear()
	tv.flex.
		SetDirection(tview.FlexRow).
		AddItem(tv.outputView, 0, 10, true).
		AddItem(tv.dropDown, 1, 1, true).
		AddItem(tv.textArea, 0, 3, true).
		AddItem(tv.submitButton, 1, 1, true)
}
