Steps to run this example:

1) Download the Go language (I'm assuming you've already done that).

2) You will need to download Postgres 10.1 for the database.
https://www.postgresql.org/

3) Configure the webserver.go script for your specific Postgres database connextion parameters.  This is in the "const" block starting at line 15.  This is 
optional, and should only be done if your Psotgres server cannot be connected via localhost (127.0.0.1).

4) from the Go command like, download the appropriate libaries (or packages) to the directory you're using in your GOPATH environment variable:
    $ go get github.com/lib/pq
    $ go get github.com/gorilla/mux
    $ go get github.com/nu7hatch/gouuid

5) In the directory that is defined in your GOPATH environment variable, you will need to run the database initialization scripts.  These exist in the sql directory 
of the GOPATH folder.  Using the Postgres administrative tool, or the command-line utility, These exist in the sql directory of the GOPATH folder.  They will create 
a database in your Protgres sql server called stockxtest.  They will also create the tables shoesizes and shoenames.  You need to follow the following procedure:

    1) Run the createdb.sql script.
    2) Connect to the stockxtest DB on the Protgres server.
    3) Run the createtables.sql script.

6) Using Go run the webserver.go script from the command line.

  $ go run webserver.go

Here are the commands I used for updating the database and reading back the average.  As per your instructions, I got the values you specified - with one exception.
The second step - adding a "2" to the list of numbers - gave me 2.53333333333332.  I'm assuming this is a database rounding issue as I used the Ave() function.


To add a size for a manufacturer:
http://localhost:8000/append?name=Yeezy&size=1

To get a current manufacturer's TrueToSize value:
http://localhost:8000/truetosize/Yeezy


I designed these as REST-ish services, the problem being that the append is not really covered clearly with a lot of the documentation I've researched.  Plus, 
I wanted to make the test as brainless to run as I could.  So...the function has its own identifier, they are both implemented as GETs, and the update uses
query parameters.  Normally you'd probably throw XML of JSON in the body of the data, but that would require making another utility, and this is only a small
technical test.

Let me know if you need anything else.