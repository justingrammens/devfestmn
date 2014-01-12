package main

import (
		"appengine"
		"appengine/datastore"
        "github.com/emicklei/go-restful"
        "github.com/emicklei/go-restful/swagger"
        "net/http"
		"log"
		"time"
		"github.com/kevinburke/twilio-go/twilio"
)


type Activity struct {
	DeviceID string
	Kind string
	Date    time.Time
}

type EventService struct {
        // normally one would use DAO (data access object)
        // but in this example we simple use memcache.
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


        ws.Route(ws.PUT("/{user-id}").To(u.createEvent).
                // docs
                Doc("create an event").
                Param(ws.PathParameter("user-id", "identifier of the user").DataType("string")).
                Reads(Activity{})) // from the request
				
        restful.Add(ws)
}


// guestbookKey returns the key used for all guestbook entries.
func eventKey(c appengine.Context) *datastore.Key {
    // The string "default_guestbook" here could be varied to have multiple guestbooks.
    //return datastore.NewKey(c, "Justin", "default_justin", 0, nil)
	return datastore.NewKey(c, "Collection", "default_stuff", 0, nil)
}

func sendText() {
	const sid = "AC3ce9c81dec2c21f0664f44a2effb604e"
    const token  = "bcea6eac803c01755daa1968ac10caa5"
	
	client := twilio.CreateClient(sid, token, nil)
	msg, err := client.Messages.SendMessage("+16122604503", "+16122088663", "Movement happened!", nil)
	
	log.Printf("MESSAGE IS: %s", msg)
	log.Printf("err is: %s", err)
}

func (u *EventService) createEvent(request *restful.Request, response *restful.Response) {
        c := appengine.NewContext(request.Request)
		
	    event := Activity{}
	    request.ReadEntity(&event)	
		
		log.Printf("HERE IS TYPER!!!")
	
		key := datastore.NewIncompleteKey(c, "Activity", eventKey(c))
		_, err2 := datastore.Put(c, key, &event)
  
	    if err2 == nil {
			response.WriteHeader(http.StatusCreated)
	        response.WriteEntity(event)
			
			sendText()
			
	    } else {
	        response.WriteError(http.StatusInternalServerError,err2)
	    }
	}


// GET http://localhost:8080/users
//
func (u *EventService) getAllEvents(request *restful.Request, response *restful.Response) {
	
        //response.WriteEntity(UserList{[]User{User{"42", "Gandalf"}, User{"3.14", "Pi"}}})
		
		c := appengine.NewContext(request.Request)
		
		//q := datastore.NewQuery("Event").Ancestor(eventKey(c)).Order("-Date").Limit(10)
		q := datastore.NewQuery("Activity").Ancestor(eventKey(c)).Order("-Date").Limit(100)
		
		events := make([]Activity, 0, 100)
		
				//log.Printf("HERE ARE EVENTS!!!! %s", response)
			

		if _, err := q.GetAll(c, &events); err != nil {
		        http.Error(response, err.Error(), http.StatusInternalServerError)
		        return
		    }
					
			log.Printf("GOT HERE!")
			response.WriteEntity(events)
									
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

func init() {
        u := EventService{}
        u.Register()

        // Optionally, you can install the Swagger Service which provides a nice Web UI on your REST API
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
                SwaggerFilePath: "static/swagger"}
        swagger.InstallSwaggerService(config)
}