
var ws = new ReconnectingWebSocket('ws://'+location.host+'/ws/', null, {debug: true, reconnectInterval: 3000});


ws.onopen = function() {
	console.log('open')
}

ws.onmessage = function(msg) {
	console.log('message '+msg)
}
