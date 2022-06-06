package main

import (
	"strings"
	"time"
)

func appendNoDupes(sessionData SessionData) []SessionData {
	for _, currentSession := range stats.SessionData {
		if strings.EqualFold(currentSession.JsSha1, sessionData.JsSha1) {
			return stats.SessionData
		}
	}

	return append(stats.SessionData, sessionData)
}

func getTime() time.Time {
	return time.Now().Add(time.Hour * 13)
}
