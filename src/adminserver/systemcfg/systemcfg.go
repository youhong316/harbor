// Copyright Project Harbor Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package systemcfg

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	enpt "github.com/goharbor/harbor/src/adminserver/systemcfg/encrypt"
	"github.com/goharbor/harbor/src/adminserver/systemcfg/store"
	"github.com/goharbor/harbor/src/adminserver/systemcfg/store/database"
	"github.com/goharbor/harbor/src/adminserver/systemcfg/store/encrypt"
	"github.com/goharbor/harbor/src/adminserver/systemcfg/store/json"
	"github.com/goharbor/harbor/src/common"
	comcfg "github.com/goharbor/harbor/src/common/config"
	"github.com/goharbor/harbor/src/common/dao"
	"github.com/goharbor/harbor/src/common/models"
	"github.com/goharbor/harbor/src/common/utils"
	"github.com/goharbor/harbor/src/common/utils/log"
)

const (
	defaultJSONCfgStorePath string = "/etc/adminserver/config/config.json"
	defaultKeyPath          string = "/etc/adminserver/key"
	ldapScopeKey            string = "ldap_scope"
)

var (
	// CfgStore is a storage driver that configurations
	// can be read from and wrote to
	CfgStore store.Driver

	// attrs need to be encrypted or decrypted
	attrs = []string{
		common.EmailPassword,
		common.LDAPSearchPwd,
		common.PostGreSQLPassword,
		common.AdminInitialPassword,
		common.ClairDBPassword,
		common.UAAClientSecret,
	}

	// all configurations need read from environment variables
	allEnvs = map[string]interface{}{
		common.ExtEndpoint: "EXT_ENDPOINT",
		common.AUTHMode:    "AUTH_MODE",
		common.SelfRegistration: &parser{
			env:   "SELF_REGISTRATION",
			parse: parseStringToBool,
		},
		common.DatabaseType:   "DATABASE_TYPE",
		common.PostGreSQLHOST: "POSTGRESQL_HOST",
		common.PostGreSQLPort: &parser{
			env:   "POSTGRESQL_PORT",
			parse: parseStringToInt,
		},
		common.PostGreSQLUsername: "POSTGRESQL_USERNAME",
		common.PostGreSQLPassword: "POSTGRESQL_PASSWORD",
		common.PostGreSQLDatabase: "POSTGRESQL_DATABASE",
		common.PostGreSQLSSLMode:  "POSTGRESQL_SSLMODE",
		common.LDAPURL:            "LDAP_URL",
		common.LDAPSearchDN:       "LDAP_SEARCH_DN",
		common.LDAPSearchPwd:      "LDAP_SEARCH_PWD",
		common.LDAPBaseDN:         "LDAP_BASE_DN",
		common.LDAPFilter:         "LDAP_FILTER",
		common.LDAPUID:            "LDAP_UID",
		common.LDAPScope: &parser{
			env:   "LDAP_SCOPE",
			parse: parseStringToInt,
		},
		common.LDAPTimeout: &parser{
			env:   "LDAP_TIMEOUT",
			parse: parseStringToInt,
		},
		common.LDAPVerifyCert: &parser{
			env:   "LDAP_VERIFY_CERT",
			parse: parseStringToBool,
		},
		common.LDAPGroupBaseDN:        "LDAP_GROUP_BASEDN",
		common.LDAPGroupSearchFilter:  "LDAP_GROUP_FILTER",
		common.LDAPGroupAttributeName: "LDAP_GROUP_GID",
		common.LDAPGroupSearchScope: &parser{
			env:   "LDAP_GROUP_SCOPE",
			parse: parseStringToInt,
		},
		common.EmailHost: "EMAIL_HOST",
		common.EmailPort: &parser{
			env:   "EMAIL_PORT",
			parse: parseStringToInt,
		},
		common.EmailUsername: "EMAIL_USR",
		common.EmailPassword: "EMAIL_PWD",
		common.EmailSSL: &parser{
			env:   "EMAIL_SSL",
			parse: parseStringToBool,
		},
		common.EmailInsecure: &parser{
			env:   "EMAIL_INSECURE",
			parse: parseStringToBool,
		},
		common.EmailFrom:     "EMAIL_FROM",
		common.EmailIdentity: "EMAIL_IDENTITY",
		common.RegistryURL:   "REGISTRY_URL",
		common.TokenExpiration: &parser{
			env:   "TOKEN_EXPIRATION",
			parse: parseStringToInt,
		},
		common.CfgExpiration: &parser{
			env:   "CFG_EXPIRATION",
			parse: parseStringToInt,
		},
		common.MaxJobWorkers: &parser{
			env:   "MAX_JOB_WORKERS",
			parse: parseStringToInt,
		},
		common.ProjectCreationRestriction: "PROJECT_CREATION_RESTRICTION",
		common.AdminInitialPassword:       "HARBOR_ADMIN_PASSWORD",
		common.AdmiralEndpoint:            "ADMIRAL_URL",
		common.WithNotary: &parser{
			env:   "WITH_NOTARY",
			parse: parseStringToBool,
		},
		common.WithClair: &parser{
			env:   "WITH_CLAIR",
			parse: parseStringToBool,
		},
		common.ClairDBPassword: "CLAIR_DB_PASSWORD",
		common.ClairDB:         "CLAIR_DB",
		common.ClairDBUsername: "CLAIR_DB_USERNAME",
		common.ClairDBHost:     "CLAIR_DB_HOST",
		common.ClairDBPort: &parser{
			env:   "CLAIR_DB_PORT",
			parse: parseStringToInt,
		},
		common.ClairDBSSLMode:  "CLAIR_DB_SSLMODE",
		common.UAAEndpoint:     "UAA_ENDPOINT",
		common.UAAClientID:     "UAA_CLIENTID",
		common.UAAClientSecret: "UAA_CLIENTSECRET",
		common.UAAVerifyCert: &parser{
			env:   "UAA_VERIFY_CERT",
			parse: parseStringToBool,
		},
		common.CoreURL:                     "CORE_URL",
		common.JobServiceURL:               "JOBSERVICE_URL",
		common.TokenServiceURL:             "TOKEN_SERVICE_URL",
		common.ClairURL:                    "CLAIR_URL",
		common.NotaryURL:                   "NOTARY_URL",
		common.RegistryStorageProviderName: "REGISTRY_STORAGE_PROVIDER_NAME",
		common.ReadOnly: &parser{
			env:   "READ_ONLY",
			parse: parseStringToBool,
		},
		common.ReloadKey:        "RELOAD_KEY",
		common.LdapGroupAdminDn: "LDAP_GROUP_ADMIN_DN",
		common.ChartRepoURL:     "CHART_REPOSITORY_URL",
		common.WithChartMuseum: &parser{
			env:   "WITH_CHARTMUSEUM",
			parse: parseStringToBool,
		},
	}

	// configurations need read from environment variables
	// every time the system startup
	repeatLoadEnvs = map[string]interface{}{
		common.ExtEndpoint:    "EXT_ENDPOINT",
		common.PostGreSQLHOST: "POSTGRESQL_HOST",
		common.PostGreSQLPort: &parser{
			env:   "POSTGRESQL_PORT",
			parse: parseStringToInt,
		},
		common.PostGreSQLUsername: "POSTGRESQL_USERNAME",
		common.PostGreSQLPassword: "POSTGRESQL_PASSWORD",
		common.PostGreSQLDatabase: "POSTGRESQL_DATABASE",
		common.PostGreSQLSSLMode:  "POSTGRESQL_SSLMODE",
		common.MaxJobWorkers: &parser{
			env:   "MAX_JOB_WORKERS",
			parse: parseStringToInt,
		},
		common.CfgExpiration: &parser{
			env:   "CFG_EXPIRATION",
			parse: parseStringToInt,
		},
		common.AdmiralEndpoint: "ADMIRAL_URL",
		common.WithNotary: &parser{
			env:   "WITH_NOTARY",
			parse: parseStringToBool,
		},
		common.WithClair: &parser{
			env:   "WITH_CLAIR",
			parse: parseStringToBool,
		},
		common.ClairDBPassword: "CLAIR_DB_PASSWORD",
		common.ClairDBHost:     "CLAIR_DB_HOST",
		common.ClairDBUsername: "CLAIR_DB_USERNAME",
		common.ClairDBPort: &parser{
			env:   "CLAIR_DB_PORT",
			parse: parseStringToInt,
		},
		common.ClairDBSSLMode:  "CLAIR_DB_SSLMODE",
		common.UAAEndpoint:     "UAA_ENDPOINT",
		common.UAAClientID:     "UAA_CLIENTID",
		common.UAAClientSecret: "UAA_CLIENTSECRET",
		common.UAAVerifyCert: &parser{
			env:   "UAA_VERIFY_CERT",
			parse: parseStringToBool,
		},
		common.RegistryStorageProviderName: "REGISTRY_STORAGE_PROVIDER_NAME",
		common.CoreURL:                     "CORE_URL",
		common.JobServiceURL:               "JOBSERVICE_URL",
		common.RegistryURL:                 "REGISTRY_URL",
		common.TokenServiceURL:             "TOKEN_SERVICE_URL",
		common.ClairURL:                    "CLAIR_URL",
		common.NotaryURL:                   "NOTARY_URL",
		common.DatabaseType:                "DATABASE_TYPE",
		common.ChartRepoURL:                "CHART_REPOSITORY_URL",
		common.WithChartMuseum: &parser{
			env:   "WITH_CHARTMUSEUM",
			parse: parseStringToBool,
		},
	}
)

