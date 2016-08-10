function getWebsocketHost() {
	if (location.protocol == 'https:') {
		return 'wss://'+location.host;
	}else{
		return 'ws://'+location.host;
	}
}

function AccountShell(account, shell, element) {
	this.term = new Terminal({
      cols: 80,
      rows: 24
    });

    this.term.open(element);

    this.element = element;

    var url = getWebsocketHost() + '/api/v1/accounts/'+account+'/shell?env='+shell;
	this.websocket = new WebSocket(url);
	var that=this;
	this.websocket.onopen = function() {
	  console.log('shell opened');

	  element.className = 'active';

	  that.term.attach(that.websocket);
	}
}

AccountShell.prototype.stop = function() {
	this.websocket.close();
	this.term.destroy();

	this.element.className = '';

	while (this.term.element.children.length) {
    	this.term.element.removeChild(this.term.element.children[0]);
  	}
}
