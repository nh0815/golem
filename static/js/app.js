
var app = angular.module('Golem', [
	'nvd3'
]);

app.controller('GolemController', ['$scope',
	function($scope){

		/* chart code */

		$scope.cpuChart = {
			data: [{
				values: [],
				key: 'User'
			}, {
				values: [],
				key: 'Nice'
			}, {
				values: [],
				key: 'Syste,'
			}],
			options: {
				chart: {
					type: 'lineChart',
					height: 450
				},
				title: {
					enable: true,
					text: 'cpu'
				}
			}
		};

		var addCpuData = function(cpu, timestamp){
			var user = {
				x: timestamp,
				y: cpu.user
			};

			var nice = {
				x: timestamp,
				y: cpu.nice
			};

			var system = {
				x: timestamp,
				y: cpu.system
			};
			$scope.cpuChart.data[0].values.push(user);
			$scope.cpuChart.data[1].values.push(nice);
			$scope.cpuChart.data[2].values.push(system);
		};

		var addMemoryData = function(memory){
			
		};

		var addNetworkData = function(network){
			
		};

		var addDiskData = function(data){
			
		};

		/* end chart code */

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

		/* socket code */

		var socket = new SockJS(window.location.origin + '/ws');

		socket.onopen = function(){
			console.log('connection open');
		};

		socket.onmessage = function(e){
			var systemInfo = JSON.parse(e.data);
			$scope.cpu = systemInfo.cpu;
			addCpuData(systemInfo.cpu, Date.parse(systemInfo.timestamp));
			$scope.memory = systemInfo.memory;
			$scope.network = systemInfo.network;
			$scope.disk = systemInfo.disk;
			$scope.$apply();
		};

		socket.onclose = function(){
			console.log('connection closed');
		};

		/* end socket code */
	}
]);