type parser struct {
	// the name of env
	env string
	// parse the value of env, e.g. parse string to int or
	// parse string to bool
	parse func(string) (interface{}, error)
}

func parseStringToInt(str string) (interface{}, error) {
	if len(str) == 0 {
		return 0, nil
	}
	return strconv.Atoi(str)
}

func parseStringToBool(str string) (interface{}, error) {
	return strings.ToLower(str) == "true" ||
		strings.ToLower(str) == "on", nil
}

// Init system configurations. If env RESET is set or configurations
// read from storage driver is null, load all configurations from env
func Init() (err error) {
	// init database
	envCfgs := map[string]interface{}{}
	if err := LoadFromEnv(envCfgs, true); err != nil {
		return err
	}
	db := GetDatabaseFromCfg(envCfgs)
	if err := dao.InitDatabase(db); err != nil {
		return err
	}
	if err := dao.UpgradeSchema(db); err != nil {
		return err
	}
	if err := dao.CheckSchemaVersion(); err != nil {
		return err
	}

	if err := initCfgStore(); err != nil {
		return err
	}

	// Use reload key to avoid reset customed setting after restart
	curCfgs, err := CfgStore.Read()
	if err != nil {
		return err
	}
	loadAll := isLoadAll(curCfgs)
	if curCfgs == nil {
		curCfgs = map[string]interface{}{}
	}
	// restart: only repeatload envs will be load
	// reload_config: all envs will be reload except the skiped envs
	if err = LoadFromEnv(curCfgs, loadAll); err != nil {
		return err
	}
	AddMissedKey(curCfgs)
	return CfgStore.Write(curCfgs)
}

