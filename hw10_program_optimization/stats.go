package hw10programoptimization

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users []User

func getUsers(r io.Reader) (result users, err error) {
	content, err := io.ReadAll(r)
	if err != nil {
		return
	}

	lines := strings.Split(string(content), "\n")
	result = make(users, 0, len(lines))
	for _, line := range lines {
		var user User
		if err = json.Unmarshal([]byte(line), &user); err != nil {
			return
		}
		result = append(result, user)
	}
	return
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)
	domainSuffix := "." + domain

	for _, user := range u {
		email := strings.ToLower(user.Email)
		if strings.HasSuffix(email, domainSuffix) {
			parts := strings.SplitN(email, "@", 2)
			if len(parts) == 2 {
				domainPart := parts[1]
				result[domainPart]++
			}
		}
	}
	return result, nil
}
