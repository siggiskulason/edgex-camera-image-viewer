

echo "copying a secrets token to '/tmp' to use. This only works if using the edgexfoundry snap"
echo "also, note that the token has a TTL of 1 hr, so you must have started up EdgeX less than 1 hour ago"

sudo cp /var/snap/edgexfoundry/current/secrets/app-rules-engine/secrets-token.json /tmp
sudo chmod a+rw /tmp/secrets-token.json 

make
google-chrome-stable --new-window --app="file://$(pwd)/image.html" &
./app-service

