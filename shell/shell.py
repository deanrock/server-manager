#!/usr/bin/python
from subprocess import Popen, PIPE
import os
import sys

shells = ["php56", "python27", "python34", "java8"]
users = [o for o in os.listdir("/home/") if os.path.isdir(os.path.join("/home/",o))]

if len(sys.argv) == 3:
    s = sys.argv[1]
    u = sys.argv[2]

    if not s in shells:
        print("shell doesnt exist!")
        sys.exit()

    if not u in users:
        print("user doesnt exist!")
        sys.exit()

    shell = 'manager/%s-base-shell' % s
    sp = Popen(
        'id -u %s' % u,
        stdout=PIPE, stderr=PIPE, shell=True)

    uid, err = sp.communicate()

    os.execv("/usr/bin/docker", ['/usr/bin/docker',
                                 'run',
                                 '-it',
                                 '--rm',
                                 '-v', '/home/%s:/home/%s' % (u, u),
                                 '-e', 'USER=%s' % u,
                                 '-e', 'USERID=%s' % uid.rstrip(),
                                 '--add-host', 'mysql:172.17.42.1',
                                 shell,
                                 '/mystuff/start.sh'])
else:
    print("help\n--------\nshells: ")

    for shell in shells:
        print(" - %s" % shell)

    print("\nuse: shell <shell> <username>")

