angular.module('trackinglog').factory('UserAgents',function($resource) {
	var listUrl = '/api/useragents',
		entityUrl = listUrl + '/:key',
		userAgents = $resource(entityUrl, {
			key: '@key'
		}, {
			query: {method:'GET', url:listUrl, isArray:true}
		});
	return userAgents;
});
