from manager.settings.base import *

# SECURITY WARNING: keep the secret key used in production secret!
SECRET_KEY = '%r*#rhgz52x#my8(fx69$mz31()g_$@l7!4nz23e6m&fxet@jz'

# SECURITY WARNING: don't run with debug turned on in production!
DEBUG = True

TEMPLATE_DEBUG = True

MYSQL_ROOT_PASSWORD = 'password'

DATABASES['manager'] = {
    'NAME': 'manager',
    'ENGINE': 'django.db.backends.mysql',
    'USER': 'manager',
    'PASSWORD': 'ret56786vew9enr89'
}
