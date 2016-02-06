
var app = angular.module('Golem', []);

app.controller('GolemController', ['$scope',
	function($scope){

		$scope.cpu = {
			user: -1,
			nice: -1,
			system: -1,
			idle: -1,
			iowait: -1
		};

		$scope.memory = {
			total: -1,
			free: -1
		};

		$scope.network = {
			interfaces: []	
		};

		$scope.disk = {
			disks: []	
		};

		var socket = new SockJS(window.location.origin + '/ws');

		socket.onopen = function(){
			console.log('connection open');
		};

		socket.onmessage = function(e){
			var systemInfo = JSON.parse(e.data);
			$scope.cpu = systemInfo.cpu;
			$scope.memory = systemInfo.memory;
			$scope.network = systemInfo.network;
			$scope.disk = systemInfo.disk;
			$scope.$apply();
		};

		socket.onclose = function(){
			console.log('connection closed');
		};
	}
]);
