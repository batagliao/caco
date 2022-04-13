package actions

import "caco/services"

var AvailableActions = []*services.DialogFlowAction{
	TeamMRsAction,
	FallbackAction,
	BombeiroAction,
}

// GetAction ...
func GetAction(actionName string) *services.DialogFlowAction {
	for _, a := range AvailableActions {
		if a.Name == actionName {
			return a
		}
	}
	return nil
}
