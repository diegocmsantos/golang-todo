# golang-todo

Code for this [medium tutorial](https://medium.com/better-programming/build-a-simple-todolist-app-in-golang-82297ec25c7d) about how to build a simple Todo application using Golang and MySql

Cloning the project
```shellscript
$ git clone https://github.com/diegocmsantos/golang-todo.git
```

Creating MySql docker container
```shellscript
$ docker run -d -p 3306:3306 --name mysql -e MYSQL_ROOT_PASSWORD=root mysql
```

Creating a 'todolist' database
```shellscript
$ docker exec -it mysql mysql -uroot -proot -e 'CREATE DATABASE todolist'
```

Building and running the app
```shellscript
$ go build && ./todolist-mysql-go
```

Open the index.html file into your browser and you're going to see the app running.
