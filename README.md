# PhotoSync

This is a small tool to allow semi-seamless synchronization of photos between Tinybeans and Aura Frames.  Tinybeans is
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
Add application configuration to a `config.yaml` file in the directory running the program:

```yaml
live: {{boolean - send emails}}

mitm:
  host: {{hostname}}
  port: {{port}}

sender:
  from: {{email address from}}
  gmail:
    credentials: {{relative path to credential file}}
  smtp:
    host: {{smtp host}}
    username: {{username}}
    password: {{password}}
    port: {{port}}

tinybeans:
  username: {{username}}
  password: {{password}}
  journals:
    - id: {{journalID}}
      users:
        - id: {{tinybeansUserID}}
          frames:
            - {{AuraFrameEmail #1}}
            - {{AuraFrameEmail #2}}
        - id: {{differentUserID}}
          frames:
            - {{AuraFrameEmail #1}}
```