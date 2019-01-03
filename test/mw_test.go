package soajsGoTest

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/soajs/soajs.golang"
	"github.com/stretchr/testify/assert"
	"gopkg.in/h2non/gock.v1"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

const soajsinjectobj = "{\"tenant\":{\"id\":\"5551aca9e179c39b760f7a1a\",\"code\":\"DBTN\"},\"key\":{\"config\":{\"mail\":{\"from\":\"soajstest@soajs.org\",\"transport\":{\"type\":\"smtp\",\"options\":{\"host\":\"secure.emailsrvr.com\",\"port\":\"587\",\"ignoreTLS\":true,\"secure\":false,\"auth\":{\"user\":\"soajstest@soajs.org\",\"pass\":\"p@ssw0rd\"}}}},\"oauth\":{\"loginMode\":\"urac\"}},\"iKey\":\"38145c67717c73d3febd16df38abf311\",\"eKey\":\"d44dfaaf1a3ba93adc6b3368816188f96134dfedec7072542eb3d84ec3e3d260f639954b8c0bc51e742c1dff3f80710e3e728edb004dce78d82d7ecd5e17e88c39fef78aa29aa2ed19ed0ca9011d75d9fc441a3c59845ebcf11f9393d5962549\"},\"application\":{\"product\":\"DSBRD\",\"package\":\"DSBRD_MAIN\",\"appId\":\"5512926a7a1f0e2123f638de\"},\"package\":{},\"device\":\"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/60.0.3112.101 Safari/537.36\",\"geo\":{\"ip\":\"127.0.0.1\"},\"awareness\":{\"host\":\"127.0.0.1\",\"port\":4000},\"urac\":{\"_id\": \"59a538becc083eecf37149df\", \"username\": \"owner\", \"firstName\": \"owner\", \"lastName\": \"owner\", \"email\": \"owner@soajs.org\", \"groups\": [ \"owner\" ], \"tenant\": { \"id\":\"5551aca9e179c39b760f7a1a\", \"code\": \"DBTN\" },\"profile\": {},\"acl\": null, \"acl_AllEnv\": null},\"param\":{\"id\":\"5551aca9e179c39b760f7a1a\"}}"

