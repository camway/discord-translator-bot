# Usage

How to use the bot once it's in your discord server.

## Discord Commands

* !commands
* !list languages
* !list groups
* !create GROUPNAME
* !join GROUPNAME LANGUAGECODE
* !leave GROUPNAME
* !delete GROUPNAME

## Concept

Groups in this project are a collection of channels

These channels are considered to be in the same "room"

Any message from any of these rooms will be translated to the other member channels

## Setup

* Create a discord server
* Create a channel for each language you want to support (make sure it's in the list below)
* In a discord channel the bot can see, run: `!create mychat`
* Get the language code below for each language you want to support
* In each language's channel, run: `!join mychat LANG_CODE` where LANG_CODE is the characters next to the language below
* Start chatting

# Supported Languages

Latest here: https://cloud.google.com/translate/docs/languages

* Amharic: am
* Arabic: ar
* Basque: eu
* Bengali: bn
* Chinese: (PRC): zh-CN
* Chinese (Taiwan): zh-TW
* English (UK): en-GB
* Bulgarian: bg
* Catalan: ca
* Cherokee: chr
* Croatian: hr
* Czech: cs
* Danish: da
* Dutch: nl
* English (US): en
* Estonian: et
* Filipino: fil
* Finnish: fi
* French: fr
* German: de
* Greek: el
* Gujarati: gu
* Hebrew: iw
* Hindi: hi
* Hungarian: hu
* Icelandic: is
* Indonesian: id
* Italian: it
* Japanese: ja
* Kannada: kn
* Korean: ko
* Latvian: lv
* Lithuanian: lt
* Malay: ms
* Malayalam: ml
* Marathi: mr
* Norwegian: no
* Polish: pl
* Portuguese (Brazil): pt-BR
* Portuguese (Portugal): pt-PT
* Romanian: ro
* Russian: ru
* Serbian: sr
* Slovak: sk
* Slovenian: sl
* Spanish: es
* Swahili: sw
* Swedish: sv
* Tamil: ta
* Telugu: te
* Thai: th
* Turkish: tr
* Urdu: ur
* Ukrainian: uk
* Vietnamese: vi
* Welsh: cy

# Deployment

Once this gets further along, I'll flush this out more.

For now, the docker-compose.yml is a good place to look for what is required (ports, services, networks, etc). The environment variables required are detailed in the developer section.

# Developers

## TODO

* Code is very much POC. Needs a refactor badly
* Add testing suite once code is structured better
* Add hooks into discord to handle potential desyncs (ex channel/server delete)
* Add commands so users can join just the channels for their language (prevent notification spam from other languages)
* Add user tracking and permissions
* Add root user(s) (Basically admin on initial boot)


## Setup

Copy the `.env.example` file to `.env`, and fill in the values

`DISCORD_BOT_TOKEN`

Your discord bot token. You can get one here: https://discord.com/developers/applications

`GOOGLE_API_JSON_TOKEN_PATH`

Your google API json file. You can get one here: https://console.cloud.google.com/apis/credentials

`GOOGLE_API_PROJECT_ID`

Your google API Project ID. You can create one here: https://console.cloud.google.com/projectcreate

`MONGO_USERNAME`
`MONGO_PASSWORD`

MongoDB credentials.

`MONGO_DATABASE`

Database name to use for storage.
