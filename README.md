# cargoshortener 🩳

A threadsafe URL shortener built with Golang and React. [Demo available here!](https://cargoshortener.herokuapp.com/)

This project uses time hashing to avoid collisions of shortened URL's and goroutines in order to ping and validate the generated URL's are working correctly. 

I built this project to familiarize myself with Go and containerization. I learned a great deal about routing in Go and necessary steps in deploying a web app. I followed [Polyglot](https://www.youtube.com/watch?v=OVBvOuxbpHA) to get a grasp of the basic structure, then made changes to improve.


## Usage

The project is dockerized but to run in development environments: 

To launch the go server: 
`go run main.go`

To launch the client with: 
`npm start `