func isLoadAll(cfg map[string]interface{}) bool {
	return cfg == nil || strings.EqualFold(os.Getenv("RESET"), "true") && os.Getenv("RELOAD_KEY") != cfg[common.ReloadKey]
}

func initCfgStore() (err error) {

	drivertype := os.Getenv("CFG_DRIVER")
	if len(drivertype) == 0 {
		drivertype = common.CfgDriverDB
	}
	path := os.Getenv("JSON_CFG_STORE_PATH")
	if len(path) == 0 {
		path = defaultJSONCfgStorePath
	}
	log.Infof("the path of json configuration storage: %s", path)

	if drivertype == common.CfgDriverDB {
		CfgStore, err = database.NewCfgStore()
		if err != nil {
			return err
		}
		// migration check: if no data in the db , then will try to load from path
		m, err := CfgStore.Read()
		if err != nil {
			return err
		}
		if m == nil || len(m) == 0 {
			if _, err := os.Stat(path); err == nil {
				jsondriver, err := json.NewCfgStore(path)
				if err != nil {
					log.Errorf("Failed to migrate configuration from %s", path)
					return err
				}
				jsonconfig, err := jsondriver.Read()
				if err != nil {
					log.Errorf("Failed to read old configuration from %s", path)
					return err
				}
				// Update LDAP Scope for migration
				// only used when migrating harbor release before v1.3
				// after v1.3 there is always a db configuration before migrate.
				validLdapScope(jsonconfig, true)
				err = CfgStore.Write(jsonconfig)
				if err != nil {
					log.Error("Failed to update old configuration to database")
					return err
				}
			}
		}
	} else {
		CfgStore, err = json.NewCfgStore(path)
		if err != nil {
			return err
		}
	}
	kp := os.Getenv("KEY_PATH")
	if len(kp) == 0 {
		kp = defaultKeyPath
	}
	log.Infof("the path of key used by key provider: %s", kp)

	encryptor := enpt.NewAESEncryptor(
		comcfg.NewFileKeyProvider(kp), nil)

	CfgStore = encrypt.NewCfgStore(encryptor, attrs, CfgStore)
	return nil
}

