package rdcfg

/*
1. package rdcfg is designed to read and parse configures

2. the config file must follow FORMAT 1 or FORMAT 2

FORMAT 1
========================================
# notes
# default section
ip =   127.0.0.1
port=1234


FORMAT 2:
=========================================
#default
 ip = 127.0.0.1

port  =  7890

# for broker
[broker]
  listen_port = 7777
time = 10

# for logger
[logger]
listen_port = 1888

Usage:
=========================================
* use cfg, err := NewRDCFG(configpath) to create a RDCFG obj
* use value, err := cfg.GetDefault("ip") to get string value in form1 (or use cfg.Get("", "ip") intead)
* use value, err := cfg.Get("broker", "listen_port") to get string value in form2
*
* other functions for use: GetInt、GetFloat64、 GetIntDefault、GetFloat64、Default
*/

import (
	"errors"
	"os"
	"regexp"
	"strconv"
	"strings"
)

const (
	RDCFG_DEFAULT_SECTION_NAME = " "
)

type RDCFG struct {
	m map[string]*section
}

func NewRDCFG(config string) (*RDCFG, error) {
	cfg := &RDCFG{
		m: make(map[string]*section, 10),
	}

	fi, err := os.Stat(config)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(config)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buff := make([]byte, fi.Size())
	f.Read(buff)
	str := string(buff)
	if !strings.HasSuffix(str, "\n") {
		return nil, errors.New("Config file does not end with a newline character.")
	}

	err = cfg.loadCfg(string(buff))
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (cfg *RDCFG) Get(sectionName string, key string) (string, error) {
	if sectionName == "" {
		sectionName = RDCFG_DEFAULT_SECTION_NAME
	}

	if section, err := cfg.getSection(sectionName); err != nil {
		return "", err
	} else {
		return section.get(key)
	}
}

func (cfg *RDCFG) GetInt(sectionName string, key string) (int, error) {
	sv, err := cfg.Get(sectionName, key)
	if err != nil {
		return 0, err
	}

	value, err := strconv.Atoi(sv)
	if err != nil {
		return 0, err
	}

	return value, err
}

func (cfg *RDCFG) GetFloat64(sectionName string, key string) (float64, error) {
	sv, err := cfg.Get(sectionName, key)

	if err != nil {
		return 0, err
	}

	value, err := strconv.ParseFloat(sv, 64)
	if err != nil {
		return 0.0, err
	}

	return value, err
}

func (cfg *RDCFG) GetDefault(key string) (string, error) {
	return cfg.Get("", key)
}

func (cfg *RDCFG) GetDefaultInt(sectionName string, key string) (int, error) {
	return cfg.GetInt("", key)
}

func (cfg *RDCFG) GetDefaultFloat64(sectionName string, key string) (float64, error) {
	return cfg.GetFloat64("", key)
}

/********************* private functions ************************************/

func (cfg *RDCFG) set(sectionName string, key string, value string) error {
	if sectionName == "" {
		sectionName = RDCFG_DEFAULT_SECTION_NAME
	}

	if section, err := cfg.getSection(sectionName); err != nil {
		return err
	} else {
		return section.set(key, value)
	}
}

func (cfg *RDCFG) setDefault(section string, key string, value string) error {
	return cfg.set("", key, value)
}

// convert file buf to RDCFG
func (cfg *RDCFG) loadCfg(filebuf string) error {
	//trim all space
	var r = strings.NewReplacer(" ", "")
	news := r.Replace(filebuf)

	// filter
	re := regexp.MustCompile("\\[.*\\]\n|.+=.+\\n")
	allLines := re.FindAllString(news, -1)

	nowSectionName := RDCFG_DEFAULT_SECTION_NAME
	cfg.addSection(nowSectionName)

	for _, line := range allLines {
		if strings.HasPrefix(line, "#") {
			continue
		}

		line := strings.TrimRight(line, "\n")

		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			tmp := strings.Trim(line, "[]")

			// section name must not be "" or " " ...
			if tmp == "" {
				return errors.New("section name is null")
			}

			if nowSectionName != tmp {
				nowSectionName = tmp
				cfg.addSection(nowSectionName)
			} else {
				return errors.New("repeat section name")
			}
		} else if pos := strings.Index(line, "="); pos >= 0 {
			k := line[:pos]
			v := line[pos+1:]
			err := cfg.set(nowSectionName, k, v)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (cfg *RDCFG) getSection(secname string) (*section, error) {
	if sec, ok := cfg.m[secname]; ok {
		return sec, nil
	} else {
		return nil, errors.New("specified section not found")
	}
}

func (cfg *RDCFG) addSection(secname string) {
	sec := newSection(secname)
	cfg.m[secname] = sec
}

func (cfg *RDCFG) addDefaultSection() {
	cfg.addSection(RDCFG_DEFAULT_SECTION_NAME)
}

// section define
type section struct {
	name string
	m    map[string]string
}

func newSection(name string) *section {
	sec := &section{
		name: name,
		m:    make(map[string]string, 100),
	}

	return sec
}

func (s *section) get(key string) (string, error) {
	if value, ok := s.m[key]; ok {
		return value, nil
	} else {
		return "", errors.New("could not find value by key")
	}
}

// may fail if key already exists
// when succeed, return nil
func (s *section) set(key string, value string) error {
	if _, ok := s.m[key]; ok {
		return errors.New("key already exists")
	} else {
		s.m[key] = value
		return nil
	}
}

// add value whether key exists or not
func (s *section) setForce(key string, value string) {
	s.m[key] = value
}
