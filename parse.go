package app

import (
	"errors"
	"fmt"
	"log"
	"regexp"
	"strconv"
)

const CreateReminderFromMessagePattern = `(?i)new reminder (?P<name>.*): (?P<text>.*): (?P<time>.*)`
const UpdateReminderFromMessagePattern = `(?i)update (?P<referenceId>.*): (?P<time>.*)`
const ReminderHoursMinutesMessagePattern = `(?i)((?P<hours>[0-9])+H)?\s?((?P<minutes>[0-9])+M)?`
const ReminderTimeMessagePattern = `(?i)(?P<year>[0-9]{4})(?P<month>[0-9]{2})(?P<day>[0-9]{2}) (?P<hour>[0-9]{1,2}):(?P<minute>[0-9]{2}) (?P<tz>[a-z]*\/[a-z]*)`

func ParseCreateReminderMessage(message string) (string, string, int, error) {
	// Messages requesting the creation of a reminder are formatted as follows:
	// "New Reminder <Reminder Name>: <Reminder Text>: <#H #M"
	log.Printf("parseCreateReminderMessage %s", message)
	var name, text string
	var nMinutes int

	match, err := regexp.Compile(CreateReminderFromMessagePattern)
	if err != nil {
		return name, text, nMinutes, err
	}
	result, err := getNamedCaptureGroups(match, message)
	if err != nil {
		return name, text, nMinutes, err
	}
	name = result["name"]
	text = result["text"]
	messageTime := result["time"]
	nMinutes, err = getReminderNMinutesFromMessage(messageTime)
	return name, text, nMinutes, err
}

func ParseUpdateReminderMessage(message string) (string, int, error) {
	// Messages requesting the update of a reminder are formatted as follows:
	// "Update <Reference ID>: <#H #M"
	log.Printf("parseUpdateReminderMessage %s", message)
	var referenceId string
	var nMinutes int

	match, err := regexp.Compile(UpdateReminderFromMessagePattern)
	if err != nil {
		return referenceId, nMinutes, err
	}
	result, err := getNamedCaptureGroups(match, message)
	if err != nil {
		return referenceId, nMinutes, err
	}
	referenceId = result["referenceId"]
	messageTime := result["time"]
	nMinutes, err = getReminderNMinutesFromMessage(messageTime)
	return referenceId, nMinutes, err
}

func getNamedCaptureGroups(r *regexp.Regexp, str string) (map[string]string, error) {
	match := r.FindStringSubmatch(str)
	results := make(map[string]string)
	if len(match) == 0 {
		return results, errors.New(fmt.Sprintf("Unable to calculate requested reminder time from %s", str))
	}
	for i, name := range r.SubexpNames() {
		if i != 0 {
			results[name] = match[i]
		}
	}
	return results, nil
}

func ReminderParseError(messageTime string) error {
	return errors.New(fmt.Sprintf("Unable to calculate requested reminder time from %s", messageTime))
}

func getReminderNMinutesFromMessage(messageTime string) (int, error) {
	var nMinutes int
	hmMatch, _ := regexp.Compile(ReminderHoursMinutesMessagePattern)

	if hmMatch != nil {
		result, err := getNamedCaptureGroups(hmMatch, messageTime)
		if err != nil {
			return nMinutes, ReminderParseError(messageTime)
		}
		hours, err := strconv.Atoi(result["hours"])
		if err != nil {
			return nMinutes, ReminderParseError(messageTime)
		}
		minutes, err := strconv.Atoi(result["minutes"])
		if err != nil {
			return nMinutes, ReminderParseError(messageTime)
		}
		return hours*60 + minutes, nil
	}
	return nMinutes, ReminderParseError(messageTime)
}
