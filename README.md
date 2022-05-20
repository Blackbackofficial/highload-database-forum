# Highload database forum

## Swagger

[Documentation API](https://github.com/mailcourses/technopark-dbms-forum)

## Start
```` javascript
docker build -t highloadforum . && docker run -p 5000:5000 -p 9432:5432 highloadforum
````

## Test

[Test program](https://github.com/mailcourses/technopark-dbms-forum)


## Balancing

```
Running 15m test @ http://5.188.141.200:80/api/service/status
  8 threads and 100 connections
    Thread Stats   Avg      Stdev     Max   +/- Stdev
    Latency    14.89ms   15.26ms 592.24ms   96.99%
    Req/Sec     0.89k   147.01     1.21k    75.79%
  Latency Distribution
     50%   12.23ms
     75%   15.25ms
     90%   20.39ms
     99%   63.47ms
  2461176 requests in 5.80m, 347.38MB read
  Socket errors: connect 0, read 92, write 0, timeout 0
Requests/sec:   7066.69
Transfer/sec:      1.00MB
```