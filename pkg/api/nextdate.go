package api

import (
	"13spr/pkg/utils"
	"errors"
	"fmt"
	"strings"
	"time"
)

// NextDate вычисляет следующую дату по правилу повторения
func NextDate(now time.Time, dstart, repeat string) (string, error) {
	if repeat == "" {
		return "", errors.New("empty repeat rule")
	}

	// Парсим начальную дату
	start, err := time.Parse(utils.DateFormat, dstart)
	if err != nil {
		return "", fmt.Errorf("invalid start date: %w", err)
	}

	// Разбиваем правило на части
	parts := strings.Split(repeat, " ")
	if len(parts) == 0 {
		return "", errors.New("invalid repeat rule")
	}

	var interval int
	var isDay, isYear, isWeek, isMonth bool
	var weekDays []int
	var monthDays []int
	var months []int

	// Определяем тип правила
	switch parts[0] {
	case "d":
		isDay = true
		if len(parts) < 2 {
			return "", errors.New("missing interval for 'd'")
		}
		interval, err = parseInterval(parts[1])
		if err != nil {
			return "", err
		}
		if interval < 1 || interval > 400 {
			return "", errors.New("interval out of range [1, 400]")
		}
	case "y":
		isYear = true
		interval = 365 // год
	case "w":
		isWeek = true
		if len(parts) < 2 {
			return "", errors.New("missing days for 'w'")
		}
		weekDays, err = parseWeekDays(parts[1])
		if err != nil {
			return "", err
		}
	case "m":
		isMonth = true
		if len(parts) < 2 {
			return "", errors.New("missing days for 'm'")
		}
		monthParts := strings.Split(parts[1], ",")
		for _, p := range monthParts {
			d, err := parseInterval(p)
			if err != nil {
				return "", err
			}
			monthDays = append(monthDays, d)
		}
		// Опционально: вторая последовательность месяцев
		if len(parts) >= 3 {
			for _, p := range strings.Split(parts[2], ",") {
				m, err := parseInterval(p)
				if err != nil {
					return "", err
				}
				months = append(months, m)
			}
		}
	default:
		return "", errors.New("unsupported repeat rule")
	}

	// Теперь ищем ближайшую дату после now
	date := start
	for !utils.AfterNow(date, now) {
		switch {
		case isDay:
			date = date.AddDate(0, 0, interval)
		case isYear:
			date = date.AddDate(1, 0, 0)
		case isWeek:
			date = findNextWeekDay(date, weekDays)
		case isMonth:
			date = findNextMonthDay(date, monthDays, months)
		}
	}

	return date.Format(utils.DateFormat), nil
}

// parseInterval преобразует строку в число
func parseInterval(s string) (int, error) {
	if s == "" {
		return 0, errors.New("empty interval")
	}
	if s[0] == '-' {
		if len(s) == 1 {
			return 0, errors.New("invalid interval")
		}
		var n int
		_, err := fmt.Sscanf(s, "%d", &n)
		if err != nil || n < -31 || n > 31 {
			return 0, errors.New("invalid day in month")
		}
		return n, nil
	}
	var n int
	_, err := fmt.Sscanf(s, "%d", &n)
	if err != nil || n < 1 || n > 400 {
		return 0, errors.New("invalid interval")
	}
	return n, nil
}

// parseWeekDays парсит список дней недели (1–7)
func parseWeekDays(s string) ([]int, error) {
	var days []int
	for _, part := range strings.Split(s, ",") {
		day, err := parseInterval(part)
		if err != nil || day < 1 || day > 7 {
			return nil, errors.New("invalid weekday")
		}
		days = append(days, day)
	}
	return days, nil
}

// findNextWeekDay находит ближайший день недели из списка
func findNextWeekDay(date time.Time, days []int) time.Time {
	next := date
	for {
		next = next.AddDate(0, 0, 1)
		weekday := int(next.Weekday())
		if weekday == 0 {
			weekday = 7
		}
		for _, d := range days {
			if d == weekday {
				return next
			}
		}
	}
}

// findNextMonthDay находит ближайшую дату в месяце
func findNextMonthDay(date time.Time, days []int, months []int) time.Time {
	next := date
	for {
		next = next.AddDate(0, 0, 1)
		_, month, day := next.Date()
		// Проверяем, является ли день одним из заданных
		for _, d := range days {
			if d == day {
				// Проверяем месяц, если указаны месяцы
				if len(months) == 0 || contains(months, int(month)) {
					return next
				}
			}
		}
	}
}

func contains(slice []int, item int) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}
