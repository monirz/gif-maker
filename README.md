# gif-maker 

## A service for creating gif image from multiple jpeg, jpg images


### Ussage

First run the server and upload the files using a `HTTP POST` request.
&nbsp;

`curl -X POST -F 'delay=1' -F 'file[]=@/path/to/pictures/img1.jpeg' -F 'file[]=@/path/to/pictures/img2.jpeg' http://localhost:8090/create`

When you get the successful response, browse to [http://localhost:8090/file](http://localhost:8090/file) 


**TODO** 
- Add support for configurable image width and height, currently it takes first image's width height as default size if all the provided images are not in same width and height. 

   
