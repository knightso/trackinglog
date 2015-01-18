angular.module('trackinglog').factory('TrackingLogs',function($resource) {
	var listUrl = '/api/useragents/:userAgentKey/trackinglogs',
		entityUrl = listUrl + '/:key',
		trackingLogs = $resource(entityUrl, {
			key: '@key'
		}, {
			query: {method:'GET', url:listUrl, isArray:true}
		});
	return trackingLogs;
});
