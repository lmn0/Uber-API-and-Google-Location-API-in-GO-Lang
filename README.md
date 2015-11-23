~~~CURL COMMANDS~~~

INSERTING

curl -H "Content-Type: application/json" -X POST -d '{"Name":"John Smit","Address":"950 Mason St","City":"San Francisco","State":"CA","Zip":"94108"}' http://localhost:8083/locations

curl -H "Content-Type: application/json" -X POST -d '{"Name":"John Smit","Address":"Golden Gate Bridge","City":"San Francisco","State":"CA","Zip":"94129"}' http://localhost:8083/locations

curl -H "Content-Type: application/json" -X POST -d '{"Name":"John Smit","Address":"Beach St. at the Embarcadero","City":"San Francisco","State":"CA","Zip":"94133"}' http://localhost:8083/locations

curl -H "Content-Type: application/json" -X POST -d '{"Name":"John Smit","Address":"Golden Gate Park","City":"San Francisco","State":"CA","Zip":"94129"}' http://localhost:8083/locations

curl -H "Content-Type: application/json" -X POST -d '{"Name":"John Smit","Address":"501 Twin Peaks Blvd","City":"San Francisco","State":"CA","Zip":"94114"}' http://localhost:8083/locations


TRIPS

curl -H "Content-Type: application/json" -X POST -d '{"starting_from_location_id": "1",
    "location_ids" : [ "2", "3", "4", "5", "1" ] }' http://localhost:8083/trips


