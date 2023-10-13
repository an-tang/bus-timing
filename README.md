## Bus-timing service
  A service to get bus information, like: bus-line, bus top and estimate arrival-time/distance from bus to bus-stop.
#### Candidate notes:
- How to run app:
1. Run app with docker-compose:
	run: `docker-compose up`
2. Run app with `go run`:
    run: `go run main.go`

- API documents: https://documenter.getpostman.com/view/7947267/2s9YR3dbXX#43175143-d380-46f2-b921-30f877b1509a
#### Approach:
1. Each bus line has their own journey, and all of positions they pass over will be called paths.
2. Bus stop stay at a position on the bus line's path.
3. Bus of the bus line will go from the start to the end of bus lines path, and they will pass bus stop in their way.
4. At any one time during the day, some (or none) buses move along bus routes, providing information such as current position, traffic density, and so on. We have an ERD below based on them
![alt text](https://github.com/an-tang/bus-timing/blob/main/images/ERD.png?raw=true)

##### How to estimate arrival time:
Suppose bus line from A to G, it has 3 running buses on the way and bus stop X like image below:
![alt text](https://github.com/an-tang/bus-timing/blob/main/images/map.png?raw=true) 
- Firstly, we need to find the nearest bus in the bus line with bus stop:
    - Distance between two points depending on longitudes and latitudes, as calculated by the Haversine_formula, refer: https://en.wikipedia.org/wiki/Haversine_formula
    - For e.g, in the above image, by using the Haversine_formula, we determine that the nearest bus stop is `bus_3`
- Secondly, find the bus stop in which path of bus line
    - For e.g, in the above image, bus stop is in path `EF`
- Next, find the bus position (found in step 1) in which path of bus line
    - For e.g, in the above image, bus is in path `CD`
- Calculate distance between bus to bus stop (follow Haversine formula above):
    - In the above image, distance from bus to bus stop = `bus_3`_D + `DE`+ `EX`
- Time arrival = `distance/speed`

*How to check a position in which part*:
- Suppose we have 2 points A, B with coordinates A(lat1, lng1), B(lat2, lng2), we need to check point X(latX,lngX) is between A & B or not.
   - Firstly, we check A, B, X are straight line or not
       - If yes, if AX + XB = AB => X between A and B
       - Else we need to check angel AXB by vectors, 
      ![alt text](https://github.com/an-tang/bus-timing/blob/main/images/vector.png?raw=true) 
       from this we can find the current path of a position.

#### Pros and Cons in my approach
1. Pros:
 - This is a simple way to simulate the requirements.
2. Cons:
- The distance is not actually accurate, just assumptions that bus will go straight in bus line from the start to the end.
- Hard code speed based on crowd level of running bus position.
- Ignore factors that can affect the estimation, like: traffic lights, crowd level, speed, etc