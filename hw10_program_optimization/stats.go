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
	decoder := json.NewDecoder(r)
	result = make(users, 0, 100_000)

	for {
		var user User
		if err = decoder.Decode(&user); err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("failed to decode JSON: %w", err)
		}
		result = append(result, user)
	}

	return result, nil
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
