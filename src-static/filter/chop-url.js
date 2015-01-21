angular.module('trackinglog').filter('chopUrl', function() {
	return function(input,arg) {
		return input.replace(/https?:\/\/[^\/]+\//, '/');
	};
});
