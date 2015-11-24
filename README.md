This simple go program will scrape all article from 
[liputan6](http://www.liputan6.com).

## Setup & Usage
Make sure you setup [Go workspace](https://golang.org/doc/code.html) and install [godep](https://github.com/tools/godep)

    git clone https://github.com/pyk/liputan6-scraper.git
    cd liputan6-scraper
    godep go build

Usage
    
    ./liputan6-scraper -tag="banjir-jakarta" -maxpage=10

it will fetch all articles in the first 10 pages of [Banjir Jakarta](http://www.liputan6.com/tag/banjir-jakarta) tag.