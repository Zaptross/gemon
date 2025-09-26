package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/pion/rtp"
	"github.com/pion/webrtc/v3/pkg/media"
	"github.com/pion/webrtc/v3/pkg/media/oggwriter"
	"github.com/samber/lo"
	"github.com/zaptross/discordgo"
)

// Variables used for command line parameters
var (
	Token       string
	ChannelID   string
	GuildID     string
	OutputDir   string
	ProgressDir string
)

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.StringVar(&GuildID, "g", "", "Guild in which voice channel exists")
	flag.StringVar(&ChannelID, "c", "", "Voice channel to connect to")
	flag.StringVar(&OutputDir, "o", "", "Directory to save finished recordings to")
	flag.StringVar(&ProgressDir, "r", "", "Directory to save in-progress recordings to")
	flag.Parse()
}

func createPionRTPPacket(p *discordgo.Packet) *rtp.Packet {
	return &rtp.Packet{
		Header: rtp.Header{
			Version: 2,
			// Taken from Discord voice docs
			PayloadType:    0x78,
			SequenceNumber: p.Sequence,
			Timestamp:      p.Timestamp,
			SSRC:           p.SSRC,
		},
		Payload: p.Opus,
	}
}

type TimeSeparatedFile struct {
	ssrc      uint32
	username  string
	startTime int64
	lastEdit  int64
	file      media.Writer
}

func (tsf *TimeSeparatedFile) Close() {
	tsf.file.Close()
}
func (tsf *TimeSeparatedFile) IsStale() bool {
	return (time.Now().UnixMilli() - tsf.lastEdit) > MAX_SILENCE_GAP
}
func (tsf *TimeSeparatedFile) FileName() string {
	name := tsf.username
	if name == "" {
		name = fmt.Sprintf(("%d"), tsf.ssrc)
	}

	return fmt.Sprintf("%d_%s.ogg", tsf.startTime, name)
}

type TSFMap map[uint32]TimeSeparatedFile

var (
	MAX_SILENCE_GAP = (1 * time.Second).Milliseconds()
)

func handleVoice(c chan *discordgo.Packet, ssrcToUser *map[int]string) {
	// files := make(map[uint32]media.Writer)
	timeSepFiles := make(TSFMap)
	for p := range c {
		// file, ok := files[p.SSRC]
		tsf, ok := timeSepFiles[p.SSRC]
		if !ok || tsf.IsStale() {
			fmt.Printf("Starting recording for SSRC %d\n", p.SSRC)
			if ok {
				// Close the old file if it exists and is stale
				fmt.Printf("Closing stale file for SSRC %d\n", p.SSRC)
				tsf.Close()

				go func(filename string) {
					// Move the file to the output directory for processing
					oldPath := fmt.Sprintf("%s/%s", ProgressDir, filename)
					newPath := fmt.Sprintf("%s/%s", OutputDir, filename)
					err := os.Rename(oldPath, newPath)

					if err != nil {
						fmt.Printf("failed to move file %s to %s: %v\n", oldPath, newPath, err)
					}
				}(tsf.FileName())
			}
			tsf = TimeSeparatedFile{
				ssrc:      p.SSRC,
				username:  (*ssrcToUser)[int(p.SSRC)],
				startTime: time.Now().UnixMilli(),
				lastEdit:  time.Now().UnixMilli(),
			}

			var err error
			tsf.file, err = oggwriter.New(fmt.Sprintf("%s/%s", ProgressDir, tsf.FileName()), 48000, 2)
			if err != nil {
				fmt.Printf("failed to create file %d.ogg, giving up on recording: %v\n", p.SSRC, err)
				return
			}
			// files[p.SSRC] = file
			timeSepFiles[p.SSRC] = tsf
		}
		// Construct pion RTP packet from DiscordGo's type.
		rtp := createPionRTPPacket(p)
		err := tsf.file.WriteRTP(rtp)
		if err != nil {
			fmt.Printf("failed to write to file %d.ogg, giving up on recording: %v\n", p.SSRC, err)
		}
		tsf.lastEdit = time.Now().UnixMilli()
		timeSepFiles[p.SSRC] = tsf // Update the map entry to reflect lastEdit change
	}

	// Once we made it here, we're done listening for packets. Close all files, and move any in-progress ones to the output directory.
	for _, tsf := range timeSepFiles {
		tsf.Close()

		go func(filename string) {
			// Move the file to the output directory for processing
			oldPath := fmt.Sprintf("%s/%s", ProgressDir, filename)
			newPath := fmt.Sprintf("%s/%s", OutputDir, filename)
			err := os.Rename(oldPath, newPath)
			if err != nil {
				fmt.Printf("failed to move file %s to %s: %v\n", oldPath, newPath, err)
			}
		}(tsf.FileName())
	}
}

func main() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	s, err := discordgo.New("Bot " + Token)
	if err != nil {
		fmt.Println("error creating Discord session:", err)
		return
	}
	defer s.Close()

	// We only really care about receiving voice state updates.
	intents := discordgo.IntentsGuildVoiceStates | discordgo.IntentsGuildMembers | discordgo.IntentsGuilds | discordgo.IntentsGuildPresences
	s.Identify.Intents = discordgo.MakeIntent(intents)

	println("Opening connection to Discord...")
	err = s.Open()
	if err != nil {
		fmt.Println("error opening connection:", err)
		return
	} else {
		fmt.Println("Bot is now running. Press CTRL-C to exit.")
	}

	members, err := s.GuildMembers(GuildID, "", 1000)

	if err != nil {
		fmt.Println("failed to get guild members:", err)
		return
	}

	ssrcToUser := make(map[int]string)

	vuh := func(vc *discordgo.VoiceConnection, vs *discordgo.VoiceSpeakingUpdate) {
		if vs.UserID == s.State.User.ID {
			// Ignore updates about ourselves
			return
		}

		if vs.SSRC != 0 && vs.UserID != "" {
			member, ok := lo.Find(members, func(m *discordgo.Member) bool {
				return m.User.ID == vs.UserID
			})

			username := "unknown"

			if ok {
				username = member.DisplayName()
			} else {
				username = vs.UserID
			}
			ssrcToUser[vs.SSRC] = username
		}
	}

	println("Joining voice channel...")
	v, err := s.ChannelVoiceJoin(GuildID, ChannelID, true, false, vuh)
	if err != nil {
		fmt.Println("failed to join voice channel:", err)
		return
	}

	done := make(chan struct{})

	go func() {
		select {
		case <-time.After(90 * time.Second):
			println("Recording timeout reached, shutting down...")
			close(done)
			v.Disconnect()
		case <-sigChan:
			println("Received termination signal, shutting down...")
			close(done)
			v.Disconnect()
		}
	}()

	go func() {
		<-done
		close(v.OpusRecv)
		v.Close()
	}()

	handleVoice(v.OpusRecv, &ssrcToUser)
}
