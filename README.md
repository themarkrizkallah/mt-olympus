# exchange
This is my implementation of an exchange platform powered by a performant matching engine.

## Getting Started

These instructions will get you a copy of the project up & running on your local machine.

### Prerequisites

You need to have [Docker](https://docs.docker.com/install/) & [docker-compose](https://docs.docker.com/compose/install/) installed.

You need to have *protoc* installed. Instructions can be found [here](https://github.com/golang/protobuf) and [here](https://google.github.io/proto-lens/installing-protoc.html).

### Installation
Clone the repo and run the build script.
```
$ git clone https://github.com/themarkrizkallah/exchange.git
$ ./build.sh
```

### Running the project
```
$ docker-compose up --build
```

Once the containers are up & running, you can make API calls to apollo on [localhost:8080](http://localhost:8080/).

### How to Use
Refer to this [README](./apollo/README.md) for all the available endpoints & the API Spec.


## Comments
The project is in its early stages. Expect a lot of the code & documentation to change.

