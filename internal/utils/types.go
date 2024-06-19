package utils

import (
	"fmt"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// TimeRegex Regular expression for HH:MM:SS format.
var TimeRegex = regexp.MustCompile(`^([01]?[0-9]|2[0-3]):[0-5][0-9]:[0-5][0-9]$`)

func AttrValueToString(val attr.Value) string {
	switch value := val.(type) {
	case types.String:
		return value.ValueString()
	default:
		return value.String()
	}
}

func DayTimeFormatToSeconds(timeStr string) (int64, error) {
	timeComponents := strings.Split(timeStr, ":")

	hours, err := strconv.Atoi(timeComponents[0])
	if err != nil {
		return 0, err
	}

	minutes, err := strconv.Atoi(timeComponents[1])
	if err != nil {
		return 0, err
	}

	seconds, err := strconv.Atoi(timeComponents[2])
	if err != nil {
		return 0, err
	}

	totalSeconds := hours*3600 + minutes*60 + seconds

	return int64(totalSeconds), nil
}

func SecondsToDayTimeFormat(seconds int) string {
	duration := time.Duration(seconds) * time.Second

	formattedDuration := fmt.Sprintf("%02d:%02d:%02d", int(duration.Hours()), int(duration.Minutes())%60, int(duration.Seconds())%60)

	return formattedDuration
}

func ConvertInterfaceToListOfString(input interface{}) (*[]string, error) {
	var output []string
	switch assertedInput := input.(type) {
	case []interface{}:
		for _, item := range assertedInput {
			if stringItem, ok := item.(string); !ok {
				return nil, fmt.Errorf("input is not a list of strings")
			} else {
				output = append(output, stringItem)
			}
		}
	case string:
		output = append(output, assertedInput)
	case nil:
		return nil, nil
	default:
		return nil, fmt.Errorf("input is not a list of strings")
	}

	return &output, nil
}

func ConvertStringArray(items []types.String) []string {
	var result []string
	for _, item := range items {
		result = append(result, item.ValueString())
	}
	return result
}
