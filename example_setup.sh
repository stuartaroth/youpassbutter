curl "http://localhost:8080?q=createAuthor&p=Stephen&p=King"
curl "http://localhost:8080?q=createAuthor&p=Richard&p=Matheson"
curl "http://localhost:8080?q=createAuthorWithMiddle&p=William&p=Peter&p=Blatty"

curl "http://localhost:8080?q=createAuthor&p=J.R.R.&p=Tolkien"
curl "http://localhost:8080?q=createAuthor&p=Robert&p=Jordan"
curl "http://localhost:8080?q=createAuthorWithMiddle&p=George&p=R.R.&p=Martin"

curl "http://localhost:8080?q=createAuthor&p=Frank&p=Herbert"
curl "http://localhost:8080?q=createAuthor&p=Isaac&p=Asimov"
curl "http://localhost:8080?q=createAuthorWithMiddle&p=Philip&p=K&p=Dick"

curl "http://localhost:8080?q=createGenre&p=Horror"
curl "http://localhost:8080?q=createGenre&p=Fantasy"
curl "http://localhost:8080?q=createGenre&p=Science+Fiction"

curl "http://localhost:8080?q=createBookWithGenre&p=1&p=1&p=It&p=1986"
curl "http://localhost:8080?q=createBookWithGenre&p=1&p=1&p=Misery&p=1987"
curl "http://localhost:8080?q=createBookWithGenre&p=1&p=1&p=The+Stand&p=1978"

curl "http://localhost:8080?q=createBookWithGenre&p=2&p=1&p=Hell+House&p=1971"
curl "http://localhost:8080?q=createBookWithGenre&p=2&p=1&p=I+Am+Legend&p=1954"

curl "http://localhost:8080?q=createBookWithGenre&p=3&p=1&p=The+Exorcist&p=1971"

curl "http://localhost:8080?q=createBookWithGenre&p=4&p=2&p=The+Fellowship+of+the+Ring&p=1954"
curl "http://localhost:8080?q=createBookWithGenre&p=4&p=2&p=The+Two+Towers&p=1954"
curl "http://localhost:8080?q=createBookWithGenre&p=4&p=2&p=The+Return+of+the+King&p=1955"
curl "http://localhost:8080?q=createBookWithGenre&p=4&p=2&p=The+Hobbit&p=1937"

curl "http://localhost:8080?q=createBookWithGenre&p=5&p=2&p=The+Eye+of+the+World&p=1990"
curl "http://localhost:8080?q=createBookWithGenre&p=5&p=2&p=The+Great+Hunt&p=1990"
curl "http://localhost:8080?q=createBookWithGenre&p=5&p=2&p=The+Dragon+Reborn&p=1991"
curl "http://localhost:8080?q=createBookWithGenre&p=5&p=2&p=The+Shadow+Rising&p=1992"
curl "http://localhost:8080?q=createBookWithGenre&p=5&p=2&p=The+Fires+of+Heaven&p=1993"

curl "http://localhost:8080?q=createBookWithGenre&p=6&p=2&p=A+Game+of+Thrones&p=1996"
curl "http://localhost:8080?q=createBookWithGenre&p=6&p=2&p=A+Clash+of+Kings&p=1998"
curl "http://localhost:8080?q=createBookWithGenre&p=6&p=2&p=A+Storm+of+Swords&p=2000"
curl "http://localhost:8080?q=createBookWithGenre&p=6&p=2&p=A+Feast+for+Crows&p=2005"
curl "http://localhost:8080?q=createBookWithGenre&p=6&p=2&p=A+Dance+with+Dragons&p=2011"

curl "http://localhost:8080?q=createBookWithGenre&p=7&p=3&p=Dune&p=1965"
curl "http://localhost:8080?q=createBookWithGenre&p=7&p=3&p=Destination:+Void&p=1966"
curl "http://localhost:8080?q=createBookWithGenre&p=7&p=3&p=The+White+Plague&p=1982"

curl "http://localhost:8080?q=createBookWithGenre&p=8&p=3&p=Pebble+in+the+sky&p=1950"
curl "http://localhost:8080?q=createBookWithGenre&p=8&p=3&p=Foundation&p=1951"
curl "http://localhost:8080?q=createBookWithGenre&p=8&p=3&p=The+Caves+of+Steel&p=1953"

curl "http://localhost:8080?q=createBookWithGenre&p=9&p=3&p=Do+Androids+Dream+Of+Electric+Sheep&p=1968"
curl "http://localhost:8080?q=createBookWithGenre&p=9&p=3&p=Ubik&p=1969"
curl "http://localhost:8080?q=createBookWithGenre&p=9&p=3&p=Radio+Free+Albemuth&p=1968"

curl "http://localhost:8080/?q=getAuthors&p=100&p=0"
curl "http://localhost:8080/?q=getBooks&p=100&p=0"
curl "http://localhost:8080/?q=searchBooksByGenre&p=fantasy&p=100&p=0"
curl "http://localhost:8080/?q=searchBooksByTitle&p=the&p=100&p=0"
curl "http://localhost:8080/?q=getBooksByAuthorId&p=6&p=100&p=0"
