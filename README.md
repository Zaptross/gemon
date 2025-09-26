# gemon
An LLM powered D&amp;D notetaking application, intended to integrate into discord.

## Setup

1. Install ffmpeg from: ffmpeg.org/download.html
2. Install golang from: https://go.dev/doc/install
3. Setup python virtual env
    * `python -m venv .venv`
    * `./.venv/Script/activate`

### Why "gemon"?

[**Ge**offery of **Mon**mouth](geoffery-monmouth) was a cleric and author of dubious reliability of some tales of King Arthur and arthurian legend.


### References

[This blog post](pavlin-blog) by Pavlin about Local Speech to Text formed the basis of the language model work.

[The voice receive example](discordgo-voice) by the maintainers of discordgo formed the base of the voice recording work.

[pavlin-blog]: https://www.pavlinbg.com/posts/python-speech-to-text-guide
[discordgo-voice]: https://github.com/zaptross/discordgo/blob/master/examples/voice_receive/main.go
[geoffery-monmouth]: https://en.wikipedia.org/wiki/Geoffrey_of_Monmouth