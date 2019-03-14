# gif-maker 

## A service for creating gif image from multiple jpeg, jpg images


### Ussage

`$ go mod download`<br/>
`$ go run main.go`

**With Docker** 

`$ docker docker build -t gif-maker . `<br/>
`$ docker run -p 8090:8090  gif-maker` 


After running the server upload the files using a `HTTP POST` request.
&nbsp;

`curl -X POST -F 'delay=1' -F 'file[]=@/path/to/pictures/img1.jpeg' -F 'file[]=@/path/to/pictures/img2.jpeg' http://localhost:8090/create`

Or you can browse and upload it through your web browser `http://localhost:8090`,    

When you get the successful response, browse to [http://localhost:8090/file](http://localhost:8090/file) 

### See the live demo [Here](https://cryptic-gorge-21126.herokuapp.com) 


**TODO** 
- Add support for configurable image width and height, currently it takes first image's width height as default size if all the provided images are not in same width and height. 

   
