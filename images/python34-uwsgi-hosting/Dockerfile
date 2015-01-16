FROM manager/python34-base-hosting

RUN echo "#user#:x:#uid#:" >> /etc/group && echo "#user#:x:#uid#:#uid#:,,,:/home/#user#:/bin/bash" >> /etc/passwd

USER #user#

cd #variable_chdir_path#

if [ ! -d #variable_virtalenv_path# ]; then
	/usr/local/opt/python-3.4/bin/pyvenv #variable_virtalenv_path#
fi

source #variable_virtalenv_path#

pip install -r #variable_requirements_file#
pip install uwsgi

EXPOSE 9000

CMD uwsgi --socket 0.0.0.0:9000 --wsgi-file #variable_wsgi_file# --master --processes #variable_processes# --threads #variable_threads#
