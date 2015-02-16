cd #variable_chdir_path#

if [ ! -d #variable_virtualenv_path# ];
then
    /usr/local/opt/python-3.4/bin/pyvenv #variable_virtualenv_path#
fi

source #variable_virtualenv_path#bin/activate

pip install -r #variable_requirements_file#

mkdir -p /home/#user#/apps/#appname#/logs && chown -R #user#:#user# /home/#user#/apps/#appname#/logs

#command#
