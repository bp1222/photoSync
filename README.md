# PhotoSync

A tool to allow semi-seamless synchronization of photos between Tinybeans and Aura Frames.  Tinybeans is
a depository of photos for children, with family members allowed to view the countless uploads of kiddo pictures.
Aura Frames are a digital picture frame that work as wonderful portals into grandma and grandpas house to share
pictures of the little ones.

Rather than needing to maintain two sets of picture data, this tools allows me to upload way too many pictures
into tinybeans, and when I "Like" a picture, this tool will ship that photo off to the Aura frames of the grandparents.
Thus allowing me a simple means of pushing only the best to be displayed at Grandma and Grandpas house.

## Installation

```shell
go get github.com/bp1222/photoSync
```

## Usage

### Configuration
#### Environment
Add your application configuration to your `.env` file in the root of your project:

```shell
TINYBEANS_USERNAME=""
TINYBEANS_PASSWORD=""

EMAIL_HOST=""
EMAIL_USERNAME=""
EMAIL_PASSWORD=""
EMAIL_PORT=""
```

#### Users and Frames
Within a `userFrameConfig.yaml` file add your journal ID's, user ID's and frame addresses

```yaml
journals:
  - id: {{journalID}}
    users:
      - id: {{userID}} 
        frames:
          - {{AuraFrameEmail #1}}
          - {{AuraFrameEmail #2}}
      - id: 2709241
        frames:
          - {{AuraFrameEmail #1}}
```


### Running
Add execution of the script to cron for a daily pull
