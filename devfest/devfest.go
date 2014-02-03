package main

import (
	"appengine"
	"appengine/datastore"
	"appengine/memcache"
	"appengine/urlfetch"
	"encoding/json"
	"fmt"
	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful/swagger"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

var twilioAccountSid string
var twilioAuthToken string
var twilioFrom string
var twilioTo string

type Activity struct {
	DeviceID string
	Kind     string
	Date     time.Time
}

type Sms struct {
	Value string
}

type EventService struct {
	// service API resource
}

func (u EventService) Register() {
	ws := new(restful.WebService)

	ws.
		Path("/events").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML) // you can specify this per route as well

	ws.Route(ws.GET("").To(u.getAllEvents).
		// docs
		Doc("update a user").
		Reads(Activity{})) // from the request

	ws.Route(ws.POST("").To(u.createEvent).
		// docs
		Doc("create an event").
		Param(ws.PathParameter("user-id", "identifier of the user").DataType("string")).
		Reads(Activity{})) // from the request

	restful.Add(ws)

	ws2 := new(restful.WebService)

	ws2.
		Path("/sms").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML) // you can specify this per route as well

	ws2.Route(ws2.GET("").To(u.getSMSValue).
		// docs
		Doc("get value").
		Reads(Sms{})) // from the request

	ws2.Route(ws2.PUT("").To(u.setSMSValue).
		// docs
		Doc("create an event").
		Param(ws2.PathParameter("value", "identifier of the user").DataType("string")).
		Reads(Sms{})) // from the request

	restful.Add(ws2)

}

func eventKey(c appengine.Context) *datastore.Key {
	return datastore.NewKey(c, "Collection", "default_stuff", 0, nil)
}

func postToTwilio(r *http.Request) error {
	formValues := url.Values{}

	formValues.Set("Body", "Motion detected")
	formValues.Set("From", twilioFrom)
	formValues.Set("To", twilioTo)
	
	log.Printf("SID VALUE %s", twilioAccountSid )
	log.Printf("twilioAuthToken %s", twilioAuthToken )
	log.Printf("FROM %s", twilioFrom )
	log.Printf("TO %s", twilioTo )
	

	req, err := http.NewRequest("POST", "https://api.twilio.com/2010-04-01/Accounts/" + twilioAccountSid + "/SMS/Messages.json", strings.NewReader(formValues.Encode()))
	if err != nil {
		return err
	}

	req.SetBasicAuth(twilioAccountSid, twilioAuthToken)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	c := appengine.NewContext(r)

	client := urlfetch.Client(c)
	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return err
	}

	_, err = ioutil.ReadAll(resp.Body)
	return err
}

// GET http://localhost:8080/sms
func (u *EventService) getSMSValue(request *restful.Request, response *restful.Response) {
	c := appengine.NewContext(request.Request)
	sms := new(Sms)
	_, err := memcache.Gob.Get(c, "sms", &sms)
	if err != nil {
		response.WriteErrorString(http.StatusNotFound, "SMS value not found.")
	} else {
		response.WriteEntity(sms)
	}
}

// PUT http://localhost:8080/sms
func (u *EventService) setSMSValue(request *restful.Request, response *restful.Response) {
	c := appengine.NewContext(request.Request)

	Smsobj := Sms{}
	request.ReadEntity(&Smsobj)

	err := request.ReadEntity(&Smsobj)
	if err == nil {
		item := &memcache.Item{
			Key:    "sms",
			Object: &Smsobj,
		}

		err = memcache.Gob.Set(c, item)

		if err != nil {
			response.WriteError(http.StatusInternalServerError, err)
			return
		}
		response.WriteHeader(http.StatusCreated)
		response.WriteEntity(Smsobj)
	} else {
		response.WriteError(http.StatusInternalServerError, err)
	}

}

// GET http://localhost:8080/events
func (u *EventService) getAllEvents(request *restful.Request, response *restful.Response) {

	c := appengine.NewContext(request.Request)
	q := datastore.NewQuery("Activity").Ancestor(eventKey(c)).Order("-Date").Limit(100)

	events := make([]Activity, 0, 100)

log.Printf("SID VALUE %s", twilioAccountSid )

	if _, err := q.GetAll(c, &events); err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Printf("GOT HERE!")
	response.WriteEntity(events)

}

// POST http://localhost:8080/events
func (u *EventService) createEvent(request *restful.Request, response *restful.Response) {
	c := appengine.NewContext(request.Request)

	event := Activity{}
	request.ReadEntity(&event)

	event.Date = time.Now()

	key := datastore.NewIncompleteKey(c, "Activity", eventKey(c))
	_, err2 := datastore.Put(c, key, &event)

	if err2 == nil {
		response.WriteHeader(http.StatusCreated)
		response.WriteEntity(event)

		sms := new(Sms)
		_, err := memcache.Gob.Get(c, "sms", &sms)
		if err == nil {
			if sms.Value == "1" {
				log.Println("SEND TEXT!")
				postToTwilio(request.Request)
			} else {
				log.Println("DONT SEND TEXT!")
			}
		}

	} else {
		response.WriteError(http.StatusInternalServerError, err2)
	}
}

func getGaeURL() string {
	if appengine.IsDevAppServer() {
		return "http://localhost:8080"
	} else {
		/**
		 * Include your URL on App Engine here.
		 * I found no way to get AppID without appengine.Context and this always
		 * based on a http.Request.
		 */
		return "http://localtone-gae.appspot.com"
	}
}

type jsonobject struct {
	Twilio struct {
		Sid   string
		Token string
		From string
		To string
	}
}

func init() {

	u := EventService{}
	u.Register()

	file, e := ioutil.ReadFile("./config.json")
	if e != nil {
		fmt.Printf("File error: %v\n", e)
		os.Exit(1)
	}

	var jsontype jsonobject
	json.Unmarshal(file, &jsontype)

	twilioAccountSid = jsontype.Twilio.Sid
	twilioAuthToken = jsontype.Twilio.Token
	twilioFrom = jsontype.Twilio.From
	twilioTo = jsontype.Twilio.To
	
	//fmt.Printf("Results: %v\n", jsontype)
	
	log.Printf("Results: %v\n", jsontype)

	// Install the Swagger Service which provides a nice Web UI on your REST API
	// You need to download the Swagger HTML5 assets and change the FilePath location in the config below.
	// Open <your_app_id>.appspot.com/apidocs and enter http://<your_app_id>.appspot.com/apidocs.json in the api input field.
	config := swagger.Config{
		WebServices:    restful.RegisteredWebServices(), // you control what services are visible
		WebServicesUrl: getGaeURL(),
		ApiPath:        "/apidocs.json",

		// Optionally, specifiy where the UI is located
		SwaggerPath: "/apidocs/",
		// GAE support static content which is configured in your app.yaml.
		// This example expect the swagger-ui in static/swagger so you should place it there :)
		SwaggerFilePath: "static/swagger",
	}
	swagger.InstallSwaggerService(config)
}
