function getWebsocketHost() {
	if (location.protocol == 'https') {
		return 'wss://'+location.host;
	}else{
		return 'ws://'+location.host;
	}
}


var ws = new ReconnectingWebSocket(getWebsocketHost() + '/ws/', null, {debug: true, reconnectInterval: 3000});


ws.onopen = function() {
	console.log('open')
}

ws.onmessage = function(msg) {
	console.log('message '+msg)
}

function AccountShell(account, shell, element) {
	this.term = new Terminal({
      cols: 80,
      rows: 24,
      screenKeys: true
    });

    this.term.open(element);

    var url = getWebsocketHost() + '/api/v1/account/'+account+'/shell?env='+shell;
	this.websocket = new WebSocket(url);
	var that=this;
	this.websocket.onopen = function() {
	  console.log('shell opened');

	  that.term.on('data', function(data) {
	    that.websocket.send(data);
	  });
	}

	this.websocket.onmessage = function(data) {
	  that.term.write(data.data);
	}
}

AccountShell.prototype.stop = function() {
	this.websocket.close();
	this.term.destroy();
}
