package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/json-iterator/go"
)

var json = jsoniter.ConfigFastest

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type UserEmail struct {
	Email string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	scanner := bufio.NewScanner(r)
	result := make(DomainStat)
	domainSuffix := "." + domain

	for scanner.Scan() {
		line := scanner.Bytes()
		var user UserEmail
		if err := json.Unmarshal(line, &user); err != nil {
			return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
		}

		email := strings.ToLower(user.Email)
		if strings.HasSuffix(email, domainSuffix) {
			parts := strings.SplitN(email, "@", 2)
			if len(parts) == 2 {
				domainPart := parts[1]
				result[domainPart]++
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scanner error: %w", err)
	}

	return result, nil
}
