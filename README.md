SOLUTION

My solution is a Go script that uses the "compress/gzip" library to stream the large .gz input file to avoid having to
unzip it all at once. Since all we care about is the array of reporting structure objects, it reads through the JSON looking 
for the beginning of that array, then loops through each element and unmarshals it into a ReportingStructure struct. Once unmarshaled,
each reporting structure object is examined and each plan name inside of it is examined to see if it contains either "PPO NY" or "NY PPO", which
after examining some of the data and looking at a couple of files linked to New York businesses on the Anthem portal, seems to
identify New York PPO plans. Once an appropriate plan name is found, the url is saved to a map to avoid outputting duplicates. When the entire
file has been processed, the urls are collected and output to a file called "urls.json".

RUNNING THE SCRIPT AND DISCUSSION

To run the program, all one needs is Go 1.21. The command is `go run main.go <input file path>`, which takes the path to the input file as the single arg. The entire
script takes a bit under 4 minutes to run and took a little under 2 hours to write. The main tradeoffs I made were to skip further optimizations and code structure improvements
as well as ensuring that all New York PPO plans were indeed captured by way of unit testing as well as deeper sleuthing.
