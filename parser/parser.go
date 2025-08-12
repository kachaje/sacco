package parser

const (
	INITIAL_SCREEN = "initialScreen"
	INPUT_SCREEN   = "inputScreen"
	QUIT_SCREEN    = "quitScreen"
)

type WorkFlow struct {
	Tree map[string]any
	Data *map[string]any

	CurrentScreen     string
	CurrentType       string
	CurrentText       string
	CurrentOptions    []map[string]any
	CurrentIdentifier string
	NextScreen        string
	PreviousScreen    string
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
