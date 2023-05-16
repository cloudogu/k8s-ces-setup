package context

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"os"

	"github.com/sirupsen/logrus"
)

// Naming settings such as fqdn, hostname and domain.
type Naming struct {
	// Fqdn contains the complete fully qualified domain name of the Cloudogu EcoSystem.
	Fqdn string `json:"fqdn"`
	// Domain is primarily used to send emails from within the EcoSystem.
	Domain string `json:"domain"`
	// CertificateType is the type of certificate used to connect to the EcoSystem.
	CertificateType string `json:"certificateType"`
	// Certificate is a PEM-formatted certificate used to connect to the EcoSystem.
	// This is only necessary if CertificateType is set to "external".
	Certificate string `json:"certificate"`
	// CertificateKey is a PEM-formatted certificate key for the EcoSystem.
	// This is only necessary if CertificateType is set to "external".
	CertificateKey string `json:"certificateKey"`
	// RelayHost over which mails get sent from the EcoSystem.
	RelayHost string `json:"relayHost"`
	// MailAddress is used by all dogus to send mail.
	MailAddress string `json:"mailAddress"`
	// Completed indicates that the Naming step should not be shown in the UI of the setup.
	Completed bool `json:"completed"`
	// UseInternalIp configures if InternalIp should be used.
	UseInternalIp bool `json:"useInternalIp"`
	// InternalIp is useful if an external loadbalancer with its own IP is configured in front of the Cloudogu EcoSystem.
	// It can be set to let dogus communicate directly within the Cloudogu EcoSystem without the detour over the load balancer.
	InternalIp string `json:"internalIp"`
}

// UserBackend contains configuration for the directory service.
type UserBackend struct {
	// DsType is the type of the UserBackend. If set to "embedded", the ldap dogu will be installed and used as a user backend.
	// If set to "external", the credentials for an external user backend have to be set.
	DsType string `json:"dsType"`
	// Server contains the type of user backend server. Can either be "activeDirectory" or "custom".
	// This is only necessary if DsType is set to "external".
	Server string `json:"server"`
	// AttributeID contains the name of the attribute describing the user id in the user backend.
	// Must be "uid" if DsType is "embedded". Must be "sAMAccountName", if DsType is "external" and Server is "activeDirectory".
	AttributeID string `json:"attributeID"`
	// AttributeGivenName contains the name of the attribute describing the given name of a user.
	// This is only necessary if DsType is set to "external".
	AttributeGivenName string `json:"attributeGivenName"`
	// AttributeSurname contains the name of the attribute describing the surname of a user.
	// This is only necessary if DsType is set to "external".
	AttributeSurname string `json:"attributeSurname"`
	// AttributeFullname contains the name of the attribute describing the full name of a user.
	// Must be "cn" if DsType is "embedded" or Server is "activeDirectory".
	AttributeFullname string `json:"attributeFullname"`
	// AttributeMail contains the name of the attribute describing the mail address of a user.
	// Must be "mail" if DsType is "embedded" or Server is "activeDirectory".
	AttributeMail string `json:"attributeMail"`
	// AttributeGroup contains the name of the attribute managing the membership of the user to a particular group.
	// Must be "memberOf" if DsType is "embedded" or Server is "activeDirectory".
	AttributeGroup string `json:"attributeGroup"`
	// BaseDN is the distinguished name from which the server is searched for users.
	// This is only necessary if DsType is set to "external".
	BaseDN string `json:"baseDN"`
	// SearchFilter is restricting which object classes should be searched.
	// Must be "(objectClass=person)" if DsType is "embedded" or Server is "activeDirectory".
	SearchFilter string `json:"searchFilter"`
	// ConnectionDN is the distinguished name of a user that is authorized to read in the user backend.
	// This is only necessary if DsType is set to "external".
	ConnectionDN string `json:"connectionDN"`
	// Password of the user in ConnectionDN.
	// This is only necessary if DsType is set to "external".
	Password string `json:"password"`
	// Host address of the external user backend.
	// This is only necessary if DsType is set to "external".
	// Must be "ldap" if DsType is "embedded".
	Host string `json:"host"`
	// Port of the external user backend.
	// This is only necessary if DsType is set to "external".
	// Must be "389" if DsType is "embedded".
	Port          string `json:"port"`
	LoginID       string `json:"loginID"`
	LoginPassword string `json:"loginPassword"`
	// Encryption determines if and how communication with the user backend should be encrypted.
	// Can be "none", "ssl", "sslAny", "startTLS" or "startTLSAny".
	// This is only necessary if DsType is set to "external".
	Encryption string `json:"encryption"`
	// Completed indicates that the UserBackend step should not be shown in the UI of the setup.
	Completed bool `json:"completed"`

	// GroupBaseDN is the distinguished name for the group mapping.
	// This is only necessary if DsType is set to "external".
	GroupBaseDN string `json:"groupBaseDN"`
	// GroupSearchFilter is restricting which object classes should be searched for the group mapping.
	// This is only necessary if DsType is set to "external".
	GroupSearchFilter string `json:"groupSearchFilter"`
	// GroupAttributeName contains the name of the attribute of the group name.
	// This is only necessary if DsType is set to "external".
	GroupAttributeName string `json:"groupAttributeName"`
	// GroupAttributeDescription contains the name of the attribute for the group description.
	// This is only necessary if DsType is set to "external".
	GroupAttributeDescription string `json:"groupAttributeDescription"`
	// GroupAttributeMember contains the name of the attribute for the group members.
	// This is only necessary if DsType is set to "external".
	GroupAttributeMember string `json:"groupAttributeMember"`
}

