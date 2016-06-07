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
    elif action == 'sync-users':
        logs = actions.sync_accounts()
    elif action == 'sync-databases':
        logs = actions.sync_databases()

    data = json.dumps(logs.logs)
    return HttpResponse(data, content_type='application/json')
