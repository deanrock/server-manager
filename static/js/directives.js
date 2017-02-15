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

drcs.directive('aceEditor', function() {
    return {
        restrict: 'A',
        scope: {
            addEditor: "&"
        },
        link: function(scope, element, attrs) {
            var textarea_id = element[0].id;
            var ace_id = 'ace_'+textarea_id;

            $('#' + textarea_id).after('<pre id="'+ace_id+'" style="height:400px;width:100%"></pre>');
            var editor = ace.edit(ace_id);
            editor.setTheme("ace/theme/chrome");
            editor.setFontSize("12px")
            editor.getSession().setMode("ace/mode/markdown");
            editor.setAutoScrollEditorIntoView(true);
            editor.setOption("maxLines", 3000);
            editor.setOption("minLines", 20);

            scope.addEditor({d: {name: element[0].id, editor: editor}});

            var textarea = $('#' + textarea_id).hide();
            editor.getSession().setValue(textarea.val());
            editor.getSession().on('change', function(){
            textarea.val(editor.getSession().getValue());
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

drcs.directive('formTextInput', function() {
    return {
        restrict: 'E',
        scope: {
            model: '=',
            label: '@',
            errors: '=',
            name: '@model'
        },
        templateUrl: 'static/templates/directives/formTextInput.html'
    }
});

drcs.directive('formCheckbox', function() {
    return {
        restrict: 'E',
        scope: {
            model: '=',
            label: '@',
            errors: '=',
            name: '@model'
        },
        templateUrl: 'static/templates/directives/formCheckbox.html'
    }
});

drcs.directive('userForm', function() {
    return {
        restrict: 'E',
        scope: {
            user: '=',
            errors: '=',
            save: '='
        },
        templateUrl: 'static/templates/directives/userForm.html'
    }
});

drcs.filter('unsafe', function($sce) {
    return function(val) {
        return $sce.trustAsHtml(val);
    };
});
