package main

import (
	"log"
	"os"

	"github.com/anacrolix/tagflag"
	"github.com/anacrolix/torrent/bencode"
	"github.com/anacrolix/torrent/metainfo"

	"github.com/CortexFoundation/torrentfs/params"
)

var (
	builtinAnnounceList = [][]string{
		//{"udp://tracker.cortexlabs.ai:5008"},
		//{"udp://tracker.openbittorrent.com:80"},
		//{"udp://tracker.publicbt.com:80"},
		//{"udp://tracker.istole.it:6969"},
		params.MainnetTrackers,
	}
)

func main() {
	log.SetFlags(log.Flags() | log.Lshortfile)
	var args struct {
		AnnounceList      []string      `name:"a" help:"extra announce-list tier entry"`
		EmptyAnnounceList bool          `name:"n" help:"exclude default announce-list entries"`
		Comment           string        `name:"t" help:"comment"`
		CreatedBy         string        `name:"c" help:"created by"`
		PieceLength       tagflag.Bytes `name:"p" help:"piece size"`
		tagflag.StartPos
		Root string
	}
	tagflag.Parse(&args, tagflag.Description("Creates a torrent metainfo for the file system rooted at ROOT, and outputs it to stdout."))
	mi := metainfo.MetaInfo{
		AnnounceList: builtinAnnounceList,
	}
	if args.EmptyAnnounceList {
		mi.AnnounceList = make([][]string, 0)
	}
	for _, a := range args.AnnounceList {
		mi.AnnounceList = append(mi.AnnounceList, []string{a})
	}
	mi.SetDefaults()
	if len(args.Comment) > 0 {
		mi.Comment = args.Comment
	}
	if len(args.CreatedBy) > 0 {
		mi.CreatedBy = args.CreatedBy
	} else {
		mi.CreatedBy = "github.com/CortexFoundation/torrentfs"
	}
	if args.PieceLength.Int64() == 0 {
		args.PieceLength = 4096
	}
	info := metainfo.Info{
		PieceLength: args.PieceLength.Int64(),
	}
	err := info.BuildFromFilePath(args.Root)
	if err != nil {
		log.Fatal(err)
	}
	mi.InfoBytes, err = bencode.Marshal(info)
	if err != nil {
		log.Fatal(err)
	}
	err = mi.Write(os.Stdout)
	if err != nil {
		log.Fatal(err)
	}
}
