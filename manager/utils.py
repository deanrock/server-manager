
import tempfile
import shutil
from subprocess import Popen, PIPE

from django.conf import settings


def get_temp_folder():
    tempdir = tempfile.mkdtemp(prefix='manager-',
                               dir=settings.TMP_FOLDER)

    return tempdir

class Logs:
    def __init__(self):
        self.logs = []

    def add(self, obj, flag=None):
        try:
            self.logs.append({'message': unicode(obj), 'flag': flag})
        except:
            self.logs.append({'message': unicode(obj, 'utf-8'), 'flag': flag})

    def append(self, logs):
        for x in logs.logs:
            self.logs.append(x)

def delete_temp_folder(temp_folder_path):
    try:
        shutil.rmtree(temp_folder_path)
    except Exception as e:
        print('cant delete %s - %s' % (temp_folder_path, e))

def exec_command(logs, cmd):
    logs.add("CMD: %s" % cmd)
    sp = Popen(
        cmd,
        stdout=PIPE, stderr=PIPE, shell=True)

    out, err = sp.communicate()

    if err:
        logs.add("ERROR: %s" % err)

    if out:
        logs.add("OUTPUT: %s" % out)

    return out, err
