from django.contrib.auth.decorators import login_required
from django.core.urlresolvers import reverse
from django.forms.models import modelform_factory
from django.http.response import HttpResponse, HttpResponseRedirect
from django.shortcuts import render_to_response
from django.template import RequestContext
import docker
from manager import actions, docker_api
from manager.forms import AppForm, DomainForm, DatabaseForm
from manager.models import App, Account, Domain, Database


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
def account_apps_edit(request, name, app=None):
    account = Account.objects.filter(name=name).first()
    af = modelform_factory(App, form=AppForm)

    if request.method == 'POST':
        formset = af(request.POST, request.FILES,
                          instance=account.apps.filter(id=app).first())

        if formset.is_valid():
            obj = formset.save(commit=False)
            obj.account = account
            obj.added_by = request.user
            obj.save()

            return HttpResponseRedirect(reverse('manager.views.account_apps', kwargs={'name': account.name}))
    else:
        formset = af(instance=account.apps.filter(id=app).first())

    return render_to_response('account/apps_edit.html',
                              {
            'account': account,
            'formset': formset
        },
                              context_instance=RequestContext(request))

@login_required
def account_apps(request, name):
    account = Account.objects.filter(name=name).first()
    apps = account.apps.all()

    containers = docker_api.cli.containers()
    mapping = {}
    for c in containers:
        try:
            name = c['Names'][0].replace('/', '')

            mapping[name] = c['Status']
        except Exception as e:
            print(e)

    print mapping

    for app in apps:
        if app.container_name() in mapping:
            app.status = mapping[app.container_name()]

            if 'Up ' in app.status:
                app.up = True

    return render_to_response('account/apps.html',
                              {
            'account': account,
            'apps': apps,
        },
                              context_instance=RequestContext(request))

@login_required
def containers(request):
    apps = App.objects.all()

    containers = docker_api.cli.containers()
    mapping = {}
    for c in containers:
        try:
            name = c['Names'][0].replace('/', '')

            mapping[name] = c['Status']
        except Exception as e:
            print(e)

    print mapping

    for app in apps:
        if app.container_name() in mapping:
            app.status = mapping[app.container_name()]

            if 'Up ' in app.status:
                app.up = True

    return render_to_response('containers.html',
                              {
            'account': account,
            'apps': apps,
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
def account_domains_edit(request, name, domain=None):
    account = Account.objects.filter(name=name).first()
    df = modelform_factory(Domain, form=DomainForm)
    domain = account.domains.filter(id=domain).first()

    if request.method == 'POST':
        formset = df(request.POST, request.FILES,
                          instance=domain)

        if formset.is_valid():
            obj = formset.save(commit=False)
            obj.account = account
            obj.added_by = request.user
            obj.save()

            return HttpResponseRedirect(reverse('manager.views.account_domains', kwargs={'name': account.name}))
    else:
        formset = df(instance=domain)

    return render_to_response('account/domains_edit.html',
                              {
            'account': account,
            'formset': formset,
            'domain': domain
        },
                              context_instance=RequestContext(request))

@login_required
def account_domains_delete(request, name, domain):
    account = Account.objects.filter(name=name).first()

    domain = account.domains.filter(id=domain).first()

    if request.method == 'POST' and 'confirmation' in request.POST and request.POST['confirmation'] == 'yes':
        domain.delete()
        return HttpResponseRedirect(reverse('manager.views.account_domains', kwargs={'name': account.name}))

    return render_to_response('account/domains_delete.html',
                              {
            'account': account,
            'domain': domain
        },
                              context_instance=RequestContext(request))


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
def update_nginx_config(request):
    logs = actions.update_nginx_config()

    return render_to_response('sync_action.html',
        {
            'logs': logs.logs,
            'action': 'Update nginx config'
        })