// User account for a Cloudogu EcoSystem instance.
type User struct {
	Username string `json:"username"`
	Mail     string `json:"mail"`
	Password string `json:"password"`
	// AdminGroup is the name of the group in the user backend that should gain admin privileges.
	AdminGroup string `json:"adminGroup"`
	// Completed indicates that this step should not be shown in the UI of the setup.
	Completed bool `json:"completed"`
	// AdminMember determines if this user should become a member of the AdminGroup.
	AdminMember     bool `json:"adminMember"`
	SendWelcomeMail bool `json:"sendWelcomeMail"`
}

// Dogus struct defines which dogus are installed and which one is the default.
type Dogus struct {
	// DefaultDogu is the dogu that a call to the EcoSystem in the browser should be redirected to.
	DefaultDogu string `json:"defaultDogu"`
	// Install contains a list of all dogus that should be installed during the setup.
	// Entries may contain a version. If they do not, the latest version will be used.
	Install []string `json:"install"`
	// Completed indicates that this step should not be shown in the UI of the setup.
	Completed bool `json:"completed"`
}

// CustomKeyValue is a map of string -> map pairs.
type CustomKeyValue map[string]map[string]interface{}

// SetupConfiguration is the main struct for the configuration of the setup.
type SetupConfiguration struct {
	// Naming configures for example FQDN, mail and certificate configuration of the EcoSystem.
	Naming Naming `json:"naming"`
	// Dogus configures the installed dogus.
	Dogus Dogus `json:"dogus"`
	// Admin configures the admin user of the EcoSystem.
	Admin User `json:"admin"`
	// UserBackend configures where and how users are stored.
	UserBackend UserBackend `json:"userBackend"`
	// RegistryConfig contains custom registry configuration that is to be applied to the EcoSystem.
	RegistryConfig CustomKeyValue `json:"registryConfig"`
	// RegistryConfigEncrypted also contains custom registry configuration but with encrypted values.
	RegistryConfigEncrypted CustomKeyValue `json:"registryConfigEncrypted"`
}

// IsCompleted checks if a SetupConfiguration is completed.
func (conf *SetupConfiguration) IsCompleted() bool {
	return conf.Naming.Completed && conf.Dogus.Completed && conf.Admin.Completed && conf.UserBackend.Completed
}

// ReadSetupConfigFromCluster reads the setup configuration from the configmap
func ReadSetupConfigFromCluster(client kubernetes.Interface, namespace string) (*SetupConfiguration, error) {
	configMap, err := client.CoreV1().ConfigMaps(namespace).Get(context.Background(), SetupStartUpConfigMap, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get setup.json configmap: %w", err)
	}

	config := &SetupConfiguration{}
	stringData := configMap.Data["setup.json"]
	err = json.Unmarshal([]byte(stringData), config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal setup configuration from configmap: %w", err)
	}

	return config, nil
}

// ReadSetupConfigFromFile reads the setup configuration from a setup json file.
func ReadSetupConfigFromFile(path string) (*SetupConfiguration, error) {
	config := &SetupConfiguration{}

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
