angular.module('trackinglog', ['ui.bootstrap','ui.utils','ui.router','ngResource', 'ngAnimate']);

angular.module('trackinglog').config(function($stateProvider, $urlRouterProvider) {

    $stateProvider.state('home', {
        url: '/home',
        templateUrl: 'partial/home/home.html'
    });
    $stateProvider.state('ua', {
        url: '/ua/:uaKey',
        templateUrl: 'partial/ua/ua.html'
    });
    /* Add New States Above */
    $urlRouterProvider.otherwise('/home');

});

angular.module('trackinglog').run(function($rootScope) {

    $rootScope.safeApply = function(fn) {
        var phase = $rootScope.$$phase;
        if (phase === '$apply' || phase === '$digest') {
            if (fn && (typeof(fn) === 'function')) {
                fn();
            }
        } else {
            this.$apply(fn);
        }
    };

});
