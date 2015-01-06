if (!$) {
    $ = django.jQuery;
}


$(document).ready(function() {
    //ace
    $('body').append('<script src="/static/ace/ace.js" type="text/javascript" charset="utf-8"></script>');

    var useAce = function(textarea_id) {
        var ace_id = 'ace_'+textarea_id;
        $('#' + textarea_id).after('<pre id="'+ace_id+'" style="height:400px;max-width:700px"></pre>');
        var editor = ace.edit(ace_id);
        editor.setTheme("ace/theme/chrome");
        editor.setFontSize("12px")
        editor.getSession().setMode("ace/mode/markdown");

        var textarea = $('#' + textarea_id).hide();
        editor.getSession().setValue(textarea.val());
        editor.getSession().on('change', function(){
        textarea.val(editor.getSession().getValue());
        });
    }

    if ($('#id_nginx_config').length) {
        useAce('id_nginx_config');
        useAce('id_apache_config');
    }
});
