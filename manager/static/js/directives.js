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

drcs.directive('forceReload',function($location,$route){
    return function(scope, element, attrs) {
        element.bind('click',function(){
            if(element[0] && element[0].href && element[0].href === $location.absUrl()){
                $route.reload();
            }
        });
    }   
});

drcs.filter('unsafe', function($sce) {
    return function(val) {
        return $sce.trustAsHtml(val);
    };
});