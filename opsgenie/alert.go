package opsgenie

import (
	"fmt"
	"github.com/opsgenie/opsgenie-go-sdk-v2/alert"
)

func (o *Opsgenie) Alert() {
	alertClient, err := alert.NewClient(o.config)
	if err != nil {
		fmt.Println("Opesgenie error: ", err)
		return
	}

	createResult, err := alertClient.Create(nil, &alert.CreateAlertRequest{
		Message:     "message1",
		Alias:       "alias1",
		Description: "alert description1",
		Responders: []alert.Responder{
			{Type: alert.EscalationResponder, Name: "TeamA_escalation"},
			{Type: alert.ScheduleResponder, Name: "TeamB_schedule"},
		},
		VisibleTo: []alert.Responder{
			{Type: alert.UserResponder, Username: "testuser@gmail.com"},
			{Type: alert.TeamResponder, Name: "admin"},
		},
		Actions: []string{"action1", "action2"},
		Tags:    []string{"tag1", "tag2"},
		Details: map[string]string{
			"key1": "value1",
			"key2": "value2",
		},
		Entity:   "entity2",
		Source:   "source2",
		Priority: alert.P1,
		User:     "testuser@gmail.com",
		Note:     "alert note2",
	})

	if err != nil {
		fmt.Println("Opesgenie error: ", err)
		return
	}

	fmt.Println("createResult: ", createResult)
}
