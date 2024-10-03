# Server

## `/aquarium`

Serve the aquarium frontend.

## Photo Upload

`/aquarium/<aquariumID>`

Upload a fish image to the aquarium.

## Subscribe Fish Changes

`/aquarium/<aquariumID>/sse`

Send new Fishes or Fish deletions (over Admin Panel).

Messages:

```
event: ping
data: {}

event: fishleft
data: {"id":"<fishID>","aquarium_id":"<aquariumID>","name":"<fishName>","filename":"<filename>"}

event: fishjoin
data: {"id":"<fishID>","aquarium_id":"<aquariumID>","name":"<fishName>","filename":"<filename>"}
```

## Get Fish Image

`/aquarium/<aquariumID>/fishes/<fishID>.png`

Serve the fish image file.

## Admin Panel

- Show Aquariums
- Delete Fishes

`/admin`
