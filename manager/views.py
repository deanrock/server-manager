from django.contrib.auth.decorators import login_required
from django.core.serializers import serialize
import json
from django.core.urlresolvers import reverse
from django.forms.models import modelform_factory, model_to_dict
from django.http.response import HttpResponse, HttpResponseRedirect, HttpResponseNotFound
from django.shortcuts import render_to_response
from django.template import RequestContext
import docker
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
def rebuild_base_image(request):
    return render_to_response('sync_action.html',
        {
            'action': 'Rebuilding base image',
            'action_url': reverse('manager.views.action_ajax', kwargs={'action': 'rebuild-base-image'}),
        })

@login_required
def sync_users(request):
    return render_to_response('sync_action.html',
        {
            'action': 'Syncing users',
            'action_url': reverse('manager.views.action_ajax', kwargs={'action': 'sync-users'}),
        })


@login_required
def sync_databases(request):
    return render_to_response('sync_action.html',
        {
            'action': 'Syncing databases',
            'action_url': reverse('manager.views.action_ajax', kwargs={'action': 'sync-databases'}),
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
    app = account.apps.filter(id=app).first()
    af = modelform_factory(App, form=AppForm)

    if request.method == 'POST':
        formset = af(request.POST, request.FILES,
                          instance=app)

        if formset.is_valid():
            obj = formset.save(commit=False)
            obj.account = account
            obj.added_by = request.user
            obj.save()

            vars = obj.image.variables.all()

            for v in vars:
                prev = obj.variables.filter(name=v.name).first()

                if prev:
                    prev.delete()

                field = 'id_variable_%s' % v.name
                if field in request.POST and request.POST[field].rstrip() != '':
                    prev = AppImageVariable()
                    prev.app = obj
                    prev.name = v.name
                    prev.value = request.POST[field]
                    prev.save()

            return HttpResponseRedirect(reverse('manager.views.account_apps', kwargs={'name': account.name}))
    else:
        formset = af(instance=app)

    variables = {}
    if app:
        for v in app.variables.all():
            variables[v.name] = v.value

    return render_to_response('account/apps_edit.html',
                              {
            'account': account,
            'formset': formset,
            'variables': json.dumps(variables),
        },
                              context_instance=RequestContext(request))


@login_required
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
def account_apps_action(request, name, app, action):
    account = Account.objects.filter(name=name).first()
    app = account.apps.filter(id=app).first()

    action_text = ''

    if action == 'redeploy':
        action_text = 'Redeploying %s for %s ...' % (app.name, account.name)
    elif action == 'stop':
        action_text = 'Stopping %s for %s ...' % (app.name, account.name)
    elif action == 'start':
        action_text = 'Starting %s for %s ...' % (app.name, account.name)
    else:
        return HttpResponseNotFound()

    return render_to_response("account/apps_action.html",
        {
            'account': account,
            'action_text': action_text,
            'action_url': reverse('manager.views.account_apps_action_ajax', kwargs={'name': name, 'app': app.id, 'action': action}),
            'app': app,
        })


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

    #ports
    apps = account.apps.all()

    ports=[]

    for app in apps:
        for port in app.image.ports.all():
            ports.append('#%s_%d_ip#' % (app.name, port.port))

    #form
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
            'domain': domain,
            'ports': ports,
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
    return render_to_response('sync_action.html',
        {
            'action': 'Updating nginx and apache config ...',
            'action_url': reverse('manager.views.action_ajax', kwargs={'action': 'update-nginx-config'}),
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


@login_required
def api_account_app(request, name, id):
    account = Account.objects.filter(name=name).first()
    app = account.apps.filter(id=id).first()

    dict = model_to_dict(app)
    dict['variables'] = []

    for v in app.variables.all():
        dict['variables'].append(model_to_dict(v))

    data = json.dumps(dict)
    return HttpResponse(data, content_type='application/json')


@login_required
def api_image(request, id):
    image = Image.objects.filter(id=id).first()

    dict = model_to_dict(image)
    dict['variables'] = []

    for v in image.variables.all():
        dict['variables'].append(model_to_dict(v))

    data = json.dumps(dict)
    return HttpResponse(data, content_type='application/json')

