

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



NEWS SOURCES

FOR PROCESSING
- add check that there are at least 10 sentences

npr.org
newsweek (https://www.newsweek.com/rss)
associated press (https://apnews.com/index.rss)
reason (https://reason.com/latest/feed/)
CBS (https://www.cbsnews.com/latest/rss/main)



bbc.com (https://feeds.bbci.co.uk/news/world/rss.xml)
reuters (https://www.reutersagency.com/feed/?best-sectors=economy&post_type=best)
- "story-collection__list-item__j4SQe" (just go through all of these and get summary, publish, link)
WILL PROBABLY HAVE TO WEBSCRAPE FRONT PAGE

Why this selection?
They're pretty trustworthy sources (though with various bias) they're free to access. no paywalls. And, most importantly, I like them <3

How do you get article content?
Webscraping, each source I have a slightly different script for getting their articles.


Article processing

The article is then
- checked if more than 10 sentences 
- given a priority score based on its words
- summarized


types of news:
- disaster
- policy
- economy
- war
- research

HIGH PRIORITY (+6)
- Ukraine, Kyiv, Russia, Gaza, Israel, Palestine, Zionist, Genocide, War, Ethnic cleansing

Disaster (+4)
- Tornado, Disaster, Hurriance, Ruin, Wreckage, 

Priority (+3)
- Union, Inflation, Economics, Congress, Bill, Law, AI, Supreme Court, Protest
- Biden, Trump, Xi jinping, AOC, Racism, Pandemic, LGBT, Minimum Wage, Trade, Military
- European Union, Middle East, Migrant, immigrant, Abortion, crisis, strike, administration,
- climate change, artificial intelligence

Low Priority (+1)
- Study, Inflation, Journal, Research, Scientific, Net Neutrality, unemployment, doctor, AI

DEPRIORITY (-5)
- editorial or opinion pieces
- celebrities, royalty
    - taylor swift, kim kardashian, drake, kanye (ye), britney spears, elon musk, mark zuckerberg
    - King Charles, Prince, Queen, Buckingham Palace, 
Title Depriority (-50, if in title)
- photo, video
- so we're not summarizing non-articles

I don't care, I don't want to hear about them.

If a word is in the title, then it gets 10x the points (either positive or negative)
If priority is negative, it won't show up

Most of this is catered to what I want to see, sorry! This is my bias, although if you think there's a topic thats ellided you think should be there, or if you think a priority is unjustifiably off, then email me and I will consider.  