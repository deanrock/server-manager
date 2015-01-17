cd #variable_chdir_path#

if [ ! -d #variable_virtualenv_path# ];
then
    /usr/local/opt/python-2.7/bin/pyvenv #variable_virtualenv_path#
fi

source #variable_virtualenv_path#bin/activate

pip install -r #variable_requirements_file#
pip install uwsgi

mkdir -p /home/#user#/apps/#appname#/logs && chown -R #user#:#user# /home/#user#/apps/#appname#/logs

uwsgi --socket 0.0.0.0:9000 --wsgi-file #variable_wsgi_file# --master --processes #variable_processes# --threads #variable_threads# --logto /home/#user#/apps/#appname#/logs/uwsgi.log
