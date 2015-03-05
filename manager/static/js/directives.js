var drcs = angular.module('managerDirectives', []);

drcs.directive('autoResizingIframe', function() {
    return {
        // Restrict it to be an attribute in this case
        restrict: 'A',
        // responsible for registering DOM listeners as well as updating the DOM
        link: function(scope, element, attrs) {
            console.log('yo');
            console.log($('iframe.auto-height'));
            $('iframe.auto-height').iframeAutoHeight({
                minHeight: 800,
                callback: function(x) {
                	console.log(x)
                	$('iframe.auto-height').css('visibility', 'visible');
                }
            });
        }
    };
});
