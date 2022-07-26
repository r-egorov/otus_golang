package hw10programoptimization

import (
	"bufio"
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

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)
	regexPattern, err := regexp.Compile("\\." + domain)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)

	var p fastjson.Parser
	var v *fastjson.Value
	for scanner.Scan() {
		if v, err = p.Parse(scanner.Text()); err != nil {
			return nil, err
		}
		email := string(v.GetStringBytes("Email"))

		if regexPattern.MatchString(email) {
			result[strings.ToLower(strings.SplitN(email, "@", 2)[1])]++
		}
	}

	return result, nil
}
