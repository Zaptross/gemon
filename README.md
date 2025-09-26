# gemon
An LLM powered D&amp;D notetaking application, intended to integrate into discord.

## Setup

1. Install ffmpeg from: ffmpeg.org/download.html
2. Install golang from: https://go.dev/doc/install
3. Setup python virtual env
    * `python -m venv .venv`

### Bot Setup

To setup your discord bot, go to https://discord.com/developers/applications and create a new application.

Once created, go to the Bot tab and configure a bot, making sure to copy out the token.

Scroll down to the intents and enable `PRESENCE INTENT` and `SERVER MEMBERS INTENT`.

Lastly, use the following link and your application's client ID to add your bot to the server of your choice.

https://discordapi.com/permissions.html#36508273664

### Running locally

To run `stt` the transcriber as a folder-watcher, run:
1. `./.venv/Script/activate`
2. `python .\stt\main.py -i $inputDir -o $outputDir`
    * where `-i` is for the input directory to watch and;
    * `-o` is the output directory to move audio and transcripts to

To run `bot` the discord bot:
1. `go run .\cmd\bot\ -t $discordToken -g $guildID -c $channelID -r $recordingDir -o $outputDir`
    * where `-t` is the token of your discord bot
    * `-g` is the ID of the server to connect to
    * `-c` is the ID of the voice channel to join
    * `-r` is the recording directory to store in-progress audio samples
    * `-o` is the output directory to move finalised audio samples to

### Why "gemon"?

[**Ge**offery of **Mon**mouth](geoffery-monmouth) was a cleric and author of dubious reliability of some tales of King Arthur and arthurian legend.


### References

[This blog post](pavlin-blog) by Pavlin about Local Speech to Text formed the basis of the language model work.

[The voice receive example](discordgo-voice) by the maintainers of discordgo formed the base of the voice recording work.

[pavlin-blog]: https://www.pavlinbg.com/posts/python-speech-to-text-guide
[discordgo-voice]: https://github.com/zaptross/discordgo/blob/master/examples/voice_receive/main.go
[geoffery-monmouth]: https://en.wikipedia.org/wiki/Geoffrey_of_Monmouth