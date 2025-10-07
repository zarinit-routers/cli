package systemctl

import (
	"regexp"
)

const ServiceNameRegex = `^[a-zA-Z0-9-_\.:\\]+$`

var compiledRegex *regexp.Regexp

func init() {
	reg, err := regexp.Compile(ServiceNameRegex)
	if err != nil {
		log.Fatal("Failed compile regular expression", "error", err, "expression", ServiceNameRegex)
	}
	compiledRegex = reg
}

type Service struct {
	name string
}

func NewService(name string) Service {
	match := compiledRegex.MatchString(name)
	if !match {
		log.Fatal("Invalid service name", "name", name)
	}
	return Service{
		name: name,
	}
}
