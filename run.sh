docker build -t my-broker . 
docker run -p 9000:9000 -p 5100:5100 my-broker