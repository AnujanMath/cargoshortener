# cargoshortener 🩳

A threadsafe URL shortener built with Golang and React [DEMO](http://cargoshortener.herokuapp.com/)

This project uses time hashing to avoid collisions of shortened URL's and goroutines in order to ping and validate the generated URL's are working correctly. 

I built this project to familiarize myself with Go and dockerization. I learned a great deal about routing in Go and necessary steps in deploying a web app.


## Usage

The project is dockerized but to run in development environments: 

To launch the go server: 
`go run main.go`

To launch the client with: 
`npm start `


## Development  
Download Deps
`go mod download`  
