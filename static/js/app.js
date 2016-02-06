
var app = angular.module('Golem', ['nvd3']);

app.controller('GolemController', ['$scope',
	function($scope){

		/* chart code */

		$scope.cpu = {
			options: {
				chart: {
					type: 'lineChart',
					height: 450,
				},
				title: {
					enable: true,
					text: 'cpu'
				}
			},
			data: [
				{
					values: [
						{x:0, y:0},
						{x:1, y:1},
						{x:2,y:2}
					],
					key: '0-5'
				}
			]
		};

		$scope.memory = {
			options: {
				chart: {
					type: 'lineChart',
					height: 450
				},
				title: {
					enable: true,
					text: 'memory'
				}
			},
			data: [
				{
					values: [
						{x: 100, y:-100},
						{x: 50, y:60},
						{x: -10, y:99}
					],
					key: 'numbers',
					area: true
				}
			]
		};

		var addCpuData = function(cpu, timestamp){};

		var addMemoryData = function(memory, timestamp){};

		var addNetworkData = function(network, timestamp){};

		var addDiskData = function(data, timestamp){};

		/* end chart code */

		/* socket code */

		var socket = new SockJS(window.location.origin + '/ws');

		socket.onopen = function(){
			console.log('connection open');
		};

		socket.onmessage = function(e){
			var systemInfo = JSON.parse(e.data);
			var time = Date.parse(systemInfo.timestamp);
			addCpuData(systemInfo.cpu, time);
			addMemoryData(systemInfo.memory, time);
			$scope.$apply();
		};

		socket.onclose = function(){
			console.log('connection closed');
		};

		/* end socket code */
	}
]);
