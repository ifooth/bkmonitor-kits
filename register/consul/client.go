package consul

import (
	"github.com/hashicorp/consul/api"
)

// Client client 封装了 consul 部分读写操作
type Client struct {
	agent   *api.Agent
	address string // address IP:Port
}

// NewClient 传入的 address 应符合 IP:Port 的结构，例如: 127.0.0.1:8080
func NewClient(address string) (*Client, error) {
	client := new(Client)
	client.address = address

	conf := api.DefaultConfig()
	conf.Address = address

	apiClient, err := api.NewClient(conf)
	if err != nil {
		return nil, err
	}

	client.agent = apiClient.Agent()
	return client, nil
}

// ServiceDeregister 注销 service, name 既是 ID
func (bc *Client) ServiceDeregister(serviceID string) error {
	return bc.agent.ServiceDeregister(serviceID)
}

func (bc *Client) GetOrCreateService(serviceID, serviceName string, tags []string, address string, port int) error {
	_, options, err := bc.agent.AgentHealthServiceByID(serviceID)
	if options == nil && err == nil {
		option := new(api.AgentServiceRegistration)
		option.Name = serviceName
		option.ID = serviceID
		option.Tags = tags
		option.Address = address
		option.Port = port

		return bc.agent.ServiceRegister(option)
	}

	return nil
}

// CheckRegister 注册检查器
func (bc *Client) CheckRegister(serviceID, checkID string, ttl string) error {
	option := new(api.AgentCheckRegistration)
	option.Name = checkID
	option.TTL = ttl
	option.ServiceID = serviceID

	return bc.agent.CheckRegister(option)
}

// CheckDeregister 注销检查器
func (bc *Client) CheckDeregister(checkID string) error {
	return bc.agent.CheckDeregister(checkID)
}

// CheckFail 标记状态为 fail
func (bc *Client) CheckFail(checkID, note string) error {
	return bc.agent.FailTTL(checkID, note)
}

// CheckPass 标记状态为 pass
func (bc *Client) CheckPass(checkID, note string) error {
	return bc.agent.PassTTL(checkID, note)
}
