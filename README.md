# "Word of Wisdom" TCP-server with protection from DDOS based on Proof of Work

## 1. Description
This project is a solution for some interview question on Golang.

## 2. Getting started
### 2.1 Requirements
+ [Go 1.17+](https://go.dev/dl/) installed (to run tests, start server or client without Docker)
+ [Docker](https://docs.docker.com/engine/install/) installed (to run docker-compose)

### 2.2 Start server and client by docker-compose:
```
make start
```

### 2.3 Start only server:
```
make start-server
```

### 2.4 Start only client:
```
make start-client
```

### 2.5 Launch tests:
```
make test
```

### 2.6 Example
[![asciicast](https://asciinema.org/a/Fe0tHCl5x4Arzb6WCm5HZsmrp.svg)](https://asciinema.org/a/Fe0tHCl5x4Arzb6WCm5HZsmrp)

## 3. Problem description
Design and implement “Word of Wisdom” tcp server. 
TCP server should be protected from DDOS attacks with the [Proof of Work](https://en.wikipedia.org/wiki/Proof_of_work), 
the challenge-response protocol should be used.  
The choice of the PoW algorithm should be explained.  
After Prof Of Work verification, server should send one of the quotes from “word of wisdom” book or any other collection of the quotes.  
Docker file should be provided both for the server and for the client that solves the PoW challenge.

## 4. Protocol definition
This solution uses TCP-based protocol. Each messsage uses delimiter \n and consists of two parts divided by symbol |:
+ header - integer number to indicate, which type of request was sent (analogue of URL in HTTP-protocol)
+ payload - optional string, that also could be json of some struct (depends on type of request)

### 4.1 Types of requests
Solution supports 5 types of requests, switching by header:
+ 0 - Quit - signal to other side to close connection
+ 1 - RequestChallenge - from client to server - request new challenge from server
+ 2 - ResponseChallenge - from server to client - message with challenge for client
+ 3 - RequestResource - from client to server - message with solved challenge
+ 4 - ResponseResource - from server to client - message with useful info is solution is correct, or with error if not

### 4.2 Examples of protocol message
Here i provide examples for all types of requests:
+ 0 Quit: ```0|```
+ 1 RequestChallenge: ```1|```
+ 2 ResponseChallenge: ```2|{"Version":1,"ZerosCount":3,"Date":1647252023,"Resource":"client1","Rand":"OTgwODE=","Counter":0}```
+ 3 RequestResource: ```3|{"Version":1,"ZerosCount":3,"Date":1647252023,"Resource":"client1","Rand":"OTgwODE=","Counter":126231}```
+ 4 ResponseResource ```4|some cool quote from Word of Wisdom```

## 5. Proof of Work
Idea of Proof of Work for DDOS protection is that client, which wants to get some resource from server, 
should firstly solve some challenge from server. 
This challenge should require more computational work on client side and verification of challenge's solution - much less on the server side.

### 5.1 Selection of an algorithm
There is some different algorithms of Proof Work. 
I compared next three algorithms as more understandable and having most extensive documentation:
+ [Merkle tree](https://en.wikipedia.org/wiki/Merkle_tree)
+ [Hashcash](https://en.wikipedia.org/wiki/Hashcash)
+ [Guided tour puzzle](https://en.wikipedia.org/wiki/Guided_tour_puzzle_protocol)

After comparison, I chose Hashcash. Other algorithms have next disadvantages:
+ In Merkle tree server should do too much work to validate client's solution. For tree consists of 4 leaves and 3 depth server will spend 3 hash calculations.
+ In guided tour puzzle client should regularly request server about next parts of guide, that complicates logic of protocol.

Hashcash, instead has next advantages:
+ simplicity of implementation
+ lots of documentation and articles with description
+ simplicity of validation on server side
+ possibility to dynamically manage complexity for client by changing required leading zeros count

Of course Hashcash also has disadvantages like:

1. Compute time depends on power of client's machine. 
For example, very weak clients possibly could not solve challenge, or too powerful computers could implement DDOS-attackls.
But complexity of challenge could be dynamically solved by changing of required zeros could from server.
2. Pre-computing challenges in advance before DDOS-attack. 
Some clients could parse protocol and compute many challenges to apply all of it in one moment.
It could be solved by additional validation of hashcash's params on server. 
For example, on creating challenge server could save **rand** value to Redis cache and check it's existence on verify step.

But all of those disadvantages could be solved in real production environment. 

## 6. Structure of the project
Project structure implements [Go-layout](https://github.com/golang-standards/project-layout) pattern.
Existing directories:
+ cmd/client - main.go for client
+ cmd/server - main.go for server
+ config - config files for both server and client
+ internal/client - all logic of client
+ internal/server - all logic of server
+ internal/pkg/config - logic of parsing config file and env variable
+ internal/pkg/pow - logic of chosen PoW algorithm (Hashcash)
+ internal/pkg/protocol - constants, models and parsing logic for implemented TCP-protocol
+ internal/pkg/clock - internal struct for getting current time (and be easily mocked in tests)

## 7. Ways to improve
Of course, every project could be improved. This project also has some ways to improve:
+ add dynamic management of Hashcash complexity based on server's overload 
(to improve DDOS protection)
+ move array of quotes to SQLite or PostgreSQL database (to get closer to real production applications)
+ add integration with Redis cache to additional verification of client's challenges (check **rand** value existence)
+ add integration tests for simulate DDOS attack - spawn more client instances