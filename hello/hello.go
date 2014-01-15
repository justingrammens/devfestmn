package main

import (
		"appengine"
		"appengine/datastore"
		"appengine/urlfetch"
        "github.com/emicklei/go-restful"
        "github.com/emicklei/go-restful/swagger"
        "net/http"
		"log"
		"time"
		"github.com/kevinburke/twilio-go/twilio"
		"appengine/memcache"
		"fmt"
)


type Activity struct {
	DeviceID string
	Kind string
	Date    time.Time
}

type Sms struct {
	Value string
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
		
		
		ws2 := new(restful.WebService)
		
        ws2.
                Path("/sms").
                Consumes(restful.MIME_XML, restful.MIME_JSON).
                Produces(restful.MIME_JSON, restful.MIME_XML) // you can specify this per route as well

        ws2.Route(ws.GET("").To(u.getSMSValue).
                // docs
                Doc("get value").
				Reads(Sms{})) // from the request

        ws2.Route(ws.PUT("/{value}").To(u.setSMSValue).
                // docs
                Doc("create an event").
                Param(ws.PathParameter("value", "identifier of the user").DataType("string")).
				Reads(Sms{})) // from the request
				
        restful.Add(ws2)
		
}


// guestbookKey returns the key used for all guestbook entries.
func eventKey(c appengine.Context) *datastore.Key {
    // The string "default_guestbook" here could be varied to have multiple guestbooks.
    //return datastore.NewKey(c, "Justin", "default_justin", 0, nil)
	return datastore.NewKey(c, "Collection", "default_stuff", 0, nil)
}

func sendText(request *restful.Request) {
	const sid = "AC3ce9c81dec2c21f0664f44a2effb604e"
    const token  = "bcea6eac803c01755daa1968ac10caa5"
	
	client := twilio.CreateClient(sid, token, nil, request.Request)
	msg, err := client.Messages.SendMessage("+16122604503", "+16122088663", "Movement happened!", nil)
	
	log.Printf("MESSAGE IS: %s", msg)
	log.Printf("err is: %s", err)
}

func (u *EventService) getSMSValue(request *restful.Request, response *restful.Response) {
	/*c := appengine.NewContext(request.Request)
	       
	        sms := new(Sms)
	        _, err := memcache.Gob.Get(c, "sms", &sms)
	        if err != nil {
	                response.WriteErrorString(http.StatusNotFound, "SMS value not found.")
	        } else {
	                response.WriteEntity(sms)
	        }
	*/
	
	
	c := appengine.NewContext(request.Request)
	client := urlfetch.Client(c)
	resp, err := client.Get("http://www.google.com/")
	if err != nil {
	    http.Error(response, err.Error(), http.StatusInternalServerError)
	    return
	 }
	 
	 	log.Printf("WE HAVE SET THIS!!!! %s", client)
		
	 fmt.Fprintf(response, "HTTP GET returned status %v", resp.Status)
				
}

func (u *EventService) setSMSValue(request *restful.Request, response *restful.Response) {
	c := appengine.NewContext(request.Request)
	
	Smsobj := Sms{Value:request.PathParameter("value")}
	
	err := request.ReadEntity(&Smsobj)
	 if err == nil {
		 item := &memcache.Item{
			 Key:   "sms",
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



func (u *EventService) createEvent(request *restful.Request, response *restful.Response) {
        c := appengine.NewContext(request.Request)
		
	    event := Activity{}
	    request.ReadEntity(&event)
		event.Date = time.Now()	
		
		log.Printf("HERE IS TYPER!!!")
	
		key := datastore.NewIncompleteKey(c, "Activity", eventKey(c))
		_, err2 := datastore.Put(c, key, &event)
  
	    if err2 == nil {
			response.WriteHeader(http.StatusCreated)
	        response.WriteEntity(event)
			
			log.Printf("SAYS THAT IT WAS CREATED????????????")
			
			sms := new (Sms)
			_, err := memcache.Gob.Get(c, "sms", &sms)
			if err == nil {
				if sms.Value == "1" {
					log.Println("SEND TEXT!")
					sendText(request)
				} else {
					log.Println("DONT SEND TEXT!")
				}
			}
			
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