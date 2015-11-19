package main

import (
    "fmt"
    "github.com/julienschmidt/httprouter"
    "net/http"
    "encoding/json"
    "strings"
    "io/ioutil"
    "gopkg.in/mgo.v2"
    "gopkg.in/mgo.v2/bson"
    "strconv"
    "github.com/anweiss/uber-api-golang/uber"
    //"bytes"
)

//func hello(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
//    fmt.Fprintf(rw, "Hello, %s!\n", p.ByName("name"))
//}

// MongoLab Auth : mongodb://tjs:password@ds039684.mongolab.com:39684/mongo

type reqObj struct{
Id int
Name string `json:"Name"`
Address string `json:"Address"`
City string `json:"City"`
State string `json:"State"`
Zip string `json:"Zip"`
Coordinates struct{
    Lat float64
    Lng float64
}
}

var id int;
type Responz struct {
    Results []struct {
        AddressComponents []struct {
            LongName  string   `json:"long_name"`
            ShortName string   `json:"short_name"`
            Types     []string `json:"types"`
        } `json:"address_components"`
        FormattedAddress string `json:"formatted_address"`
        Geometry         struct {
            Location struct {
                Lat float64 `json:"lat"`
                Lng float64 `json:"lng"`
            } `json:"location"`
            LocationType string `json:"location_type"`
            Viewport     struct {
                Northeast struct {
                    Lat float64 `json:"lat"`
                    Lng float64 `json:"lng"`
                } `json:"northeast"`
                Southwest struct {
                    Lat float64 `json:"lat"`
                    Lng float64 `json:"lng"`
                } `json:"southwest"`
            } `json:"viewport"`
        } `json:"geometry"`
        PartialMatch bool     `json:"partial_match"`
        PlaceID      string   `json:"place_id"`
        Types        []string `json:"types"`
    } `json:"results"`
    Status string `json:"status"`
}



type resObj struct{
Greeting string
}

func createlocation(rw http.ResponseWriter, req *http.Request, p httprouter.Params) {
    id=id+1;


    decoder := json.NewDecoder(req.Body)
    var t reqObj 
    t.Id = id; 
    err := decoder.Decode(&t)
    if err != nil {
        fmt.Println("Error")
    }


    //lstring := strings.Split(t.Loc," ");
    st:=strings.Join(strings.Split(t.Address," "),"+");
    fmt.Println(st);
    constr := []string {strings.Join(strings.Split(t.Address," "),"+"),strings.Join(strings.Split(t.City," "),"+"),t.State}
    lstringplus := strings.Join(constr,"+")
    locstr := []string{"http://maps.google.com/maps/api/geocode/json?address=",lstringplus}
    //fmt.Println(strings.Join(locstr,""));
    resp, err := http.Get(strings.Join(locstr,""))
    //fmt.Println(resp);
    
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
       fmt.Println("Error: Wrong address");
     }
     var data Responz
    err = json.Unmarshal(body, &data)
    fmt.Println(data.Status)
    // n := bytes.IndexByte(body, 0)
    // stz := string(body[:n])
    // fmt.Println(stz);

 //    s := []string{"Hello, ",t.Name}
 //    g := resObj{strings.Join(s,"")}
    t.Coordinates.Lat=data.Results[0].Geometry.Location.Lat;
    t.Coordinates.Lng=data.Results[0].Geometry.Location.Lng;


//Mongo Persistence

 conn, err := mgo.Dial("mongodb://localhost:27017/mongo")

    // Check if connection error, is mongo running?
    if err != nil {
        panic(err)
    }
    defer conn.Close();

conn.SetMode(mgo.Monotonic,true);
c:=conn.DB("test").C("details");
err = c.Insert(t);

//Response
    js,err := json.Marshal(t)
    if err != nil{
	   fmt.Println("Error")
	   return
	}
    rw.Header().Set("Content-Type","application/json")
    rw.Write(js)
}

func getloc(rw http.ResponseWriter, req *http.Request, p httprouter.Params){
fmt.Println(p.ByName("locid"));
id ,err1:= strconv.Atoi(p.ByName("locid"))
if err1 != nil {
        panic(err1)
    }
 conn, err := mgo.Dial("mongodb://localhost:27017/mongo")

    // Check if connection error, is mongo running?
    if err != nil {
        panic(err)
    }
    defer conn.Close();

conn.SetMode(mgo.Monotonic,true);
c:=conn.DB("test").C("details");
result:=reqObj{}
err = c.Find(bson.M{"id":id}).One(&result)
if err != nil {
                fmt.Println(err)
        }

        //fmt.Println("Name:", result.Name)

        //Response
        js,err := json.Marshal(result)
    if err != nil{
       fmt.Println("Error")
       return
    }
    rw.Header().Set("Content-Type","application/json")
    rw.Write(js)
}

type modReqObj struct{
    Address string `json:"address"`
    City string `json:"city"`
    State string `json:"state"`
    Zip string `json:"zip"`
}

