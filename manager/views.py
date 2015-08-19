from django.contrib.auth.decorators import login_required
from django.core.serializers import serialize
import json
from django.core.urlresolvers import reverse
from django.forms.models import modelform_factory, model_to_dict
from django.http.response import HttpResponse, HttpResponseRedirect, HttpResponseNotFound
from django.shortcuts import render_to_response
from django.template import RequestContext
import docker
from django.views.decorators.csrf import csrf_exempt
from manager import actions, docker_api
from manager.forms import AppForm, DomainForm, DatabaseForm, UserSSHKeyForm
from manager.models import App, Account, Domain, Database, UserSSHKey, Image, AppImageVariable


@login_required
def action_ajax(request, action):
    if action == 'rebuild-base-image':
        logs = actions.rebuild_base_image()
    elif action == 'sync-users':
        logs = actions.sync_accounts()
    elif action == 'sync-databases':
        logs = actions.sync_databases()
    elif action == 'update-nginx-config':
        logs = actions.update_nginx_config()

    data = json.dumps(logs.logs)
    return HttpResponse(data, content_type='application/json')


@login_required
@csrf_exempt
def account_apps_action_ajax(request, name, app, action):
    account = Account.objects.filter(name=name).first()
    app = account.apps.filter(id=app).first()

    logs = actions.Logs()

    logs.add("app %s, account %s, container name %s" % (app.name, app.account.name, app.container_name()))

    if action == 'redeploy':
        logs.append(app.redeploy())
    elif action == 'stop':
        logs.add(app.stop())
    elif action == 'start':
        logs.add(app.start())

    data = json.dumps(logs.logs)
    return HttpResponse(data, content_type='application/json')


@login_required
def account_apps(request, name):
    account = Account.objects.filter(name=name).first()
    apps = account.apps.all()

    containers = docker_api.cli.containers()
    mapping = {}
    for c in containers:
        try:
            name = c['Names'][-1].replace('/', '')

            mapping[name] = (c['Status'], c['Id'])
        except Exception as e:
            print(e)

    dapps = []

    for app in apps:
        app.up = False
        app.status = None
        app.container_id = None
        if app.container_name() in mapping:
            app.status = mapping[app.container_name()][0]
            app.container_id = mapping[app.container_name()][1]

            if 'Up ' in app.status:
                app.up = True

        dapps.append({
            'id': app.id,
            'memory': app.memory,
            'up': app.up,
            'status': app.status,
            'account': app.account.id,
            'account_name': app.account.name,
            'container_id': app.container_id,
            'name': app.name,
            'image': app.image.name,
            })

    data = json.dumps(dapps)
    return HttpResponse(data, content_type='application/json')

@login_required
def containers(request):
    apps = App.objects.all()

    containers = docker_api.cli.containers()
    mapping = {}
    for c in containers:
        try:
            name = c['Names'][-1].replace('/', '')

            mapping[name] = c['Status']
        except Exception as e:
            print(e)

    dapps = []

    print (mapping)

    for app in apps:
        app.up = False
        app.status = ''

        if app.container_name() in mapping:
            app.status = mapping[app.container_name()]

            if 'Up ' in app.status:
                app.up = True

        dapps.append({
            'id': app.id,
            'memory': app.memory,
            'up': app.up,
            'status': app.status,
            'account': app.account.id,
            'account_name': app.account.name,
            'name': app.name,
            'image': app.image.name,
            })

    


    data = json.dumps(dapps)
    return HttpResponse(data, content_type='application/json')

@login_required
def account_databases(request, name):
    return render_to_response('account/databases.html',
                              {
            'account': Account.objects.filter(name=name).first()
        },
                              context_instance=RequestContext(request))


@login_required
def account_databases_edit(request, name, database=None):
    account = Account.objects.filter(name=name).first()
    df = modelform_factory(Database, form=DatabaseForm)

    if request.method == 'POST':
        formset = df(request.POST, request.FILES,
                          instance=account.databases.filter(id=database).first())

        if formset.is_valid():
            obj = formset.save(commit=False)
            obj.account = account
            obj.added_by = request.user
            obj.save()

            return HttpResponseRedirect(reverse('manager.views.account_databases', kwargs={'name': account.name}))
    else:
        formset = df(instance=account.databases.filter(id=database).first())

    return render_to_response('account/databases_edit.html',
                              {
            'account': account,
            'formset': formset
        },
                              context_instance=RequestContext(request))

@login_required
def profile_sshkeys(request):
    return render_to_response('profile/ssh_keys.html',
        {},
                              context_instance=RequestContext(request))


@login_required
def profile_sshkeys_edit(request, key=None):
    df = modelform_factory(UserSSHKey, form=UserSSHKeyForm)

    if request.method == 'POST':
        formset = df(request.POST, request.FILES,
                          instance=request.user.ssh_keys.filter(id=key).first())

        if formset.is_valid():
            obj = formset.save(commit=False)
            obj.added_by = request.user
            obj.user = request.user
            obj.save()

            return HttpResponseRedirect(reverse('manager.views.profile_sshkeys'))
    else:
        formset = df(instance=request.user.ssh_keys.filter(id=key).first())

    return render_to_response('profile/ssh_keys_edit.html',
                              {
            'formset': formset
        },
                              context_instance=RequestContext(request))

@login_required
def profile_sshkeys_delete(request, key):
    key = request.user.ssh_keys.filter(id=key).first()

    if request.method == 'POST' and 'confirmation' in request.POST and request.POST['confirmation'] == 'yes':
        key.delete()
        return HttpResponseRedirect(reverse('manager.views.profile_sshkeys'))

    return render_to_response('profile/ssh_keys_delete.html',
                              {
            'account': account,
            'key': key
        },
                              context_instance=RequestContext(request))
