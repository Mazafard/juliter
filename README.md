# Juliter

Check tweet of Users and calculate words they used.

# Requirement

- [golang +1.8](https://getgb.io/)
- [gb as A project based build tool for the Go programming language.](https://golang.org/) 


# installation

Installation as much as drink a glass of water:


    git clone git@github.com:Mazafard/juliter.git juliter
    cd juliter
    export CONSUMER_KEY=Mazafard
    export CONSUMER_SECRET= asal
    gb vendor restore
    gb build server
    ./bin/server -u Mazafard,asal,shib,baam


This won't produce the expected output either, but, you will now have a
terms table fully populated where you will be able to play with
twitter's term counts:

    cd $HOME
    sqlite3 .juli-db
    sqlite> SELECT a.term, a.count (*) b.count FROM tweet_terms a INNER
      JOIN tweet_terms b ON a.term = b.term AND a.username <> b.username
      ORDER by a.count*b.count desc;
