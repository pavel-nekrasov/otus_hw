package hw10programoptimization

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"

	"github.com/mailru/easyjson"
)

//easyjson:json
type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

var ErrNilReader = errors.New("nil reader")

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users []*User

func getUsers(r io.Reader) (users, error) {
	if r == nil {
		return nil, ErrNilReader
	}
	// читаем сканнером - не надо загружать сразу весь файл в память
	// scanner читает переиспользуя внутренний буффер - экономим память
	scanner := bufio.NewScanner(r)
	// читаем в построчном режиме
	scanner.Split(bufio.ScanLines)
	result := make(users, 0)

	for scanner.Scan() {
		user := User{}
		// scanner.Bytes() - не надо конвертировать строки в байты - сразу передаем байты (экономим на памяти)
		// ну и читаем easyjson - он быстрее + экономит память(меньше аллокаций)
		if err := easyjson.Unmarshal(scanner.Bytes(), &user); err != nil {
			return nil, err
		}
		// т.к. слайс указателей - данные копировать не надо - в слайс заносим только указатель
		result = append(result, &user)
	}
	return result, nil
}

func countDomains(u users, domain string) (DomainStat, error) {
	result := make(DomainStat)
	// один раз компилируем regex
	r, err := regexp.Compile("\\." + domain)
	if err != nil {
		return nil, err
	}

	for _, user := range u {
		matched := r.Match([]byte(user.Email))

		if matched {
			// так вроде немного меньше памяти жрет. В SplitN  - аллоцируется две строки первая из которых нам заведомо не нужна
			idx := strings.Index(user.Email, "@")
			if idx < 0 {
				continue
			}
			domainPart := strings.ToLower(user.Email[idx+1:])

			// domainPart := strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
			num := result[domainPart]
			num++
			result[domainPart] = num
		}
	}
	return result, nil
}
