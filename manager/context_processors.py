from django.conf import settings # import the settings file

def server_friendly_name(request):
    # return the value you want as a dictionnary. you may add multiple values in there.
    return {'server_friendly_name': settings.NAME}