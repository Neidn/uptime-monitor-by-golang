package workflows

import (
	"github.com/Neidn/uptime-monitor-by-golang/config"
	"github.com/google/go-github/v59/github"
	"log"
	"math"
	"time"
)

type DownPercent struct {
	Day              string
	Week             string
	Month            string
	Year             string
	All              string
	DailyMinutesDown map[string]int
}

type Downtime struct {
	Day              int
	Week             int
	Month            int
	Year             int
	All              int
	DailyMinutesDown map[string]int
}

type TimeOverlap struct {
	Start time.Time
	End   time.Time
}

func GetUptimePercentForSite(
	slugName string,
	issues []*github.Issue,
) (downPercent DownPercent, err error) {
	siteHistory, err := ReadSiteHistory(slugName)
	if err != nil {
		return DownPercent{}, err
	}

	// Calculate Start Time
	startTime, err := time.Parse(config.TimeFormat, siteHistory.StartTime)
	if err != nil {
		return DownPercent{}, err
	}

	// Calculate Total Seconds
	// Now - Start Time
	totalSeconds := int(time.Since(startTime).Seconds())
	log.Println("Total Seconds: ", totalSeconds)

	// Calculate Down Seconds
	downtimeSeconds, err := GetDowntimeFromSiteIssues(issues)
	if err != nil {
		return DownPercent{}, err
	}

	log.Println("Downtime Seconds: ", downtimeSeconds)

	downPercent.Day = math.Max(0, 100-float64(downtimeSeconds.Day)/float64(86400)*100)

	return
}

func GetDowntimeFromSiteIssues(
	issues []*github.Issue,
) (downtime Downtime, err error) {
	day := 0
	week := 0
	month := 0
	year := 0

	for _, issue := range issues {
		dailyMinutesDown := map[string]int{}
		end := time.Now()

		// Calculate Issue Downtime
		var issueDowntime int
		if issue.ClosedAt == nil {
			issueDowntime = int(time.Since(*issue.CreatedAt.GetTime()).Seconds())
		} else {
			issueDowntime = int(issue.ClosedAt.GetTime().Sub(*issue.CreatedAt.GetTime()).Seconds())
		}
		downtime.All += issueDowntime

		log.Println("Issue Downtime: ", issueDowntime)

		issueOverlap := TimeOverlap{
			Start: *issue.CreatedAt.GetTime(),
			End:   *issue.ClosedAt.GetTime(),
		}

		// Calculate Issue Downtime by Day
		for day := 0; day < 365; day++ {
			date := time.Now().AddDate(0, 0, -day)
			tmpOverlap := TimeOverlap{
				Start: time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, time.UTC),
				End:   time.Date(date.Year(), date.Month(), date.Day(), 23, 59, 59, 0, time.UTC),
			}

			overlap := checkOverlap(issueOverlap, tmpOverlap)
			if overlap <= 0 {
				continue
			}

			log.Println("Issue Downtime by Day: ", date, overlap)

			dateKey := date.Format("2006-01-02")
			if _, ok := dailyMinutesDown[dateKey]; !ok {
				dailyMinutesDown[dateKey] = 0
			} else {
				dailyMinutesDown[dateKey] += overlap / 60
			}
			dailyMinutesDown[date.Format("2006-01-02")] += overlap
		}

		day += checkOverlap(issueOverlap, TimeOverlap{
			// Start: one day ago
			Start: time.Now().AddDate(0, 0, -1),
			End:   end,
		})

		week += checkOverlap(issueOverlap, TimeOverlap{
			// Start: one week ago
			Start: time.Now().AddDate(0, 0, -7),
			End:   end,
		})

		month += checkOverlap(issueOverlap, TimeOverlap{
			// Start: one month ago
			Start: time.Now().AddDate(0, -1, 0),
			End:   end,
		})

		year += checkOverlap(issueOverlap, TimeOverlap{
			// Start: one year ago
			Start: time.Now().AddDate(-1, 0, 0),
			End:   end,
		})
	}

	downtime.Day = day
	downtime.Week = week
	downtime.Month = month
	downtime.Year = year

	return downtime, nil
}

func checkOverlap(
	firstOverlap TimeOverlap,
	secondOverlap TimeOverlap,
) int {
	var minOverlap TimeOverlap
	var maxOverlap TimeOverlap

	if firstOverlap.Start.Before(secondOverlap.Start) {
		minOverlap = firstOverlap
		maxOverlap = secondOverlap
	} else {
		minOverlap = secondOverlap
		maxOverlap = firstOverlap
	}

	// minOverlap.End < maxOverlap.Start
	// return 0
	if minOverlap.End.Before(maxOverlap.Start) {
		return 0
	}

	// minOverlap.End - maxOverlap.Start
	// return int(seconds)
	return int(minOverlap.End.Sub(maxOverlap.Start).Seconds())
}
