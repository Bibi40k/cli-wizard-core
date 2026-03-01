package wizard

import "fmt"

type ManageAction string

const (
	ManageCreate     ManageAction = "create"
	ManageEdit       ManageAction = "edit"
	ManageDelete     ManageAction = "delete"
	ManageSetDefault ManageAction = "set_default"
	ManageCancel     ManageAction = "cancel"
)

type ManageSelection[T any] struct {
	Action ManageAction
	Item   *T
}

type ManageOptions[T any] struct {
	Items         []T
	ItemLabel     func(item T) string
	Select        func(options []string, defaultOption, message string) string
	MainMessage   string
	ItemMessage   string
	DefaultLabel  string
	CreateLabel   string
	EditLabel     string
	DeleteLabel   string
	SetDefaultLbl string
	CancelLabel   string
}

func Manage[T any](opts ManageOptions[T]) (ManageSelection[T], error) {
	if opts.ItemLabel == nil {
		return ManageSelection[T]{}, fmt.Errorf("item label function is required")
	}
	if opts.Select == nil {
		return ManageSelection[T]{}, fmt.Errorf("select function is required")
	}

	createLabel := strOr(opts.CreateLabel, "Create new")
	editLabel := strOr(opts.EditLabel, "Edit existing")
	deleteLabel := strOr(opts.DeleteLabel, "Delete")
	setDefaultLabel := strOr(opts.SetDefaultLbl, "Set default")
	cancelLabel := strOr(opts.CancelLabel, "Cancel")
	mainMessage := strOr(opts.MainMessage, "Action:")
	itemMessage := strOr(opts.ItemMessage, "Select item:")

	actionOptions := []string{createLabel}
	if len(opts.Items) > 0 {
		actionOptions = []string{editLabel, createLabel, setDefaultLabel, deleteLabel}
	}
	actionOptions = append(actionOptions, cancelLabel)

	defaultAction := createLabel
	if len(opts.Items) > 0 {
		defaultAction = editLabel
	}
	action := opts.Select(actionOptions, defaultAction, mainMessage)
	switch action {
	case cancelLabel:
		return ManageSelection[T]{Action: ManageCancel}, nil
	case createLabel:
		return ManageSelection[T]{Action: ManageCreate}, nil
	case editLabel, deleteLabel, setDefaultLabel:
		if len(opts.Items) == 0 {
			return ManageSelection[T]{Action: ManageCancel}, nil
		}
	default:
		return ManageSelection[T]{Action: ManageCancel}, nil
	}

	labels := make([]string, 0, len(opts.Items))
	indexByLabel := make(map[string]int, len(opts.Items))
	defaultItem := ""
	for i, item := range opts.Items {
		label := opts.ItemLabel(item)
		labels = append(labels, label)
		indexByLabel[label] = i
		if opts.DefaultLabel != "" && label == opts.DefaultLabel {
			defaultItem = label
		}
	}
	if defaultItem == "" && len(labels) > 0 {
		defaultItem = labels[0]
	}
	itemOptions := append([]string{}, labels...)
	itemOptions = append(itemOptions, cancelLabel)
	chosen := opts.Select(itemOptions, defaultItem, itemMessage)
	if chosen == cancelLabel {
		return ManageSelection[T]{Action: ManageCancel}, nil
	}
	idx, ok := indexByLabel[chosen]
	if !ok {
		return ManageSelection[T]{Action: ManageCancel}, nil
	}
	item := opts.Items[idx]

	switch action {
	case editLabel:
		return ManageSelection[T]{Action: ManageEdit, Item: &item}, nil
	case deleteLabel:
		return ManageSelection[T]{Action: ManageDelete, Item: &item}, nil
	case setDefaultLabel:
		return ManageSelection[T]{Action: ManageSetDefault, Item: &item}, nil
	default:
		return ManageSelection[T]{Action: ManageCancel}, nil
	}
}

func strOr(v, fallback string) string {
	if v != "" {
		return v
	}
	return fallback
}