var registryOutputByte = []byte(`{"result":true,"ts":1538662943210,"service":{"service":"CONTROLLER","type":"rest","route":"/reloadRegistry"},"data":{"timeLoaded":1538662943199,"name":"dev","environment":"dev","profileOnly":false,"coreDB":{"provision":{"name":"core_provision","prefix":"","servers":[{"host":"127.0.0.1","port":27017}],"credentials":null,"streaming":{"batchSize":10000,"colName":{"batchSize":10000}},"URLParam":{"bufferMaxEntries":0},"registryLocation":{"l1":"coreDB","l2":"provision","env":"dev"},"timeConnected":1538662943201},"session":{"name":"core_session","prefix":"","store":{},"collection":"sessions","stringify":false,"expireAfter":1209600000,"registryLocation":{"l1":"coreDB","l2":"session","env":"dev"},"cluster":"dash_cluster","servers":[{"host":"192.168.30.31","port":27017}],"credentials":{},"URLParam":{"bufferMaxEntries":0,"maxPoolSize":5},"extraParam":{"db":{"native_parser":true,"bufferMaxEntries":0},"server":{}},"streaming":{}}},"tenantMetaDB":{"urac":{"prefix":"","cluster":"dash_cluster","servers":[{"host":"192.168.30.31","port":27017}],"credentials":{},"URLParam":{"bufferMaxEntries":0,"maxPoolSize":5},"extraParam":{"db":{"native_parser":true,"bufferMaxEntries":0},"server":{}},"streaming":{},"name":"#TENANT_NAME#_urac"}},"domain":"soajs.org","apiPrefix":"dev-api","sitePrefix":"dev","protocol":"http","port":80,"serviceConfig":{"awareness":{"cacheTTL":3600000,"healthCheckInterval":5000,"autoRelaodRegistry":1000,"maxLogCount":5,"autoRegisterService":true},"agent":{"topologyDir":"/opt/soajs/"},"logger":{"src":true,"level":"debug","formatter":{"levelInString":true,"outputMode":"long"}},"cors":{"enabled":true,"origin":"*","credentials":"true","methods":"GET,HEAD,PUT,PATCH,POST,DELETE","headers":"key,soajsauth,Keep-Alive,User-Agent,X-Requested-With,If-Modified-Since,Cache-Control,Content-Type,Authorization","maxage":1728000},"oauth":{"grants":["password","refresh_token"],"debug":false,"accessTokenLifetime":7200,"refreshTokenLifetime":1209600},"ports":{"controller":4000,"maintenanceInc":1000,"randomInc":100},"cookie":{"secret":"this is a secret sentence"},"session":{"name":"soajsID","secret":"server","cookie":{"path":"/","httpOnly":true,"secure":false,"maxAge":null},"resave":false,"saveUninitialized":false,"rolling":false,"unset":"keep"},"key":{"algorithm":"aes256","password":"soajs"}},"deployer":{"type":"manual","selected":"manual","manual":{"nodes":"127.0.0.1"},"container":{"docker":{"local":{"nodes":"127.0.0.1","socketPath":"/var/run/docker.sock"},"remote":{"nodes":"","apiProtocol":"","auth":{"token":""}}},"kubernetes":{"local":{"nodes":"127.0.0.1","namespace":"%namespace%","auth":{"token":"%kubetoken%"}},"remote":{"nodes":"127.0.0.1","namespace":"%namespace%","auth":{"token":"%kubetoken%"}}}}},"custom":{"tester":{"_id":"5bb758970bc295fbe40ae9d5","name":"tester","locked":false,"plugged":true,"shared":false,"value":{"test":true},"created":"DASHBOARD","author":"owner"}},"resources":{"cluster":{"dash_cluster":{"_id":"5b169de7a5d5f82f25f8db9a","name":"dash_cluster","type":"cluster","category":"mongo","created":"DEV","author":"owner","locked":true,"plugged":true,"shared":true,"config":{"servers":[{"host":"192.168.30.31","port":27017}],"credentials":{},"URLParam":{"bufferMaxEntries":0,"maxPoolSize":5},"extraParam":{"db":{"native_parser":true,"bufferMaxEntries":0},"server":{}},"streaming":{}}}}},"services":{"controller":{"group":"controller","maxPoolSize":100,"authorization":true,"port":4000,"requestTimeout":30,"requestTimeoutRenewal":null,"hosts":{"1":["127.0.0.1"],"latest":1}},"example01":{"group":"SOAJS Example Service","port":4010,"versions":{"1":{"extKeyRequired":false,"awareness":true,"urac":false,"urac_Profile":false,"urac_ACL":false,"provision_ACL":false,"apis":[{"l":"Test Get","v":"/testGet"},{"l":"Test Delete","v":"/testDel"},{"l":"Build Name","v":"/buildName"},{"l":"Test Post","v":"/testPost"},{"l":"Test Put","v":"/testPut"}]}},"requestTimeoutRenewal":5,"requestTimeout":30,"version":1,"extKeyRequired":false,"oauth":false}},"daemons":{"testDaemon":{"group":"tests","port":4501,"versions":{"1":{"jobs":{"buildCatalog":{}}}}}}}}`)

const (
	registryApi     = "127.0.0.1"
	registryApiPort = "5000"
	envCode         = "dev"
)

func StartRegistryServer(t *testing.T) {
	log.Println("Starting registry mock server ...")

	var registryOutput map[string]interface{}
	err := json.Unmarshal(registryOutputByte, &registryOutput)
	assert.NoError(t, err)

	gock.New("http://"+registryApi+":"+registryApiPort).
		Get("/getRegistry").
		Persist().
		MatchParam("env", "dev").
		MatchParam("serviceName", "mux").
		Reply(200).
		JSON(registryOutput)
	gock.EnableNetworking()
}

func StartTestServer(t *testing.T, handlerFunc func(*testing.T) http.HandlerFunc) *httptest.Server {
	jsonFile, err := os.Open("soa.json")
	if err != nil {
		panic(err)
	}
	defer jsonFile.Close()
	byteValue, _ := ioutil.ReadAll(jsonFile)
	var result map[string]interface{}
	json.Unmarshal([]byte(byteValue), &result)

	soajsMiddleware := soajsGo.InitMiddleware(result)
	testServer := httptest.NewServer(soajsMiddleware(handlerFunc(t)))

	return testServer
}

