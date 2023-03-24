#!/bin/bash
rm -rf ./dist
npm run build
rm dist.zip
zip -r dist.zip dist/
echo "Uploading the dist.zip file to the remote server..."
scp dist.zip ubuntu@app.judgefest.com:~/
user="ubuntu"
server="plailist.app"

# Remove the directory and its contents
echo "Removing the ~/dist directory and its contents..."

ssh ${user}@${server} 'rm -rf ~/dist'

echo "Unzipping the dist.zip file, moving the contents, and setting ownership..."
# Unzip the dist.zip file, move the contents, and set ownership
ssh ${user}@${server} '
  unzip -q ~/dist.zip -d ~ &&
  sudo rm -rf /var/www/plailist/* &&
  sudo mv ~/dist/* /var/www/plailist/ &&
  sudo chown -R www-data:www-data /var/www/plailist;
'

echo "Operation completed successfully."