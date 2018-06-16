#!/bin/sh

pm2 stop ehelp
rm ehelp.exe
go build -o ehelp.exe
pm2 start ehelp.exe

while true; do
    
    if [[ $(git pull origin master) == *up-to-date* ]]; 
    then
        echo "no change"
    else
        echo "detect changes"
        sleep 2s
        echo "stop ehelp"
		pm2 stop ehelp
		rm ehelp.exe
        go build -o ehelp.exe
        pm2 restart ehelp
    fi

    echo "sleep 30s"

    sleep 30s        

done
update sh tren remote
