#!/bin/bash
user="ubuntu"
server="plailist.app"
# Get the directory where the script resides
DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

# Get the parent directory of the script
PARENT_DIR="$(dirname "$DIR")"

cd $PARENT_DIR/server
mkdir build
env GOOS=linux GOARCH=arm64 go build -o build/chat_arm64
scp chat_arm64 ubuntu@plailist.app:~/

ssh ${user}@${server} '
  sudo systemctl stop chat &&
  sudo mv chat_arm64 /usr/local/bin/ &&
  sudo systemctl start chat
'

echo "Removing the build folder."
rm -rf build;
