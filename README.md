      ____           _ _ _     _          _   
     / ___| ___   __| (_) |   | |    ___ | |_ 
    | |  _ / _ \ / _` | | |   | |   / _ \| __|
    | |_| | (_) | (_| | | |___| |__| (_) | |_ 
     \____|\___/ \__,_|_|_____|_____\___/ \__|
                                       

Simple interface to Monit for monitoring multiple instances (lile M/Monit), using Go language

####Usage
Type `godillot` from the command line

####Config
Configuration with a yaml file, named "godillot.yaml"

    outputfile: index.html #html file created by godillot with the collected info
    servers: #list of servers
      - server: Watchdog #Name
        url: http://admin:monit@watchdog:2812/_status?format=xml #url (don't forget _status?formal=xml, and to enable monit web server)
      - server: Obelix  
        url: http://admin:monit@obelix:2812/_status?format=xml
      - server: Domify  
        url: http://admin:monit@domify:2812/_status?format=xml

####Example
See godillot in action on [godillot.zeneffy.fr](http://godillot.zeneffy.fr)