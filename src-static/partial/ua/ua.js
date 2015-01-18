angular.module('trackinglog').controller('UaCtrl',function($scope, $stateParams, UserAgents, TrackingLogs){

  $scope.uaKey = $stateParams.uaKey;

  $scope.userAgent = UserAgents.get({key: $scope.uaKey});

  TrackingLogs.query({userAgentKey: $scope.uaKey}).$promise.then(function(logs) {
    $scope.trackingLogs = _.map(logs, function(log) {
      log.createdAt = Date.parse(log.createdAt);
      log.updatedAt = Date.parse(log.updatedAt);
      return log;
    });
  });
});
