### Jaeger Tracing Example 

Overview
- jaeger "all-in-one"
    - UI exposed on port 16686 
- go-service
    - http server with Go standard lib
    - calls the py-service
- py-service
    - http server with FastAPI
    - handles requests from go-service
    

Goal(s)
- Simple demo of Jaeger Tracing
- Best practices for Tags and Logs
- Trace requests/spans across services
    
