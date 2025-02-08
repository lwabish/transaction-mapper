package config

import (
	"github.com/lwabish/transaction-mapper/pkg/transaction"
	"gopkg.in/yaml.v3"
	"strings"
)

type Config struct {
	DualLevel            bool                  `yaml:"dualLevel"`
	LevelSplitter        string                `yaml:"levelSplitter"`
	Rules                map[string][]string   `yaml:"rules"`
	Default              string                `yaml:"default"`
	TransferAccountRules []transferAccountRule `yaml:"transferAccountRules"`
	keywordsToCategory   map[string]string
}

type transferAccountRule struct {
	Account       string `yaml:"account"`
	Keyword       string `yaml:"keyword"`
	ToAccountType string `yaml:"toAccountType"`
	ToAccountName string `yaml:"toAccountName"`
}

func (c *Config) InferCategory(t transaction.Transaction) (string, string) {
	return c.inferByRules(t.Description)
}

func (c *Config) InferTransferToAccount(t transaction.Transaction, ai transaction.AccountInfo) (string, string) {
	for _, rule := range c.TransferAccountRules {
		if rule.Account == ai.Name {
			if strings.Contains(t.Description, rule.Keyword) {
				return rule.ToAccountType, rule.ToAccountName
			}
		}
	}
	return "", ""
}

func (c *Config) inferByRules(desc string) (string, string) {
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

func NewConfig(bs []byte) (*Config, error) {
	c := &Config{
		Rules:              make(map[string][]string),
		keywordsToCategory: make(map[string]string),
	}
	if err := yaml.Unmarshal(bs, c); err != nil {
		return nil, err
	}
	for key, value := range c.Rules {
		for _, v := range value {
			c.keywordsToCategory[v] = key
		}
	}
	return c, nil
}

type Loader struct {
	config *Config
}

func (r *Loader) LoadConf(c *Config) {
	r.config = c
}

func (r *Loader) GetConf() *Config {
	return r.config
}