func CallTestServer(t *testing.T, testServer *httptest.Server) *http.Response {
	// construct the url
	var u bytes.Buffer
	u.WriteString(string(testServer.URL))
	u.WriteString("/")

	// init the http client
	httpClient := &http.Client{}
	req, _ := http.NewRequest("GET", u.String(), nil)

	// add soajsinjectobj header
	req.Header.Set("soajsinjectobj", soajsinjectobj)

	res, err := httpClient.Do(req)
	assert.NoError(t, err)
	if res != nil {
		defer res.Body.Close()
	}

	return res
}

func GetTestHandler(t *testing.T) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		log.Println("test handler reached!")

		soajs := r.Context().Value("soajs").(soajsGo.SOAJSObject)
		assert.NotEqual(t, soajs, soajsGo.SOAJSObject{})

		assert.Equal(t, soajs.Tenant.Id, "5551aca9e179c39b760f7a1a")
		assert.Equal(t, soajs.Tenant.Code, "DBTN")
		assert.Equal(t, soajs.Tenant.Key.IKey, "38145c67717c73d3febd16df38abf311")

		assert.Equal(t, soajs.Urac.Id, "59a538becc083eecf37149df")

		assert.Equal(t, soajs.Awareness.Host, "127.0.0.1")
		assert.Equal(t, soajs.Awareness.Port, 4000)

		_, err := soajs.Reg.GetDatabases()
		assert.Equal(t, errors.New("Environment registry not found"), err)
	}
	return http.HandlerFunc(fn)
}

func GetDatabaseOpsHandler(t *testing.T) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		log.Println("Database operations handler loaded ...")

		soajs := r.Context().Value("soajs").(soajsGo.SOAJSObject)
		assert.NotEqual(t, soajs, soajsGo.SOAJSObject{})

		time.Sleep(1000 * time.Millisecond)

		// get all databases
		dbs, err := soajs.Reg.GetDatabases()
		assert.NoError(t, err)
		assert.NotEmpty(t, dbs)

		// missing database name in param
		_, err = soajs.Reg.GetDatabase("")
		assert.Equal(t, errors.New("Database name is required"), err)

		// get one database
		oneDb, err := soajs.Reg.GetDatabase("provision")
		assert.NoError(t, err)
		assert.NotEmpty(t, oneDb)
	}
	return http.HandlerFunc(fn)
}

func GetServiceConfigOpsHandler(t *testing.T) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		log.Println("Service config operations handler loaded ...")

		soajs := r.Context().Value("soajs").(soajsGo.SOAJSObject)
		assert.NotEqual(t, soajs, soajsGo.SOAJSObject{})

		time.Sleep(1000 * time.Millisecond)

		// get service config
		serviceConfig, err := soajs.Reg.GetServiceConfig()
		assert.NoError(t, err)
		assert.NotEmpty(t, serviceConfig)
	}
	return http.HandlerFunc(fn)
}

func GetCustomRegistryOpsHandler(t *testing.T) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		log.Println("Custom registry operations handler loaded ...")

		soajs := r.Context().Value("soajs").(soajsGo.SOAJSObject)
		assert.NotEqual(t, soajs, soajsGo.SOAJSObject{})

		time.Sleep(1000 * time.Millisecond)

		// get custom registry
		customRegistry, err := soajs.Reg.GetCustom()
		assert.NoError(t, err)
		assert.NotEmpty(t, customRegistry)
	}
	return http.HandlerFunc(fn)
}

func GetResourcesOpsHandler(t *testing.T) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		log.Println("Resources operations handler loaded ...")

		soajs := r.Context().Value("soajs").(soajsGo.SOAJSObject)
		assert.NotEqual(t, soajs, soajsGo.SOAJSObject{})

		time.Sleep(1000 * time.Millisecond)

		// get all resources
		resources, err := soajs.Reg.GetResources()
		assert.NoError(t, err)
		assert.NotEmpty(t, resources)

		// missing resource name in param
		_, err = soajs.Reg.GetResource("")
		assert.Equal(t, errors.New("Resource name is required"), err)

		// get one resource
		oneResource, err := soajs.Reg.GetResource("dash_cluster")
		assert.NoError(t, err)
		assert.NotEmpty(t, oneResource)
	}
	return http.HandlerFunc(fn)
}