func updateloc(rw http.ResponseWriter, req *http.Request, p httprouter.Params){
    //fmt.Println("sdf");
 id ,err1:= strconv.Atoi(p.ByName("locid"))
 //fmt.Println(id);
 if err1 != nil {
         panic(err1)
     }
  conn, err := mgo.Dial("mongodb://localhost:27017/mongo")

//     // Check if connection error, is mongo running?
     if err != nil {
         panic(err)
     }
     defer conn.Close();

conn.SetMode(mgo.Monotonic,true);
 c:=conn.DB("test").C("details");


     decoder := json.NewDecoder(req.Body)
     var t modReqObj  
     err = decoder.Decode(&t)
     if err != nil {
         fmt.Println("Error")
     }


     colQuerier := bson.M{"id": id}
     change := bson.M{"$set": bson.M{"address": t.Address, "city":t.City,"state":t.State,"zip":t.Zip}}
     err = c.Update(colQuerier, change)
     if err != nil {
         panic(err)
     }

}

func deleteloc(rw http.ResponseWriter, req *http.Request, p httprouter.Params){
     id ,err1:= strconv.Atoi(p.ByName("locid"))
 //fmt.Println(id);
 if err1 != nil {
         panic(err1)
     }
  conn, err := mgo.Dial("mongodb://localhost:27017/mongo")
  conn.SetMode(mgo.Monotonic,true);
c:=conn.DB("test").C("details");

//     // Check if connection error, is mongo running?
     if err != nil {
         panic(err)
     }
     defer conn.Close();
     err=c.Remove(bson.M{"id":id})
     if err != nil { fmt.Printf("Could not find kitten %s to delete", id)}
    rw.WriteHeader(http.StatusNoContent)
}

type userUber struct {
    LocationIds            []string `json:"location_ids"`
    StartingFromLocationID string   `json:"starting_from_location_id"`
}

func plantrip(rw http.ResponseWriter, req *http.Request, p httprouter.Params){

    decoder := json.NewDecoder(req.Body)
    var uUD userUber 
    err := decoder.Decode(&uUD)
    if err != nil {
        fmt.Println("Error")
    }

        fmt.Println(uUD.StartingFromLocationID);


///UBERRRRRR !!!!
    var options uber.RequestOptions;
    options.ServerToken= "S37TkXJu1TBNbDea22MxgIAjoM1C__fJ3r6vbQ-5";
    options.ClientId= "5-BNiHDpt1CZvQoWd2G2vV2GSvSnIu2j";
    options.ClientSecret= "P5qyGJI-sJw5m-s2kFHljzg59kccexZ8qkbaL44P";
    options.AppName= "CMPE273-A3";
    options.BaseUrl= "https://sandbox-api.uber.com/v1/";
    

    client := uber.Create(&options);

//Quering for the locations: start and the rest
        sid ,err1:= strconv.Atoi(uUD.StartingFromLocationID)
 //fmt.Println(id);
 if err1 != nil {
         panic(err1)
     }

    conn, err := mgo.Dial("mongodb://localhost:27017/mongo");

    // Check if connection error, is mongo running?
    if err != nil {
        panic(err)
    }
    defer conn.Close();

    conn.SetMode(mgo.Monotonic,true);
    c:=conn.DB("test").C("details");
    result:=reqObj{}
    err = c.Find(bson.M{"id":sid}).One(&result)
    if err != nil {
                fmt.Println(err)
        }

    // distance := []int{};
    // price :=[]float64{};
    index:=0;
    // totalPrice := 0.0;
    totalDistance :=0.0;
    // totalDuration :=0.0;

    for _,ids := range uUD.LocationIds{
    
        lid,err1:= strconv.Atoi(ids)
            //fmt.Println(id);
        if err1 != nil {
            panic(err1)
        }
        

        resultLID:=reqObj{}
        err = c.Find(bson.M{"id":lid}).One(&resultLID)
        if err != nil {
             fmt.Println(err)
        }
        pe := &uber.PriceEstimates{}
        pe.StartLatitude = result.Coordinates.Lat;
        pe.StartLongitude = result.Coordinates.Lng;
        pe.EndLatitude = resultLID.Coordinates.Lat;
        pe.EndLongitude = resultLID.Coordinates.Lng;

        if e := client.Get(pe); e != nil {
            fmt.Println(e);
        }
        fmt.Println(result.Address);
        fmt.Println("to");
        fmt.Println(resultLID.Address);
        totalDistance=totalDistance+pe.Prices[0].Distance;
    for _, price := range pe.Prices {
        fmt.Println(price.DisplayName + ": "+strconv.Itoa(price.LowEstimate) + "; Surge: " + strconv.FormatFloat(price.Distance , 'f', 2, 32))
    }
        index=index+1;
    }
    

     fmt.Println(totalDistance);




    }

func main() {
    mux := httprouter.New()
    //mux.GET("/hello/:name", hello)

    /// Uber
    

    ///
    id=0;
    mux.POST("/locations",createlocation)
    mux.POST("/trips",plantrip)
    mux.GET("/locations/:locid",getloc)
    mux.PUT("/locations/:locid",updateloc)
    mux.DELETE("/locations/:locid",deleteloc)
    server := http.Server{
            Addr:        "0.0.0.0:8083",
            Handler: mux,
    }

    server.ListenAndServe()
}
