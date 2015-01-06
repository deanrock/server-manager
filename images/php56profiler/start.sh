#!/bin/bash
echo $USER
echo "$USER:x:$UID:" >> /etc/group
echo "$USER:x:$UID:$UID:,,,:/home/$USER:/bin/bash" >> /etc/passwd

sed -i s/ENV_USER/$USER/ /etc/php5/fpm/php-fpm.conf
sed -i s/ENV_USER/$USER/ /etc/php5/fpm/pool.d/www.conf

mkdir /mystuff
chown $USER:$USER /mystuff

#export
#su - $USER -c $@

export > /home/$USER/blah.log
echo "yo\n" >> /home/$USER/blah.log
$@
