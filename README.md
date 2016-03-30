      ____  ___      _ _ _ _       _   
     / ___|/ _ \  __| (_) | | ___ | |_ 
    | |  _| | | |/ _` | | | |/ _ \| __|
    | |_| | |_| | (_| | | | | (_) | |_ 
     \____|\___/ \__,_|_|_|_|\___/ \__|
                                   

Simple interface to Monit using Go language

####Usage
Type `godillot` from the command line

####Config
Configuration with a yaml file, named "godillot.yaml"

    outputfile: index.html
    servers:
      - server: Watchdog
        url: http://admin:monit@watchdog:2812/_status?format=xml
      - server: Obelix  
        url: http://admin:monit@obelix:2812/_status?format=xml
      - server: Domify  
        url: http://admin:monit@domify:2812/_status?format=xml