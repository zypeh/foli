# foli
Simple Behance Photo View

### Installation
Just download it and go build it. I don't know there is any good dependency tools for golang too.

### How it works ?
This program accepts OS env `API` as API key to fetch images from [Behance](https://www.behance.net/dev/api/endpoints/)

```bash
API=xxx ./main
```

And you will find there is a directory called `images` had created, including images that fetched on programme startup.
