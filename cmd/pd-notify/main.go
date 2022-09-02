package main

import (
	"context"
	_ "embed"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"time"

	pd "github.com/PagerDuty/go-pagerduty"
	"github.com/jmbaur/pd-notify/notifications"
)

//go:embed alert.png
var alertIconContents []byte

type Config struct {
	ApiKey string
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

func logic() error {
	overrideUser := flag.String("user", "", "Name of user to listen for (default is current user)")
	flag.Parse()

	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return err
	}
	alertIconPath := filepath.Join(cacheDir, "pd-notify-alert-icon.png")
	if _, err := os.Stat(alertIconPath); errors.Is(err, os.ErrNotExist) {
		alertIconFile, err := os.Create(alertIconPath)
		if err != nil {
			return err
		}
		if _, err := alertIconFile.Write(alertIconContents); err != nil {
			return err
		}
		alertIconFile.Close()
	} else if err != nil {
		return err
	}

	apiKey, ok := os.LookupEnv("PD_API_KEY")
	if !ok {
		return errors.New("could not find API key, make sure PD_API_KEY is set")
	}
	config := &Config{
		ApiKey: apiKey,
	}

	client := pd.NewClient(config.ApiKey)
	currentUser, err := client.GetCurrentUser(pd.GetCurrentUserOptions{})
	if err != nil {
		return fmt.Errorf("failed to get user: %v", err)
	}

	var userID string
	var userName string

	teamIDs := []string{}
	for _, team := range currentUser.Teams {
		teamIDs = append(teamIDs, team.ID)
		if userID == "" && *overrideUser != "" {
			members := paginate(func(next uint) ([]pd.Member, bool, uint) {
				response, err := client.ListTeamMembers(context.Background(), team.ID, pd.ListTeamMembersOptions{})
				if err != nil {
					return []pd.Member{}, false, 0
				}
				return response.Members, response.More, response.Total
			})

			for _, member := range members {
				if member.User.Summary == *overrideUser {
					userID = member.User.ID
					userName = member.User.Summary
				}
			}
		}
	}

	if userID == "" {
		userID = currentUser.ID
		userName = currentUser.Name
	}

	escPolicies := paginate(func(next uint) ([]pd.EscalationPolicy, bool, uint) {
		response, err := client.ListEscalationPoliciesWithContext(context.Background(), pd.ListEscalationPoliciesOptions{
			UserIDs: []string{userID},
			Offset:  next,
		})
		if err != nil {
			return []pd.EscalationPolicy{}, false, 0
		}
		return response.EscalationPolicies, response.More, response.Total
	})

	fmt.Printf("Active user: %s\n", userName)

	if len(escPolicies) == 0 {
		fmt.Printf("Looks like %s don't have any escalation policies.\nGoodbye!\n", userName)
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
			UserIDs:             []string{userID},
			EscalationPolicyIDs: escalationPolicyIDs,
		})
		if err != nil {
			return []pd.OnCall{}, false, 0
		}
		return response.OnCalls, response.More, response.Total
	})

	timeUntilOncall := time.Duration(0)
	for _, oncall := range oncalls {
		start, err := time.Parse(time.RFC3339, oncall.Start)
		if err != nil {
			log.Println("got invalid time format from pagerduty")
			continue
		}

		timeUntilOncall = time.Until(start)
	}

	if len(oncalls) == 0 || timeUntilOncall > 1*time.Hour {
		fmt.Printf("Looks like %s doesn't have any oncalls starting soon.\nGoodbye!\n", userName)
		os.Exit(0)
	} else if timeUntilOncall > time.Duration(0) {
		fmt.Printf("%s is starting oncall in %s. Waiting until then.", userName, timeUntilOncall)
		time.Sleep(timeUntilOncall)
	}

	fmt.Println("Listening for incidents...")
	for {
		incidents := paginate(func(next uint) ([]pd.Incident, bool, uint) {
			response, err := client.ListIncidentsWithContext(context.Background(), pd.ListIncidentsOptions{
				ServiceIDs: serviceIDs,
				UserIDs:    []string{userID},
				TeamIDs:    teamIDs,
				Offset:     next,
			})
			if err != nil {
				return []pd.Incident{}, false, 0
			}
			return response.Incidents, response.More, response.Total
		})

		for _, incident := range incidents {
			if err := notifications.Notify(incident.Description, incident.Summary, alertIconPath); err != nil {
				log.Println("failed to send notification", err)
				fmt.Println("NEW INCIDENT")
				fmt.Println(incident.Description)
				fmt.Println(incident.Summary)
				log.Println(strings.Repeat("=", 80))
			}
		}
		time.Sleep(5 * time.Minute)
	}
}

func main() {
	if err := logic(); err != nil {
		log.Fatal(err)
	}
}
