angular.module('trackinglog').controller('HomeCtrl',function($scope, UserAgents){

  $scope.userAgents = UserAgents.query();

  UserAgents.query().$promise.then(function(uas) {
    $scope.userAgents = _.map(uas, function(ua) {
      ua.createdAt = Date.parse(ua.createdAt);
      ua.updatedAt = Date.parse(ua.updatedAt);
      return ua;
    });
  });

});
