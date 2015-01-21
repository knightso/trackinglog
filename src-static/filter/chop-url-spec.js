describe('chopUrl', function() {

	beforeEach(module('trackinglog'));

	it('should ...', inject(function($filter) {

        var filter = $filter('chopUrl');

		expect(filter('input')).toEqual('output');

	}));

});