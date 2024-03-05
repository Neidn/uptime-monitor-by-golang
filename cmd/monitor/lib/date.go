package lib

import (
	"fmt"
	"github.com/google/go-github/v59/github"
	"log"
	"time"
)

// ConvertDateToHumanReadableTimeDifference converts a date to a human-readable time difference
// param dateTime time.Time: the time to convert
// returns string: the human-readable time difference
// e.g. 1 day, 2 hours, 3 minutes
// e.g. 1 hour, 2 minutes
func ConvertDateToHumanReadableTimeDifference(dateTime github.Timestamp) string {
	var result string

	// Get the current time
	currentTime := github.Timestamp{Time: time.Now()}

	log.Println("currentTime", currentTime)

	// Get the time difference
	timeDifference := currentTime.Time.Sub(dateTime.Time)

	// Get the number of days
	days := int(timeDifference.Hours() / 24)

	// Get the number of hours subtracting the days
	hours := int(timeDifference.Hours()) - (days * 24)

	// Get the number of minutes subtracting the days and hours
	minutes := int(timeDifference.Minutes()) - (days * 24 * 60) - (hours * 60)

	// add the days to the result
	if days > 0 {
		result += fmt.Sprintf("%d", days) + " day"
		if days > 1 {
			result += "s"
		}
	}

	// add the hours to the result
	if hours > 0 {
		if result != "" {
			result += ", "
		}
		result += fmt.Sprintf("%d", hours) + " hour"
		if hours > 1 {
			result += "s"
		}
	}

	// add the minutes to the result
	if minutes > 0 {
		if result != "" {
			result += ", "
		}
		result += fmt.Sprintf("%d", minutes) + " minute"
		if minutes > 1 {
			result += "s"
		}
	}

	return result
}
