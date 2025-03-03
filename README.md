A simple project using Go to create an web-server timer. The objective is creating a timer for each post request, that when finished, the server will make an request to the desirable destination.
Each timer can be reseted by making another request with the same timer id.
The request needs to be a POST to "\reset-timer" request with a body like:
```
{
  id: [String]
  time: [Number]
  request: [Object]
}

```
