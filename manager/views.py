from django.contrib.auth.decorators import login_required
from django.http.response import HttpResponse, HttpResponseRedirect
from django.shortcuts import render_to_response
from manager import actions
from manager.models import App


@login_required
def rebuild_base_image(request):
    logs = actions.rebuild_base_image()

    return render_to_response('sync_action.html',
        {
            'logs': logs.logs,
            'action': 'Rebuild base image'
        })

@login_required
def sync_users(request):
    logs = actions.sync_accounts()

    return render_to_response('sync_action.html',
        {
            'logs': logs.logs,
            'action': 'Sync users'
        })


@login_required
def sync_databases(request):
    logs = actions.sync_databases()

    return render_to_response('sync_action.html',
        {
            'logs': logs.logs,
            'action': 'Sync databases'
        })

@login_required
def actions_container(request, id, action):
    app = App.objects.get(id=id)

    logs = actions.Logs()

    logs.add("app %s, account %s, container name %s" % (app.name, app.account.name, app.container_name()))

    logs.append(app.redeploy())

    return render_to_response("sync_action.html",
        {
            "action": app.name,
            "logs": logs.logs,
        })

@login_required
def index(request):
    return HttpResponseRedirect('/admin/')

@login_required
def update_nginx_config(request):
    logs = actions.update_nginx_config()

    return render_to_response('sync_action.html',
        {
            'logs': logs.logs,
            'action': 'Update nginx config'
        })
