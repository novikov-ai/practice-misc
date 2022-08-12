package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/valyala/fastjson"
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

const UserFieldEmail = "Email"

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users [100_000]User

func getUsers(r io.Reader) (result users, err error) {
	reader := bufio.NewScanner(r)

	parser := fastjson.Parser{}
	i := 0

	for reader.Scan() {
		content := reader.Bytes()

		v, err := parser.ParseBytes(content)
		if err != nil {
			return result, err
		}

		result[i] = User{Email: string(v.GetStringBytes(UserFieldEmail))}
		i++
	}

	return
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)

	for _, user := range u {
		if strings.HasSuffix(user.Email, "."+domain) {
			domainName := strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
			result[domainName]++
		}
	}
	return result, nil
}
