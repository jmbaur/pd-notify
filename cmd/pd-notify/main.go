package main

import (
	"context"
	_ "embed"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	pd "github.com/PagerDuty/go-pagerduty"
	"github.com/gen2brain/beeep"
)

type Config struct {
	ApiKey string
}

//go:embed alert.png
var alertIconContents string

func main() {
	flag.Parse()

	alertIconFile, err := ioutil.TempFile("", "pd-notify")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(alertIconFile.Name())

	if _, err := alertIconFile.WriteString(alertIconContents); err != nil {
		log.Fatal(err)
	}

	config := &Config{
		ApiKey: os.Getenv("PD_API_KEY"),
	}

	client := pd.NewClient(config.ApiKey)
	user, err := client.GetCurrentUser(pd.GetCurrentUserOptions{})
	if err != nil {
		log.Fatal("failed to get user", err)
	}

	teamIDs := []string{}
	for _, team := range user.Teams {
		teamIDs = append(teamIDs, team.ID)
	}

	escPolicies := paginate(func(next uint) ([]pd.EscalationPolicy, bool, uint) {
		response, err := client.ListEscalationPoliciesWithContext(context.Background(), pd.ListEscalationPoliciesOptions{
			UserIDs: []string{user.ID},
			Offset:  next,
		})
		if err != nil {
			return []pd.EscalationPolicy{}, false, 0
		}
		return response.EscalationPolicies, response.More, response.Total
	})

	fmt.Printf("Hello, %s!\n", user.Name)

	if len(escPolicies) == 0 {
		fmt.Printf("Looks like you don't have any escalation policies, good for you!\nGoodbye.\n")
		os.Exit(0)
	}

	escalationPolicyIDs := []string{}
	serviceIDs := []string{}
	for _, policy := range escPolicies {
		escalationPolicyIDs = append(escalationPolicyIDs, policy.ID)
		for _, service := range policy.Services {
			serviceIDs = append(serviceIDs, service.ID)
		}
	}

	oncalls := paginate(func(next uint) ([]pd.OnCall, bool, uint) {
		response, err := client.ListOnCallsWithContext(context.Background(), pd.ListOnCallOptions{
			UserIDs:             []string{user.ID},
			EscalationPolicyIDs: escalationPolicyIDs,
		})
		if err != nil {
			return []pd.OnCall{}, false, 0
		}
		return response.OnCalls, response.More, response.Total
	})

	currentlyOncall := false
	hasOncallStartingSoon := false
	for _, oncall := range oncalls {
		start, err := time.Parse(time.RFC3339, oncall.Start)
		if err != nil {
			log.Println("got invalid time format from pagerduty")
			continue
		}

		until := time.Until(start)
		if until < time.Duration(0) {
			currentlyOncall = true
		}

		if until < time.Duration(1*time.Hour) && until > time.Duration(0) {
			hasOncallStartingSoon = true
		}
	}

	if !currentlyOncall && !hasOncallStartingSoon {
		fmt.Printf("Looks like you don't have any oncalls starting soon.\nGoodbye!\n")
		os.Exit(0)
	}

	_ = beeep.Alert("foo", "bar", alertIconFile.Name())

	fmt.Println("Checking for incidents...")
	for {
		incidents := paginate(func(next uint) ([]pd.Incident, bool, uint) {
			response, err := client.ListIncidentsWithContext(context.Background(), pd.ListIncidentsOptions{
				ServiceIDs: serviceIDs,
				UserIDs:    []string{user.ID},
				TeamIDs:    teamIDs,
				Offset:     next,
			})
			if err != nil {
				return []pd.Incident{}, false, 0
			}
			return response.Incidents, response.More, response.Total
		})

		for _, incident := range incidents {
			if err := beeep.Alert(incident.Description, incident.Summary, alertIconFile.Name()); err != nil {
				log.Println("failed to send notification", err)
				fmt.Println("NEW INCIDENT")
				fmt.Println(incident.Description)
				fmt.Println(incident.Summary)
			}
		}
		time.Sleep(5 * time.Minute)
	}
}

func paginate[T any](f func(next uint) (t []T, more bool, nnext uint)) []T {
	allItems := []T{}
	var next uint = 0
	for {
		items, more, nnext := f(next)
		next = nnext
		allItems = append(allItems, items...)
		if !more {
			break
		}
	}

	return allItems
}
