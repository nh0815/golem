
var socket = new SockJS(window.location.origin + '/ws');

socket.onopen = function(){
	console.log('connection open');
};

socket.onmessage = function(e){
	console.log(e);
};

socket.onclose = function(){
	console.log('connection closed');
}
