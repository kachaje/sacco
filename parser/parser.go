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

	CurrentScreen   string
	NextScreen      string
	PreviousScreen  string
	CurrentLanguage string
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

func (w *WorkFlow) InputIncluded(input string, options []map[string]any) bool {
	found := false

	for _, option := range options {
		if option["position"] != nil {
			if val, ok := option["position"].(int); ok && fmt.Sprint(val) == input {
				found = true
				break
			}
		}
	}

	return found
}

func (w *WorkFlow) NodeOptions(input string) []string {
	options := []string{}

	node := w.GetNode(input)
	if node != nil && node["options"] != nil {
		opts, ok := node["options"].([]any)
		if ok {
			for _, row := range opts {
				optVal, ok := row.(map[string]any)
				if ok {
					position := fmt.Sprintf("%v", optVal["position"])

					val, ok := optVal["label"].(map[string]any)
					if ok {
						if val["all"] != nil {
							entry := fmt.Sprintf("%s. %s", position, val["all"])

							options = append(options, entry)
						} else if w.CurrentLanguage != "" && val[w.CurrentLanguage] != nil {
							entry := fmt.Sprintf("%s. %s", position, val[w.CurrentLanguage])

							options = append(options, entry)
						}
					}
				}
			}
		}
	}

	return options
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
