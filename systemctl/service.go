package systemctl

import (
	"regexp"
)

const ServiceNameRegex = `^[a-zA-Z0-9-_\.:\\]+$`

type Service string

func NewService(name string) Service {
	match, err := regexp.MatchString(ServiceNameRegex, name)
	if err != nil {
		log.Fatal("Failed to match service name", "error", err, "name", name)
	}
	if !match {
		log.Fatal("Invalid service name", "name", name)
	}
	return Service(name)
}
