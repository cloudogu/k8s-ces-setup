package context

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/sirupsen/logrus"
)

// Naming settings such as fqdn, hostname and domain.
type Naming struct {
	Fqdn            string `json:"fqdn"`
	Domain          string `json:"domain"`
	CertificateType string `json:"certificateType"`
	Certificate     string `json:"certificate"`
	CertificateKey  string `json:"certificateKey"`
	RelayHost       string `json:"relayHost"`
	MailAddress     string `json:"mailAddress"`
	Completed       bool   `json:"completed"`
	UseInternalIp   bool   `json:"useInternalIp"` //todo look at this
	InternalIp      string `json:"internalIp"`
}

// UserBackend contains configuration for the directory service.
type UserBackend struct {
	DsType             string `json:"dsType"`
	Server             string `json:"server"`
	AttributeID        string `json:"attributeID"`
	AttributeGivenName string `json:"attributeGivenName"`
	AttributeSurname   string `json:"attributeSurname"`
	AttributeFullname  string `json:"attributeFullname"`
	AttributeMail      string `json:"attributeMail"`
	AttributeGroup     string `json:"attributeGroup"`
	BaseDN             string `json:"baseDN"`
	SearchFilter       string `json:"searchFilter"`
	ConnectionDN       string `json:"connectionDN"`
	Password           string `json:"password"`
	Host               string `json:"host"`
	Port               string `json:"port"`
	LoginID            string `json:"loginID"`
	LoginPassword      string `json:"loginPassword"`
	Encryption         string `json:"encryption"`
	Completed          bool   `json:"completed"`

	GroupBaseDN               string `json:"groupBaseDN"`
	GroupSearchFilter         string `json:"groupSearchFilter"`
	GroupAttributeName        string `json:"groupAttributeName"`
	GroupAttributeDescription string `json:"groupAttributeDescription"`
	GroupAttributeMember      string `json:"groupAttributeMember"`
}

// User account for a ces instance.
type User struct {
	Username        string `json:"username"`
	Mail            string `json:"mail"`
	Password        string `json:"password"`
	AdminGroup      string `json:"adminGroup"`
	Completed       bool   `json:"completed"`
	AdminMember     bool   `json:"adminMember"`
	SendWelcomeMail bool   `json:"sendWelcomeMail"`
}

// Dogus struct defines which dogus are installed and which one is the default.
type Dogus struct {
	DefaultDogu string   `json:"defaultDogu"`
	Install     []string `json:"install"`
	Completed   bool     `json:"completed"`
}

// CustomKeyValue is a map of string -> map pairs.
type CustomKeyValue map[string]map[string]interface{}

// SetupConfiguration is the main struct for the configuration of the setup.
type SetupConfiguration struct {
	Naming                  Naming         `json:"naming"`
	Dogus                   Dogus          `json:"dogus"`
	Admin                   User           `json:"admin"`
	UserBackend             UserBackend    `json:"userBackend"`
	RegistryConfig          CustomKeyValue `json:"registryConfig"`
	RegistryConfigEncrypted CustomKeyValue `json:"registryConfigEncrypted"`
}

// IsCompleted checks if a SetupConfiguration is completed.
func (conf *SetupConfiguration) IsCompleted() bool {
	return conf.Naming.Completed && conf.Dogus.Completed && conf.Admin.Completed && conf.UserBackend.Completed
}

// ReadSetupConfig reads the setup configuration from a setup json file.
func ReadSetupConfig(path string) (SetupConfiguration, error) {
	config := SetupConfiguration{}

	fileInfo, err := os.Stat(path)

	// Kubernetes mounts not existent optional config maps as empty dirs, hence, an empty dir indicates that no
	// setup.json is available.
	if os.IsNotExist(err) || (fileInfo != nil && fileInfo.IsDir()) {
		logrus.Print("Found no setup.json")
		return config, nil
	} else if err != nil {
		return config, fmt.Errorf("could not find file at %s", path)
	}

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return config, fmt.Errorf("failed to read setup configuration %s: %w", path, err)
	}

	err = json.Unmarshal(data, &config)
	if err != nil {
		return config, fmt.Errorf("failed to unmarshal setup configuration %s: %w", path, err)
	}

	return config, nil
}
