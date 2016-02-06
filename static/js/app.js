
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
					values: [],
					key: 'user'
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
					values: [],
					key: 'memory',
					area: true
				}
			]
		};

		var addCpuData = function(cpu, timestamp){
			$scope.cpu.data[0].values.push({x:timestamp, y:cpu.user});
		};

		var addMemoryData = function(memory, timestamp){
			var memoryUsage = (memory.total-memory.free) / memory.total;
			$scope.memory.data[0].values.push({x: timestamp, y: memoryUsage});
		};

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
