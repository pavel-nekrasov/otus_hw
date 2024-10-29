package hw10programoptimization

import (
	"bufio"
	"errors"
	"fmt"
	"io"
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
	result, err := countDomains(r, domain)
	if err != nil {
		return nil, fmt.Errorf("count domains error: %w", err)
	}
	return result, nil
}

func countDomains(r io.Reader, domain string) (DomainStat, error) {
	if r == nil {
		return nil, ErrNilReader
	}

	result := make(DomainStat)

	domainSuffix := "." + domain
	// читаем сканнером - не надо загружать сразу весь файл в память
	// scanner читает переиспользуя внутренний буффер - экономим память
	scanner := bufio.NewScanner(r)
	// читаем в построчном режиме
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		user := User{}
		// scanner.Bytes() - не надо конвертировать строки в байты - сразу передаем байты (экономим на памяти)
		// ну и читаем easyjson - он быстрее + экономит память(меньше аллокаций)
		if err := easyjson.Unmarshal(scanner.Bytes(), &user); err != nil {
			return nil, err
		}

		if strings.HasSuffix(user.Email, domainSuffix) {
			// так вроде немного меньше памяти жрет. В SplitN  - аллоцируется две строки первая из которых нам заведомо не нужна
			idx := strings.Index(user.Email, "@")
			if idx < 0 {
				continue
			}
			domainPart := strings.ToLower(user.Email[idx+1:])
			result[domainPart]++
		}
	}
	return result, nil
}
