# PhotoSync

A tool to allow semi-seamless synchronization of photos between Tinybeans and Aura Frames.  Tinybeans is
a depository of photos of children, to share with family members allowed to view the countless uploads of kiddo pictures.
Aura Frames is a digital picture frame that work as wonderful portals into grandma and grandpas house to share
pictures of the little ones.

Rather than needing to maintain two sets of picture data, this tools allows me to upload way too many pictures
into tinybeans, and when I "like" a picture, this tool will ship that photo off to the Aura frames of the grandparents.
Thus allowing me a simple means of pushing only the best to be displayed at Grandma and Grandpas house.

It is configurable to allow custom members of the Tinybeans journal to track, and to what frames their likes should be
sent to.

## Installation
```shell
go get github.com/bp1222/photoSync
```

## Usage
### Configuration
#### Environment
Add application configuration to a `.env` file in the directory running the program:

```shell
TINYBEANS_USERNAME=""
TINYBEANS_PASSWORD=""

EMAIL_HOST=""
EMAIL_USERNAME=""
EMAIL_PASSWORD=""
EMAIL_PORT=""

# Enable proxying of transport through a MITM Proxy to watch traffic
DEBUG_MITM_PROXY=true/false

# Defaults to false, needs to be enabled for actual sending to happen
LIVE_SEND_MAIL=true/false
```

#### Users and Frames
Within a `userFrameConfig.yaml` file add your journal ID's, user ID's and frame addresses

```yaml
journals:
  - id: {{journalID}}
    users:
      - id: {{tinybeansUserID}} 
        frames:
          - {{AuraFrameEmail #1}}
          - {{AuraFrameEmail #2}}
      - id: 2709241
        frames:
          - {{AuraFrameEmail #1}}
```