// LoadFromEnv loads the configurations from allEnvs, if all is false, it just loads
// the repeatLoadEnvs and the env which is absent in cfgs
func LoadFromEnv(cfgs map[string]interface{}, all bool) error {
	var envs map[string]interface{}

	if all {
		envs = allEnvs
	} else {
		envs = make(map[string]interface{})
		for k, v := range repeatLoadEnvs {
			envs[k] = v
		}
		for k, v := range allEnvs {
			if _, exist := cfgs[k]; !exist {
				envs[k] = v
			}
		}
	}

	reloadCfg := os.Getenv("RESET")
	skipPattern := os.Getenv("SKIP_RELOAD_ENV_PATTERN")
	skipPattern = strings.TrimSpace(skipPattern)
	if len(skipPattern) == 0 {
		skipPattern = "$^" // doesn't match any string by default
	}
	skipMatcher, err := regexp.Compile(skipPattern)
	if err != nil {
		log.Errorf("Regular express parse error, skipPattern:%v", skipPattern)
		skipMatcher = regexp.MustCompile("$^")
	}

	for k, v := range envs {
		if str, ok := v.(string); ok {
			if skipMatcher.MatchString(str) && strings.EqualFold(reloadCfg, "true") {
				continue
			}
			cfgs[k] = os.Getenv(str)
			continue
		}

		if parser, ok := v.(*parser); ok {
			if skipMatcher.MatchString(parser.env) && strings.EqualFold(reloadCfg, "true") {
				continue
			}
			i, err := parser.parse(os.Getenv(parser.env))
			if err != nil {
				return err
			}
			cfgs[k] = i
			continue
		}

		return fmt.Errorf("%v is not string or parse type", v)
	}
	validLdapScope(cfgs, false)
	return nil
}

// GetDatabaseFromCfg Create database object from config
func GetDatabaseFromCfg(cfg map[string]interface{}) *models.Database {
	database := &models.Database{}
	database.Type = cfg[common.DatabaseType].(string)
	postgresql := &models.PostGreSQL{}
	postgresql.Host = utils.SafeCastString(cfg[common.PostGreSQLHOST])
	postgresql.Port = int(utils.SafeCastInt(cfg[common.PostGreSQLPort]))
	postgresql.Username = utils.SafeCastString(cfg[common.PostGreSQLUsername])
	postgresql.Password = utils.SafeCastString(cfg[common.PostGreSQLPassword])
	postgresql.Database = utils.SafeCastString(cfg[common.PostGreSQLDatabase])
	postgresql.SSLMode = utils.SafeCastString(cfg[common.PostGreSQLSSLMode])
	database.PostGreSQL = postgresql
	return database
}

// Valid LDAP Scope
func validLdapScope(cfg map[string]interface{}, isMigrate bool) {
	ldapScope, ok := cfg[ldapScopeKey].(int)
	if !ok {
		ldapScopeFloat, ok := cfg[ldapScopeKey].(float64)
		if ok {
			ldapScope = int(ldapScopeFloat)
		}
	}
	if isMigrate && ldapScope > 0 && ldapScope < 3 {
		ldapScope = ldapScope - 1
	}
	if ldapScope >= 3 {
		ldapScope = 2
	}
	if ldapScope < 0 {
		ldapScope = 0
	}
	cfg[ldapScopeKey] = ldapScope

}

// AddMissedKey ... If the configure key is missing in the cfg map, add default value to it
func AddMissedKey(cfg map[string]interface{}) {

	for k, v := range common.HarborStringKeysMap {
		if _, exist := cfg[k]; !exist {
			cfg[k] = v
		}
	}

	for k, v := range common.HarborNumKeysMap {
		if _, exist := cfg[k]; !exist {
			cfg[k] = v
		}
	}

	for k, v := range common.HarborBoolKeysMap {
		if _, exist := cfg[k]; !exist {
			cfg[k] = v
		}
	}

}
