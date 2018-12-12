set -ex

pushd /tmp
test -e zos_linux ||  wget https://github.com/threefoldtech/zos/releases/download/0.1.6/zos_linux

chmod +x zos_linux

./zos_linux configure --name=examplemachine --address="$ZOSMACHINEADDRESS" --setdefault

# zos container new --root=https://hub.grid.tf/tf-autobuilder/threefoldfoundation-tfchain-master_bridged.flist
# NEEDS EXPLICIT MERGE with ubuntu flist (https://hub.grid.tf/tf-autobuilder/threefoldfoundation-tfchain-master_bridged.flist)
./zos_linux container new --root=https://hub.grid.tf/thabet/ubuntu_bridged.flist  
