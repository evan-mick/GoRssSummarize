

IDEALLY
go server manages whats displayed as well
only js stuff will be for animations n what not


Structure

Main file sets up go routines

2 parts of application
- one checks for API calls, when a user wants some articles, it sends a list with summaries
- another runs twice a day, it goes to selected websites, 
    webscrapes articles, summarizes them, then stores them on a database to be accessed online


https://cpanel.infinityfree.com/panel/indexpl.php?option=mysql&ttt=-1439843499561137152
thanks infinity free for hosting <3



Templ renders html, puts it in static folder
