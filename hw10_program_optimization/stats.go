package hw10programoptimization

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
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

const (
	UserFieldId       = "Id"
	UserFieldName     = "Name"
	UserFieldUsername = "Username"
	UserFieldEmail    = "Email"
	UserFieldPhone    = "Phone"
	UserFieldPassword = "Password"
	UserFieldAddress  = "Address"
)

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

		user := getUserFromParsedValue(v)

		result[i] = *user
		i++
	}

	return
}

func getUserFromParsedValue(value *fastjson.Value) *User {
	return &User{
		ID:       value.GetInt(UserFieldId),
		Name:     string(value.GetStringBytes(UserFieldName)),
		Username: string(value.GetStringBytes(UserFieldUsername)),
		Email:    string(value.GetStringBytes(UserFieldEmail)),
		Phone:    string(value.GetStringBytes(UserFieldPhone)),
		Password: string(value.GetStringBytes(UserFieldPassword)),
		Address:  string(value.GetStringBytes(UserFieldAddress)),
	}
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)

	regExpr, err := regexp.Compile("\\." + domain)
	if err != nil {
		return nil, err
	}

	for _, user := range u {
		matched := regExpr.MatchString(user.Email)
		if matched {
			domainName := strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
			result[domainName]++
		}
	}
	return result, nil
}
