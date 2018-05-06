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

### How to use ?

After it complete fetching images from Behance, you may dump all the fetched images by accesing the `/` root route.

```
GET localhost:8080/

[
    {
        "id": 1,
        "title": "Soumkine Notebooks",
        "description": "For almost 20 years I work as an illustrator and graphic designer. \nEach new project starts with sketches in my notebooks. To record any idea you need a good support to write or sketch out.\n\nOver this time I have tried all kinds of notebooks; everything from cheap school notepads to premium class diaries. One day I realised no notebooks existed to meet these conditions, so I created my own. The Soumkine Notebook.\n\nFurthermore, I created a stationery company Soumkine Notebooks that focuses on high quality hand-made notebooks with premium paper.\n\nSoumkine — It's my last name in the French manner. It pronounces as [sum’kin].\n\nP.S. Next time I will tell more about the Soumkine identity and especially about the logo. But today I can't wait to present the last collection of \"A5 Slim\" notebooks. By the way, they are all available in our online shop Soumkine.com, so don't miss out! ",
        "filename": "d595f541911437.Y3JvcCwxMjcyLDk5Niw2NSww.jpg",
        "src": "https://mir-s3-cdn-cf.behance.net/projects/original/d595f541911437.Y3JvcCwxMjcyLDk5Niw2NSww.jpg"
    },
    ...
]
```

For the query system, you have to feed it JSON to query. It have to use array typed, because of the multiple queries support.

##### Query fields

| Fields | Type | Description |
| ------ | ---- | ----------- |
| title | string | The title of the project in Behance |
| description | string | Yes, we do query the description too, regex not supported :( |
| filename | string | Filename of the image covers fetched from Behance |
| src | string | The original source address of the image covers from Behance. If you want to access from localhost, do using the `/imgs` route |

```
POST localhost:8080/q

[
    {
        "title": "Soumkine Notebooks"
    },
    {
        "src": "https://mir-s3-cdn-cf.behance.net/projects/original/d595f541911437.Y3JvcCwxMjcyLDk5Niw2NSww.jpg"
    }
]
```

#### Statics
To access the image covers fetched from Behance locally, do use `localhost:8080/imgs/` and your filename after it.

```
GET localhost:8080/imgs/your_filename
```