func GetServicesOpsHandler(t *testing.T) http.HandlerFunc {
	fn := func(w http.ResponseWriter, r *http.Request) {
		log.Println("Services operations handler loaded ...")

		soajs := r.Context().Value("soajs").(soajsGo.SOAJSObject)
		assert.NotEqual(t, soajs, soajsGo.SOAJSObject{})

		time.Sleep(1000 * time.Millisecond)

		// get all services
		services, err := soajs.Reg.GetServices()
		assert.NoError(t, err)
		assert.NotEmpty(t, services)

		// missing service name in param
		_, err = soajs.Reg.GetService("")
		assert.Equal(t, errors.New("Service name is required"), err)

		// get one service
		oneService, err := soajs.Reg.GetService("example01")
		assert.NoError(t, err)
		assert.NotEmpty(t, oneService)
	}
	return http.HandlerFunc(fn)
}

func TestSoajsMiddleware(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		URL         string
		Method      string
		Description string
	}{
		{
			URL:         "/tidbit/hello",
			Method:      "GET",
			Description: "Testing middleware behavior with GET requests",
		},
		{
			URL:         "/tidbit/hello",
			Method:      "POST",
			Description: "Testing middleware behavior with POST requests",
		},
	}

	testServer := httptest.NewServer(soajsGo.SoajsMiddleware(GetTestHandler(t)))
	defer testServer.Close()

	for _, testCase := range tests {
		// construct the url
		var u bytes.Buffer
		u.WriteString(string(testServer.URL))
		u.WriteString(testCase.URL)

		// init the http client
		httpClient := &http.Client{}
		req, _ := http.NewRequest(testCase.Method, u.String(), nil)

		// add soajsinjectobj header
		req.Header.Set("soajsinjectobj", soajsinjectobj)

		res, err := httpClient.Do(req)
		assert.NoError(err)
		if res != nil {
			defer res.Body.Close()
		}
	}
}

func TestExecRegistry(t *testing.T) {
	assert := assert.New(t)
	defer gock.Off()

	var registryOutput map[string]interface{}
	err := json.Unmarshal(registryOutputByte, &registryOutput)
	assert.NoError(err)

	gock.New("http://"+registryApi+":"+registryApiPort).
		Get("/getRegistry").
		MatchParam("env", "dev").
		MatchParam("serviceName", "mux").
		Reply(200).
		JSON(registryOutput)

	os.Setenv("SOAJS_ENV", envCode)
	os.Setenv("SOAJS_REGISTRY_API", registryApi+":"+registryApiPort)

	regObj, err := soajsGo.ExecRegistry(map[string]string{"envCode": "dev", "serviceName": "mux"})
	assert.NoError(err)
	assert.NotEmpty(regObj)

	// invalid port passed
	os.Setenv("SOAJS_REGISTRY_API", registryApi+":invalidPort")
	_, err = soajsGo.ExecRegistry(map[string]string{"envCode": "dev", "serviceName": "mux"})
	assert.Equal("Port must be an integer [invalidPort]", err.Error())

	// invalid registry api format
	os.Setenv("SOAJS_REGISTRY_API", "127.0.0.1")
	_, err = soajsGo.ExecRegistry(map[string]string{"envCode": "dev", "serviceName": "mux"})
	assert.Equal("Invalid format for SOAJS_REGISTRY_API [hostname:port]: 127.0.0.1", err.Error())
}

func TestGetters(t *testing.T) {
	log.Println("Testing getters ...")

	StartRegistryServer(t)
	defer gock.Off()

	os.Setenv("SOAJS_ENV", envCode)
	os.Setenv("SOAJS_REGISTRY_API", registryApi+":"+registryApiPort)

	handlers := map[string]func(t *testing.T) http.HandlerFunc{
		"databases":      GetDatabaseOpsHandler,
		"serviceConfig":  GetServiceConfigOpsHandler,
		"customRegistry": GetCustomRegistryOpsHandler,
		"resources":      GetResourcesOpsHandler,
		"services":       GetServicesOpsHandler,
	}

	for handlerName, handlerFunc := range handlers {
		log.Println("Testing " + handlerName + " ...")

		testServer := StartTestServer(t, handlerFunc)

		CallTestServer(t, testServer)
		testServer.Close()
	}
}
