Build the container;
```
$ docker build . -t distil-api:latest
```

Run container;
```
$ docker run -d --restart always -p 127.0.0.1:3000:3000 -d distil-api:latest
```