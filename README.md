# P2P file system of cortex full node
```
go get github.com/CortexFoundation/torrentfs

or

make
```
#### torrent-create : to create torrent file
./torrent-create file > torrent
#### torrent-magnet : load info hash from torrent file
./torrent-magnet < torrent
#### seeding : to seed file to dht
./seeding -dataDir=
#### torrent : to download file
./torrent download $magnet
