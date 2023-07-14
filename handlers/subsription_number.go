package handlers

import (
	"time"
)

type SubscriptionNumberPackageInterface interface {
	GetTimePeriod() (time.Time, time.Time)
}

type OneMonthSubscriptionNumberPackage struct {
	StartDate time.Time
	EndDate   time.Time
}

type TwoMonthSubscriptionNumberPackage struct {
	StartDate time.Time
	EndDate   time.Time
}

func (s *OneMonthSubscriptionNumberPackage) GetTimePeriod() (time.Time, time.Time) {
	s.EndDate = s.StartDate.AddDate(0, 1, 0)

	return s.StartDate, s.EndDate
}

func (s *TwoMonthSubscriptionNumberPackage) GetTimePeriod() (time.Time, time.Time) {
	s.EndDate = s.StartDate.AddDate(0, 2, 0)

	return s.StartDate, s.EndDate
}
