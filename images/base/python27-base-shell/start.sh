#!/bin/bash

echo "$USER:x:$USERID:" >> /etc/group && echo "$USER:x:$USERID:$USERID:,,,:/home/$USER:/bin/bash" >> /etc/passwd

su - $USER
