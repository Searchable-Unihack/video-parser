# video-parser

Takes a collection of .png images, and connects to Azure to retrieve the text contained within.

To use, specify your key via the environment variable SUBSCRIPTION_KEY.

For example:

SUBSCRIPTION_KEY=<key> ./video-parser *.png
