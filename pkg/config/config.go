package config

import (
	"github.com/lwabish/transaction-mapper/pkg/transaction"
	"gopkg.in/yaml.v3"
	"log"
	"os"
	"strings"
)

type config struct {
	DualLevel          bool                `yaml:"dualLevel"`
	LevelSplitter      string              `yaml:"levelSplitter"`
	Rules              map[string][]string `yaml:"rules"`
	Default            string              `yaml:"default"`
	keywordsToCategory map[string]string
}

var (
	Config = &config{
		Rules:              make(map[string][]string),
		keywordsToCategory: make(map[string]string),
	}
)

func init() {
	log.Println("load config rules from file")
	bs, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalln(err)
	}
	err = yaml.Unmarshal(bs, &Config)
	if err != nil {
		log.Fatal(err)
	}
	for key, value := range Config.Rules {
		for _, v := range value {
			Config.keywordsToCategory[v] = key
		}
	}
}

func (c *config) InferCategory(t transaction.Transaction) (string, string) {
	return c.inferByRules(t.Description)
}

func (c *config) inferByRules(desc string) (string, string) {
	for k, v := range c.keywordsToCategory {
		if strings.Contains(desc, k) {
			if c.DualLevel {
				parts := strings.Split(v, c.LevelSplitter)
				return parts[0], parts[1]
			}
			return "", v
		}
	}
	if c.DualLevel {
		parts := strings.Split(c.Default, c.LevelSplitter)
		return parts[0], parts[1]
	}
	return "", c.Default
}
