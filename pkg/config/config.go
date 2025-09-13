package config

import (
	"strings"

	"github.com/lwabish/transaction-mapper/pkg/transaction"
	"github.com/samber/lo"
	"gopkg.in/yaml.v3"
)

type Config struct {
	DualLevel     bool                `yaml:"dualLevel"`
	LevelSplitter string              `yaml:"levelSplitter"`
	Rules         map[string][]string `yaml:"rules"`
	Default       struct {
		Outcome string `yaml:"outcome"`
		Income  string `yaml:"income"`
	} `yaml:"default"`
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
	return c.inferByRules(t)
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

func (c *Config) inferByRules(t transaction.Transaction) (string, string) {
	desc := t.Description
	defaultCategory := lo.Ternary(t.Amount > 0, c.Default.Outcome, c.Default.Income)
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
		parts := strings.Split(defaultCategory, c.LevelSplitter)
		return parts[0], parts[1]
	}
	return "", defaultCategory
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
