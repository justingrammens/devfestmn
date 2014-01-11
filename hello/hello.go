package main

import (
		"html/template"
		"appengine"
		"appengine/datastore"
        "github.com/emicklei/go-restful"
        "github.com/emicklei/go-restful/swagger"
        "net/http"
		"log"
		"time"
)

// This example is functionally the same as ../restful-user-service.go
// but it`s supposed to run on Goole App Engine (GAE)
//
// contributed by ivanhawkes

type Greeting struct {
    Author  string
    Content string
    Date    time.Time
}


type Activity struct {
	DeviceID string
	Kind string
	Date    time.Time
}

/*type ActivityDevice struct {
    DeviceID  string
    Kind string
    Date    time.Time
}*/



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


/*        ws.Route(ws.GET("/{user-id}").To(u.findEvent).
                // docs
                Doc("get a user").
                Param(ws.PathParameter("user-id", "identifier of the user").DataType("string")).
                Writes(User{})) // on the response
*/
        ws.Route(ws.PUT("/{user-id}").To(u.createEvent).
                // docs
                Doc("create an event").
                Param(ws.PathParameter("user-id", "identifier of the user").DataType("string")).
                Reads(Activity{})) // from the request
				/*
        ws.Route(ws.DELETE("/{user-id}").To(u.removeEvent).
                // docs
                Doc("delete a user").
                Param(ws.PathParameter("user-id", "identifier of the user").DataType("string")))
*/
        restful.Add(ws)
}

// GET http://localhost:8080/users/1
//
/*
func (u EventService) findEvent(request *restful.Request, response *restful.Response) {
        c := appengine.NewContext(request.Request)
        id := request.PathParameter("user-id")
        usr := new(User)
		
		
        _, err := memcache.Gob.Get(c, id, &usr)
        if err != nil || len(usr.Id) == 0 {
                response.WriteErrorString(http.StatusNotFound, "User could not be found.")
        } else {
                response.WriteEntity(usr)
        }
}

// PATCH http://localhost:8080/users
// <User><Id>1</Id><Name>Melissa Raspberry</Name></User>
//
func (u *UserService) updateEvent(request *restful.Request, response *restful.Response) {
        c := appengine.NewContext(request.Request)
        evt := new(Event)
        err := request.ReadEntity(&evt)
        if err == nil {
                item := &memcache.Item{
                        Key:    evt.Id,
                        Object: &usr,
                }
                err = memcache.Gob.Set(c, item)
                if err != nil {
                        response.WriteError(http.StatusInternalServerError, err)
                        return
                }
                response.WriteEntity(usr)
        } else {
                response.WriteError(http.StatusInternalServerError, err)
        }
}

// PUT http://localhost:8080/users/1
// <User><Id>1</Id><Name>Melissa</Name></User>
//
*/

// guestbookKey returns the key used for all guestbook entries.
func eventKey(c appengine.Context) *datastore.Key {
    // The string "default_guestbook" here could be varied to have multiple guestbooks.
    //return datastore.NewKey(c, "Justin", "default_justin", 0, nil)
	return datastore.NewKey(c, "Collection", "default_stuff", 0, nil)
}

func (u *EventService) createEvent(request *restful.Request, response *restful.Response) {
        c := appengine.NewContext(request.Request)
		
		
//	    event := new(Greeting) ??? WTF this didn't work? but below did...
	    event := Activity{}
	    request.ReadEntity(&event)	
		
		log.Printf("HERE IS TYPER!!!")
	
		 // persist the event
/*		event := Event{
		        Id: time.Now() ,
		        Type: typer,
		}
*/		
		key := datastore.NewIncompleteKey(c, "Activity", eventKey(c))
		_, err2 := datastore.Put(c, key, &event)
  
	    if err2 == nil {
			response.WriteHeader(http.StatusCreated)
	        response.WriteEntity(event)
	    } else {
	        response.WriteError(http.StatusInternalServerError,err2)
	    }
			
        //usr := User{Id: request.PathParameter("user-id")}
        //err := request.ReadEntity(&usr)
        /*if err == nil {
                item := &memcache.Item{
                        Key:    usr.Id,
                        Object: &usr,
                }
                err = memcache.Gob.Add(c, item)
                if err != nil {
                        response.WriteError(http.StatusInternalServerError, err)
                        return
                }
                response.WriteHeader(http.StatusCreated)
                response.WriteEntity(usr)
        } else {
                response.WriteError(http.StatusInternalServerError, err)
        }*/
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
			
		    /*if err := guestbookTemplate.Execute(response, events); err != nil {
		        http.Error(response, err.Error(), http.StatusInternalServerError)
		    }*/
			
			/*
		    if _, err := q.GetAll(c, &events); err != nil {
		        http.Error(response, err.Error(), http.StatusInternalServerError)
		        return
		    }
			*/
			//log.Printf("HERE ARE EVENTS!!!! %s", events)
			
			log.Printf("GOT HERE!")
			response.WriteEntity(events)
						
			
/*		    if err := guestbookTemplate.Execute(w, greetings); err != nil {
		        http.Error(w, err.Error(), http.StatusInternalServerError)
		    }
			*/
			
}


var guestbookTemplate = template.Must(template.New("book").Parse(guestbookTemplateHTML))

const guestbookTemplateHTML = `
<html>
<body>
    {{range .}}
      {{with .Kind}}
        <p><b>{{.}}</b> wrote:</p>
      {{else}}
        <p>An anonymous person wrote:</p>
      {{end}}
      <pre>{{.Date}}</pre>
    {{end}}
    <form action="/sign" method="post">
      <div><textarea name="content" rows="3" cols="60"></textarea></div>
      <div><input type="submit" value="Sign Guestbook"></div>
    </form>
  </body>
</html>
`


// DELETE http://localhost:8080/users/1
//
/*
func (u *UserService) removeUser(request *restful.Request, response *restful.Response) {
        c := appengine.NewContext(request.Request)
        id := request.PathParameter("user-id")
        err := memcache.Delete(c, id)
        if err != nil {
                response.WriteError(http.StatusInternalServerError, err)
        }
}
*/

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