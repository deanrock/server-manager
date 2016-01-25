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
