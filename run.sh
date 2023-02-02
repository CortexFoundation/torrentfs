./build/bin/torrent-create workspace/data -p=4096 > workspace/test-torrent
./build/bin/torrent-magnet < workspace/test-torrent
mkdir -p mnt/9196320d998fdab966bcb3a08f3f087e1f993c12/data
cp workspace/test-torrent mnt/9196320d998fdab966bcb3a08f3f087e1f993c12/torrent
cp -r workspace/data/* mnt/9196320d998fdab966bcb3a08f3f087e1f993c12/data
./build/bin/seeding -dataDir=mnt
