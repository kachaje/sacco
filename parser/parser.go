package parser

import "fmt"

const (
	INITIAL_SCREEN = "initialScreen"
	INPUT_SCREEN   = "inputScreen"
	QUIT_SCREEN    = "quitScreen"
)

type WorkFlow struct {
	Tree map[string]any
	Data *map[string]any

	CurrentScreen  string
	NextScreen     string
	PreviousScreen string
}

func NewWorkflow(tree map[string]any) *WorkFlow {
	return &WorkFlow{
		Tree:          tree,
		Data:          &map[string]any{},
		CurrentScreen: INITIAL_SCREEN,
	}
}

func (w *WorkFlow) GetNode(screen string) map[string]any {
	if w.Tree[screen] != nil {
		node, ok := w.Tree[screen].(map[string]any)
		if ok {
			return node
		}
	}

	return nil
}

func (w *WorkFlow) NextNode(input string) map[string]any {
	var node map[string]any
	var nextScreen string
	var ok bool

	if w.CurrentScreen == INITIAL_SCREEN {
		nextScreen, ok = w.Tree[INITIAL_SCREEN].(string)
		if ok {
			node = w.GetNode(nextScreen)
		}
	} else {
		node = w.GetNode(w.CurrentScreen)

		fmt.Println("##########")
	}

	w.PreviousScreen = w.CurrentScreen
	w.CurrentScreen = nextScreen

	return node
}
