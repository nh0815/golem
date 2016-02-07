
var app = angular.module('Golem', ['nvd3']);

app.controller('GolemController', ['$scope',
	function($scope){

		var DATA_LIMIT = 100;

		var previousCpu = {
			user: 0,
			nice: 0,
			system: 0,
			idle: 0,
			iowait: 0,
			irq: 0,
			irqSoft: 0,
			steal: 0,
			guest: 0,
			guestNice: 0
		};

		var previousNetwork = null;

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

		$scope.disk = {};

		$scope.network = {
			options: {
				chart: {
					type: 'stackedAreaChart',
					height: 450
				},
				title: {
					enable: true,
					text: 'network'
				}
			},
			data: [
				{
					values: [],
					key: 'transmit'
				},
				{
					values: [],
					key: 'receive'
				}
			]
		};

		var addCpuData = function(cpu, timestamp){
			$scope.cpu.data[0].values.push({x:timestamp, y:cpu});
			var length = $scope.cpu.data[0].values.length;
			if(length > DATA_LIMIT){
				$scope.cpu.data[0].values = $scope.cpu.data[0].values.slice(length - DATA_LIMIT, length);
			}
		};

		var getCpuUsage = function(cpu){
			var previousIdle = previousCpu.idle + previousCpu.iowait;
			var idle = cpu.idle + cpu.iowait;

			var previousNonIdle = previousCpu.user + previousCpu.nice + previousCpu.system + previousCpu.irq + previousCpu.irqSoft + previousCpu.steal;
			var nonIdle = cpu.user + cpu.nice + cpu.system + cpu.irq + cpu.irqSoft + cpu.steal;

			var previousTotal = previousIdle + previousNonIdle;
			var total = idle + nonIdle;

			var totalDelta = total - previousTotal;
			var idleDelta = idle - previousIdle;

			return Math.round(100 * (totalDelta - idleDelta) / totalDelta);

		};

		var addMemoryData = function(memory, timestamp){
			var memoryUsage = Math.round(100 * (memory.total-memory.free) / memory.total);
			$scope.memory.data[0].values.push({x: timestamp, y: memoryUsage});
			var length = $scope.memory.data[0].values.length;
			if(length > DATA_LIMIT){
				$scope.memory.data[0].values = $scope.memory.data[0].values.slice(length - DATA_LIMIT, length);
			}
		};

		var addNetworkData = function(network, timestamp){
			if(previousNetwork){
				$scope.network.data[0].values.push({x:timestamp, y:network.transmitBytes-previousNetwork.transmitBytes});
				$scope.network.data[1].values.push({x:timestamp, y:network.receiveBytes-previousNetwork.receiveBytes});
				var length1 = $scope.network.data[0].values.length;
				if(length1 > DATA_LIMIT){
					$scope.network.data[0].values = $scope.network.data[0].values.slice(length1 - DATA_LIMIT, length1);
				}
				var length2 = $scope.network.data[1].values.length;
				if(length2 > DATA_LIMIT){
					$scope.network.data[1].values = $scope.network.data[1].values.slice(length2 - DATA_LIMIT, length2);
				}
			}
			previousNetwork = network;
		};

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
			addCpuData(getCpuUsage(systemInfo.cpu), time);
			addMemoryData(systemInfo.memory, time);
			addNetworkData(systemInfo.network.interfaces[0], time);
			$scope.$apply();
		};

		socket.onclose = function(){
			console.log('connection closed');
		};

		/* end socket code */
	}
]);
