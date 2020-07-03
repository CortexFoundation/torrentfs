# P2P file system of cortex full node
```
go get github.com/CortexFoundation/torrentfs

or

make
```
```cd build/bin```

```echo "Hello torrent" > file ```
#### torrent-create : to create torrent file
```./torrent-create file > torrent```
#### torrent-magnet : load info hash from torrent file
```./torrent-magnet < torrent```
#### seeding : to seed file to dht
```./seeding -dataDir=store```
under store folder
```
f3bc013c3cda4bb9f74a6f7696e66efa7be92a9a
├── data
└── torrent
```
#### torrent : to download file
```./torrent download 'infohash:f3bc013c3cda4bb9f74a6f7696e66efa7be92a9a' ```
