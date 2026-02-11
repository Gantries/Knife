// Package nacos contains nacos config service utility.
package nacos

import (
	"reflect"

	"github.com/gantries/knife/pkg/maps"
	"github.com/gantries/knife/pkg/national"
	"github.com/gantries/knife/pkg/tel"
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
	"gopkg.in/yaml.v3"
)

var logger = tel.Logger("pkg/nacos")

// Data contains a configuration data identifier.
type Data struct {
	Group string `yaml:"group" json:"group"`
	ID    string `yaml:"id" json:"id"`
}

// Configuration contains nacos configuration structure
type Configuration struct {
	ServerAddress   string `yaml:"server_address" json:"server_address"`
	Timeout         int    `yaml:"timeout" json:"timeout"`
	Port            int    `yaml:"port" json:"port"`
	DataID          string `yaml:"data_id" json:"data_id"`
	ContextPath     string `yaml:"context_path" json:"context_path"`
	Group           string `yaml:"group" json:"group"`
	Namespace       string `yaml:"namespace" json:"namespace"`
	Scheme          string `yaml:"scheme" json:"scheme"`
	Username        string `yaml:"username" json:"username"`
	Password        string `yaml:"password" json:"password"`
	Languages       []Data `yaml:"languages" json:"languages"`
	DefaultLanguage string `yaml:"default_language" json:"default_language"`
}

// Watcher contains definition to watch remote configuration changes
type Watcher struct {
	client config_client.IConfigClient
}

// NewWatcher is the constructor to create Watcher instances.
func NewWatcher(config Configuration) (*Watcher, error) {
	var client config_client.IConfigClient
	var err error
	if config.ServerAddress != "" {
		serverConfigs := []constant.ServerConfig{
			{
				IpAddr:      config.ServerAddress,
				Port:        uint64(config.Port), // #nosec G115 - type conversion is safe for port numbers
				ContextPath: config.ContextPath,
				Scheme:      config.Scheme,
			},
		}
		clientConfig := constant.ClientConfig{
			NamespaceId:         config.Namespace,
			TimeoutMs:           uint64(config.Timeout), // #nosec G115 - type conversion is safe for timeout values
			NotLoadCacheAtStart: true,
			Username:            config.Username,
			Password:            config.Password,
		}

		client, err = clients.CreateConfigClient(map[string]interface{}{
			"serverConfigs": serverConfigs,
			"clientConfig":  clientConfig,
		})
		if err != nil {
			return nil, err
		}
	} else {
		client = nil
	}

	return &Watcher{client: client}, nil
}

// Get should be used to get specific configuration from nacos server.
func (w *Watcher) Get(id string, group string, output any) error {
	if w.client == nil {
		logger.Info("Skip getting config for empty client")
		return nil
	}
	data, err := w.client.GetConfig(vo.ConfigParam{
		DataId: id,
		Group:  group,
	})
	if err != nil {
		logger.Error("Unable to get configuration", "group", group, "id", id, "error", err)
		return err
	}
	err = yaml.Unmarshal([]byte(data), output)
	if err != nil {
		logger.Error("Unable to deserialize config", "group", group, "id",
			id, "type", reflect.TypeOf(output), "error", err)
		return err
	}
	return nil
}

// Watch should be used to watch a configuration changes.
func (w *Watcher) Watch(id string, group string, action func(ns, group, id, cfg string)) error {
	if w.client == nil {
		logger.Info("Skip watching config for empty client")
		return nil
	}
	return w.client.ListenConfig(vo.ConfigParam{
		DataId:   id,
		Group:    group,
		OnChange: action,
	})
}

// DefaultLang returns default language setting.
func (c Configuration) DefaultLang() string {
	return c.DefaultLanguage
}

// SetupInternational is used to prepare i18n support.
func (w *Watcher) SetupInternational(c Configuration) {
	if len(c.DefaultLanguage) > 0 {
		national.Prepare(c)
	}
	if len(c.Languages) > 0 {
		for _, lang := range c.Languages {
			messages := maps.Map[string, maps.Map[string, string]]{}
			if err := w.Get(lang.ID, lang.Group, messages); err != nil {
				logger.Warn("Unable to load messages", "error", err)
				continue
			}
			national.LoadMessages(messages)
			if err := w.Watch(lang.ID, lang.Group, func(ns, group, id, cfg string) {
				national.LoadMessagesFromString(cfg)
			}); err != nil {
				logger.Warn("Unable to watch config", "error", err)
			}
		}
	}
}
