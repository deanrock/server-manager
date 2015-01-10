from django.contrib.auth.decorators import login_required
from django.forms.models import modelform_factory
from django.http.response import HttpResponse, HttpResponseRedirect
from django.shortcuts import render_to_response
from django.template import RequestContext
from manager import actions
from manager.models import App, Account


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

    if action == 'redeploy':
        logs.append(app.redeploy())
    elif action == 'stop':
        logs.add(app.stop())
    elif action == 'start':
        logs.add(app.start())

    return render_to_response("sync_action.html",
        {
            "action": action,
            "logs": logs.logs,
        })

@login_required
def index(request):
    return render_to_response('index.html',
        {
            'accounts': Account.objects.all(),
        },
                              context_instance=RequestContext(request))

@login_required
def account(request, name):
    return render_to_response('account/overview.html',
        {
            'account': Account.objects.filter(name=name).first()
        },
                              context_instance=RequestContext(request))

@login_required
def account_apps_edit(request, name, app):
    account = Account.objects.filter(name=name).first()
    AppForm = modelform_factory(App)
    return render_to_response('account/apps_edit.html',
                              {
            'account': account,
            'formset': AppForm(instance=account.apps.filter(id=app).first())
        },
                              context_instance=RequestContext(request))

@login_required
def account_apps(request, name):
    account = Account.objects.filter(name=name).first()
    return render_to_response('account/apps.html',
                              {
            'account': account,
            'apps': account.apps.all(),
        },
                              context_instance=RequestContext(request))

@login_required
def account_domains(request, name):
    return render_to_response('account/domains.html',
                              {
            'account': Account.objects.filter(name=name).first()
        },
                              context_instance=RequestContext(request))

@login_required
def account_databases(request, name):
    return render_to_response('account/overview.html',
                              {
            'account': Account.objects.filter(name=name).first()
        },
                              context_instance=RequestContext(request))


@login_required
def update_nginx_config(request):
    logs = actions.update_nginx_config()

    return render_to_response('sync_action.html',
        {
            'logs': logs.logs,
            'action': 'Update nginx config'
        })
